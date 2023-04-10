package gogpt

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
)

// initCookieJarAndHttpClient initialises the autoFillingCookieJar and http.Client instances inside the current *gpt instance
func (g *gpt) initCookieJarAndHttpClient() error {
	if g.cookieJar == nil {
		cookieJar, err := createNewAutoFillingCookieJar(baseURL, g.getUserCookiesSupplier(baseURL))
		if err != nil {
			return err
		}
		g.cookieJar = cookieJar
		g.httpClient = &http.Client{
			Jar: g.cookieJar,
		}
		return nil
	}
	if g.httpClient == nil {
		g.httpClient = &http.Client{
			Jar: g.cookieJar,
		}
	}
	return nil
}

// initSession initializes the session of the current gpt instance.
// It returns an error if something goes wrong while unmarshalling the api response
func (g *gpt) initSession() error {
	err := g.cookieJar.setExpiredCookies()
	if err != nil {
		return err
	}
	resp, err := g.httpClient.Get("https://chat.openai.com/api/auth/session")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Debug("session token response", zap.String("session-request-response", string(body)))
	s, err := unmarshalGPTSessionResponseJSON(body)
	if err != nil {
		logger.Error("Error while unmarshalling session response")
		return err
	}
	g.session = s

	return nil
}

// refreshSession verifies if there's a session exists. If there's no session exists, creates one using initSession
// if there's an existing session verifies if the session is expired using isExpired function. If the session is expired
// recreates the session using initSession
func (g *gpt) refreshSession() error {
	if g.session == nil {
		return g.initSession()
	}
	isExpired, err := g.session.isExpired()
	if err != nil {
		logger.Error("Error while checking if the existing session is expired", zap.String("expiration-date-string", g.session.Expires))
		return g.initSession()
	}
	if isExpired {
		return g.initSession()
	}
	return nil
}

// prepareRequest prepares the cookies and the user session to use in each http request. This function should be called
// before each http request to ensure that the request will not be blocked
func (g *gpt) prepareRequest() error {
	err := g.initCookieJarAndHttpClient()
	if err != nil {
		return err
	}
	err = g.cookieJar.setExpiredCookies()
	if err != nil {
		return err
	}
	err = g.refreshSession()
	if err != nil {
		return err
	}
	return nil
}

func createAPIURL(endpoint string) string {
	return fmt.Sprintf("%s/backend-api/%s", baseURL, endpoint)
}

func (g *gpt) createRequest(method string, endpoint string, body io.Reader) (*http.Request, error) {
	err := g.prepareRequest()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, createAPIURL(endpoint), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", g.session.AccessToken))
	req.Header.Set("content-type", "/application/json")
	return req, nil
}

func runAPIRequest[T any](g *gpt, method, endpoint string, requestBody io.Reader) (*T, error) {
	req, err := g.createRequest(method, endpoint, requestBody)
	if err != nil {
		return nil, err
	}
	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("run http %s request on  %s with %+v failed with status code %d", method, endpoint, body, resp.StatusCode)
	}
	var response T
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (g *gpt) getConversationHistory(offset, limit uint) (*ConversationsResponse, error) {
	return runAPIRequest[ConversationsResponse](g, "GET", fmt.Sprintf("conversations?offset=%d&limit=%d", offset, limit), nil)
}
