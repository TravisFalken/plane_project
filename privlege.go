package main

//Privilege struct is for the privileges of each item
type Privilege struct {
	PrivilegeID string
	Username    string
	PlanID      string
	Write       bool
}
