package db

import (
	"fmt"

	"github.com/gocql/gocql"
	"github.com/studingprojects/bookstore_oauth-api/src/client/cassandra"
	"github.com/studingprojects/bookstore_oauth-api/src/domain/access_token"
	errors "github.com/studingprojects/bookstore_utils-go/rest_errors"
)

const (
	queryGetAccessTokenQuery = "SELECT access_token, user_id, client_id, expires FROM access_token WHERE access_token=?;"
	queryCreateAccessToken   = "INSERT INTO access_token (access_token, user_id, client_id, expires) VALUES (?, ?, ?, ?);"
)

type DbRepository interface {
	GetById(string) (*access_token.AccessToken, errors.RestErr)
	Create(access_token.AccessToken) errors.RestErr
}

type dbRepository struct{}

func NewRepository() DbRepository {
	return &dbRepository{}
}

func (r *dbRepository) GetById(id string) (*access_token.AccessToken, errors.RestErr) {
	var result access_token.AccessToken
	if err := cassandra.GetSession().Query(queryGetAccessTokenQuery, id).Consistency(gocql.One).Scan(&result.AccessToken, &result.UserId, &result.ClientId, &result.Expires); err != nil {
		if err == gocql.ErrNotFound {
			return nil, errors.NewNotFounfError("token not found")
		}
		return nil, errors.NewInternalServerError("DB error", err)
	}
	return &result, nil
}

func (r *dbRepository) Create(at access_token.AccessToken) errors.RestErr {
	if err := cassandra.GetSession().Query(queryCreateAccessToken,
		at.AccessToken,
		at.UserId,
		at.ClientId,
		at.Expires,
	).Exec(); err != nil {
		return errors.NewInternalServerError(
			fmt.Sprintf("error when trying to save access token in database: %s", err.Error()),
			err)
	}
	return nil
}
