package queries

import (
	"path/filepath"
	"runtime"
)

// Will return the path to the database file
func getDBPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(b)) // Adjust path based on actual structure
	return filepath.Join(basePath, "sqlite", "social_network.db")
}

//Return the path of the upload folder
func getUploadPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(b)) // Adjust path based on actual structure
	return filepath.Join(basePath)
}
