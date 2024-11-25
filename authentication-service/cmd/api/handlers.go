package main

import (
	"errors"
	"fmt"
	"net/http"
)

func (c *Config) Authenticate(w http.ResponseWriter, r *http.Request) {

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := c.readJSON(w, r, &request)

	if err != nil {
		c.errorJson(w, err)
		return
	}

	user, err := c.Models.User.GetByEmail(request.Email)
	if err != nil {
		c.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	valid, err := user.PasswordMatches(request.Password)
	if err != nil || !valid {
		c.errorJson(w, errors.New("invalid credentials"), http.StatusUnauthorized)
		return
	}

	payload := jsonResponse{
		Error:   false,
		Message: fmt.Sprintf("user logged in %q", user.Email),
		Data:    user,
	}

	c.writeJSON(w, http.StatusOK, payload)

}
