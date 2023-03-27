package gogpt

import (
	"github.com/playwright-community/playwright-go"
	"math/rand"
	"net/http"
	"time"
)

const baseURL = "https://chat.openai.com"


//randomTimeOut returns a random float64 value used as a timeout between 1000 and 10,000
func randomTimeOut() float64 {
	r := rand.New(rand.NewSource(time.Now().UnixMicro()))
	return float64(r.Intn(10000) + 1000)
}

func playwrightSameSiteAttributeToHttpSameSite(attribute *playwright.SameSiteAttribute) http.SameSite {
	switch attribute {
	case playwright.SameSiteAttributeStrict: return http.SameSiteStrictMode
	case playwright.SameSiteAttributeLax: return http.SameSiteLaxMode
	case playwright.SameSiteAttributeNone: return http.SameSiteNoneMode
	default: return http.SameSiteDefaultMode
	}
}

func playwrightToHttpCookie (playwrightCookie *playwright.BrowserContextCookiesResult) *http.Cookie {
	return &http.Cookie{
		Name:       playwrightCookie.Name,
		Value:      playwrightCookie.Value,
		Path:       playwrightCookie.Path,
		Domain:     playwrightCookie.Domain,
		Expires:    time.Unix(int64(playwrightCookie.Expires), 0),
		Secure:     playwrightCookie.Secure,
		HttpOnly:   playwrightCookie.HttpOnly,
		SameSite:   playwrightSameSiteAttributeToHttpSameSite(&playwrightCookie.SameSite),
	}
}

func playwrightCookiesToHttpCookies(cookies []*playwright.BrowserContextCookiesResult) []*http.Cookie {
	var result []*http.Cookie = make([]*http.Cookie, len(cookies), len(cookies))
	for i,cookie := range cookies {
		result[i] = playwrightToHttpCookie(cookie)
	}
	return result
}