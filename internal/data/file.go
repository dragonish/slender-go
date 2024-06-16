package data

import (
	"errors"
	"io/fs"
	"os"
	"slender/internal/logger"
)

// IsPathExists returns true when the path exists.
func IsPathExists(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, fs.ErrNotExist)
}

// DeleteFile deletes specified file.
//
// No operation when the file does not exist.
func DeleteFile(path string) error {
	log := logger.New("path", path)

	if IsPathExists(path) {
		err := os.Remove(path)
		if err != nil {
			return log.Err("delete file error", err)
		}
	} else {
		log.Debug("the path does not exist")
	}

	return nil
}
