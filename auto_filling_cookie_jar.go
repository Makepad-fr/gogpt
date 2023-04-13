package gogpt

import (
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type httpCookieSupplier func() ([]*http.Cookie, error)

// autoFillingCookieJar embeds *cookiejar.Jar and adds a custom method that supplies fresh cookies
type autoFillingCookieJar struct {
	*cookiejar.Jar
	u                 *url.URL
	newCookieSupplier httpCookieSupplier
}

// setExpiredCookies checks for expired cookies and sets new ones using newCookieSupplier function
func (c *autoFillingCookieJar) setExpiredCookies() error {
	if c.newCookieSupplier == nil {
		return errors.New("NewCookiesSupplier is empty")
	}

	cookies := c.Cookies(c.u)
	for _, cookie := range cookies {
		if cookie.Expires.Before(time.Now()) {
			newCookies, err := c.newCookieSupplier()
			if err != nil {
				return err
			}
			c.SetCookies(c.u, newCookies)
		}
	}
	return nil
}

// createNewAutoFillingCookieJar creates a new cookie jar related to the given url string and with given httpCookieSupplier
func createNewAutoFillingCookieJar(urlString string, supplier httpCookieSupplier) (*autoFillingCookieJar, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	cookies, err := supplier()
	if err != nil {
		return nil, err
	}
	cj := &autoFillingCookieJar{
		Jar:               jar,
		newCookieSupplier: supplier,
		u:                 u,
	}
	cj.SetCookies(cj.u, cookies)
	return cj, nil
}
