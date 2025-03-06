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

func GetAllPosts(userID string) ([]datamodels.Post, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM posts WHERE user_id = ? OR privacy = 'public' OR allowedUsers LIKE ?", userID, "%"+userID+"%")
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
