package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/bhagatnilesh11/fullstack/api/auth"
	"github.com/bhagatnilesh11/fullstack/api/models"
	"github.com/bhagatnilesh11/fullstack/api/responses"
	"github.com/bhagatnilesh11/fullstack/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) GetEmployees(w http.ResponseWriter, r *http.Request) {
	employee := models.Employee{}
	employees, err := employee.FindAllEmployees(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, employees)

}

func (server *Server) GetEmployee(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	uid, err := strconv.ParseUint(vars["id"], 10, 32)

	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	employee := models.Employee{}
	getUser, err := employee.FindEmployeeByID(server.DB, uint32(uid))
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, getUser)

}

func (server *Server) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	employee := models.Employee{}
	err = json.Unmarshal(body, &employee)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}
	employee.Prepare()
	empCreated, err := employee.CreateEmployee(server.DB)
	if err != nil {
		formaterror := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formaterror)
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, empCreated.ID))
	responses.JSON(w, http.StatusOK, empCreated)
}

func (server *Server) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uid, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
	}

	employee := models.Employee{}
	err = json.Unmarshal(body, &employee)

	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	tokenid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New("unauthorised"))
		return
	}
	if tokenid != uint32(uid) {
		responses.ERROR(w, http.StatusUnprocessableEntity, errors.New(http.StatusText(http.StatusUnauthorized)))
	}
	employee.Prepare()
	//err = employee.Validate()
	updateduser, err := employee.UpdateEmployee(server.DB, int32(uid))
	if err != nil {
		formaterror := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formaterror)
		return
	}
	responses.JSON(w, http.StatusOK, updateduser)
}

func (server *Server) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	uid, err := strconv.ParseUint(params["id"], 10, 32)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	employee := models.Employee{}
	_, err = employee.DeleteEmployee(server.DB, int32(uid))

	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", uid))
	responses.JSON(w, http.StatusNoContent, "")

}
