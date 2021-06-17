package mysql

import (
	"github.com/VividCortex/mysqlerr"
	"github.com/go-sql-driver/mysql"
	"github.com/nhatnhanchiha/bookstore_utils-go/rest_errors"
	"strings"
)

const (
	errorNoRows = "no rows in result set"
)

func ParseError(err error) *rest_errors.RestErr {
	sqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		if strings.Contains(err.Error(), errorNoRows) {
			return rest_errors.NewNotFoundError("no record matching")
		}
		return rest_errors.NewInternalServerError("error parsing database response", rest_errors.NewError(""))
	}

	switch sqlErr.Number {
	case mysqlerr.ER_DUP_ENTRY:
		return rest_errors.NewBadRequestError("invalid data")
	}

	return rest_errors.NewInternalServerError("error processing request", rest_errors.NewError(""))
}
