package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestUserServiceDoNotResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8088/users/login",
		ReqBody:      `{"email":"xuan.vo@outlook.com","password":"123123"}`,
		RespHTTPCode: -1,
		RespBody:     "",
	})
	userRepository := &userRepository{}
	restUser, err := userRepository.Login("xuan.vo@outlook.com", "123123")

	assert.Nil(t, restUser)
	assert.NotNil(t, err)
	assert.Equal(t, "user service: invalid parameters", err.Message)
	assert.Equal(t, 400, err.Status)
}

func TestUserServiceResponseErrorResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8088/users/login",
		ReqBody:      `{"email":"xuan.vo@outlook.com","password":"123123"}`,
		RespHTTPCode: 400,
		RespBody:     `{"message":"something wrong","status":400,"error":"bad_parameters"}`,
	})
	userRepository := &userRepository{}
	restUser, err := userRepository.Login("xuan.vo@outlook.com", "123123")

	assert.Nil(t, restUser)
	assert.NotNil(t, err)
	assert.Equal(t, "something wrong", err.Message)
	assert.Equal(t, "bad_parameters", err.Error)
	assert.Equal(t, 400, err.Status)
}

func TestUserServiceUnmarshalResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8088/users/login",
		ReqBody:      `{"email":"xuan.vo@outlook.com","password":"123123"}`,
		RespHTTPCode: 400,
		RespBody:     `{other:other_value}`,
	})
	userRepository := &userRepository{}
	restUser, err := userRepository.Login("xuan.vo@outlook.com", "123123")

	assert.Nil(t, restUser)
	assert.NotNil(t, err)
	assert.Equal(t, "user service: could not parse user response", err.Message)
	assert.Equal(t, 500, err.Status)
}

func TestUserServiceFailedResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8088/users/login",
		ReqBody:      `{"email":"xuan.vo@outlook.com","password":"123123"}`,
		RespHTTPCode: 200,
		RespBody:     "",
	})
	userRepository := &userRepository{}
	restUser, err := userRepository.Login("xuan.vo@outlook.com", "123123")

	assert.Nil(t, restUser)
	assert.NotNil(t, err)
	assert.Equal(t, "user service: could not parse user response", err.Message)
	assert.Equal(t, 500, err.Status)
}

func TestUserServiceSuccessResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "http://localhost:8088/users/login",
		ReqBody:      `{"email":"xuan.vo@outlook.com","password":"123123"}`,
		RespHTTPCode: 200,
		RespBody:     `{"id":10,"firstName":"Xuan","lastName":"Vo","email":"xuan.vo02@lazada.com"}`,
	})
	userRepository := &userRepository{}
	restUser, err := userRepository.Login("xuan.vo02@lazada.com", "123123")

	assert.NotNil(t, restUser)
	assert.Nil(t, err)
	assert.Equal(t, "xuan.vo02@lazada.com", restUser.Email)
	assert.Equal(t, "Xuan", restUser.FirstName)
	assert.Equal(t, "Vo", restUser.LastName)
}
