package spell

import (
	"github.com/justinas/nosurf"
	uuid "github.com/satori/go.uuid"
	"net/http"
)

type Context struct {
	w        http.ResponseWriter
	r        *http.Request
	id       string
	Flash    Flash
	Session  Session
	ViewData map[string]interface{}
	options  ContextOptions
	Logger   Logger
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
) (*Context, error) {
	id := markRequest(r)
	flash, err := getFlash(r)
	if err != nil {
		logger.Logf(ErrorLevel, "Error getting flash: %s", err)
		return nil, err
	}

	session, err := getSession(r)
	if err != nil {
		logger.Logf(ErrorLevel, "Error getting session: %s", err)
		return nil, err
	}

	viewData := make(map[string]interface{})
	viewData["flash"] = flash

	if options.csrfToken {
		token := nosurf.Token(r)
		viewData[CSRFToken] = token
		viewData["renderToken"] = renderToken(token)
	}

	return &Context{
		id:       id,
		w:        w,
		r:        r,
		options:  options,
		Flash:    NewFlash(),
		Session:  session,
		ViewData: viewData,
		Logger:   logger,
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

func getFlash(r *http.Request) (Flash, error) {
	flash := NewFlash()

	value, err := getCookieValue(r, FlashCookie)
	if err != nil {
		return flash, err
	}

	err = flash.decode(value)

	return flash, err
}

func getSession(r *http.Request) (Session, error) {
	session := NewSession()

	value, err := getCookieValue(r, SessionCookie)
	if err != nil {
		return session, err
	}

	err = session.decode(value)

	return session, err
}

func (c *Context) makeResponse() error {
	c.Logger.Logf(DebugLevel, "Making response for context with id: %s", c.id)

	flash, err := c.Flash.encode()
	if err != nil {
		c.Logger.Logf(ErrorLevel, "Error encoding flash: %s", err)
		return err
	}

	setCookies(c.w, FlashCookie, flash, c.options.cookies)

	session, err := c.Session.encode()
	if err != nil {
		c.Logger.Logf(ErrorLevel, "Error encoding session: %s", err)
		return err
	}

	setCookies(c.w, SessionCookie, session, c.options.cookies)

	return nil
}
