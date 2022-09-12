package spell

import (
	"github.com/gorilla/securecookie"
	"net/http"
)

type CookiesOptions struct {
	Path     string
	Domain   string
	SameSite http.SameSite
	Secure   bool
	HttpOnly bool
}

type CookieManager interface {
	SetCookies(w http.ResponseWriter, name string, value any, options CookiesOptions) error
	GetCookieValue(r *http.Request, name string, dst any) error
}

type DefaultCookieManager struct {
	secure *securecookie.SecureCookie
}

var _ CookieManager = &DefaultCookieManager{}

func NewCookieManager(hashKey, blockKey string) *DefaultCookieManager {
	return &DefaultCookieManager{
		secure: securecookie.New([]byte(hashKey), []byte(blockKey)),
	}
}

func (m *DefaultCookieManager) SetCookies(
	w http.ResponseWriter,
	name string,
	value any,
	options CookiesOptions,
) error {
	encoded, err := m.secure.Encode(name, value)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Path:     options.Path,
		Value:    encoded,
		Secure:   options.Secure,
		HttpOnly: options.HttpOnly,
		SameSite: options.SameSite,
	})

	return nil
}

func (m *DefaultCookieManager) GetCookieValue(r *http.Request, name string, dst any) error {
	cookie, err := r.Cookie(name)

	if err != nil {
		return err
	}

	return m.secure.Decode(name, cookie.Value, dst)
}

func (m *DefaultCookieManager) Encode(name string, value any) (string, error) {
	return m.secure.Encode(name, value)
}
