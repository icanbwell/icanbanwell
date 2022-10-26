package icanbanwell_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/icanbwell/icanbanwell"
	"github.com/stretchr/testify/assert"
)

func TestDisabled(t *testing.T) {
	cfg := icanbanwell.CreateConfig()
	cfg.Enabled = false

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := icanbanwell.New(ctx, next, cfg, "icanbanwell-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
}

func TestEmptyHeader(t *testing.T) {
	cfg := icanbanwell.CreateConfig()
	cfg.Enabled = true
	cfg.Bans = make(map[string]string)
	cfg.Bans["1.2.3.4"] = time.Now().Add(5 * time.Minute).Format(time.RFC3339)

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := icanbanwell.New(ctx, next, cfg, "icanbanwell-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Result().StatusCode)
}

func TestFutureBan(t *testing.T) {
	cfg := icanbanwell.CreateConfig()
	cfg.Enabled = true
	cfg.Bans = make(map[string]string)
	cfg.Bans["1.2.3.4"] = time.Now().Add(5 * time.Minute).Format(time.RFC3339)

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := icanbanwell.New(ctx, next, cfg, "icanbanwell-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("X-Forwarded-For", "1.2.3.4,2.3.4.5")

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusForbidden, recorder.Result().StatusCode)
}

func TestExpiredBan(t *testing.T) {
	cfg := icanbanwell.CreateConfig()
	cfg.Enabled = true
	cfg.Bans = make(map[string]string)
	cfg.Bans["1.2.3.4"] = time.Now().Add(-5 * time.Minute).Format(time.RFC3339)

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := icanbanwell.New(ctx, next, cfg, "icanbanwell-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	req.Header.Add("X-Forwarded-For", "1.2.3.4,2.3.4.5")

	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
}
