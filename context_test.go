package spell

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockFlash() Flash {
	f := NewFlash()
	f.Error("error")
	f.Success("success")

	return f
}

func mockSession() Session {
	s := NewSession()
	s["foo"] = "bar"

	return s
}

func mockContext(flash Flash, session Session) (*Context, *http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", "https://example.com/foo", nil)
	w := httptest.NewRecorder()

	options := ContextOptions{
		csrfToken: true,
		cookies:   defaultCookies,
	}
	cookieManager := NewCookieManager("qTvsuYaLZTOmaRz8", "qTvsuYaLZTOmaRz8")

	if flash != nil {
		encoded, err := cookieManager.Encode(FlashCookie, flash)
		if err != nil {
			panic(err)
		}

		req.AddCookie(&http.Cookie{
			Name:     FlashCookie,
			Value:    encoded,
			Path:     options.cookies.Path,
			Secure:   options.cookies.Secure,
			HttpOnly: options.cookies.HttpOnly,
			SameSite: options.cookies.SameSite,
			Domain:   options.cookies.Domain,
		})
	}

	if session != nil {
		encoded, err := cookieManager.Encode(SessionCookie, session)
		if err != nil {
			panic(err)
		}

		req.AddCookie(&http.Cookie{
			Name:     SessionCookie,
			Value:    encoded,
			Path:     options.cookies.Path,
			Secure:   options.cookies.Secure,
			HttpOnly: options.cookies.HttpOnly,
			SameSite: options.cookies.SameSite,
			Domain:   options.cookies.Domain,
		})
	}

	c, err := NewContext(w, req, options, NewDefaultLogger(), cookieManager)
	if err != nil {
		panic(err)
	}

	return c, req, w
}

func TestFlash(t *testing.T) {
	f := mockFlash()
	c, _, _ := mockContext(f, nil)

	flash := c.ViewData["flash"].(Flash)

	if flash["error"] != "error" {
		t.Errorf("Expected flash error to be 'error', got '%s'", flash["error"])
	}
	if flash["success"] != "success" {
		t.Errorf("Expected flash success to be 'success', got '%s'", flash["success"])
	}
}

func TestSession(t *testing.T) {
	s := mockSession()
	c, _, _ := mockContext(nil, s)

	if c.Session["foo"] != "bar" {
		t.Errorf("Expected session foo to be 'bar', got '%s'", c.Session["foo"])
	}
}

func TestMakeResponse(t *testing.T) {
	c, _, w := mockContext(nil, nil)

	c.Flash.Error("error")
	c.Session["foo"] = "bar"

	err := c.makeResponse()
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	//Check Set-Cookie headers
	cookies := w.Header()["Set-Cookie"]
	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookies, got %d", len(cookies))
	}

	//Check flash cookie
	flashCookie := cookies[0]
	if !strings.Contains(flashCookie, "flash=") {
		t.Errorf("Expected flash cookie, got '%s'", flashCookie)
	}

	//Check session cookie
	sessionCookie := cookies[1]
	if !strings.Contains(sessionCookie, "session=") {
		t.Errorf("Expected session cookie, got '%s'", sessionCookie)
	}
}

func TestMakeResponseNoFlash(t *testing.T) {
	c, _, w := mockContext(mockFlash(), nil)

	err := c.makeResponse()
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	// check Set-Cookie headers
	cookies := w.Header()["Set-Cookie"]
	if len(cookies) != 2 {
		t.Errorf("Expected 2 cookie, got %d", len(cookies))
	}
}

func TestRedirect(t *testing.T) {
	c, _, w := mockContext(nil, nil)

	c.Redirect("/foo")

	if w.Code != http.StatusFound {
		t.Errorf("Expected status code %d, got %d", http.StatusFound, w.Code)
	}

	if w.Header().Get("Location") != "/foo" {
		t.Errorf("Expected location to be '/foo', got '%s'", w.Header().Get("Location"))
	}
}

func TestCsrfToken(t *testing.T) {
	c, _, _ := mockContext(nil, nil)

	if c.ViewData[CSRFToken] == nil {
		t.Error("Expected csrfToken to be set")
	}

	if c.ViewData["renderToken"] == nil {
		t.Error("Expected renderToken to be set")
	}
}
