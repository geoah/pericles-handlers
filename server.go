package main

import (
	"log"
	"net/http"
	"github.com/codegangsta/martini-contrib/encoder"

	"github.com/codegangsta/martini"
)

var m *martini.Martini

func init() {
	// Martini instance
	m = martini.New()

	// Setup middleware
	m.Use(martini.Recovery())
	m.Use(martini.Logger())

	// Setup routes
	r := martini.NewRouter()
	r.Get(`/payloads`, GetPayloads)
	r.Get(`/payloads/:id`, GetPayload)
	r.Post(`/payloads`, AddPayload)

	m.Use(func(c martini.Context, w http.ResponseWriter) {
		// Inject JSON Encoder
		c.MapTo(encoder.JsonEncoder{}, (*encoder.Encoder)(nil))
		// Force Content-Type
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	})
	// Inject database
	m.MapTo(store, (*Store)(nil))
	// Add the router action
	m.Action(r.Handle)
}

func main() {
	// Startup HTTP server
	if err := http.ListenAndServe(":8000", m); err != nil {
		log.Fatal(err)
	}
}
