package users

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/krish8learn/bookstore_user_api/domain/users"
	"github.com/krish8learn/bookstore_user_api/services"
	"github.com/krish8learn/bookstore_user_api/utils/crypto_utils"
	"github.com/krish8learn/bookstore_user_api/utils/errors"
)

//4th layer
//function of these file takes input, send them to service
//error thrown from this file is invalid input error
func CreateUser(c *gin.Context) {
	var user users.User
	if err := c.ShouldBindJSON(&user); err != nil {
		//error while converting data from the body struct(json from struct)
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}
	result, saveErr := services.CreateUser(user)
	if saveErr != nil {
		c.JSON(saveErr.Status, saveErr)
		//handle the error of user creation
		return
	}
	c.JSON(http.StatusCreated, result.PMarshall(c.GetHeader("X-Public") == "true"))
}

func GetUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id, should be number")
		c.JSON(err.Status, err)
		return
	}

	user, getErr := services.GetUser(userId)
	if getErr != nil {
		c.JSON(getErr.Status, getErr)
		//handle the error of user creation
		return
	}
	c.JSON(http.StatusCreated, user.PMarshall(c.GetHeader("X-Public") == "true"))
}

func UpdateUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("invalid user id, should be number")
		c.JSON(err.Status, err)
		return
	}

	var user users.User

	if err := c.ShouldBindJSON(&user); err != nil {
		//error while converting data from the body struct(json from struct)
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	user.Id = userId
	user.Password = crypto_utils.GetMd5(user.Password)

	isPartial := c.Request.Method == http.MethodPatch

	result, err := services.UpdateUser(isPartial, user)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result.PMarshall(c.GetHeader("X-Public") == "true"))
	c.JSON(http.StatusOK, map[string]string{"status": "updated"})
}

func DeleteUser(c *gin.Context) {
	userId, userErr := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if userErr != nil {
		err := errors.NewBadRequestError("user id should be number")
		c.JSON(err.Status, err)
		return
	}

	var user users.User
	user.Id = userId
	result, err := services.DeleteUser(user)
	if err != nil {
		c.JSON(err.Status, result)
	}
	c.JSON(http.StatusOK, result)
	c.JSON(http.StatusOK, map[string]string{"status": "deleted"})

}

func Search(c *gin.Context) {
	status := c.Query("status")
	users, err := services.Search(status)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}

	result := make([]interface{}, len(users))
	for index, val := range users {
		result[index] = val.PMarshall(c.GetHeader("X-Public") == "true")
	}
	c.JSON(http.StatusOK, result)
}

func Login(c *gin.Context) {
	var req users.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		restErr := errors.NewBadRequestError("invalid json body")
		c.JSON(restErr.Status, restErr)
		return
	}

	result, err := services.LoginUser(req)
	if err != nil {
		c.JSON(err.Status, err)
		return
	}
	c.JSON(http.StatusOK, result.PMarshall(c.GetHeader("X-Public") == "true"))
}
