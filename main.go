package icanbanwell

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

// Config the plugin configuration.
type Config struct {
	Enabled bool `json:"enabled,omitempty"`

	// bans: ["1.2.3.4": "RFC3339 Timestamp", ... ]
	Bans map[string]string `json:"bans,omitempty"`

	// whitelist: ["2.3.4.5", ...]
	Whitelist []string `json:"whitelist,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		Bans:    make(map[string]string),
		Enabled: false,
	}
}

// Demo a Demo plugin.
type ICanBanwell struct {
	next    http.Handler
	bans    map[string]string
	enabled bool
	name    string
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	return &ICanBanwell{
		enabled: config.Enabled,
		bans:    config.Bans,
		next:    next,
		name:    name,
	}, nil
}

func (b *ICanBanwell) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if !b.enabled {
		b.next.ServeHTTP(rw, req)
		return
	}

	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xForwardedFor == "" {
		rw.WriteHeader(http.StatusForbidden)
		return
	}

	now := time.Now()
	for _, ip := range strings.Split(xForwardedFor, ",") {
		expireString, ok := b.bans[ip]
		if ok {
			expireTime, err := time.Parse(time.RFC3339, expireString)
			if err != nil {
				log.Warn().
					Str("timestamp", expireString).
					Msg("Expire timestamp cannot be parsed, use an RFC3339 style timestamp.")
				delete(b.bans, ip)
			} else {
				if now.Before(expireTime) {
					rw.WriteHeader(http.StatusForbidden)
					return
				} else {
					// Expire has expired, remove the ban
					delete(b.bans, ip)
				}
			}
		}
	}

	b.next.ServeHTTP(rw, req)
}
