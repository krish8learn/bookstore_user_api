package services

import (
	"github.com/krish8learn/bookstore_user_api/domain/users"
	"github.com/krish8learn/bookstore_user_api/utils/crypto_utils"
	"github.com/krish8learn/bookstore_user_api/utils/errors"
)

//5th layer
//functions of this file apply the business logic
//error thrown from this file comes from the user_dao file
func GetUser(userID int64) (*users.User, *errors.RestErr) {
	if userID <= 0 {
		return nil, errors.NewBadRequestError("invalid user id")
	}

	result := &users.User{Id: userID}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return result, nil
}

func CreateUser(user users.User) (*users.User, *errors.RestErr) {

	if err := user.Validate(); err != nil {
		return nil, err
	}
	user.Status = users.StatusActive
	user.Password = crypto_utils.GetMd5(user.Password)
	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateUser(isPartial bool, user users.User) (*users.User, *errors.RestErr) {
	current, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	/*if err := user.Validate(); err != nil {
		return nil, err
	}*/

	if isPartial {
		if user.FirstName != "" {
			current.FirstName = user.FirstName
		}
		if user.LastName != "" {
			current.LastName = user.LastName
		}
		if user.Email != "" {
			current.Email = user.Email
		}
		if user.Password != "" {
			current.Status = user.Status
		}
		if user.Password != "" {
			current.Password = user.Password
		}
	} else {
		current.Id = user.Id
		current.FirstName = user.FirstName
		current.LastName = user.LastName
		current.Email = user.Email
		current.Status = user.Status
		current.Password = user.Password
	}

	if err := current.Update(); err != nil {
		return nil, err
	}
	return current, nil

}

func DeleteUser(user users.User) (*users.User, *errors.RestErr) {
	//fetching the data before deleting
	current, err := GetUser(user.Id)
	if err != nil {
		return nil, err
	}

	if err := current.Delete(); err != nil {
		return nil, err
	}
	return current, nil
}

func Search(status string) ([]users.User, *errors.RestErr) {
	findstatususer := &users.User{}
	users, err := findstatususer.FindUserByStatus(status)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func LoginUser(req users.LoginRequest) (*users.User, *errors.RestErr) {
	dao := &users.User{
		Email:    req.Email,
		Password: crypto_utils.GetMd5(req.Password),
	}
	if err := dao.FindByEmailAndPassword(); err != nil {
		return nil, err
	}
	return dao, nil
}
