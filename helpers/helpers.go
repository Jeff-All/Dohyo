package helpers

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// HigherLogLevel - returns the higher logrus.Level of the two log levels submitted
func HigherLogLevel(config string, cli string) (log.Level, error) {
	config = strings.ToLower(config)
	cli = strings.ToLower(cli)
	if config == "debug" || cli == "debug" {
		return log.DebugLevel, nil
	}
	if config == "info" || cli == "info" {
		return log.InfoLevel, nil
	}
	if config == "warn" || cli == "warn" {
		return log.WarnLevel, nil
	}
	if config == "error" || cli == "error" {
		return log.ErrorLevel, nil
	}
	return log.ErrorLevel, fmt.Errorf("invalid log level: %s && %s", config, cli)
}

// DeleteFile - Deletes the file
func DeleteFile(filename string) error {
	return os.Remove(filename)
}

// AppendOrCreateFile - Opens an existing file to be appended or creates a new one if
// it does not already exist
func AppendOrCreateFile(filename string) (*os.File, error) {
	if exists, err := FileExists(filename); err != nil {
		return nil, err
	} else if exists {
		return os.OpenFile(filename, os.O_WRONLY|os.O_APPEND, 0666)
	}
	if err := CreateDirectoryForFile(filename); err != nil {
		return nil, err
	}
	return os.Create(filename)
}

// CreateDirectoryForFile - Creates the directory for the given file
func CreateDirectoryForFile(filename string) error {
	err := os.MkdirAll(filepath.Dir(filename), os.ModeAppend)
	return err
}

// FileExists - Returns if the file exists
func FileExists(filename string) (bool, error) {
	if _, err := os.Stat(filename); err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		fmt.Printf("%s !exists\n", filename)
		return false, nil
	} else {
		return false, err
	}
}

// SplitFileName - Returns the directory, filename and extension in 3 seperate strings
func SplitFileName(filename string) (string, string, string) {
	var dir = filepath.Dir(filename)
	var longName = filepath.Base(filename)
	var extension = filepath.Ext(filename)
	var shortName = strings.TrimSuffix(longName, extension)
	return dir, shortName, extension[1:]
}
