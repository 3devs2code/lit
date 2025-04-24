package initialise

import (
	"fmt"
	"os"
)

func InitRepo() error {
	if _, err := os.Stat(".lit"); !os.IsNotExist(err) {
		return fmt.Errorf("directory already exists")
	}

	err := os.MkdirAll(".lit/objects", 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}

	if err := os.MkdirAll(".lit/refs/heads", 0755); err != nil {
		return fmt.Errorf("failed to create refs/heads directory: %v", err)
	}

	headContent := []byte("ref: refs/heads/master\n")
	if err := os.WriteFile(".lit/HEAD", headContent, 0644); err != nil {
		return fmt.Errorf("failed to create HEAD file: %v", err)
	}

	fmt.Println("Initialized empty Lit repository in .lit/")
	return nil
}
