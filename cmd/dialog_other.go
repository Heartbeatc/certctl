// +build !windows

package cmd

import "fmt"

func selectFolderWindows(title string) (string, error) {
	return "", fmt.Errorf("not supported on this platform")
}
