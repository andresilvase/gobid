package api

import (
	"errors"
	"net/http"

	"github.com/andresilvase/gobid/internal/jsonutils"
	"github.com/andresilvase/gobid/internal/services"
	"github.com/andresilvase/gobid/internal/usecase/user"
)

func (a *Api) handleSignup(w http.ResponseWriter, r *http.Request) {

	data, problems, err := jsonutils.DecodeValidJson[user.CreateUserReq](r)

	if err != nil {
		_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, problems)
		return
	}

	id, err := a.UserService.CreateUser(
		r.Context(),
		data.UserName,
		data.Email,
		data.Password,
		data.Bio,
	)

	if err != nil {
		if errors.Is(err, services.ErrDuplicateEmailOrUsername) {
			_ = jsonutils.EncodeJson(w, r, http.StatusUnprocessableEntity, map[string]any{
				"error": "email or username already exists",
			})
			return
		}
	}

	_ = jsonutils.EncodeJson(w, r, http.StatusCreated, map[string]any{
		"user_id": id,
	})

}

func (a *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login user"))
}

func (a *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout user"))
}
