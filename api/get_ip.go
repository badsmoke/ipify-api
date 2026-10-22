// ipify-api/api
//
// This package holds our API handlers which we use to service REST API
// requests.

package api

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/rdegges/ipify-api/models"
)

var validJSONPCallback = regexp.MustCompile(`^[A-Za-z_$][A-Za-z0-9_$]*(\.[A-Za-z_$][A-Za-z0-9_$]*)*$`)

// GetIP returns a user's public facing IP address (IPv4 OR IPv6).
//
// By default, it will return the IP address in plain text, but can also return
// data in both JSON and JSONP if requested to.
func GetIP(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid request form", http.StatusBadRequest)
		return
	}

	ip := clientIP(r)
	if ip == "" {
		http.Error(w, "Could not determine client IP address", http.StatusBadRequest)
		return
	}

	// If the user specifies a 'format' querystring, we'll try to return the
	// user's IP address in the specified format.
	if format, ok := r.Form["format"]; ok && len(format) > 0 {
		jsonBytes, err := json.Marshal(models.IPAddress{IP: ip})
		if err != nil {
			http.Error(w, "Could not encode response", http.StatusInternalServerError)
			return
		}

		switch format[0] {
		case "json":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write(jsonBytes)
			return
		case "jsonp":
			// If the user specifies a 'callback' parameter, we'll use that as
			// the name of our JSONP callback.
			callback := "callback"
			if val, ok := r.Form["callback"]; ok && len(val) > 0 {
				callback = val[0]
			}
			if !validJSONPCallback.MatchString(callback) {
				http.Error(w, "Invalid JSONP callback", http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/javascript")
			_, _ = fmt.Fprintf(w, "%s(%s);", callback, jsonBytes)
			return
		}
	}

	// If no 'format' querystring was specified, we'll default to returning the
	// IP in plain text.
	w.Header().Set("Content-Type", "text/plain")
	_, _ = fmt.Fprint(w, ip)
}

// clientIP uses the first forwarded address when the service is behind a
// proxy that overwrites X-Forwarded-For, then falls back to the TCP peer.
func clientIP(r *http.Request) string {
	forwardedFor := strings.TrimSpace(strings.SplitN(r.Header.Get("X-Forwarded-For"), ",", 2)[0])
	if ip := net.ParseIP(forwardedFor); ip != nil {
		return ip.String()
	}

	remoteHost, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteHost = r.RemoteAddr
	}
	if ip := net.ParseIP(strings.Trim(remoteHost, "[]")); ip != nil {
		return ip.String()
	}

	return ""
}
