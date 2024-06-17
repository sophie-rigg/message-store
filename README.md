# message-store

## Description
This is a simple http server for storing messages.

## Usage
To start the server, run the following command:
```go run ./cmd/message-store/main.go```

The server will start on port 8080 by default.

## Endpoints

### POST /messages
This endpoint is used to upload a plain text message to the server and return json.
The message should be sent in the body of the request.
The response will be in the form of:
```json
{
    "id": "1"
}
```
The ID of the message will be included in the response. It increases by 1 for each message uploaded.

### GET /messages/{id}
This endpoint is used to retrieve a message from the server and return it in the body.
The ID of the message should be included in the URL and should be an integer.

## Testing
To run the tests, run the following command:
```go test ./...```

## Future Improvements
- Store messages in a database