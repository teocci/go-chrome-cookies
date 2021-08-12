// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package browser

import (
	"github.com/teocci/go-chrome-cookies/core/data"
	"github.com/teocci/go-chrome-cookies/logger"
	"strings"
)

type Firefox struct {
	name        string
	profilePath string
	keyPath     string
}

// NewFirefox return firefox browser interface
func NewFirefox(profile, key, name, storage string) (Browser, error) {
	return &Firefox{profilePath: profile, keyPath: key, name: name}, nil
}

// GetAllItems return all item with firefox
func (f *Firefox) GetAllItems() ([]data.Item, error) {
	var items []data.Item
	for item, choice := range firefoxItems {
		var (
			sub, main string
			err       error
		)
		if choice.subFile != "" {
			sub, err = getItemPath(f.profilePath, choice.subFile)
			if err != nil {
				logger.Debugf("%s find %s file failed, ERR:%s", f.name, item, err)
				continue
			}
		}
		main, err = getItemPath(f.profilePath, choice.mainFile)
		if err != nil {
			logger.Debugf("%s find %s file failed, ERR:%s", f.name, item, err)
			continue
		}
		i := choice.newItem(main, sub)
		logger.Debugf("%s find %s file success", f.name, item)
		items = append(items, i)
	}
	return items, nil
}

func (f *Firefox) GetItem(itemName string) (data.Item, error) {
	itemName = strings.ToLower(itemName)
	if item, ok := firefoxItems[itemName]; ok {
		var (
			sub, main string
			err       error
		)
		if item.subFile != "" {
			sub, err = getItemPath(f.profilePath, item.subFile)
			if err != nil {
				logger.Debugf("%s find %s file failed, ERR:%s", f.name, item.subFile, err)
			}
		}
		main, err = getItemPath(f.profilePath, item.mainFile)
		if err != nil {
			logger.Debugf("%s find %s file failed, ERR:%s", f.name, item.mainFile, err)
		}
		i := item.newItem(main, sub)
		logger.Debugf("%s find %s file success", f.name, item.mainFile)
		return i, nil
	} else {
		return nil, errItemNotSupported
	}
}

func (f *Firefox) GetName() string {
	return f.name
}

// GetSecretKey for firefox is always nil
// this method used to implement Browser interface
func (f *Firefox) GetSecretKey() []byte {
	return nil
}

// InitSecretKey for firefox is always nil
// this method used to implement Browser interface
func (f *Firefox) InitSecretKey() error {
	return nil
}
