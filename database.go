package main

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

//Connect to database
func connectDatabase() (db *sql.DB) {
	db, err := sql.Open("postgres", "user=postgres password=password dbname=planeCheck sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	return db
}

/////////////VALIDATION////////////////////////////////

//Validate username is unique
func validateUsername(username string) bool {
	//Connect to database
	db := connectDatabase()
	var result string
	//Prepare statement
	stmt, err := db.Prepare("SELECT username FROM _user WHERE username = $1;")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Query statement
	err = stmt.QueryRow(username).Scan(&result)
	if err == sql.ErrNoRows {
		return true
	}

	if err != nil {
		log.Panic(err)
		return false
	}

	if result == username {
		return false
	}

	return true
}

//////////////////////////////USER CRUD///////////////////////////////////////////////////

/////////////////////////////          //////////////////////////////////////////////////

//Create User
func createUser(user User) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare("INSERT INTO _user(username, user_given_name, user_family_name, user_email, user_password) VALUES($1,$2,$3,$4,$5);")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute statment
	_, err = stmt.Exec(user.Username, user.UserGivenName, user.UserFamilyName, user.UserEmail, user.UserPassword)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true

}

//get user
func getUser(filteredBy string, value string) (user User, flag bool) {
	//Connect to Database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare("SELECT username, user_given_name, user_family_name, user_email, user_password  FROM _user WHERE " + filteredBy + " = $1;")
	if err != nil {
		log.Panic(err)
		return user, false
	}

	//Query Database
	err = stmt.QueryRow(value).Scan(&user.Username, &user.UserGivenName, &user.UserFamilyName, &user.UserEmail, &user.UserPassword)
	if err == sql.ErrNoRows {
		log.Println("No user found") // for testing
		return user, false
	}
	if err != nil {
		log.Panic(err)
		return user, false
	}

	return user, true
}

//update user
func updateUser(user User) bool {
	//Connect to database
	db := connectDatabase()

	//Prepare statement
	stmt, err := db.Prepare("UPDATE _user SET user_given_name = $1, user_family_name = $2, user_email = $4, user_password = $5, username = $5 WHERE user_id = $6;")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute statement
	_, err = stmt.Exec(user.UserGivenName, user.UserFamilyName, user.UserEmail, user.UserPassword, user.Username, user.UserID)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true

}

//////////////////////////////PLAN CRUD////////////////////////////////////////

/////////////////////////////         /////////////////////////////////////////

//Add new plan
func addPlan(plan Plan) (planID string, flag bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("INSERT INTO _plan(plan_title, plan_owner, plan_completed, plan_date_created, plan_last_updated, plan_percentage) VALUES($1,$2,$3,$4,$5,$6) RETURNING plan_id;")
	if err != nil {
		log.Panic(err)
		return "-5", false
	}

	//Execute statement
	err = stmt.QueryRow(plan.Title, plan.PlanOwner, plan.Completed, plan.CreatedDate, plan.LastUpdated, plan.Percentage).Scan(&planID)
	if err != nil {
		log.Panic(err)
		return "-5", false
	}

	return planID, true
}

//Get the item from the database
func getPlan(filteredBy string, value string) (plan Plan, flag bool) {
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("SELECT plan_title, plan_owner, plan_completed, plan_id, plan_last_updated, plan_date_created, plan_percentage FROM _plan WHERE " + filteredBy + " = $1;")
	if err != nil {
		log.Panic(err)
		return plan, false
	}

	//Query statement
	err = stmt.QueryRow(value).Scan(&plan.Title, &plan.PlanOwner, &plan.Completed, &plan.PlanID, &plan.LastUpdated, &plan.CreatedDate, &plan.Percentage)
	if err != nil {
		log.Panic(err)
	}

	return plan, true
}

//Update the plan on the database
func updatePlan(plan Plan, planID string) bool {
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("UPDATE _plan SET plan_title = $1, plan_last_updated = $2, plan_percentage = $3, plan_completed = $4 WHERE plan_id = $5;")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute statement
	_, err = stmt.Exec(plan.Title, plan.LastUpdated, plan.Percentage, plan.Completed, planID)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true
}

//Deletes a plan
func deletePlan(planID int) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("DELETE FROM _plan WHERE plan_id = $1;")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute statement
	result, err := stmt.Exec(planID)
	if err != nil {
		log.Panic(err)
		return false
	}

	//Validate that plan was deleted
	_, err = result.RowsAffected()
	if err == sql.ErrNoRows {
		log.Println(err)
		return false
	}

	if err != nil {
		log.Panic(err)
		return false
	}

	return true

}

/////////////////////////////////ITEM CRUD///////////////////////////////////////////

////////////////////////////////          //////////////////////////////////////////

//Create Item
func AddItem(itemOwner string, item Item) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare the statement
	stmt, err := db.Prepare("INSERT INTO _item(item_title, item_description, item_owner, item_completed) VALUES($1,$2,$3,$4)")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute statement
	_, err = stmt.Exec(item.Title, item.Description, item.ItemOwner, item.Completed)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true
}

//get an item
func getItem(itemID string) (item Item) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("SELECT item_title, item_description, item_completed, item_owner FROM _item WHERE item_id = $1;")
	if err != nil {
		log.Panic(err)
		return item
	}

	//Query statment
	err = stmt.QueryRow(itemID).Scan(&item.Title, &item.Description, &item.Completed, &item.ItemOwner)
	if err != nil {
		log.Panic(err)
		return item
	}
	item.ItemID = itemID
	return item
}

//get all its belonging to plan
func getAllItems(planID string) (items []Item, flag bool) {
	//connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("SELECT item_title, item_description, item_completed, item_owner FROM _item WHERE plan_id = $1;")
	if err != nil {
		log.Panic(err)
		return items, false
	}

	//Query statement
	rows, err := stmt.Query(planID)
	if err != nil {
		log.Panic(err)
		return items, false
	}

	var item Item
	//scan the rows
	for rows.Next() {
		err = rows.Scan(&item.Title, &item.Description, &item.Completed, &item.ItemOwner)
		if err != nil {
			log.Panic(err)
			return items, false
		}
		items = append(items, item)
	}

	return items, true
}

////////////////////update an item/////////////////////////////
func updateItem(item Item) (successfull bool) {
	//connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("UPDATE _item SET item_description = $1, item_title = $2, item_completed = $3 WHERE item_id = $4;")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute the statement
	_, err = stmt.Exec(item.Description, item.Title, item.Completed, item.ItemID)
	if err != nil {
		return false
	}

	return true
}

func deleteItem(itemID string) (successful bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//prepare statement
	stmt, err := db.Prepare("DELETE FROM _item WHERE item_id = $1")
	if err != nil {
		log.Panic(err)
		return false
	}

	//execute statement
	_, err = stmt.Exec(itemID)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true
}

/////////////////////////////GROUP CRUD//////////////////////////////////////////////////

//////////////////////////////PRIVILEGES CRUD////////////////////////////////////////

///Add privileges to database : look at adding return the privilege id
func addPrivilege(newPriv Privilege) (success bool) {
	//connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("INSERT INTO _privileges(plan_id, username, write) VALUES($1,$2,$3);")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute statement
	_, err = stmt.Exec(newPriv.PlanID, newPriv.Username, newPriv.Write)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true
}

//Read privilege 		: Look into making more flexible (searching in different ways)
func readPrivilege(typeOfSearch string, key string) (success bool, privilege Privilege) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statment
	stmt, err := db.Prepare("SELECT plan_id, privilege_id, username, write FROM _privileges WHERE $2 = $1;")
	if err != nil {
		log.Panic(err)
		return false, privilege
	}

	//Query statement
	result, err := stmt.Query(typeOfSearch, key)
	if err != nil {
		log.Panic(err)
		return false, privilege
	}
	err = result.Scan(&privilege.PlanID, &privilege.PrivilegeID, &privilege.Username, &privilege.Write)
	if err != nil {
		log.Fatal(err)
	}

	return true, privilege
}

//Get all of the privileges from a plan
func getAllPrivileges(planID string) (success bool, privileges []Privilege) {
	//connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("SELECT privilege_id, username, write FROM _privileges WHERE plan_id = $1;")
	if err != nil {
		log.Panic(err)
		return false, privileges
	}

	//Query database
	result, err := stmt.Query(planID)
	if err != nil {
		log.Panic(err)
		return false, privileges
	}

	var privilege Privilege
	//Scan the result
	for result.Next() {
		err = result.Scan(&privilege.PrivilegeID, &privilege.Username, &privilege.Write)
		if err != nil {
			log.Panic(err)
		}
		privilege.PlanID = planID
		privileges = append(privileges, privilege)
	}

	return true, privileges
}

//update a privilege
func updatePrivilege(newPrivilege Privilege) (success bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("UPDATE _privileges SET write = $1 WHERE privilege_id = $2;")
	if err != nil {
		log.Panic(err)
		return false
	}

	//Execute statement
	_, err = stmt.Exec(newPrivilege.Write, newPrivilege.PrivilegeID)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true

}

//Delete Privilege from database
func deletePrivilege(privilegeToDelete Privilege) (success bool) {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("DELETE FROM _privileges WHERE privilege_id = $1;")

	//Execute statement
	_, err = stmt.Exec(privilegeToDelete.PrivilegeID)
	if err != nil {
		log.Panic(err)
		return false
	}

	return true
}

///LOGIN SECTION/////////
func addSessionToDatabase(sessionID string, username string) bool {
	//Connect to database
	db := connectDatabase()
	defer db.Close()

	//Prepare statement
	stmt, err := db.Prepare("UPDATE _user SET session_id = $1 WHERE username = $2;")
	if err != nil {
		log.Println("Entered Prepate panic add session") // for testing
		log.Panic(err)
		return false
	}

	//Execute statement
	_, err = stmt.Exec(sessionID, username)
	if err != nil {
		log.Println("Ented exec panic in add session") // for testing
		log.Panic(err)
		return false
	}

	return true
}
