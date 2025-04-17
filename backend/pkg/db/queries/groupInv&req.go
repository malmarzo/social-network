package queries

import ("log"
"database/sql"
datamodels "social-network/pkg/dataModels"
 //"fmt"
// "social-network/pkg/utils"
// "sort"
)




func InviteUser(GroupID int, UserID, InvitedBy string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("INSERT INTO group_members (group_id, user_id, invited_by, status, type) VALUES (?, ?, ?, 'pending','invitation')",GroupID, UserID, InvitedBy)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil

}



func RequestToJoin(GroupID int, UserID, groupCreator string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("INSERT INTO group_members (group_id, user_id, invited_by, status, type) VALUES (?, ?, ?, 'pending', 'request')",GroupID, UserID, groupCreator)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil

}




func AcceptInvitation(GroupID int, UserID, InvitedBy string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ? AND invited_by = ?",
        "accepted", GroupID, UserID, InvitedBy,)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil

}

func DeclineInvitation(GroupID int, UserID, InvitedBy string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("UPDATE group_members SET status = ? WHERE group_id = ? AND user_id = ? AND invited_by = ?",
        "declined", GroupID, UserID, InvitedBy,)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil
}




func AcceptRequest( InvitedBy string, GroupID int, UserID string) error {
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()

	_, err2 := db.Exec("UPDATE group_members SET status = ? WHERE invited_by = ? AND group_id = ? AND user_id = ?",
		"accepted", InvitedBy, GroupID, UserID) // Set invited_by to a provided value

	if err2 != nil {
		log.Println(err2)
		return err2
	}
	return nil
}


func DeclineRequest( InvitedBy string, GroupID int, UserID string)error{
	dbPath := getDBPath()
	db, err1 := sql.Open("sqlite3", dbPath)
	if err1 != nil {
		log.Println(err1)
		return err1
	}
	defer db.Close()
    _, err2 := db.Exec("UPDATE group_members SET status = ? WHERE invited_by = ? AND group_id = ? AND user_id = ?",
		"declined", InvitedBy, GroupID, UserID)
    if err2 != nil {
        log.Println(err2)
		return err2
    }
	return nil
}


func GetPendingInvitations(userID string) ([]datamodels.Invite, error) {
	var invites []datamodels.Invite
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT group_id, user_id, invited_by 
		FROM group_members 
		WHERE user_id = ? AND status = 'pending' AND type = 'invitation'`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error querying pending invites:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var invite datamodels.Invite
		if err := rows.Scan(&invite.GroupID, &invite.UserID, &invite.InvitedBy); err != nil {
			log.Println("Error scanning invite row:", err)
			return nil, err
		}
		invites = append(invites, invite)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over invite rows:", err)
		return nil, err
	}
	

	return invites, nil
}


func GetPendingRequests(userID string) ([]datamodels.Request, error) {
	var requests []datamodels.Request
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT group_id, invited_by, user_id
		FROM group_members 
		WHERE invited_by = ? AND status = 'pending' AND type = 'request'`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error querying pending invites:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var request datamodels.Request
		if err := rows.Scan(&request.GroupID, &request.GroupCreator,&request.UserID); err != nil {
			log.Println("Error scanning invite row:", err)
			return nil, err
		}
		requests = append(requests, request)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over invite rows:", err)
		return nil, err
	}

	return requests, nil
}



func GetUserInvitationList(userID string, groupID int) ([]datamodels.User, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("DB connection error:", err)
		return nil, err
	}
	defer db.Close()

	query := `
	SELECT id, nickname
	FROM users
	WHERE id != ?
	  AND id NOT IN (
		SELECT user_id FROM group_members
		WHERE group_id = ? AND status IN ('pending', 'accepted')
	  )
	  AND id != (
		SELECT creator_id FROM groups WHERE id = ?
	  );
	`

	rows, err := db.Query(query, userID, groupID, groupID)
	if err != nil {
		log.Println("Query error:", err)
		return nil, err
	}
	defer rows.Close()

	var users []datamodels.User

	for rows.Next() {
		var user datamodels.User
		if err := rows.Scan(&user.ID, &user.Nickname); err != nil {
			log.Println("Row scan error:", err)
			continue
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Println("Rows error:", err)
		return nil, err
	}

	return users, nil
}
