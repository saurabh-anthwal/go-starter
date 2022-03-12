package models

import "github.com/saurabh-anthwal/dummy/pkg/config"


type User struct {

	FullName string
	Password string
	Email 	 string
	Username  string

}


func init() {
	config.RegisterModel(&User{})
}


