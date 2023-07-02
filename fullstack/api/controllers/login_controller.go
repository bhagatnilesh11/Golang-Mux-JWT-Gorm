package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/bhagatnilesh11/fullstack/api/auth"
	"github.com/bhagatnilesh11/fullstack/api/models"
	"github.com/bhagatnilesh11/fullstack/api/responses"
	"github.com/bhagatnilesh11/fullstack/api/utils/formaterror"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	user := models.User{}
	// `err = json.Unmarshal(body, &user)` is decoding the JSON data in the `body` variable and storing it
	// in the `user` variable. It is converting the JSON data into a Go struct.
	err = json.Unmarshal(body, &user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// `user.Prepare()` is a method call on the `user` object. It is a function defined in the `models`
	// package that prepares the user object before further processing. The exact implementation of the
	// `Prepare()` method is not shown in the code snippet provided, but it is likely performing tasks such
	// as trimming whitespace, converting fields to lowercase, or any other necessary data preparation
	// steps before validating or saving the user object.
	user.Prepare()

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	token, err := server.SignIn(user.Email, user.Password)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusUnprocessableEntity, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, token)
}

func (server *Server) SignIn(email, password string) (string, error) {

	var err error

	user := models.User{}

	err = server.DB.Debug().Model(models.User{}).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return "", err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}
	return auth.CreateToken(user.ID)
}
