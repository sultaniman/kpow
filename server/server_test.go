package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/sultaniman/kpow/config"
)

func loadTestConfig(t *testing.T) *config.Config {
	t.Helper()
	cfg, err := config.GetConfig("config.toml")
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	return cfg
}

func newTestServer(t *testing.T, cfg *config.Config) *echo.Echo {
	t.Helper()
	h, err := NewHandler(cfg)
	if err != nil {
		t.Fatalf("failed to create handler: %v", err)
	}
	e, err := CreateServer(cfg, h)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	return e
}

func findCSRFCookie(cookies []*http.Cookie) *http.Cookie {
	for _, c := range cookies {
		if strings.Contains(strings.ToLower(c.Name), "csrf") {
			return c
		}
	}
	return nil
}

func TestRateLimiting(t *testing.T) {
	cfg := loadTestConfig(t)
	cfg.RateLimiter = &config.RateLimiter{RPM: 1, Burst: 1, CooldownSeconds: 60}

	e := newTestServer(t, cfg)

	// initial GET to obtain csrf token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	csrfCookie := findCSRFCookie(rec.Result().Cookies())
	if csrfCookie == nil {
		t.Fatal("csrf cookie not found")
	}

	form := url.Values{}
	form.Set("subject", "hello")
	form.Set("content", "world")
	form.Set("csrf", csrfCookie.Value)
	body := strings.NewReader(form.Encode())

	postReq := httptest.NewRequest(http.MethodPost, "/", body)
	postReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	postReq.AddCookie(csrfCookie)
	postRec := httptest.NewRecorder()
	e.ServeHTTP(postRec, postReq)
	assert.Equal(t, http.StatusOK, postRec.Code)

	// second request should hit the limiter
	postReq2 := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(form.Encode()))
	postReq2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	postReq2.AddCookie(csrfCookie)
	postRec2 := httptest.NewRecorder()
	e.ServeHTTP(postRec2, postReq2)
	assert.Equal(t, http.StatusTooManyRequests, postRec2.Code)
}

func TestInvalidCSRFToken(t *testing.T) {
	cfg := loadTestConfig(t)
	cfg.RateLimiter = &config.RateLimiter{RPM: 0}

	e := newTestServer(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	csrfCookie := findCSRFCookie(rec.Result().Cookies())
	if csrfCookie == nil {
		t.Fatal("csrf cookie not found")
	}

	form := url.Values{}
	form.Set("subject", "hello")
	form.Set("content", "world")
	form.Set("csrf", "badtoken")
	body := strings.NewReader(form.Encode())

	postReq := httptest.NewRequest(http.MethodPost, "/", body)
	postReq.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	postReq.AddCookie(csrfCookie)
	postRec := httptest.NewRecorder()
	e.ServeHTTP(postRec, postReq)
	assert.Equal(t, http.StatusForbidden, postRec.Code)
}
