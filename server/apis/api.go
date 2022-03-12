package apis

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/saurabh-anthwal/dummy/pkg/config"
	"github.com/saurabh-anthwal/dummy/server/models"
	"net/http"
)

var zlog = config.GetLogger()

// HelloResponse is the JSON representation for a customized message
type HelloResponse struct {
	Message string `json:"message"`
}

// HelloWorld returns a basic "Hello World!" message
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	response := HelloResponse{
		Message: "Hello world!",
	}
	jsonResponse(w, response, http.StatusOK)
}


// HelloName returns a personalized JSON message
func HelloName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	response := HelloResponse{
		Message: fmt.Sprintf("Hello %s!", name),
	}


	jsonResponse(w, response, http.StatusOK)
}

// User api singups users
func User(w http.ResponseWriter, r *http.Request) {
	type Signup struct {
		FullName string `validate:"required" json:"full_name"`
		Password string `validate:"required"`
		Email 	 string  `validate:"required,email"`
		Username  string `validate:"required,min=10,max=10"`//phone number as username
	}
	values := &Signup{}
	if err := json.NewDecoder(r.Body).Decode(values); err != nil {
		errJsonResponse(w, fmt.Errorf("failed reading request body, error: %v", err))
		return
	}


	user := &models.User{
		FullName: values.FullName,
		Password: values.Password,
		Email:    values.Email,
		Username: values.Username,
	}

	if err := config.GetDB().Create(user).Error; err != nil {
		zlog.Errorf("error saving user: %v", err)
		errJsonResponse(w, fmt.Errorf("error signing up: %v", err))
	}

	response := HelloResponse{
		Message: fmt.Sprintf("signup success"),
	}
	jsonResponse(w, response, http.StatusOK)
}




