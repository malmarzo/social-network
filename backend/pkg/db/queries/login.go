package queries

import ("database/sql"
"log"
)



func GetUserIdByEmail( email string)(string, error){
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return "",err1
	}
	defer db.Close()
	var userID string
	err2 := db.QueryRow("SELECT id FROM users WHERE  email = ?",  email).Scan(&userID)
	if err2 != nil {
		log.Println(err2)
		return "",err2
	}
	return userID, nil
}

func GetPasswordByEmail(email string)(string,error){
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return "",err1
	}
	defer db.Close()
	var password string
	err := db.QueryRow("SELECT password FROM users WHERE email = ?", email).Scan(&password)
	if err != nil {
		log.Println( err)
		return "",err
	}
	return password,nil
}

