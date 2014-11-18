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

##### Request

    {
        cmd: "/sbin/ping",
        params: [ "-c", "5", "google.com" ],
    }

##### Response

    {
        id: "cada5f19-6ee3-4f26-81a9-4f0c4bc9ca3e",
        cmd: "/sbin/ping",
        params: [ "-c", "5", "google.com" ],
        status: "working"
    }

##### Sample calls

* Sleep for 30 seconds
  `curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -d '{ "Cmd": "/bin/sleep", "Params": ["30"] }' http://localhost:8000/payloads`

* Ping google.com 5 times
  `curl -X POST -H "Content-Type: application/json" -H "Cache-Control: no-cache" -d '{ "cmd": "/sbin/ping", "params": ["-c", "5", "google.com"] }' http://localhost:8000/payloads`

#### GET /payloads

Returns all payloads.

##### Response

    [
        {
            id: "cada5f19-6ee3-4f26-81a9-4f0c4bc9ca3e",
            cmd: "/sbin/ping",
            params: [ "-c", "5", "google.com" ],
            status: "working"
        }, {
            id: "08ca0c54-c65f-4cb1-8fdb-46fbb6445e84",
            cmd: "sleep",
            params: [ "30" ],
            status: "finished"
        }, {
            id: "edc72406-dfca-44ca-a1a3-c992a897491c",
            cmd: "/bin/sleep",
            params: [ "30" ],
            status: "finished"
        }, {
            id: "b52a37ac-665b-44f6-9956-2317bc0354a2",
            cmd: "/sbin/ping",
            params: [ "-something-wrong", "5" ],
            status: "error (exit status 64)"
        }
    ]

#### Get /payloads/:id

Returns single payload by its id.

##### Response

    {
        id: "cada5f19-6ee3-4f26-81a9-4f0c4bc9ca3e",
        cmd: "/sbin/ping",
        params: [ "-c", "5", "google.com" ],
        status: "working"
    }

## Installation

Requires Go 1.1+

    git clone ... pericles-handlers
    cd pericles-handlers

    go get github.com/codegangsta/martini
    go get github.com/codegangsta/martini-contrib/encoder
    go get github.com/twinj/uuid

## Running

    go build && ./pericles-handlers
