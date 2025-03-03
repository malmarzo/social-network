package queries

import ("log"
"database/sql"
datamodels "social-network/pkg/dataModels"
)

// func InsertGroup(title,description, creatorID string)error{
// 	dbPath := getDBPath()
// 	db, err1 := sql.Open("sqlite3", dbPath)
// 	if err1 != nil {
// 		log.Println(err1)
// 		return err1
// 	}
// 	defer db.Close()
	
//     _,err2 := db.Exec("INSERT INTO groups (title, description, creator_id RETURNING id) VALUES (?, ?, ?)",title, description, creatorID)
//     if err2 != nil {
// 		log.Println(err2)
//         return err2
//     }
// 	return nil
// }


func InsertGroup(title, description string, creatorID string) (int, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return 0, err
	}
	defer db.Close()

	var groupID int
	err = db.QueryRow("INSERT INTO groups (title, description, creator_id) VALUES (?, ?, ?) RETURNING id", title, description, creatorID).Scan(&groupID)
	if err != nil {
		log.Println("Error inserting group:", err)
		return 0, err
	}

	return groupID, nil
}

func InviteUser(GroupID int, UserID, InvitedBy string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("INSERT INTO group_members (group_id, user_id, invited_by, status) VALUES (?, ?, ?, 'pending')",GroupID, UserID, InvitedBy)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil

}

func GetUsersList()([]datamodels.User, error){
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return nil,err1
	}
	defer db.Close()
	var users []datamodels.User
	rows, err := db.Query("SELECT id, nickname FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user datamodels.User
		if err := rows.Scan(&user.ID, &user.Nickname); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}