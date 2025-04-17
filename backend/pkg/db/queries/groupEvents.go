package queries

import ("log"
"database/sql"
datamodels "social-network/pkg/dataModels"
 //"fmt"
 "social-network/pkg/utils"
 "sort"
)




func InsertEvent(groupID int, senderID, title, description, dateTime string) (int, error) {
    dbPath := getDBPath() // Get the database path
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Println("Error opening DB:", err)
        return 0, err
    }
    defer db.Close()

    // Insert the event into the table
    query := `
        INSERT INTO events (group_id, creator_id, title, description, event_date)
        VALUES (?, ?, ?, ?, ?)
    `
    result, err := db.Exec(query, groupID, senderID, title, description, dateTime)
    if err != nil {
        log.Println("Error inserting event:", err)
        return 0, err
    }

    // Get the last inserted event ID using LAST_INSERT_ROWID()
    lastInsertID, err := result.LastInsertId()
    if err != nil {
        log.Println("Error getting last insert ID:", err)
        return 0, err
    }

    return int(lastInsertID), nil
}

func InsertEventOptions(eventID int, option []datamodels.Option) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		INSERT INTO event_options (event_id, option_text, option_id)
		VALUES (?, ?, ?)
	`

	for _, option := range option {
		_, execErr := db.Exec(query, eventID, option.Text,option.ID)
		if execErr != nil {
			log.Printf("Error inserting event option (%s): %v\n", option.Text, execErr)
			return execErr
		}
	}

	return nil
}

func GetEventCreatedAt(groupID int, senderID, title, description string) (string, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return "", err
	}
	defer db.Close()

	var createdAt string

	query := `
		SELECT created_at 
		FROM events 
		WHERE group_id = ? AND creator_id = ? AND title = ? AND description = ?
	`

	err = db.QueryRow(query, groupID, senderID, title, description).Scan(&createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			log.Println("No event found with the given parameters.")
		} else {
			log.Println("Error retrieving event timestamp:", err)
		}
		return "", err
	}

	return createdAt, nil
}


func InsertEventResponse(eventID int, senderID string, optionID int) error {
    dbPath := getDBPath() // Function to get the DB path.
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Println("Error opening DB:", err)
        return err
    }
    defer db.Close()

    // Use INSERT OR REPLACE to update the row if the combination of event_id and user_id already exists
    stmt, err := db.Prepare(`
        INSERT INTO event_participation (event_id, user_id, option_id)
        VALUES (?, ?, ?)
        ON CONFLICT(event_id, user_id) DO UPDATE SET option_id = excluded.option_id
    `)
    if err != nil {
        log.Println("Error preparing statement:", err)
        return err
    }
    defer stmt.Close()

    // Execute the statement to insert or update the row
    _, execErr := stmt.Exec(eventID, senderID, optionID)
    if execErr != nil {
        log.Println("Error executing statement:", execErr)
        return execErr
    }

    log.Println("Insert or update successful!") // Confirm the action
    return nil
}




func GetEventOptionID(eventID int, optionText string) (int, error) {
    dbPath := getDBPath() // Get the database path
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        log.Println("Error opening DB:", err)
        return 0, err
    }
    defer db.Close()

    var optionID int
    query := `
        SELECT id 
        FROM event_options 
        WHERE event_id = ? AND option_text = ? 
        LIMIT 1
    `
    err = db.QueryRow(query, eventID, optionText).Scan(&optionID)
    if err != nil {
        if err == sql.ErrNoRows {
            log.Println("No matching event option found.")
            return 0, nil // No matching option found
        }
        log.Println("Error retrieving event option ID:", err)
        return 0, err
    }

    return optionID, nil
}

func OldGroupEvents(groupID int) ([]datamodels.EventMessage, error) {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT e.id, e.group_id, e.creator_id, e.title, e.description, e.event_date, e.created_at,
		       eo.id, eo.option_text
		FROM events e
		LEFT JOIN event_options eo ON e.id = eo.event_id
		WHERE e.group_id = ?
		ORDER BY e.event_date DESC
	`
	
	rows, err := db.Query(query, groupID)
	if err != nil {
		log.Println("Error querying old events:", err)
		return nil, err
	}
	defer rows.Close()

	eventMap := make(map[int]*datamodels.EventMessage)

	// Process rows and build event details with options
	for rows.Next() {
		var eventID, optionID int
		var groupID int
		var creatorID, title, description, eventDate, createdAt, optionText string

		err := rows.Scan(&eventID, &groupID, &creatorID, &title, &description, &eventDate, &createdAt, &optionID, &optionText)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}
		senderName, err := GetFirstNameById(creatorID)
		if err != nil {
			log.Println("Error retrieving sender name", err)
			return nil, err
		}

		day, err:= utils.ExtractDayFromEvents(eventDate)
		if err != nil {
		log.Println("Error extracting the day", err)
        return nil , err
		}


		// Check if event already exists in the map
		event, exists := eventMap[eventID]
		if !exists {
			event = &datamodels.EventMessage{
				EventID:     eventID,
				GroupID:     groupID,
				SenderID:    creatorID,
				Title:       title,
				Description: description,
				DateTime:    eventDate,
				CreatedAt:   createdAt,
				Options:     []datamodels.Option{},
				Day:day,
				FirstName:senderName,
			}
			eventMap[eventID] = event
		}

		// Add option if it exists
		if optionText != "" {
			event.Options = append(event.Options, datamodels.Option{
				ID:   optionID,
				Text: optionText,
			})
		}
	}

	// Convert map to slice
	var eventHistory []datamodels.EventMessage
	for _, event := range eventMap {
		eventHistory = append(eventHistory, *event)
	}

	sort.Slice(eventHistory, func(i, j int) bool {
		return eventHistory[i].DateTime > eventHistory[j].DateTime
	})
	

	return eventHistory, nil
}

func GetEventResponses(eventID int) ([]datamodels.EventResponseMessage, error) {
	log.Println("the error in this function")
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	query := `
		SELECT ep.user_id, ep.option_id, eo.option_text, ep.event_id, u.first_name
		FROM event_participation ep
		LEFT JOIN event_options eo ON ep.option_id = eo.id
		LEFT JOIN users u ON ep.user_id = u.id
		WHERE ep.event_id = ?
	`

	rows, err := db.Query(query, eventID)
	if err != nil {
		log.Println("Error querying event responses:", err)
		return nil, err
	}
	defer rows.Close()

	var eventResponses []datamodels.EventResponseMessage

	// Process rows and build event responses
	for rows.Next() {
		var userID,firstName string
		var optionID, eventID int
		var optionText string

		err := rows.Scan(&userID, &optionID, &optionText, &eventID, &firstName)
		if err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		eventResponse := datamodels.EventResponseMessage{
			EventID:   eventID,
			OptionID:  optionID,
			SenderID:  userID,
			//OptionText: optionText,
			FirstName: firstName,  // Populate the FirstName field
		}

		eventResponses = append(eventResponses, eventResponse)
	}

	return eventResponses, nil
}



func InsertEventNotification(userID string, eventID int) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		INSERT INTO event_notification (user_id, event_id, status)
		VALUES (?, ?, 'pending')
	`

	_, err = db.Exec(query, userID, eventID)
	if err != nil {
		log.Println("Error inserting event notification:", err)
		return err
	}

	log.Println("Event notification inserted successfully")
	return nil
}


func UpdateEventNotificationStatus(userID string, eventID int) error {
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return err
	}
	defer db.Close()

	query := `
		UPDATE event_notification
		SET status = 'delivered'
		WHERE user_id = ? AND event_id = ?
	`

	_, err = db.Exec(query, userID, eventID)
	if err != nil {
		log.Println("Error updating event notification status:", err)
		return err
	}

	log.Println("Event notification status updated to 'delivered'")
	return nil
}

func GetPendingEventNotifications(userID string) ([]datamodels.EventNotification, error) {
	var eventNotifications []datamodels.EventNotification
	dbPath := getDBPath() // Get the database path
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println("Error opening DB:", err)
		return nil, err
	}
	defer db.Close()

	// SQL query to get event details along with the notification details
	query := `
		SELECT en.event_id, e.title, e.creator_id
		FROM event_notification en
		JOIN events e ON en.event_id = e.id
		WHERE en.user_id = ? AND en.status = 'pending'
	`

	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error querying pending event notifications with event details:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var eventNotification datamodels.EventNotification
		var title string
		var creatorID string

		// Scan the event notification along with the event details (title and creator_id)
		if err := rows.Scan(&eventNotification.EventID, &title, &creatorID); err != nil {
			log.Println("Error scanning event notification row:", err)
			return nil, err
		}

		// Assign the title and creator_id to the event notification struct
		eventNotification.Title = title
		eventNotification.SenderID = creatorID

		eventNotifications = append(eventNotifications, eventNotification)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error iterating over event notification rows:", err)
		return nil, err
	}

	return eventNotifications, nil
}
