package gogpt

import (
	"encoding/json"
	"time"
)

const sessionExpirationTimeLayout = "2006-01-02T15:04:05.999Z"

type User struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Image   string   `json:"image"`
	Picture string   `json:"picture"`
	Groups  []string `json:"groups"`
}

type Session struct {
	User        User   `json:"user"`
	Expires     string `json:"expires"`
	AccessToken string `json:"accessToken"`
}

// unmarshalGPTSessionResponseJSON returns a pointer to a Session from the jsonData passed in parameters
// it returns an error if something goes wrong while unmarshalling the json
func unmarshalGPTSessionResponseJSON(jsonData []byte) (*Session, error) {
	var result Session
	err := json.Unmarshal(jsonData, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// isExpired function verifies if the current session is expired by using its Expires attribute
// if returns an error if the Expires string can not be pased using sessionExpirationTimeLayout
func (s *Session) isExpired() (bool, error) {
	expirationTime, err := time.Parse(sessionExpirationTimeLayout, s.Expires)
	if err != nil {
		return true, err
	}
	if expirationTime.Before(time.Now()) {
		return true, nil
	}
	return false, nil
}
