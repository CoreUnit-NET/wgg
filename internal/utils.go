package wgg

import "os/exec"

// IsCommandAvailable returns true if the command is available in the system's PATH, false otherwise.
func IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}
