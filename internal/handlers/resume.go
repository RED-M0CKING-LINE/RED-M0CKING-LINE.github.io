package handlers

import (
	"net/http"
)

func (p *Pages) Resume(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	var userPIN string = r.URL.Query().Get("pin")

	var resumePINs map[string]bool = map[string]bool{
		"MTU1885": true,
	}

	var authed bool = false

	if resumePINs[userPIN] {
		authed = true
	}

	sess, _ := p.Auth.SessionFrom(r)
	_ = p.Tpl.Render(w, "resume", p.base(r, "resume", map[string]any{
		"User":   sess,
		"Authed": authed,
	}))
}
