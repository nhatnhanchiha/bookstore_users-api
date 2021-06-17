package users

import (
	"database/sql"
	"fmt"
	"github.com/nhatnhanchiha/bookstore_users-api/datasources/mysql/user_db"
	"github.com/nhatnhanchiha/bookstore_users-api/logger"
	"github.com/nhatnhanchiha/bookstore_users-api/utils/date"
	"github.com/nhatnhanchiha/bookstore_users-api/utils/mysql"
	"github.com/nhatnhanchiha/bookstore_utils-go/rest_errors"
	"log"
)

const (
	queryInsertUser             = "insert into users(first_name, last_name, email, date_created, status, password) value(?, ?, ?, ?, ?, ?);"
	queryGetUser                = "select id, first_name, last_name, email, date_created, status from users where id = ?;"
	queryUpdateUser             = "update users set first_name = ?, last_name = ?, email = ? where id = ?;"
	queryDeleteUser             = "delete from users where id = ?;"
	queryFindByUserByStatus     = "select id, first_name, last_name, email, date_created, status from users where status = ?;"
	queryFindByEmailAndPassword = "select id, first_name, last_name, email, date_created, status from users where email = ? and password = ? and status = ?;"
)

func (user *User) Get() *rest_errors.RestErr {
	stmt, err := user_db.Client.Prepare(queryGetUser)
	if err != nil {
		logger.Error("error when trying to prepare get user by id statement", err)
		return rest_errors.NewInternalServerError("error when trying to get user", rest_errors.NewError("database error"))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("error when trying to close get user by id statement")
		}
	}(stmt)

	result := stmt.QueryRow(user.Id)

	getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)

	if getErr != nil {
		return rest_errors.NewNotFoundError("Cannot find any user matched with provided id")
	}

	return nil
}

func (user *User) Save() *rest_errors.RestErr {
	stmt, err := user_db.Client.Prepare(queryInsertUser)
	if err != nil {
		logger.Error("error when trying to prepare save user statement", err)
		return rest_errors.NewInternalServerError("error when trying to save user", rest_errors.NewError("database err"))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			logger.Error("error when trying to close save user statement", err)
		}
	}(stmt)

	user.DateCreated = date.GetNowDbFormat()

	insertResult, saveErr := stmt.Exec(user.FirstName, user.LastName, user.Email, user.DateCreated, user.Status, user.Password)
	if saveErr != nil {
		return mysql.ParseError(saveErr)
	}

	userId, err := insertResult.LastInsertId()
	if err != nil {
		logger.Error("error when trying to save user", err)
		return rest_errors.NewInternalServerError("error when trying to save user", rest_errors.NewError("database error"))
	}

	user.Id = userId

	return nil
}

func (user *User) Update() *rest_errors.RestErr {
	stmt, err := user_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		logger.Error("error when trying to prepare update user statement", err)
		return rest_errors.NewInternalServerError("error when trying to update user", rest_errors.NewError("database error"))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			logger.Error("error when trying to close update user statement", err)
		}
	}(stmt)

	_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Id)

	if err != nil {
		logger.Error("error when trying to update user", err)
		return rest_errors.NewInternalServerError("error when trying to update user", rest_errors.NewError("database error"))
	}

	return nil
}

func (user *User) Delete() *rest_errors.RestErr {
	stmt, err := user_db.Client.Prepare(queryDeleteUser)
	if err != nil {
		logger.Error("error when trying to prepare delete user by id statement", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", rest_errors.NewError("database error"))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("error when trying to close delete user by id statement")
		}
	}(stmt)

	_, err = stmt.Exec(user.Id)

	if err != nil {
		logger.Error("error when trying to delete user by id", err)
		return rest_errors.NewInternalServerError("error when trying to delete user", rest_errors.NewError("database error"))
	}

	return nil
}

func (user User) FindByStatus(status string) ([]User, *rest_errors.RestErr) {
	stmt, err := user_db.Client.Prepare(queryFindByUserByStatus)
	if err != nil {
		logger.Error("error when trying to prepare find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find users by status", rest_errors.NewError("database error"))

	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("error when trying to close find users by status statement")
		}
	}(stmt)

	rows, err := stmt.Query(status)
	if err != nil {
		logger.Error("error when trying to query find users by status statement", err)
		return nil, rest_errors.NewInternalServerError("error when trying to find users by status", rest_errors.NewError("database error"))
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println("error when trying to close find users by status statement")
		}
	}(rows)

	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status); err != nil {
			logger.Error("error when trying to scan rows were generated by find users by status statement", err)
			return nil, rest_errors.NewInternalServerError("error when trying to find users by status", rest_errors.NewError("database error"))
		}
		results = append(results, user)
	}

	if len(results) == 0 {
		return nil, rest_errors.NewNotFoundError(fmt.Sprintf("No user matching status %s", status))
	}

	return results, nil
}

func (user *User) FindByEmailAndPassword() *rest_errors.RestErr {
	stmt, err := user_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		logger.Error("error when trying to prepare get user by id statement", err)
		return rest_errors.NewInternalServerError("error when trying to find user by email and password", rest_errors.NewError("database error"))
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println("error when trying to close get user by id statement")
		}
	}(stmt)

	result := stmt.QueryRow(user.Email, user.Password, StatusActive)

	getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.DateCreated, &user.Status)

	if getErr != nil {
		return rest_errors.NewInternalServerError("error when trying to find user by email and password", rest_errors.NewError("database error"))
	}

	return nil
}
