# go-grpc

## Description
`go-grpc` is a sample Go project that implements gRPC services using [ConnectRPC](https://connectrpc.com/) for CRUD operations and streaming functionalities. This project includes a complete server implementation, with various interceptors for logging, authentication, and recovery, as well as a simple client CLI application to run all functionalities implemented.

>This repository and the code in it is meant as a source of good practices and an example for any implementations of this kind, to help anyone trying to build similar applications.

## Features
- CRUD Service: Create, Read, Update, and Delete operations.
- Stream Service: Uploading files and sending direct messages (bidi).
- Interceptors: Logging, Authentication, and Recovery.

## Installation
To install the project dependencies, run:
```bash
go mod download
```

## Usage
### Running the Server
To start the gRPC server, use:
```bash
go run server/cmd/server/main.go
```

### Client Commands
Run the following command to see the capabilities and usage of the client CLI:
```bash
go run client/main.go -h
```

### Unit/integration testing the Server implementation
A number of sample tests are part of the project for the server implementation.
These can be run with the following command:
```bash
go test -v ./...
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[GPLv3](https://www.gnu.org/licenses/)
