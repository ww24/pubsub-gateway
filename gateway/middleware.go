package gateway

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
)

// cors is middleware for cors headers.
func (g *gateway) cors(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		origin := g.defaultOrigin
		if ref := r.Header.Get("Origin"); ref != "" {
			u, err := url.Parse(ref)
			if err == nil {
				if strings.HasSuffix(u.Hostname(), g.allowOriginSuffix) {
					origin = u.Scheme + "://" + u.Host
				}
			}
		}
		header := w.Header()
		header.Set("Access-Control-Allow-Origin", origin)
		header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		header.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		header.Set("Access-Control-Expose-Headers", "Content-Length")
		header.Set("Access-Control-Allow-Credentials", "true")
		header.Set("Access-Control-Max-Age", "600")
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

func (g *gateway) authorizeIDToken(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			// passed
			h.ServeHTTP(w, r)
			return
		}

		userInfoHeader := r.Header.Get("X-Endpoint-API-UserInfo")
		if userInfoHeader == "" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		d, err := base64.RawURLEncoding.DecodeString(userInfoHeader)
		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("failed to decode base64 encoded header: %+v\n", err)
			return
		}
		log.Println("espv2-userinfo:", string(d))
		userInfo := &ESPv2UserInfo{}
		if err := json.Unmarshal(d, userInfo); err != nil {
			w.WriteHeader(http.StatusForbidden)
			log.Printf("failed to unmarshal json: %+v\n", err)
			return
		}
		if !userInfo.EmailVerified {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		for _, u := range authorizedUsers {
			if userInfo.Email == u {
				// passed
				h.ServeHTTP(w, r)
				return
			}
		}
		// failed
		w.WriteHeader(http.StatusForbidden)
	}
	return http.HandlerFunc(fn)
}

// ESPv2UserInfo is request header of ESPv2.
type ESPv2UserInfo struct {
	Iss           string `json:"iss"`
	Azp           string `json:"azp"`
	Aud           string `json:"aud"`
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AtHash        string `json:"at_hash"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Locale        string `json:"locale"`
	Iat           int64  `json:"iat"`
	Exp           int64  `json:"exp"`
	Jti           string `json:"jti"`
}
