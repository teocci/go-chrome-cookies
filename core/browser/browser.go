// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package browser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/teocci/go-chrome-cookies/core/data"
	"github.com/teocci/go-chrome-cookies/logger"
)

const (
	chromeName         = "Chrome"
	chromeBetaName     = "Chrome Beta"
	chromiumName       = "Chromium"
	edgeName           = "Microsoft Edge"
	firefoxName        = "Firefox"
	firefoxBetaName    = "Firefox Beta"
	firefoxDevName     = "Firefox Dev"
	firefoxNightlyName = "Firefox Nightly"
	firefoxESRName     = "Firefox ESR"
	speed360Name       = "360speed"
	qqBrowserName      = "qq"
	braveName          = "Brave"
	operaName          = "Opera"
	operaGXName        = "OperaGX"
	vivaldiName        = "Vivaldi"
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
}

const (
	cookie     = "cookie"
	history    = "history"
	bookmark   = "bookmark"
	download   = "download"
	password   = "password"
	creditCard = "credit-card"
)

var (
	errItemNotSupported    = errors.New(`item not supported, default is "all", choose from history|download|password|bookmark|cookie`)
	errBrowserNotSupported = errors.New("browser not supported")
	errChromeSecretIsEmpty = errors.New("chrome secret is empty")
	errDbusSecretIsEmpty   = errors.New("dbus secret key is empty")
)

var (
	chromiumItems = map[string]struct {
		mainFile string
		newItem  func(mainFile, subFile string) data.Item
	}{
		bookmark: {
			mainFile: data.ChromeBookmarkFile,
			newItem:  data.NewBookmarks,
		},
		cookie: {
			mainFile: data.ChromeCookieFile,
			newItem:  data.NewCookies,
		},
		history: {
			mainFile: data.ChromeHistoryFile,
			newItem:  data.NewHistoryData,
		},
		download: {
			mainFile: data.ChromeDownloadFile,
			newItem:  data.NewDownloads,
		},
		password: {
			mainFile: data.ChromePasswordFile,
			newItem:  data.NewCPasswords,
		},
		creditCard: {
			mainFile: data.ChromeCreditFile,
			newItem:  data.NewCCards,
		},
	}
	firefoxItems = map[string]struct {
		mainFile string
		subFile  string
		newItem  func(mainFile, subFile string) data.Item
	}{
		bookmark: {
			mainFile: data.FirefoxDataFile,
			newItem:  data.NewBookmarks,
		},
		cookie: {
			mainFile: data.FirefoxCookieFile,
			newItem:  data.NewCookies,
		},
		history: {
			mainFile: data.FirefoxDataFile,
			newItem:  data.NewHistoryData,
		},
		download: {
			mainFile: data.FirefoxDataFile,
			newItem:  data.NewDownloads,
		},
		password: {
			mainFile: data.FirefoxKey4File,
			subFile:  data.FirefoxLoginFile,
			newItem:  data.NewFPasswords,
		},
	}
)

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
	return nil, errBrowserNotSupported
}

// PickCustomBrowser pick single browser with custom browser profile path and key file path (windows only).
// If custom key file path is empty, but the current browser requires key file (chromium for windows version > 80)
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

func getItemPath(profilePath, file string) (string, error) {
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

// check key file path is exist
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

func ListBrowser() []string {
	var l []string
	for k := range browserList {
		l = append(l, k)
	}
	return l
}

func ListItem() []string {
	var l []string
	for k := range chromiumItems {
		l = append(l, k)
	}
	return l
}