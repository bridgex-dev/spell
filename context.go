package spell

import (
	"errors"
	"github.com/justinas/nosurf"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type Context struct {
	w             http.ResponseWriter
	r             *http.Request
	id            string
	Flash         Flash
	Session       Session
	ViewData      map[string]interface{}
	options       ContextOptions
	Logger        Logger
	CookieManager CookieManager
}

type ContextOptions struct {
	cookies   CookiesOptions
	csrfToken bool
}

func NewContext(
	w http.ResponseWriter,
	r *http.Request,
	options ContextOptions,
	logger Logger,
	cookieManager CookieManager,
) (*Context, error) {
	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	if cookieManager == nil {
		return nil, errors.New("cookie manager cannot be nil")
	}

	id := markRequest(r)
	flash, err := getFlash(cookieManager, r)
	if err != nil {
		logger.Logf(ErrorLevel, "Error getting flash: %s", err)
	}

	session, err := getSession(cookieManager, r)
	if err != nil {
		logger.Logf(ErrorLevel, "Error getting session: %s", err)
	}

	viewData := make(map[string]interface{})
	viewData["flash"] = flash

	if options.csrfToken {
		token := nosurf.Token(r)
		viewData[CSRFToken] = token
		viewData["renderToken"] = renderToken(token)
	}

	return &Context{
		id:            id,
		w:             w,
		r:             r,
		options:       options,
		Flash:         NewFlash(),
		Session:       session,
		ViewData:      viewData,
		Logger:        logger,
		CookieManager: cookieManager,
	}, nil
}

func (c *Context) Redirect(url string) {
	http.Redirect(c.w, c.r, url, http.StatusFound)
}

func markRequest(r *http.Request) string {
	id := uuid.NewV4().String()
	r.Header.Set(HeaderId, id)

	return id
}

func getFlash(manager CookieManager, r *http.Request) (Flash, error) {
	flash := NewFlash()

	err := manager.GetCookieValue(r, FlashCookie, &flash)

	return flash, err
}

func getSession(manager CookieManager, r *http.Request) (Session, error) {
	session := NewSession()

	err := manager.GetCookieValue(r, SessionCookie, &session)

	return session, err
}

func (c *Context) makeResponse() error {
	c.Logger.Logf(DebugLevel, "Making response for context with id: %s", c.id)

	err := c.CookieManager.SetCookies(c.w, FlashCookie, c.Flash, c.options.cookies)
	if err != nil {
		c.Logger.Logf(ErrorLevel, "Error setting flash cookie: %s", err)
		return err
	}

	err = c.CookieManager.SetCookies(c.w, SessionCookie, c.Session, c.options.cookies)
	if err != nil {
		c.Logger.Logf(ErrorLevel, "Error setting session cookie: %s", err)
		return err
	}

	return nil
}
