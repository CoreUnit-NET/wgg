package wgg

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// CleanUpOutDir removes all files in outDir that have a name that starts with "node." or "client.".
func CleanUpOutDir(outDir string) error {
	files, err := os.ReadDir(outDir)
	if err != nil {
		return errors.New("Error reading outDir: " + err.Error())
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), "node.") ||
			strings.HasPrefix(file.Name(), "client.") {
			err = os.Remove(outDir + "/" + file.Name())
			if err != nil {
				return errors.New("Error removing '" + outDir + "/" + file.Name() + "': " + err.Error())
			}
		}
	}

	return nil
}

// InitOutDir initializes the output directory for configuration files.
// It retrieves the directory path from the WGG_OUT_DIR environment variable.
// If the path is not absolute, it prefixes it with the current working directory.
// If the directory does not exist, it attempts to create it with the appropriate permissions.
// Returns the absolute path of the output directory or an error if any operation fails.
func InitOutDir() (string, string, error) {
	outDir := os.Getenv("WGG_OUT_DIR")
	if len(outDir) <= 0 {
		return "", "", errors.New("the WGG_OUT_DIR env var is not set or empty")
	} else if !strings.HasPrefix(outDir, "/") {
		outDir = FatalCwd() + "/" + outDir
	}

	keyDir := outDir + "/keys"

	_, err := os.Stat(keyDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(keyDir, 0755)
		if err != nil {
			return "", "", errors.New("Error creating outDir at '" + keyDir + "': " + err.Error())
		}
	}

	return outDir, keyDir, nil
}

func Cwd() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current working directory: %w", err)
	}

	return cwd, nil
}

func FatalCwd() string {
	cwd, err := Cwd()
	if err != nil {
		log.Fatalln(err)
	}
	return cwd
}
