package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	atDomain "github.com/studingprojects/bookstore_oauth-api/src/domain/access_token"
	"github.com/studingprojects/bookstore_oauth-api/src/domain/user"
	"github.com/studingprojects/bookstore_oauth-api/src/service/access_token"
	errors "github.com/studingprojects/bookstore_utils-go/rest_errors"
)

type AccessTokenHandler interface {
	GetById(*gin.Context)
	Login(*gin.Context)
	Create(*gin.Context)
}

type accessTokenHandler struct {
	service access_token.Service
}

func NewAccessTokenHandler(service access_token.Service) AccessTokenHandler {
	return &accessTokenHandler{
		service: service,
	}
}

func (h *accessTokenHandler) GetById(c *gin.Context) {
	token, err := h.service.GetById(c.Param("access_token_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	c.JSON(http.StatusOK, token)
	return
}

func (h *accessTokenHandler) Login(c *gin.Context) {
	var loginUser user.UserLoginRequest
	if err := c.ShouldBindJSON(&loginUser); err != nil {
		restErr := errors.NewBadRequestError("invalid json")
		c.JSON(restErr.Status, restErr)
		return
	}
	var userInfo *user.User
	var loginErr *errors.RestErr
	if userInfo, loginErr = h.service.Login(loginUser.Email, loginUser.Password); loginErr != nil {
		c.JSON(loginErr.Status, loginErr)
		return
	}
	c.JSON(http.StatusOK, userInfo)
	return
}

func (h *accessTokenHandler) Create(c *gin.Context) {
	var request atDomain.AccessTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, "invalid json")
		return
	}
	accessToken, tokenErr := h.service.Create(request)
	if tokenErr != nil {
		c.JSON(tokenErr.Status, tokenErr)
		return
	}
	c.JSON(http.StatusOK, accessToken)
	return
}
