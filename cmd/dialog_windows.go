// +build windows

package cmd

import "github.com/sqweek/dialog"

func selectFolderWindows(title string) (string, error) {
	return dialog.Directory().Title(title).Browse()
}
