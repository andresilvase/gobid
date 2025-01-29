package api

import "net/http"

func (a *Api) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create user"))
}

func (a *Api) handleLoginUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Login user"))
}

func (a *Api) handleLogoutUser(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Logout user"))
}
