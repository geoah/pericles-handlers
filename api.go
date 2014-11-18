package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/encoder"
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
	// TODO Check for errors
	id, err := store.Add(&payload)
	if err != nil {
		// TODO Provide more information on errors
		return http.StatusConflict, encoder.Must(enc.Encode(NewError(500, "Something went wrong.")))
	}
	// Once payload is stored, redirect the user to the newly create object's URI
	w.Header().Set("Location", fmt.Sprintf("/payloads/%s", id))
	return http.StatusCreated, encoder.Must(enc.Encode(payload))
}
