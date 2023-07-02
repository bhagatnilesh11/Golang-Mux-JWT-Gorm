// This code is defining a controller package in Go. It imports the "net/http" package and a custom
// package "github.com/bhagatnilesh11/fullstack/api/responses".
package controllers

import (
	"net/http"

	"github.com/bhagatnilesh11/fullstack/api/responses"
)

func (server *Server) Home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome To This Awesome API")

}
