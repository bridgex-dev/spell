package spell

import (
	"net/http"
	"net/url"
)

type CookiesOptions struct {
	Path     string
	Domain   string
	SameSite http.SameSite
	Secure   bool
	HttpOnly bool
}

var defaultCookies = CookiesOptions{
	Path:     "/",
	Domain:   "",
	SameSite: http.SameSiteStrictMode,
	Secure:   true,
	HttpOnly: true,
}

func setCookies(w http.ResponseWriter, name, value string, options CookiesOptions) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Path:     options.Path,
		Value:    encodeCookie(value),
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	})
}

func getCookieValue(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}

	return decodeCookie(cookie.Value)
}

func encodeCookie(value string) string {
	return url.QueryEscape(value)
}

func decodeCookie(value string) (string, error) {
	return url.QueryUnescape(value)
}
