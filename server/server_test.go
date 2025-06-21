package server

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/sultaniman/kpow/config"
)

func loadTestConfig(t *testing.T) *config.Config {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to resolve working dir: %v", err)
	}
	root := filepath.Dir(wd)

	cfgPath := filepath.Join(root, "config.toml")
	cfg, err := config.GetConfig(cfgPath)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	keyBytes, err := os.ReadFile(filepath.Join(root, "server/enc/testkeys/pubkey.gpg"))
	if err != nil {
		t.Fatalf("failed to read pubkey: %v", err)
	}

	keyDir := t.TempDir()
	keyFile := filepath.Join(keyDir, "pubkey.gpg")
	if err := os.WriteFile(keyFile, keyBytes, 0o644); err != nil {
		t.Fatalf("failed to write pubkey: %v", err)
	}

	cfg.Key.Path = keyFile

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

	rateLimitHit := false
	for range 100 {
		postReq2 := httptest.NewRequest(http.MethodGet, "/", nil)
		postReq2.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
		postReq2.Header.Set("X-Real-IP", "100.100.100.100")
		postReq2.AddCookie(csrfCookie)
		postRec2 := httptest.NewRecorder()
		e.ServeHTTP(postRec2, postReq2)
		if http.StatusTooManyRequests == postRec2.Code {
			rateLimitHit = true
			break
		}
	}

	assert.True(t, rateLimitHit)
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

func TestAdvertiseKeyRendering(t *testing.T) {
	cfg := loadTestConfig(t)
	cfg.RateLimiter = &config.RateLimiter{RPM: 0}

	// ensure the public key bytes are loaded for the form renderer
	if errs := cfg.Validate(); len(errs) > 0 {
		t.Fatalf("config validation failed: %v", errs)
	}

	t.Run("enabled", func(t *testing.T) {
		cfg.Key.Advertise = true
		e := newTestServer(t, cfg)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		assert.Contains(t, body, "id=\"pubkey\"")
		assert.Contains(t, body, string(cfg.Key.KeyBytes))
	})

	t.Run("disabled", func(t *testing.T) {
		cfg.Key.Advertise = false
		e := newTestServer(t, cfg)

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		assert.NotContains(t, body, "id=\"pubkey\"")
	})
}

func TestFormBannerRendering(t *testing.T) {
	cfg := loadTestConfig(t)
	cfg.RateLimiter = &config.RateLimiter{RPM: 0}
	cfg.Server.CustomBanner = "<span>test banner</span>"

	e := newTestServer(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.Contains(t, body, cfg.Server.CustomBanner)
}

func TestFormHideLogo(t *testing.T) {
	cfg := loadTestConfig(t)
	cfg.RateLimiter = &config.RateLimiter{RPM: 0}
	cfg.Server.HideLogo = true

	e := newTestServer(t, cfg)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	assert.NotContains(t, body, "class=\"logo\"")
}
