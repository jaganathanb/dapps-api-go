package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"io"

	"github.com/gorilla/mux"
	"github.com/jaganathanb/dapps-api-go/api/auth"
	"github.com/jaganathanb/dapps-api-go/api/models"
	"github.com/jaganathanb/dapps-api-go/api/responses"
	"github.com/jaganathanb/dapps-api-go/api/utils/formaterror"
)

func (server *Server) CreateGst(w http.ResponseWriter, r *http.Request) {

	body, err := io.ReadAll(io.Reader(r.Body))
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	gst := models.GST{}
	err = json.Unmarshal(body, &gst)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	gst.Prepare()
	err = gst.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	if uid != gst.ID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New(http.StatusText(http.StatusUnauthorized)))
		return
	}
	gstCreated, err := gst.SaveGST(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.URL.Path, gstCreated.ID))
	responses.JSON(w, http.StatusCreated, gstCreated)
}

func (server *Server) GetGsts(w http.ResponseWriter, r *http.Request) {

	gst := models.GST{}

	gsts, err := gst.FindAllGSTs(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, gsts)
}

func (server *Server) GetGst(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	gstin := vars["gstin"]
	if gstin == "" {
		responses.ERROR(w, http.StatusBadRequest, errors.New("GSTIN can't be empty"))
		return
	}
	gst := models.GST{}

	gstReceived, err := gst.FindGSTByID(server.DB, gstin)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, gstReceived)
}

func (server *Server) UpdateGst(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// Check if the gstin is valid
	gstin, isErr := vars["gstin"]
	if isErr {
		responses.ERROR(w, http.StatusBadRequest, errors.New("GSTIN can't be empty"))
		return
	}

	//CHeck if the auth token is valid and get the user id from it
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the gst exist
	gst := models.GST{}
	err = server.DB.Debug().Model(models.GST{}).Where("gstin = ?", gstin).Take(&gst).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("GST not found"))
		return
	}

	// Read the data gsted
	body, err := io.ReadAll(r.Body)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	// Start processing the request data
	gstUpdate := models.GST{}
	err = json.Unmarshal(body, &gstUpdate)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	gstUpdate.Prepare()
	err = gstUpdate.Validate()
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	gstUpdate.ID = gst.ID //this is important to tell the model the gst id to update, the other update field are set above

	gstUpdated, err := gstUpdate.UpdateAGST(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, gstUpdated)
}

func (server *Server) DeleteGst(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	gstin, isErr := vars["gstin"]
	if isErr {
		responses.ERROR(w, http.StatusBadRequest, errors.New("GSTIN can't be empty"))
		return
	}

	// Is this user authenticated?
	_, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check if the gst exist
	gst := models.GST{}
	err = server.DB.Debug().Model(models.GST{}).Where("gstin = ?", gstin).Take(&gst).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Unauthorized"))
		return
	}

	_, err = gst.DeleteAGST(server.DB, gstin)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", gstin)
	responses.JSON(w, http.StatusNoContent, "")
}
