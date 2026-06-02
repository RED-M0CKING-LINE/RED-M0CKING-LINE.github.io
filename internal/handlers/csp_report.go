// Handle CSP reports
package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (p *Pages) CSPReport(w http.ResponseWriter, r *http.Request) {
	ct := r.Header.Get("Content-Type")

	switch {
	case strings.HasPrefix(ct, "application/reports+json"):
		// Reporting API v1, array of report objects
		var reports []struct {
			Type string `json:"type"`
			Body struct {
				DocumentURL        string `json:"documentURL"`
				BlockedURL         string `json:"blockedURL"`
				EffectiveDirective string `json:"effectiveDirective"`
				OriginalPolicy     string `json:"originalPolicy"`
				StatusCode         int    `json:"statusCode"`
			} `json:"body"`
		}
		if err := json.NewDecoder(r.Body).Decode(&reports); err != nil {
			p.Log.Warn("csp report decode failed", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		for _, rep := range reports {
			p.Log.Warn("csp violation",
				"type", rep.Type,
				"document", rep.Body.DocumentURL,
				"blocked", rep.Body.BlockedURL,
				"directive", rep.Body.EffectiveDirective,
			)
		}

	case strings.HasPrefix(ct, "application/csp-report"):
		// Legacy report-uri format, a single wrapped object
		var report struct {
			Body struct {
				DocumentURI        string `json:"document-uri"`
				BlockedURI         string `json:"blocked-uri"`
				ViolatedDirective  string `json:"violated-directive"`
				EffectiveDirective string `json:"effective-directive"`
				OriginalPolicy     string `json:"original-policy"`
				StatusCode         int    `json:"status-code"`
			} `json:"csp-report"`
		}
		if err := json.NewDecoder(r.Body).Decode(&report); err != nil {
			p.Log.Warn("csp report decode failed", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		p.Log.Warn("csp violation",
			"document", report.Body.DocumentURI,
			"blocked", report.Body.BlockedURI,
			"directive", report.Body.EffectiveDirective,
		)

	default:
		p.Log.Warn("csp report unknown content-type", "content-type", ct)
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
