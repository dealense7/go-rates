package main

import (
	"fmt"
	"github.com/dealense7/go-rate-app/internal/utils"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func dumpDatabase(location, filename string) (string, error) {
	config := utils.LoadDBEnv()

	// Ensure the location directory exists
	if err := os.MkdirAll(location, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %v", location, err)
	}

	// Construct the full file path
	filePath := filepath.Join(location, filename)

	// Prepare mysqldump command
	cmd := exec.Command(
		"mysqldump",
		"-u", config.Username,
		fmt.Sprintf("-p%s", config.Password), // Password passed with -p
		"-h", config.Host,
		"-P", config.Port,
		config.Database,
	)

	// Create output file
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	// Set the command's output to the file
	cmd.Stdout = file

	// Run the mysqldump command
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("failed to execute mysqldump: %v", err)
	}

	return filePath, nil
}

func main() {
	path, _ := os.Getwd()
	dumpName := fmt.Sprintf("dump-%s.sql", time.Now().Format("2006-01-02"))
	str, err := dumpDatabase(path+"/static/dump", dumpName)
	if err != nil {
		fmt.Println("Dump failed:", err)
	} else {
		fmt.Println(str)
	}
}
