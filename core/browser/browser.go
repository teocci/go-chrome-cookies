// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package browser

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/teocci/go-chrome-cookies/core/data"
	"github.com/teocci/go-chrome-cookies/core/throw"
	"github.com/teocci/go-chrome-cookies/logger"
)

type Browser interface {
	// InitSecretKey is init chrome secret key, firefox's key always empty
	InitSecretKey() error

	// GetName return browser name
	GetName() string

	// GetSecretKey return browser secret key
	GetSecretKey() []byte

	// GetAllItems return all items (password|bookmark|cookie|history)
	GetAllItems() ([]data.Item, error)

	// GetItem return single one from the password|bookmark|cookie|history
	GetItem(itemName string) (data.Item, error)

	// ListItems return list of items
	ListItems() []string
}

// PickBrowser return a list of browser interface
func PickBrowser(name string) ([]Browser, error) {
	var browsers []Browser
	name = strings.ToLower(name)
	if name == "all" {
		for _, v := range browserList {
			b, err := v.New(v.ProfilePath, v.KeyPath, v.Name, v.Storage)
			if err != nil {
				logger.Error(err)
			}
			browsers = append(browsers, b)
		}
		return browsers, nil
	} else if choice, ok := browserList[name]; ok {
		b, err := choice.New(choice.ProfilePath, choice.KeyPath, choice.Name, choice.Storage)
		browsers = append(browsers, b)
		return browsers, err
	}
	return nil, throw.ErrorBrowserNotSupported()
}

// PickCustomBrowser pick single browser with custom browser profile path and key file path (Windows only).
// If custom key file path is empty, but the current browser requires key file (chromium for Windows version > 80)
// key file path will be automatically found in the profile path's parent directory.
func PickCustomBrowser(browserName, cusProfile, cusKey string) ([]Browser, error) {
	var (
		browsers []Browser
	)
	browserName = strings.ToLower(browserName)
	supportBrowser := strings.Join(ListBrowser(), "|")
	if browserName == "all" {
		return nil, fmt.Errorf("can't select all browser, pick one from %s with -b flag\n", supportBrowser)
	}
	if choice, ok := browserList[browserName]; ok {
		// if this browser need key path
		if choice.KeyPath != "" {
			var err error
			// if browser need key path and cusKey is empty, try to get key path with profile dir
			if cusKey == "" {
				cusKey, err = getKeyPath(cusProfile)
				if err != nil {
					return nil, err
				}
			}
			if err := checkKeyPath(cusKey); err != nil {
				return nil, err
			}
		}
		b, err := choice.New(cusProfile, cusKey, choice.Name, choice.Storage)
		browsers = append(browsers, b)
		return browsers, err
	} else {
		return nil, fmt.Errorf("%s not support, pick one from %s with -b flag\n", browserName, supportBrowser)
	}
}

// GetItemPath try to get item file path with the browser's profile path
// default key file path is in the parent directory of the profile dir, and name is [Local State]
func GetItemPath(profilePath, file string) (string, error) {
	p, err := filepath.Glob(filepath.Join(profilePath, file))
	if err != nil {
		return "", err
	}
	if len(p) > 0 {
		return p[0], nil
	}
	return "", fmt.Errorf("find %s failed", file)
}

// getKeyPath try to get key file path with the browser's profile path
// default key file path is in the parent directory of the profile dir, and name is [Local State]
func getKeyPath(profilePath string) (string, error) {
	if _, err := os.Stat(filepath.Clean(profilePath)); os.IsNotExist(err) {
		return "", err
	}
	parentDir := getParentDirectory(profilePath)
	keyPath := filepath.Join(parentDir, "Local State")
	return keyPath, nil
}

// checkKeyPath check if key file path exist
func checkKeyPath(keyPath string) error {
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return fmt.Errorf("secret key path not exist, please check %s", keyPath)
	}
	return nil
}

func getParentDirectory(dir string) string {
	var (
		length int
	)
	// filepath.Clean(dir) auto remove
	dir = strings.ReplaceAll(filepath.Clean(dir), `\`, `/`)
	length = strings.LastIndex(dir, "/")
	if length > 0 {
		if length > len([]rune(dir)) {
			length = len([]rune(dir))
		}
		return string([]rune(dir)[:length])
	}
	return ""
}