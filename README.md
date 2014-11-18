# Handlers Proof of Concept in Go (golang)

Creates an HTTP server to manage payloads.
Each payload is a work to be executed on the server and is assigned a unique id by which its status can be tracked.

Once a new new payload is created via `POST /payloads` its id will be returned and its status will be set to `working`.  
The status of the payload can be tracked via a `GET /payloads/:id` or from the list of all payloads of `GET /payloads`.  
Available statuses are `pending`, `working`, `finished` or `error:*` which will contain the error message.  

Payloads are currently stored in-memory so if the server is re-started no old payloads will be present or continue work.

### Available endpoints

#### POST /payloads

Creates and executes a new payload.  
Payload's structure is as follows:

    {
        "Id": "...", // Automatically generated UUID.v4
        "Cmd": "sleep" // Command to be executed
        "Params": ["30"], // Command arguments
        "Wid": "random-wid", // Not used
        "Wiid": "random-wiid", // Not used
        "Xid": "random-xid", // Not used
        "Xuri": "random-xuri", // Not used
        "Status": "..." // Provides current status of command; Read only
    }

Currently only the `Cmd` and `Params` are required to execute a new command.  

##### Sample commands

* Sleep
  `curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -d '{ "Cmd": "/bin/sleep", "Params": ["30"] }' http://localhost:8000/payloads`

* Ping
  `curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -d '{ "cmd": "/sbin/ping", "params": ["-c", "5", "google.com"] }' http://localhost:8000/payloads`

## Installation

#### GET /payloads

Returns all payloads.

#### Get /payloads/:id

Returns single payload by its id.

## Installation

Requires Go 1.1+

    git clone ... pericles-handlers
    cd pericles-handlers

    go get github.com/codegangsta/martini
    go get github.com/codegangsta/martini-contrib/encoder
    go get github.com/twinj/uuid

## Running

    go build && ./pericles-handlers
