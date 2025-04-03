package queries

import (
	"database/sql"
	"errors"
	"log"
	"mime"
	"os"
	"path/filepath"
	datamodels "social-network/pkg/dataModels"
)


func InsertNewGroupComment(comment datamodels.GroupComment) error {
	dbPath := getDBPath()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}

	defer db.Close()

	if comment.CommentText == "" && comment.CommentImage == "" {
		return errors.New("missing required fields")
	}

	// Start transaction
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO group_comments (id, user_id, post_id, content, created_at, image) VALUES (?, ?, ?, ?, ?, ?)", comment.CommentID, comment.UserID, comment.PostID, comment.CommentText, comment.CreatedAt, comment.CommentImage)
	if err != nil {
		return err
	}

	//Update comments count in posts table
	_, err = tx.Exec("UPDATE group_posts SET num_comments = num_comments + 1 WHERE id = ?", comment.PostID)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}



func GetGroupComment(commentID string) (datamodels.GroupComment, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return datamodels.GroupComment{}, err
	}

	defer db.Close()

	var comment datamodels.GroupComment
	err = db.QueryRow("SELECT id, user_id, post_id, content, created_at, image FROM group_comments WHERE id = ?", commentID).Scan(
		&comment.CommentID,
		&comment.UserID,
		&comment.PostID,
		&comment.CommentText,
		&comment.CreatedAt,
		&comment.CommentImage,
	)

	if err != nil {
		return datamodels.GroupComment{}, err
	}

	// If comment has an image, read it from uploads directory
	if comment.CommentImage != "" {
		fullPath := filepath.Join(getUploadPath(), comment.CommentImage)

		// Read the image file
		imageData, err := os.ReadFile(fullPath)
		if err != nil {
			log.Printf("Error reading comment image file: %v", err)
			return datamodels.GroupComment{}, err
		}

		// Get the extension and mime type
		ext := filepath.Ext(fullPath)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		comment.ImageDataURL = imageData
		comment.ImageMimeType = mimeType
	}

	return comment, nil
}



func GetPostGroupComments(postID string) ([]datamodels.GroupComment, error) {
	dbPath := getDBPath()
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Query comments with user data using JOIN
	rows, err := db.Query(`
        SELECT 
            c.id,
            c.user_id,
            c.post_id,
            c.content,
            c.created_at,
            c.image,
            u.nickname
        FROM group_comments c
        JOIN users u ON c.user_id = u.id
        WHERE c.post_id = ?`, postID)
	if err != nil {
		log.Printf("Query error: %v", err)
		return nil, err
	}
	defer rows.Close()

	var comments []datamodels.GroupComment
	for rows.Next() {
		var comment datamodels.GroupComment
		var imageStr sql.NullString // Handle NULL values for image

		err = rows.Scan(
			&comment.CommentID,
			&comment.UserID,
			&comment.PostID,
			&comment.CommentText,
			&comment.CreatedAt,
			&imageStr,
			&comment.UserNickname,
		)
		if err != nil {
			log.Printf("Scan error: %v", err)
			return nil, err
		}

		// Handle NULL image value
		if imageStr.Valid {
			comment.CommentImage = imageStr.String
		}

		// If comment has an image, read it from uploads directory
		if comment.CommentImage != "" {
			fullPath := filepath.Join(getUploadPath(), comment.CommentImage)

			// Read the image file
			imageData, err := os.ReadFile(fullPath)
			if err != nil {
				log.Printf("Error reading comment image file: %v", err)
				continue // Skip image if can't be read but continue with comments
			}

			// Get the extension and mime type
			ext := filepath.Ext(fullPath)
			mimeType := mime.TypeByExtension(ext)
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			comment.ImageDataURL = imageData
			comment.ImageMimeType = mimeType
		}

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows error: %v", err)
		return nil, err
	}
	return comments, nil
}
