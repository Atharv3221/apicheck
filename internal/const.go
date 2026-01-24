package internal

import (
	"os"
	"path/filepath"
)

const file = ".apicheck.yaml"

// Returns global Path for apicheck config based on os.
func GetGlobalPath() string {
	homedir, _ := os.UserHomeDir()
	globalPath := filepath.Join(homedir, file)
	return globalPath
}

func GetLocalPath() string {
	cureentDir, _ := os.Getwd()
	localPath := filepath.Join(cureentDir, file)
	return localPath
}
