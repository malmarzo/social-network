package queries
import (
	"log"
"database/sql"
datamodels "social-network/pkg/dataModels"
//"fmt"
//"social-network/pkg/utils"
"os"
"path/filepath"
	"mime"
)


func InsertNewGroupPost(post datamodels.GroupPost) error {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
        INSERT INTO group_posts (
            id, group_id, user_id, user_name, title, 
            content, image, created_at, 
            num_likes, num_dislikes, num_comments 
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		post.PostID,
		post.GroupID,
		post.UserID,
		post.UserNickname, // This field will be stored in user_name column
		post.PostTitle,
		post.Content,
		post.PostImage,
		//post.PostPrivacy,
		post.CreatedAt,
		post.NumOfLikes,
		post.NumOfDislikes,
		post.NumOfComments,
		//post.AllowedUsers,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}


func GetAllGroupPosts(userID string,groupID int, tab string) ([]datamodels.GroupPost, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rows *sql.Rows
	var errFetch error

	baseQuery := `
    SELECT gp.* FROM group_posts gp
    WHERE gp.group_id = ?
    `

	switch tab {
	case "trending":
		rows, errFetch = db.Query(baseQuery+` 
            ORDER BY (p.num_likes + p.num_dislikes + p.num_comments) DESC`,
			groupID,)
	case "my-posts":
		rows, errFetch = db.Query("SELECT * FROM group_posts WHERE user_id = ? AND group_id = ? ORDER BY created_at DESC", userID, groupID)
	default: // "latest" or empty
		rows, errFetch = db.Query(baseQuery+`ORDER BY p.created_at DESC`,
			groupID,)
	}

	if errFetch != nil {
		return nil, errFetch
	}
	defer rows.Close()

	var posts []datamodels.GroupPost
	for rows.Next() {
		var post datamodels.GroupPost
		err = rows.Scan(
			&post.PostID,
			&post.GroupID,
			&post.UserID,
			&post.PostTitle,
			&post.Content,
			&post.PostImage,
			&post.NumOfLikes,
			&post.NumOfDislikes,
			&post.NumOfComments,
			&post.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		// If post has an image, read it from uploads directory
		if post.PostImage != "" {
			fullPath := filepath.Join(getUploadPath(), post.PostImage)

			// Read the image file
			imageData, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Error reading post image file: %v", err)
				continue // Skip image if can't be read but continue with posts
			}

			// Get the extension and mime type
			ext := filepath.Ext(fullPath)
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			post.ImageDataURL = imageData
			post.ImageMimeType = mimeType
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}



func GetGroupPostInteractionStats(groupPostID string) (datamodels.GroupPostInteractions, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return datamodels.GroupPostInteractions{}, err
	}
	defer db.Close()

	var groupPostInteractions datamodels.GroupPostInteractions
	err = db.QueryRow("SELECT num_likes, num_dislikes, num_comments FROM group_posts WHERE id = ?", groupPostID).Scan(&groupPostInteractions.Likes, &groupPostInteractions.Dislikes, &groupPostInteractions.Comments)
	if err != nil {
		return datamodels.GroupPostInteractions{}, err
	}

	return groupPostInteractions, nil
}



func LikeGroupPost(postID string, userID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check current interaction status
	var currentType string
	err = tx.QueryRow("SELECT type FROM group_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&currentType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Handle different cases
	switch currentType {
	case "like": // Remove like
		_, err = tx.Exec("DELETE FROM group_likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE group_posts SET num_likes = CASE WHEN num_likes > 0 THEN num_likes - 1 ELSE 0 END WHERE id = ?", postID)

	case "dislike": // Change dislike to like
		_, err = tx.Exec("UPDATE group_likes SET type = 'like' WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
            UPDATE group_posts 
            SET num_likes = CASE WHEN num_likes >= 0 THEN num_likes + 1 ELSE 1 END,
                num_dislikes = CASE WHEN num_dislikes > 0 THEN num_dislikes - 1 ELSE 0 END 
            WHERE id = ?`, postID)

	default: // New like
		_, err = tx.Exec("INSERT INTO group_likes (user_id, post_id, type) VALUES (?, ?, 'like')", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE group_posts SET num_likes = num_likes + 1 WHERE id = ?", postID)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}



func DislikeGroupPost(postID string, userID string) error {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Check current interaction status
	var currentType string
	err = tx.QueryRow("SELECT type FROM group_likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&currentType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Handle different cases
	switch currentType {
	case "dislike": // Remove dislike
		_, err = tx.Exec("DELETE FROM group_likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE group_posts SET num_dislikes = CASE WHEN num_dislikes > 0 THEN num_dislikes - 1 ELSE 0 END WHERE id = ?", postID)

	case "like": // Change like to dislike
		_, err = tx.Exec("UPDATE group_likes SET type = 'dislike' WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
            UPDATE group_posts 
            SET num_dislikes = CASE WHEN num_dislikes >= 0 THEN num_dislikes + 1 ELSE 1 END,
                num_likes = CASE WHEN num_likes > 0 THEN num_likes - 1 ELSE 0 END 
            WHERE id = ?`, postID)

	default: // New dislike
		_, err = tx.Exec("INSERT INTO group_likes (user_id, post_id, type) VALUES (?, ?, 'dislike')", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE group_posts SET num_dislikes = num_dislikes + 1 WHERE id = ?", postID)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}
