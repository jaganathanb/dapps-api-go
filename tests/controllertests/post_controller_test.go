package controllertests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jaganathanb/dapps-api-go/api/models"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateGST(t *testing.T) {

	err := refreshUserAndGstTable()
	if err != nil {
		log.Fatal(err)
	}
	user, err := seedOneUser()
	if err != nil {
		log.Fatalf("Cannot seed user %v\n", err)
	}
	token, err := server.SignIn(user.Email, "password") //Note the password in the database is already hashed, we want unhashed
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	samples := []struct {
		inputJSON    string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   201,
			tokenGiven:   tokenString,
			title:        "The title",
			content:      "the content",
			errorMessage: "",
		},
		{
			inputJSON:    `{"title":"The title", "content": "the content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			// When no token is passed
			inputJSON:    `{"title":"When no token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "",
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			inputJSON:    `{"title":"When incorrect token is passed", "content": "the content", "author_id": 1}`,
			statusCode:   401,
			tokenGiven:   "This is an incorrect token",
			errorMessage: "Unauthorized",
		},
		{
			inputJSON:    `{"title": "", "content": "The content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			inputJSON:    `{"title": "This is a title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			inputJSON:    `{"title": "This is an awesome title", "content": "the content"}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Author",
		},
		{
			// When user 2 uses user 1 token
			inputJSON:    `{"title": "This is an awesome title", "content": "the content", "author_id": 2}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range samples {

		req, err := http.NewRequest("POST", "/gsts", bytes.NewBufferString(v.inputJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.CreateGst)

		req.Header.Set("Authorization", v.tokenGiven)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			fmt.Printf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 201 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just for both ids to have the same type
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestGetGsts(t *testing.T) {

	err := refreshUserAndGstTable()
	if err != nil {
		log.Fatal(err)
	}
	_, _, err = seedUsersAndGsts()
	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest("GET", "/gsts", nil)
	if err != nil {
		t.Errorf("this is the error: %v\n", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(server.GetGsts)
	handler.ServeHTTP(rr, req)

	var gsts []models.GST
	err = json.Unmarshal([]byte(rr.Body.String()), &gsts)

	assert.Equal(t, rr.Code, http.StatusOK)
	assert.Equal(t, len(gsts), 2)
}
func TestGetGstByID(t *testing.T) {

	err := refreshUserAndGstTable()
	if err != nil {
		log.Fatal(err)
	}
	gst, err := seedOneUserAndOneGst()
	if err != nil {
		log.Fatal(err)
	}
	gstSample := []struct {
		ID           uuid.UUID
		GSTIN        string
		TradeName    string
		RegisteredAt time.Time
		Address      string
		OwnerName    string
		statusCode   int
	}{
		{
			ID:         gst.ID,
			statusCode: 200,
			GSTIN:      gst.GSTIN,
			TradeName:  gst.TradeName,
			OwnerName:  gst.OwnerName,
		},
		{
			ID:         uuid.Nil,
			statusCode: 400,
		},
	}
	for _, v := range gstSample {

		req, err := http.NewRequest("GET", "/gsts", nil)
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}

		req = mux.SetURLVars(req, map[string]string{"gstin": v.GSTIN})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.GetGst)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)

		if err != nil {
			log.Fatalf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 200 {
			assert.Equal(t, gst.GSTIN, responseMap["GSTIN"])
			assert.Equal(t, gst.TradeName, responseMap["TradeName"])
			assert.Equal(t, gst.Address, responseMap["Address"]) //the response author id is float64
		}
	}
}

func TestUpdateGst(t *testing.T) {

	var GstUserEmail, GstUserPassword string
	var AuthGstAuthorID uint32
	var AuthGstID uuid.UUID

	err := refreshUserAndGstTable()
	if err != nil {
		log.Fatal(err)
	}
	_, gsts, err := seedUsersAndGsts()
	if err != nil {
		log.Fatal(err)
	}

	//Login the user and get the authentication token
	token, err := server.SignIn(GstUserEmail, GstUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the first gst
	for _, gst := range gsts {
		if gst.ID == uuid.Nil {
			continue
		}
		AuthGstID = gst.ID
	}
	// fmt.Printf("this is the auth gst: %v\n", AuthGstID)

	samples := []struct {
		id           uuid.UUID
		updateJSON   string
		statusCode   int
		title        string
		content      string
		author_id    uint32
		tokenGiven   string
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           AuthGstID,
			updateJSON:   `{"title":"The updated gst", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   200,
			title:        "The updated gst",
			content:      "This is the updated content",
			author_id:    AuthGstAuthorID,
			tokenGiven:   tokenString,
			errorMessage: "",
		},
		{
			// When no token is provided
			id:           AuthGstID,
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is provided
			id:           AuthGstID,
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 1}`,
			tokenGiven:   "this is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			//Note: "Title 2" belongs to gst 2, and title must be unique
			id:           AuthGstID,
			updateJSON:   `{"title":"Title 2", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   500,
			tokenGiven:   tokenString,
			errorMessage: "Title Already Taken",
		},
		{
			id:           AuthGstID,
			updateJSON:   `{"title":"", "content": "This is the updated content", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Title",
		},
		{
			id:           AuthGstID,
			updateJSON:   `{"title":"Awesome title", "content": "", "author_id": 1}`,
			statusCode:   422,
			tokenGiven:   tokenString,
			errorMessage: "Required Content",
		},
		{
			id:           AuthGstID,
			updateJSON:   `{"title":"This is another title", "content": "This is the updated content"}`,
			statusCode:   401,
			tokenGiven:   tokenString,
			errorMessage: "Unauthorized",
		},
		{
			id:         uuid.Nil,
			statusCode: 400,
		},
		{
			id:           AuthGstID,
			updateJSON:   `{"title":"This is still another title", "content": "This is the updated content", "author_id": 2}`,
			tokenGiven:   tokenString,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}

	for _, v := range samples {

		req, err := http.NewRequest("POST", "/gsts", bytes.NewBufferString(v.updateJSON))
		if err != nil {
			t.Errorf("this is the error: %v\n", err)
		}
		req = mux.SetURLVars(req, map[string]string{"gstin": v.id.String()})
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.UpdateGst)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			t.Errorf("Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, v.statusCode)
		if v.statusCode == 200 {
			assert.Equal(t, responseMap["title"], v.title)
			assert.Equal(t, responseMap["content"], v.content)
			assert.Equal(t, responseMap["author_id"], float64(v.author_id)) //just to match the type of the json we receive thats why we used float64
		}
		if v.statusCode == 401 || v.statusCode == 422 || v.statusCode == 500 && v.errorMessage != "" {
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}

func TestDeleteGST(t *testing.T) {

	var GstUserEmail, GstUserPassword string
	var GstUserID uint32
	var AuthGstID uuid.UUID

	err := refreshUserAndGstTable()
	if err != nil {
		log.Fatal(err)
	}
	_, gsts, err := seedUsersAndGsts()
	if err != nil {
		log.Fatal(err)
	}

	//Login the user and get the authentication token
	token, err := server.SignIn(GstUserEmail, GstUserPassword)
	if err != nil {
		log.Fatalf("cannot login: %v\n", err)
	}
	tokenString := fmt.Sprintf("Bearer %v", token)

	// Get only the second gst
	for _, gst := range gsts {
		if gst.ID == uuid.Nil {
			continue
		}
		AuthGstID = gst.ID
	}
	gstSample := []struct {
		id           uuid.UUID
		author_id    uint32
		tokenGiven   string
		statusCode   int
		errorMessage string
	}{
		{
			// Convert int64 to int first before converting to string
			id:           AuthGstID,
			author_id:    GstUserID,
			tokenGiven:   tokenString,
			statusCode:   204,
			errorMessage: "",
		},
		{
			// When empty token is passed
			id:           AuthGstID,
			author_id:    GstUserID,
			tokenGiven:   "",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			// When incorrect token is passed
			id:           AuthGstID,
			author_id:    GstUserID,
			tokenGiven:   "This is an incorrect token",
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
		{
			id:         uuid.Nil,
			tokenGiven: tokenString,
			statusCode: 400,
		},
		{
			id:           AuthGstID,
			statusCode:   401,
			errorMessage: "Unauthorized",
		},
	}
	for _, v := range gstSample {

		req, _ := http.NewRequest("GET", "/gsts", nil)
		req = mux.SetURLVars(req, map[string]string{"id": v.id.String()})

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.DeleteGst)

		req.Header.Set("Authorization", v.tokenGiven)

		handler.ServeHTTP(rr, req)

		assert.Equal(t, rr.Code, v.statusCode)

		if v.statusCode == 401 && v.errorMessage != "" {

			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("Cannot convert to json: %v", err)
			}
			assert.Equal(t, responseMap["error"], v.errorMessage)
		}
	}
}
