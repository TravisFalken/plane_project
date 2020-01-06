package main

type Item struct {
	ItemID      string
	Title       string
	Description string
	PlanID      string
	ItemOwner   string
	Privilges   []Privilege
	Completed   bool
}
