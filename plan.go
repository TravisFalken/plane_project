package main

import "time"

//Plan holds a slice of items and a slice of privileges
type Plan struct {
	PlanID      string
	Title       string
	PlanOwner   string
	Items       []Item
	Privilges   []Privilege
	Completed   bool
	Percentage  int
	CreatedDate string
	LastUpdated string
}

//Check if a plan is completed and also returns if the plan has been changed
func (plan *Plan) checkComplete() (Changed bool) {
	var count int //number of items completed
	// size of the items slice
	size := len(plan.Items)
	//Check that there are actually items
	if size == 0 {
		return false
	}

	//run through all the items and check if it is completed
	for _, item := range plan.Items {
		if item.Completed == true {
			count++
		}
	}

	//calculate percentage completed of plan
	percentage := (count / size) * 100
	//Check if plan has been completed
	if percentage == 100 {
		plan.Completed = true
	} else {
		plan.Completed = false
	}

	//Check that nothing has changed to the plan
	if plan.Percentage == percentage {
		return false
	}

	plan.Percentage = percentage
	return true
}

//Update last update
func (plan *Plan) setLastUpdate() bool {
	plan.LastUpdated = time.Now().Format("2006-01-02")
	return true
}
