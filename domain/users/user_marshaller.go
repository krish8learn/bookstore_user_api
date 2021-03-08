package users

import "encoding/json"

type PublicUser struct {
	Id        int64  `json:"id"`
	FirstName string `json:"first_name"`
	//LastName  string `json:"last_name"`
	//Email       string `json:"email"`
	Datecreated string `json:"date_created"`
	Status      string `json:"status"`
	//Password    string `json:"password"`
}

type PrivateUser struct {
	Id          int64  `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	Datecreated string `json:"date_created"`
	Status      string `json:"status"`
	//Password    string `json:"password"`
}

func (user *User) PMarshall(isPublic bool) interface{} {
	//1st approach to convert User to struct to PublicUser
	if isPublic {
		return PublicUser{
			Id:          user.Id,
			FirstName:   user.FirstName,
			Datecreated: user.Datecreated,
			Status:      user.Status,
		}
	}
	//2nd approach to convert User struct to PrivateUser
	userJson, _ := json.Marshal(user)
	var priUs PrivateUser
	json.Unmarshal(userJson, &priUs)
	return priUs
}
