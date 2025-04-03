package queries

import (
	"path/filepath"
	"runtime"
)

//Return the path of the upload folder
func getUploadPath() string {
	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(filepath.Dir(b)) // Adjust path based on actual structure
	return filepath.Join(basePath)
}
