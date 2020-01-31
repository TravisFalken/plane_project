package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

//Handles the splash page
func index(w http.ResponseWriter, r *http.Request) {
	if !getloggedIn(r) {
		fmt.Fprintf(w, "Splash page test") // for testing splash page
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}

}

//Shows the home page
func homePage(w http.ResponseWriter, r *http.Request) {
	if getloggedIn(r) {
		fmt.Fprintf(w, "This is the home page") // for testing
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Show the signup page
func showSignUp(w http.ResponseWriter, r *http.Request) {
	if !getloggedIn(r) {
		fmt.Fprintf(w, "This is the signup page") // for testing
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

//Show the login page
func loginPage(w http.ResponseWriter, r *http.Request) {
	if !getloggedIn(r) {
		fmt.Fprintf(w, "this is the login page") //For testing
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

//Signup the user
func signUp(w http.ResponseWriter, r *http.Request) {
	var newUser User

	//Get values from the form
	newUser.UserGivenName = r.FormValue("given_name")
	newUser.UserFamilyName = r.FormValue("family_name")
	newUser.UserEmail = r.FormValue("email")
	newUser.UserPassword = r.FormValue("password")
	newUser.Username = r.FormValue("username")

	//Validate that required inputs have been enterd
	if validateInput(newUser.UserGivenName) || validateInput(newUser.UserEmail) || validateInput(newUser.UserPassword) || validateInput(newUser.Username) {
		//Validate the username is unique
		if validateUsername(newUser.Username) {
			//Add user to the database
			if createUser(newUser) {
				fmt.Fprintf(w, "user successfuly created") // for testing
			} else {
				fmt.Fprintf(w, "The user has not successfuly created") // for testing
			}
		} else {
			fmt.Fprintf(w, "Username is not unique") // for testing
		}
	}
}

//login the user
func loginUser(w http.ResponseWriter, r *http.Request) {
	//check that user is still logged in
	if !getloggedIn(r) {
		//Get form values
		username := r.FormValue("username")
		password := r.FormValue("password")
		log.Println("Username: " + username) // for testing
		log.Println("Password: " + password) // for testing
		//check login details are correct
		if validateLoginDetails(username, password) {
			//Make sure adding session was successful
			if addSession(w, r, username) {
				http.Redirect(w, r, "/home", http.StatusSeeOther)
			} else {
				log.Panic("Failed to add session to user")
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}

		} else {
			fmt.Fprintf(w, "Incorrect login details") //For testing still need to fix
		}
	} else {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	}
}

//Log user out
func logout(w http.ResponseWriter, r *http.Request) {
	if getloggedIn(r) {
		//Remove the cookie session from user
		if removeSession(w, r) {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			log.Panic("Could not log user out!")
			fmt.Fprintf(w, "An error occured and we could not log you out!") // for testing
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

/////////////////////PLAN SECTION////////////////////////////////

////////////////////            /////////////////////////////////

//Creates a new plan
func createPlan(w http.ResponseWriter, r *http.Request) {
	if getloggedIn(r) {

		var newPlan Plan
		//get the title from the form
		newPlan.Title = r.FormValue("plan_title")
		//get the user creating the
		user, found := getUser("session_id", getSession(r))
		if found {
			//Validate that user actually put in values
			if validateInput(newPlan.Title) {
				//Get the form values
				newPlan.Title = r.FormValue("plan_title")
				newPlan.PlanOwner = user.Username
				newPlan.CreatedDate = (time.Now()).Format("2006-01-02")
				newPlan.Completed = false
				newPlan.Percentage = 0
				newPlan.setLastUpdate()
				//Add the plan to the database
				planID, success := addPlan(newPlan)
				newPlan.PlanID = planID
				if success {
					fmt.Fprintf(w, "Plan was successfully created") // for testing
				} else {
					fmt.Fprintf(w, "could not create new plan") // for testing
				}
			} else {
				fmt.Fprintf(w, "Please enter the value") // for testing
			}
		} else {
			fmt.Fprintf(w, "Could not create new Plan") //for testing
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//Create Plan Page
func createPlanPage(w http.ResponseWriter, r *http.Request) {
	if getloggedIn(r) {
		fmt.Fprintf(w, "This is the create plan page")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Shows plan
func showPlan(w http.ResponseWriter, r *http.Request) {
	if getloggedIn(r) {
		fmt.Fprintf(w, "Show plan page")
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//edits a plan
func editPlan(w http.ResponseWriter, r *http.Request) {
	if getloggedIn(r) {
		var plan Plan
		newTitle := r.FormValue("title")
		//Make sure input is valid

		if validateInput(newTitle) {
			//Make sure user is owner
			if validOwner(r) {
				plan.Title = newTitle
				plan.setLastUpdate()
				plan.checkComplete() // for testing
				id := getID(r)
				//updates the plan and returns if it is successful
				if updatePlan(plan, id) {
					fmt.Fprintf(w, "Plan has been updated") //for testing
				} else {
					fmt.Fprintf(w, "Plan could not updated") //for testing
				}
			} else {
				fmt.Fprintf(w, "You do not have the permission to change the plan")
			}
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

//Delete user plan
func deleteUserPlan(w http.ResponseWriter, r *http.Request) {
	//Make sure user is logged in
	if getloggedIn(r) {
		//Validate that user is plan owner
		if validOwner(r) {
			idString := getID(r)
			id, err := strconv.Atoi(idString)
			if err != nil {
				log.Panic(err)
			}
			//delete a plan and make sure it is successfully deleted
			if deletePlan(id) {
				fmt.Fprintf(w, "Plan was successfully deleted") // for testing
			} else {
				fmt.Fprintf(w, "Plan was not successfully deleted") // for testing
			}
		} else {
			fmt.Fprintf(w, "You cannot delete the note because you are not the plan owner") //for testing
		}
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//Validate owner of plan
func validOwner(r *http.Request) bool {
	//get user from session
	user, found := getUser("session_id", getSession(r))

	//get plan id
	planID := getID(r)
	//make sure user is found
	if found == false {
		log.Print("could not find user") //for testing
		log.Print("Failed to find user from session id")
		return false
	}
	//get the plan for the database
	plan, result := getPlan("plan_id", planID)

	//see if plan was found
	if result == false {
		return false
	}

	//check that the user is the owner
	if user.Username == plan.PlanOwner {
		return true
	}

	return false
}

//get the sesssion cookie value
func getSession(r *http.Request) (sessionID string) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		sessionID = " "
		return sessionID
	}
	sessionID = sessionCookie.Value
	return sessionID
}

//////////////////////////ITEM SECTION//////////////////////////////////

/////////////////////////             /////////////////////////////////

//Adds item to a plan
func addItemToPlan(w http.ResponseWriter, r *http.Request) {
	//check user is still logged on
	if getloggedIn(r) {
		//check user is note owner
		//need to finish up
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

//////////////CONTROLLER VALIDATE SECTION///////////////////////

//Validate login details
//Fix weird bug refusing to log in even though details are correct
func validateLoginDetails(username string, password string) bool {
	user, found := getUser("username", username)
	//check that there is a user with the username
	if found {
		//Check that passwords match
		if password == user.UserPassword {
			return true
		}
		return false
	}
	return false

}

//Validate that user has entered an input
func validateInput(input string) bool {
	if input == "" || input == " " {
		return false
	}
	return true
}

//Validate logged in
func getloggedIn(r *http.Request) bool {
	sessionID := getSession(r)
	_, found := getUser("session_id", sessionID)
	return found
}

//Add session cookie
func addSession(w http.ResponseWriter, r *http.Request, username string) bool {
	//get new session id
	sessionID := newSessionID()
	sessionCookie := &http.Cookie{
		Name:  "session",
		Value: sessionID,
	}
	http.SetCookie(w, sessionCookie)
	if addSessionToDatabase(sessionID, username) {
		return true
	}
	return false
}

//Remove the session cookie from the user
func removeSession(w http.ResponseWriter, r *http.Request) bool {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Modify the cookie
	sessionCookie = &http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
	}
	//Set the cookie
	http.SetCookie(w, sessionCookie)
	return true
}

//gets an id from the url
func getID(r *http.Request) (id string) {
	id = mux.Vars(r)["id"]
	return id
}

//func check if user has write access to the notes
func validateWriteAccess(r *http.Request) bool {
	
}
