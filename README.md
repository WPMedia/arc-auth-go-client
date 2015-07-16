# arc-auth-go-client
A client for the arc-auth-server, written in Go

## Testing
Run `go test -v` to run the client tests; a couple of the tests will use a running arc-auth-server on localhost `http://boot2docker:3000` if it is running to do real end-to-end tests of the client code.  If the boot2docker instance isn't running those end-to-end tests are just skipped.