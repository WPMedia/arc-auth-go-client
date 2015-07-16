package arcauth

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewClientNeedsServer(t *testing.T) {
	_, err := New("", "v1", "user", "pass")
	assert.Error(t, err)
}

func TestNewClientNeedsUser(t *testing.T) {
	_, err := New("whatever", "v1", "", "pass")
	assert.Error(t, err)
}

func TestNewClientNeedsPass(t *testing.T) {
	_, err := New("whatever", "v1", "user", "")
	assert.Error(t, err)
}

func TestClientWhenServerSendsGoodResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprint(w, "Hello, client")
    }

	testServer := httptest.NewServer(http.HandlerFunc(handler))
    defer testServer.Close()

    arcAuthClient, err := New(testServer.URL, "v1", "user", "pass")
    if err != nil {
		t.Errorf("unexpected error %v", err)		
	}

	body, error := arcAuthClient.Auth("FakeDemoToken")

	assert.Equal(t, "Hello, client", body)
	assert.NoError(t, error)
}

func TestClientWhenServerSendsBadResponse(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
        http.Error(w, "something failed", http.StatusInternalServerError)
    }

    testServer := httptest.NewServer(http.HandlerFunc(handler))
    defer testServer.Close()

	arcAuthClient, _ := New(testServer.URL, "v1", "user", "pass")
	_, responseErr := arcAuthClient.Auth("FakeDemoToken")

    assert.NotNil(t, responseErr, "We expect the client to propogate a server error to its caller")
}

func TestMask(t *testing.T) {
	arcAuthClient, _ := New("whatever", "v1", "user", "pass")

	assert.Equal(t, "", arcAuthClient.Mask(""))
	assert.Equal(t, "*****", arcAuthClient.Mask("a"))
	assert.Equal(t, "*****", arcAuthClient.Mask("foo"))
	assert.Equal(t, "*****", arcAuthClient.Mask("abcde"))
	assert.Equal(t, "a***ef", arcAuthClient.Mask("abcdef"))
}

/*
 * These boot2docker tests will only run if you have booted up a localhost version of the arc-auth application that's
 * available at the boot2docker:3000 port.
 * You should be able to do this by running "docker-compose build && docker-compose up" within the arc-auth/
 * project directory and as long as you have a line in your /etc/hosts file mapping the servername "boot2docker"
 * to the boot2docker ip.
 *
 * The username and password used in this test are one of the authorized "peers" in the arc-auth-server and
 * the assertion of a user "vaughant" showing up in the responseBody is based on the fixture data loaded into 
 * the arc-auth-server's fake MySQL database on startup.
 */
const localhostServer = "boot2docker:3000"

func TestClientWithGoodTokenAgainstBoot2DockerImage(t *testing.T) {
	responseBody, responseErr := runBoot2DockerTest(t, "FakeDemoToken")

	assert.Nil(t, responseErr, "Unexpected responseErr %s", responseErr)
	assert.Contains(t, responseBody, "vaughant") // that should be the user associated with the FakeDemoToken
}

func TestClientWithBadTokenAgainstBoot2DockerImage(t *testing.T) {
	responseBody, responseErr := runBoot2DockerTest(t, "No Such Token")	
	assert.Nil(t, responseErr, "Unexpected responseErr %s", responseErr)
	assert.Equal(t, responseBody, "")
}

func runBoot2DockerTest(t *testing.T, token string) (string, error) {
	// See if the boot2docker image is up where we expect, skip the test execution if it isn't reachable
	// Hint: run "go test -v" to see whether the tests are skipped
	_, connErr := net.Dial("tcp", localhostServer)
	if (connErr != nil) {
		t.Skip("This test won't run unless it can reach ", localhostServer)
	}
	arcAuthClient, arcAuthClientErr := New("http://" + localhostServer, "v1", "demo-app", "WKZd$&vk&$I7VCa@ueVl1sMMj7iFW315")
	assert.Nil(t, arcAuthClientErr, "Unexpected error constructing an arcAuthClient")

	return arcAuthClient.Auth(token)
}