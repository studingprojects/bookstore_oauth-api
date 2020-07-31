package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/studingprojects/bookstore_oauth-api/src/domain/user"
	"github.com/studingprojects/bookstore_oauth-api/src/utils/errors"
)

var (
	userRestClient = rest.RequestBuilder{
		BaseURL: "http://localhost:8088",
		Timeout: 1000 * time.Millisecond,
	}
)

type UserRepository interface {
	Login(string, string) (*user.User, *errors.RestErr)
}
type userRepository struct {
}

func NewRepository() UserRepository {
	return &userRepository{}
}

func (r *userRepository) Login(username string, password string) (*user.User, *errors.RestErr) {
	request := user.UserLoginRequest{
		Email:    username,
		Password: password,
	}
	response := userRestClient.Post("/users/login", request)
	fmt.Println(response)
	fmt.Println(response.Err)
	if response == nil || response.Response == nil {
		return nil, errors.NewBadRequestError("user service: invalid parameters")
	}
	if response.StatusCode > 299 {
		var restErr errors.RestErr
		if err := json.Unmarshal(response.Bytes(), &restErr); err != nil {
			return nil, errors.NewInternalServerError("user service: could not parse user response")
		}
		return nil, &errors.RestErr{
			Status:  restErr.Status,
			Message: fmt.Sprintf("user service: %s", restErr.Message),
			Error:   "external_service_error",
		}
	}
	var userInfo user.User
	if err := json.Unmarshal(response.Bytes(), &userInfo); err != nil {
		return nil, errors.NewInternalServerError("user service: could not parse user response")
	}

	return &userInfo, nil
}