// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
// +build !windows,!plan9,!nacl,!darwin

package browser

import (
	"crypto/sha1"
	"github.com/godbus/dbus/v5"
	"github.com/teocci/go-chrome-cookies/logger"

	keyring "github.com/ppacher/go-dbus-keyring"
	"golang.org/x/crypto/pbkdf2"
)

const (
	fireFoxProfilePath        = "/home/*/.mozilla/firefox/*.default-release*/"
	fireFoxBetaProfilePath    = "/home/*/.mozilla/firefox/*.default-beta*/"
	fireFoxDevProfilePath     = "/home/*/.mozilla/firefox/*.dev-edition-default*/"
	fireFoxNightlyProfilePath = "/home/*/.mozilla/firefox/*.default-nightly*/"
	fireFoxESRProfilePath     = "/home/*/.mozilla/firefox/*.default-esr*/"
	chromeProfilePath         = "/home/*/.config/google-chrome/*/"
	chromiumProfilePath       = "/home/*/.config/chromium/*/"
	edgeProfilePath           = "/home/*/.config/microsoft-edge*/*/"
	braveProfilePath          = "/home/*/.config/BraveSoftware/Brave-Browser/*/"
	chromeBetaProfilePath     = "/home/*/.config/google-chrome-beta/*/"
	operaProfilePath          = "/home/*/.config/opera/"
	vivaldiProfilePath        = "/home/*/.config/vivaldi/*/"
)

const (
	chromeStorageName     = "Chrome Safe Storage"
	chromiumStorageName   = "chromium Safe Storage"
	edgeStorageName       = "chromium Safe Storage"
	braveStorageName      = "Brave Safe Storage"
	chromeBetaStorageName = "Chrome Safe Storage"
	operaStorageName      = "chromium Safe Storage"
	vivaldiStorageName    = "Chrome Safe Storage"
)

var (
	browserList = map[string]struct {
		ProfilePath string
		Name        string
		KeyPath     string
		Storage     string
		New         func(profile, key, name, storage string) (Browser, error)
	}{
		"firefox": {
			ProfilePath: fireFoxProfilePath,
			Name:        firefoxName,
			New:         NewFirefox,
		},
		"firefox-beta": {
			ProfilePath: fireFoxBetaProfilePath,
			Name:        firefoxBetaName,
			New:         NewFirefox,
		},
		"firefox-dev": {
			ProfilePath: fireFoxDevProfilePath,
			Name:        firefoxDevName,
			New:         NewFirefox,
		},
		"firefox-nightly": {
			ProfilePath: fireFoxNightlyProfilePath,
			Name:        firefoxNightlyName,
			New:         NewFirefox,
		},
		"firefox-esr": {
			ProfilePath: fireFoxESRProfilePath,
			Name:        firefoxESRName,
			New:         NewFirefox,
		},
		"chrome": {
			ProfilePath: chromeProfilePath,
			Name:        chromeName,
			Storage:     chromeStorageName,
			New:         NewChromium,
		},
		"edge": {
			ProfilePath: edgeProfilePath,
			Name:        edgeName,
			Storage:     edgeStorageName,
			New:         NewChromium,
		},
		"brave": {
			ProfilePath: braveProfilePath,
			Name:        braveName,
			Storage:     braveStorageName,
			New:         NewChromium,
		},
		"chrome-beta": {
			ProfilePath: chromeBetaProfilePath,
			Name:        chromeBetaName,
			Storage:     chromeBetaStorageName,
			New:         NewChromium,
		},
		"chromium": {
			ProfilePath: chromiumProfilePath,
			Name:        chromiumName,
			Storage:     chromiumStorageName,
			New:         NewChromium,
		},
		"opera": {
			ProfilePath: operaProfilePath,
			Name:        operaName,
			Storage:     operaStorageName,
			New:         NewChromium,
		},
		"vivaldi": {
			ProfilePath: vivaldiProfilePath,
			Name:        vivaldiName,
			Storage:     vivaldiStorageName,
			New:         NewChromium,
		},
	}
)

func InitSecretKey(c *Chromium) error {
	// what is d-bus @https://dbus.freedesktop.org/
	var chromeSecret []byte
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}
	svc, err := keyring.GetSecretService(conn)
	if err != nil {
		return err
	}
	session, err := svc.OpenSession()
	if err != nil {
		return err
	}
	defer func() {
		session.Close()
	}()
	collections, err := svc.GetAllCollections()
	if err != nil {
		return err
	}
	for _, col := range collections {
		items, err := col.GetAllItems()
		if err != nil {
			return err
		}
		for _, item := range items {
			label, err := item.GetLabel()
			if err != nil {
				logger.Error(err)
				continue
			}
			if label == c.GetStorage() {
				se, err := item.GetSecret(session.Path())
				if err != nil {
					logger.Error(err)
					return err
				}
				chromeSecret = se.Value
			}
		}
	}
	if chromeSecret == nil {
		return ThrowErrorDbusSecretIsEmpty()
	}
	var chromeSalt = []byte("saltysalt")
	// @https://source.chromium.org/chromium/chromium/src/+/master:components/os_crypt/os_crypt_linux.cc
	key := pbkdf2.Key(chromeSecret, chromeSalt, 1, 16, sha1.New)
	c.SetSecretKey(key)
	return nil
}
