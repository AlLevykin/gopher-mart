package models

type User struct {
	Login    string `json:"login" type:"string"`
	Password string `json:"password" type:"string"`
}
