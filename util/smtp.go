package util

import (
	"errors"
	"net/smtp"
)

//https://stackoverflow.com/questions/57783841/how-to-send-email-using-outlooks-smtp-servers
//https://github.com/jordan-wright/email/blob/master/email.go

type loginAuth struct {
	username, password string
}

func OutLookEmailLoginAuth(username, password string) smtp.Auth {
	return &loginAuth{username, password}
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte{}, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unkown fromServer")
		}
	}
	return nil, nil
}
