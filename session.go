package gogpt

import (
	"net/http"
	"net/url"
)

type Session struct {
	Cookies []*http.Cookie
	AuthToken string
	Email string
	Avatar url.URL
}
