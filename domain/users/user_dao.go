package users

import (
	"fmt"
	"strings"

	"github.com/krish8learn/bookstore_user_api/datasources/mysql/users_db"
	"github.com/krish8learn/bookstore_user_api/utils/date_utils"
	"github.com/krish8learn/bookstore_user_api/utils/errors"
	"github.com/krish8learn/bookstore_user_api/utils/mysql_utils"
)

//6th layer
//functions of this file interact with the database
//most of the errors thrown from this file comes from mysql error
const (
	queryInsertUser             = "INSERT INTO users(first_name, last_name, email, date_created, status, password) VALUES(?,?,?,?,?,?);"
	queryGetUser                = "SELECT id, first_name, last_name, email, date_created, status, password FROM users WHERE id =?;"
	queryUpdateUser             = "UPDATE users SET first_name=?, last_name=?, email=?, status=?, password=? WHERE id=?;"
	queryDelteUser              = "DELETE FROM users WHERE id=?;"
	queryFindUserByStatus       = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE status=?;"
	queryFindByEmailAndPassword = "SELECT id, first_name, last_name, email, date_created, status FROM users WHERE email=? AND password=? AND status=?;"
)

func (user *User) Get() *errors.RestErr {
	statement, err := users_db.Client.Prepare(queryGetUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()
	result := statement.QueryRow(user.Id)
	getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Datecreated, &user.Status, &user.Password)
	if getErr != nil {
		return mysql_utils.ParseError(getErr)
	}
	return nil
}

func (user *User) Save() *errors.RestErr {
	statement, err := users_db.Client.Prepare(queryInsertUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()
	user.Datecreated = date_utils.GetNowString()
	insertResult, saveErr := statement.Exec(user.FirstName, user.LastName, user.Email, user.Datecreated, user.Status, user.Password)
	if saveErr != nil {
		return mysql_utils.ParseError(saveErr)
	}
	userId, err := insertResult.LastInsertId()
	if err != nil {
		return mysql_utils.ParseError(saveErr)
	}
	user.Id = userId
	return nil
}

func (user *User) Update() *errors.RestErr {
	statement, err := users_db.Client.Prepare(queryUpdateUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()

	_, updateErr := statement.Exec(user.FirstName, user.LastName, user.Email, user.Status, user.Password, user.Id)
	if updateErr != nil {
		return mysql_utils.ParseError(updateErr)
	}
	return nil
}

func (user *User) Delete() *errors.RestErr {
	statement, err := users_db.Client.Prepare(queryDelteUser)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()
	_, deleteErr := statement.Exec(user.Id)
	if deleteErr != nil {
		return mysql_utils.ParseError(deleteErr)
	}
	return nil
}

func (user *User) FindUserByStatus(status string) ([]User, *errors.RestErr) {
	statement, err := users_db.Client.Prepare(queryFindUserByStatus)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()
	rows, err := statement.Query(status)
	if err != nil {
		return nil, errors.NewInternalServerError(err.Error())
	}
	defer rows.Close()
	results := make([]User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Datecreated, &user.Status); err != nil {
			return nil, mysql_utils.ParseError(err)
		}
		results = append(results, user)
	}
	if len(results) == 0 {
		return nil, errors.NewNotFoundError(fmt.Sprintf("no user mathces the status of %s", status))
	}
	return results, nil
}

func (user *User) FindByEmailAndPassword() *errors.RestErr {
	statement, err := users_db.Client.Prepare(queryFindByEmailAndPassword)
	if err != nil {
		return errors.NewInternalServerError(err.Error())
	}
	defer statement.Close()

	result := statement.QueryRow(user.Email, user.Password, StatusActive)

	if getErr := result.Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Datecreated, &user.Status); getErr != nil {
		if strings.Contains(getErr.Error(), mysql_utils.ErrorNoRows) {
			return errors.NewNotFoundError("invalid user credentials")
		}
		return errors.NewInternalServerError("error when trying to find user")
	}
	return nil
}
