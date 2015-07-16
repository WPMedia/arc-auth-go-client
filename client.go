package arcauth

// TODO
// Create a transport with keep-alives and reuse for a long time
// Cache results
// Do a better job returning error strings based of HTTP response codes
// add logging
// mask the token in the logging

import (
	"fmt"
	"io/ioutil"
    "log"
    "net/http"
)

const DefaultApiVerion = "v1"
const AdmiralTokenHeader = "X-Admiral-Token"

type ArcAuthClient struct {
    Host  string
    User  string
    Pass  string
}

type ErrorResponse struct {
    Code    int     `json:"code"`
    Message string  `json:"mesage"`
}

/**
 * New constructs a new ArcAuthClient for communication with an arc-auth-server
 * server - the root FQDN of the arc-auth-server (e.g. https://arc-auth.ext.nile.works)
 * apiVersion - the version of the server-side API to communicate with.  The only supported option right now is "v1"
 * user - the user to use in BasicAuth when making requests for token authentication (this go client itself must be authenticated!)
 * pass - the password for the user when sending BasicAuth
 */
func New(server, apiVersion, user string, pass string) (*ArcAuthClient, error) {
	log.Printf("Constructing new arc-auth client")

	if server == "" {
		return nil, fmt.Errorf("Arc Auth Server cannot be empty, provide FQDN value like 'http://your.service.com'")
	}
	if apiVersion == "" {
		apiVersion = DefaultApiVerion
	}
	if user == "" {
		return nil, fmt.Errorf("You must provide a user to authenticate against the arc-auth server")
	}
	if pass == "" {
		return nil, fmt.Errorf("You must provide a password to authenticate against the arc-auth server")
	}

	return &ArcAuthClient {
		Host:	fmt.Sprintf("%s/api/%s", server, apiVersion),
		User:	user,
		Pass:	pass,
	}, nil
}

/**
 * Auth makes a request to the arc-auth-server's ".../auth" endpoint with the
 * token string set as the Header associated to the AdmiralTokenHeader key
 */
func (arcAuthClient *ArcAuthClient) Auth(token string) (string, error) {

	httpClient := &http.Client{ } // TODO move to the struct itself?
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/auth", arcAuthClient.Host), nil)
	req.SetBasicAuth(arcAuthClient.User, arcAuthClient.Pass)
	req.Header.Set(AdmiralTokenHeader, token)

	resp, err := httpClient.Do(req)

	if err != nil {
		fmt.Printf("Error : %s", err)
		return "", err
	} else if (resp.Status != "200 OK" && resp.Status != "204 No Content") {
		log.Printf("Got response code %s when authenticating token %s", resp.Status, token)
		return "", fmt.Errorf("Non-20X response code %s", resp.Status)
	} else {
		body, error := ioutil.ReadAll(resp.Body)
		return string(body), error
	}
}