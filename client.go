package arcauth

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strings"
)

const AdmiralTokenHeader = "X-Admiral-Token"

type ArcAuthClient struct {
    Host       string
    User       string
    Pass       string
    HttpClient *http.Client
}


/**
 * New constructs a new ArcAuthClient for communication with an arc-auth-server
 * server - the root FQDN of the arc-auth-server (e.g. https://arc-auth.ext.nile.works)
 * user - the user to use in BasicAuth when making requests for token authentication (this go client itself must be authenticated!)
 * pass - the password for the user when sending BasicAuth
 */
func New(server, user string, pass string) (*ArcAuthClient, error) {
    log.Printf("Constructing new arc-auth client")

    if server == "" {
        return nil, fmt.Errorf("Arc Auth Server cannot be empty, provide FQDN value like 'http://your.service.com'")
    }
    if user == "" {
        return nil, fmt.Errorf("You must provide a user to authenticate against the arc-auth server")
    }
    if pass == "" {
        return nil, fmt.Errorf("You must provide a password to authenticate against the arc-auth server")
    }

    return &ArcAuthClient {
        Host:   fmt.Sprintf("%s/api/v1", server),
        User:   user,
        Pass:   pass,
        HttpClient: &http.Client{ },
    }, nil
}

/**
 * Auth makes a request to the arc-auth-server's ".../auth" endpoint with the
 * token string set as the Header associated to the AdmiralTokenHeader key
 *
 * On a succesful connection, the raw JSON from the server is returned by this method.  Note that an invalid
 * token will still be "successful" and a 204/Empty Content from the server will result in an empty string
 * being returned to the caller
 */
func (this *ArcAuthClient) Auth(token string) (string, error) {
    request, err := http.NewRequest("GET", fmt.Sprintf("%s/auth", this.Host), nil)
    request.SetBasicAuth(this.User, this.Pass)
    request.Header.Set(AdmiralTokenHeader, token)

    response, err := this.HttpClient.Do(request)

    defer response.Body.Close()

    if err != nil {
        log.Printf("Error : %s", err)
        return "", err
    } 

    if (response.StatusCode != http.StatusOK && response.StatusCode != http.StatusNoContent) {
        log.Printf("Got response code %s when authenticating token %s", response.StatusCode, this.Mask(token))
        return "", &ErrorResponse{Code: response.StatusCode, Message: "Non-20X response code"}
    }
    
    body, error := ioutil.ReadAll(response.Body)
    return string(body), error    
}

type ErrorResponse struct {
    Code    int     `json:"code"`
    Message string  `json:"message"`
}

func (e *ErrorResponse) Error() string {
    return fmt.Sprintf("HTTP Code %s | %s", e.Code, e.Message)
}

/**
 * Invokes this.Mask() with the maskChar "*"
 */
func (this *ArcAuthClient) Mask(plaintext string) string {
    return this.MaskWithChar(plaintext, "*")
}

/**
 * Mask keeps the first 2 characters and the last character of the plaintext input string in tact but replaces
 * everything else with the mask character
 *
 * Plaintext strings less than or equal to length 5 are masked with 5 "mask" characters
 *
 * The empty input string is not masked at all and an empty string is returned
 */
func (this *ArcAuthClient) MaskWithChar(plaintext, maskChar string) string {
    if plaintext == "" {
        return ""
    }
    if len(plaintext) <= 5 {
        return strings.Repeat(maskChar, 5)
    }
    chars := strings.Split(plaintext, "")
    
    var buffer bytes.Buffer
    buffer.WriteString(chars[0])
    buffer.WriteString(strings.Repeat(maskChar, len(plaintext) - 3))
    buffer.WriteString(chars[len(plaintext) - 2])
    buffer.WriteString(chars[len(plaintext) - 1])
    return buffer.String()
}