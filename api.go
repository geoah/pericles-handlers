package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/encoder"
	"github.com/twinj/uuid"
)

func GetPayloads(r *http.Request, enc encoder.Encoder, store Store) []byte {
	// Return all payload Objects
	return encoder.Must(enc.Encode(store.GetAll()))
}

func GetPayload(enc encoder.Encoder, store Store, parms martini.Params) (int, []byte) {
	// TODO Check for parms:id
	// Get payload Object from Store
	payload := store.Get(parms["id"])
	// TODO Check if payload exists
	return http.StatusOK, encoder.Must(enc.Encode(payload))
}

func AddPayload(w http.ResponseWriter, r *http.Request, enc encoder.Encoder, store Store) (int, []byte) {
	// Decode JSON from body and create a Payload Object from it
	decoder := json.NewDecoder(r.Body)
	var payload Payload
	decoder.Decode(&payload)

	// Create and set a unique Id
	payload.Id = uuid.Formatter(uuid.NewV4(), uuid.CleanHyphen)

	// Reset the status
	payload.Status = "pending"

	// TODO Check for errors
	id, err := store.Add(&payload)
	if err != nil {
		// TODO Provide more information on errors
		return http.StatusConflict, encoder.Must(enc.Encode(NewError(500, "Something went wrong.")))
	}

	// As a test we create a goroutine to perform an async task
	go func() {
		// Just sleep for 30 seconds
		cmd := exec.Command("sleep", "30")
		// Once we start the execution set the payload's status to working
		err := cmd.Start()
		payload.Status = "working"
		if err != nil {
			// If there are any errors just push them to the error:start status
			payload.Status = fmt.Printf("error:start (%v)", err)
		} else {
			// Wait until the execution is complete
			err = cmd.Wait()
			// Once completed
			if err != nil {
				// If there are any errors just push them to the error status
				payload.Status = fmt.Printf("error (%v)", err)
			} else {
				// And mark the payload's status as finished
				payload.Status = "finished"
				store.Update(&payload)
			}
		}
	}()

	// Once payload is stored, redirect the user to the newly create object's URI
	w.Header().Set("Location", fmt.Sprintf("/payloads/%s", id))
	return http.StatusCreated, encoder.Must(enc.Encode(payload))
}
