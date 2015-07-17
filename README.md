# arc-auth-go-client
A client for the arc-auth-server, written in Go

## Using
Import the client into your code:

```
import (
	"github.com/WPMedia/arc-auth-go-client"
)
```

Create a new client:

```
arcAuthClient, err := New("https://the-arc-auth.server.url", "user", "pass")
```

Use the client to get the authorization JSON for a token:

```
json, err := arcAuthClient.Auth("FakeDemoToken")
```    


## Testing
Run `godep fo test -v` to run the client tests; a couple of the tests will use a running arc-auth-server on localhost `http://boot2docker:3000` if it is running to do real end-to-end tests of the client code.  If the boot2docker instance isn't running those end-to-end tests are just skipped.