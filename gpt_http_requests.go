package gogpt

import (
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

//initCookieJarAndHttpClient initialises the autoFillingCookieJar and http.Client instances inside the current *gpt instance
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

func (g *gpt) loadSessionDetails() error{
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	// TODO: Parse the body as session
	logger.Debug("session token response", zap.String("session-request-response", string(body)))
	return nil
}
