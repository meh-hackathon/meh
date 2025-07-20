package server

import (
	"net/http"

	"github.com/meh-hackathon/meh/auth"
	"github.com/meh-hackathon/meh/httpx"
)

func (s *Server) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUser(r.Context())
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	httpx.WriteJSON(w, http.StatusOK, user)
}

func (s *Server) LoginUser(w http.ResponseWriter, r *http.Request) {
	body, err := httpx.ParseReqBody[auth.TokenRequest](r)
	if err != nil {
		httpx.WriteError(w, err)
		return
	}

	switch body.GrantType {
	case "password":
		tokenPair, err := s.oauthHandler.HandlePasswordGrant(r.Context(), body)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, tokenPair)
		return
	case "refresh_token":
		tokenPair, err := s.oauthHandler.HandleRefreshGrant(r.Context(), body)
		if err != nil {
			httpx.WriteError(w, err)
			return
		}
		httpx.WriteJSON(w, http.StatusOK, tokenPair)
		return
	default:
		httpx.WriteError(w, auth.ErrInvalidToken.WithApiMessagef("Unsupported grant type: %s. Supported types are: password, refresh_token", body.GrantType))
	}
}
