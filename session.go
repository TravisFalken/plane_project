package main

import (
	uuid "github.com/satori/go.uuid"
)

//get new session id
func newSessionID() string {
	id, _ := uuid.NewV4()
	return id.String()
}
