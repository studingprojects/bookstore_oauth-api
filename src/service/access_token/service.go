package access_token

import (
	"fmt"

	"github.com/studingprojects/bookstore_oauth-api/src/domain/access_token"
	"github.com/studingprojects/bookstore_oauth-api/src/domain/user"
	"github.com/studingprojects/bookstore_oauth-api/src/repository/db"
	"github.com/studingprojects/bookstore_oauth-api/src/repository/rest"
	errors "github.com/studingprojects/bookstore_utils-go/rest_errors"
)

type Service interface {
	GetById(string) (*access_token.AccessToken, errors.RestErr)
	Login(string, string) (*user.User, errors.RestErr)
	Create(access_token.AccessTokenRequest) (*access_token.AccessToken, errors.RestErr)
}

type service struct {
	dbRepo   db.DbRepository
	userRepo rest.UserRepository
}

func NewService(repo db.DbRepository, userRepo rest.UserRepository) Service {
	return &service{dbRepo: repo, userRepo: userRepo}
}

func (s *service) GetById(id string) (*access_token.AccessToken, errors.RestErr) {
	fmt.Println("RUN HERE")
	return s.dbRepo.GetById(id)
}

func (s *service) Login(email string, password string) (*user.User, errors.RestErr) {
	return s.userRepo.Login(email, password)
}

func (s *service) Create(request access_token.AccessTokenRequest) (*access_token.AccessToken, errors.RestErr) {
	if err := request.Validate(); err != nil {
		return nil, err
	}
	// TODO: Support both grant type password and client_credentials
	userInfo, loginErr := s.userRepo.Login(request.Username, request.Password)
	if loginErr != nil {
		return nil, loginErr
	}
	at := access_token.GetNewAccessToken(userInfo.Id)
	at.Generate()
	fmt.Println(at)
	if err := s.dbRepo.Create(at); err != nil {
		return nil, err
	}

	return &at, nil
}
