package arcauth

import (
	"fmt"
	"testing"
)

func TestNewClient(t *testing.T) {
	arcAuthClient, err := New("http://boot2docker:3000", "v1", "demo-app", "WKZd$&vk&$I7VCa@ueVl1sMMj7iFW315")
	if err != nil {
		t.Errorf("unexpected error %v", err)		
	}

	body, error := arcAuthClient.Auth("FakeDemoToken")

	fmt.Print("body is ", body)
	fmt.Print("error is ", error)

}