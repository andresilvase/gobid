package api

import (
	"log/slog"
	"net/http"

	"github.com/andresilvase/gobid/internal/jsonutils"
	"github.com/gorilla/csrf"
)

func (api *Api) HandleGetCSRFToken(w http.ResponseWriter, r *http.Request) {
	csrfToken := csrf.Token(r)

	jsonutils.EncodeJson(w, r, http.StatusOK, map[string]any{
		"csrf_token": csrfToken,
	})
}

func (api *Api) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Auth middleware triggered",
			"cookies", r.Cookies(),
			"headers", r.Header,
			"session_cookie", r.Header.Get("Cookie"))

		// Get the session data
		userId := api.Sessions.Get(r.Context(), "AuthenticatedUserId")
		slog.Info("Session data", "userId", userId)

		if userId == nil {
			slog.Error("No user ID found in session")
			jsonutils.EncodeJson(w, r, http.StatusUnauthorized, map[string]any{
				"message": "unauthorized",
			})
			return
		}

		next.ServeHTTP(w, r)
	})
}
