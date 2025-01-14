package utils

import (
	"path"
)

func ReplaceExt(filePath, newExt string) string {
	dir := path.Dir(filePath)          // Get the directory part
	base := path.Base(filePath)        // Get the filename with extension
	ext := path.Ext(base)              // Get the current file extension
	name := base[:len(base)-len(ext)]  // Get the filename without the extension
	return path.Join(dir, name+newExt) // Construct the new file path
}
