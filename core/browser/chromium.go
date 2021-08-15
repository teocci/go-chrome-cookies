// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package browser

import (
	"github.com/teocci/go-chrome-cookies/core/data"
	"github.com/teocci/go-chrome-cookies/core/throw"
	"github.com/teocci/go-chrome-cookies/logger"
	"strings"
)

var chromiumItems = map[string]struct {
	mainFile string
	newItem  func(mainFile, subFile string) data.Item
}{
	data.ItemNameBookmark: {
		mainFile: data.ChromeBookmarkFile,
		newItem:  data.NewBookmarks,
	},
	data.ItemNameCookie: {
		mainFile: data.ChromeCookieFile,
		newItem:  data.NewCookies,
	},
	data.ItemNameHistory: {
		mainFile: data.ChromeHistoryFile,
		newItem:  data.NewHistoryData,
	},
	data.ItemNameDownload: {
		mainFile: data.ChromeDownloadFile,
		newItem:  data.NewDownloads,
	},
	data.ItemNamePassword: {
		mainFile: data.ChromePasswordFile,
		newItem:  data.NewCPasswords,
	},
	data.ItemNameCreditCard: {
		mainFile: data.ChromeCreditFile,
		newItem:  data.NewCCards,
	},
}

type Chromium struct {
	name        string
	profilePath string
	keyPath     string
	storage     string // storage use for linux and macOS, get secret key
	secretKey   []byte
}

// NewChromium return Chromium browser interface
func NewChromium(profile, key, name, storage string) (Browser, error) {
	return &Chromium{profilePath: profile, keyPath: key, name: name, storage: storage}, nil
}

func (c *Chromium) GetName() string {
	return c.name
}

func (c *Chromium) GetKeyPath() string {
	return c.keyPath
}

func (c *Chromium) GetStorage() string {
	return c.storage
}

func (c *Chromium) GetSecretKey() []byte {
	return c.secretKey
}

func (c *Chromium) SetSecretKey(secretKey []byte) {
	c.secretKey = secretKey
}

// GetAllItems return all chromium items from browser
// If it can't find the item path, log error then continue
func (c *Chromium) GetAllItems() ([]data.Item, error) {
	var items []data.Item
	for item, choice := range chromiumItems {
		m, err := GetItemPath(c.profilePath, choice.mainFile)
		if err != nil {
			logger.Debugf("%s find %s file failed, ERR:%s", c.name, item, err)
			continue
		}
		i := choice.newItem(m, "")
		logger.Debugf("%s find %s File Success", c.name, item)
		items = append(items, i)
	}
	return items, nil
}

// GetItem return single item
func (c *Chromium) GetItem(itemName string) (data.Item, error) {
	itemName = strings.ToLower(itemName)
	if item, ok := chromiumItems[itemName]; ok {
		m, err := GetItemPath(c.profilePath, item.mainFile)
		if err != nil {
			logger.Debugf("%s find %s file failed, ERR:%s", c.name, item.mainFile, err)
		}
		i := item.newItem(m, "")
		return i, nil
	} else {
		return nil, throw.ErrorItemNotSupported()
	}
}

func (c *Chromium) InitSecretKey() error {
	err := InitSecretKey(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Chromium) ListItems() []string {
	var l []string
	for k := range chromiumItems {
		l = append(l, k)
	}
	return l
}
