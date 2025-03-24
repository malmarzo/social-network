package queries

import (
	"database/sql"
	"log"
	"mime"
	"os"
	"path/filepath"
	datamodels "social-network/pkg/dataModels"
)

// Returns the number of posts for a user
func GetNumOfPosts(userID string) (int, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return 0, err
	}

	defer db.Close()

	var numOfPOsts int
	err = db.QueryRow("SELECT COUNT(*) FROM posts WHERE user_id = ?", userID).Scan(&numOfPOsts)
	if err != nil {
		return 0, err
	}

	return numOfPOsts, nil
}

func InsertNewPost(post datamodels.Post) error {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Println(err)
		return err
	}
	defer db.Close()

	_, err = db.Exec(`
        INSERT INTO posts (
            id, user_id, user_name, title, 
            content, image, privacy, created_at, 
            num_likes, num_dislikes, num_comments, 
            allowedUsers
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		post.PostID,
		post.UserID,
		post.UserNickname, // This field will be stored in user_name column
		post.PostTitle,
		post.Content,
		post.PostImage,
		post.PostPrivacy,
		post.CreatedAt,
		post.NumOfLikes,
		post.NumOfDislikes,
		post.NumOfComments,
		post.AllowedUsers,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetAllPosts(userID, tab string) ([]datamodels.Post, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rows *sql.Rows
	var errFetch error

	baseQuery := `
    SELECT p.* FROM posts p
    LEFT JOIN (
        SELECT DISTINCT following_id FROM followers 
        WHERE follower_id = ? AND status = 'accepted'
    ) f ON p.user_id = f.following_id
    WHERE 
        (p.user_id = ? OR
        p.privacy = 'public' OR
        (p.privacy = 'almost_private' AND f.following_id IS NOT NULL) OR
        (p.privacy = 'private' AND p.allowedUsers LIKE ?))
    GROUP BY p.id
    `

	switch tab {
	case "trending":
		rows, errFetch = db.Query(baseQuery+` 
            ORDER BY (p.num_likes + p.num_dislikes + p.num_comments) DESC`,
			userID, userID, "%"+userID+"%")
	case "my-posts":
		rows, errFetch = db.Query("SELECT * FROM posts WHERE user_id = ? ORDER BY created_at DESC", userID)
	default: // "latest" or empty
		rows, errFetch = db.Query(baseQuery+`ORDER BY p.created_at DESC`,
			userID, userID, "%"+userID+"%")
	}

	if errFetch != nil {
		return nil, errFetch
	}
	defer rows.Close()

	var posts []datamodels.Post
	for rows.Next() {
		var post datamodels.Post
		err = rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.UserNickname,
			&post.PostTitle,
			&post.Content,
			&post.PostImage,
			&post.PostPrivacy,
			&post.CreatedAt,
			&post.NumOfLikes,
			&post.NumOfDislikes,
			&post.NumOfComments,
			&post.AllowedUsers,
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

func GetPostInteractionStats(postID string) (datamodels.PostInteractions, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return datamodels.PostInteractions{}, err
	}
	defer db.Close()

	var postInteractions datamodels.PostInteractions
	err = db.QueryRow("SELECT num_likes, num_dislikes, num_comments FROM posts WHERE id = ?", postID).Scan(&postInteractions.Likes, &postInteractions.Dislikes, &postInteractions.Comments)
	if err != nil {
		return datamodels.PostInteractions{}, err
	}

	return postInteractions, nil
}

func LikePost(postID string, userID string) error {
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
	err = tx.QueryRow("SELECT type FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&currentType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Handle different cases
	switch currentType {
	case "like": // Remove like
		_, err = tx.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE posts SET num_likes = CASE WHEN num_likes > 0 THEN num_likes - 1 ELSE 0 END WHERE id = ?", postID)

	case "dislike": // Change dislike to like
		_, err = tx.Exec("UPDATE likes SET type = 'like' WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
            UPDATE posts 
            SET num_likes = CASE WHEN num_likes >= 0 THEN num_likes + 1 ELSE 1 END,
                num_dislikes = CASE WHEN num_dislikes > 0 THEN num_dislikes - 1 ELSE 0 END 
            WHERE id = ?`, postID)

	default: // New like
		_, err = tx.Exec("INSERT INTO likes (user_id, post_id, type) VALUES (?, ?, 'like')", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE posts SET num_likes = num_likes + 1 WHERE id = ?", postID)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

func DislikePost(postID string, userID string) error {
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
	err = tx.QueryRow("SELECT type FROM likes WHERE user_id = ? AND post_id = ?", userID, postID).Scan(&currentType)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	// Handle different cases
	switch currentType {
	case "dislike": // Remove dislike
		_, err = tx.Exec("DELETE FROM likes WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE posts SET num_dislikes = CASE WHEN num_dislikes > 0 THEN num_dislikes - 1 ELSE 0 END WHERE id = ?", postID)

	case "like": // Change like to dislike
		_, err = tx.Exec("UPDATE likes SET type = 'dislike' WHERE user_id = ? AND post_id = ?", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec(`
            UPDATE posts 
            SET num_dislikes = CASE WHEN num_dislikes >= 0 THEN num_dislikes + 1 ELSE 1 END,
                num_likes = CASE WHEN num_likes > 0 THEN num_likes - 1 ELSE 0 END 
            WHERE id = ?`, postID)

	default: // New dislike
		_, err = tx.Exec("INSERT INTO likes (user_id, post_id, type) VALUES (?, ?, 'dislike')", userID, postID)
		if err != nil {
			return err
		}
		_, err = tx.Exec("UPDATE posts SET num_dislikes = num_dislikes + 1 WHERE id = ?", postID)
	}

	if err != nil {
		return err
	}

	return tx.Commit()
}

func GetProfilePosts(profileID, userID string, isMyProfile bool) ([]datamodels.Post, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var rows *sql.Rows
	if isMyProfile {
		rows, err = db.Query(`
            SELECT * FROM posts 
            WHERE user_id = ? 
            GROUP BY id 
            ORDER BY created_at DESC`, profileID)
	} else {
		rows, err = db.Query(`
            SELECT p.* FROM posts p
            LEFT JOIN followers f ON p.user_id = f.following_id AND f.follower_id = ?
            WHERE p.user_id = ? 
            AND (
                p.privacy = 'public' 
                OR (p.privacy = 'almost_private' AND f.status = 'accepted')
                OR (p.privacy = 'private' AND p.allowedUsers LIKE ?)
            )
            GROUP BY p.id
            ORDER BY p.created_at DESC`,
			userID, profileID, "%"+userID+"%")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []datamodels.Post
	for rows.Next() {
		var post datamodels.Post
		err = rows.Scan(
			&post.PostID,
			&post.UserID,
			&post.UserNickname,
			&post.PostTitle,
			&post.Content,
			&post.PostImage,
			&post.PostPrivacy,
			&post.CreatedAt,
			&post.NumOfLikes,
			&post.NumOfDislikes,
			&post.NumOfComments,
			&post.AllowedUsers,
		)
		if err != nil {
			return nil, err
		}

		// Handle post image if exists
		if post.PostImage != "" {
			fullPath := filepath.Join(getUploadPath(), post.PostImage)
			imageData, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Error reading post image file: %v", err)
				continue
			}

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
