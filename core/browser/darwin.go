// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
// +build !windows,!plan9,!nacl,!linux

package browser

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"os/exec"

	"golang.org/x/crypto/pbkdf2"
)

const (
	fireFoxProfilePath        = "/Users/*/Library/Application Support/Firefox/Profiles/*.default-release*/"
	fireFoxBetaProfilePath    = "/Users/*/Library/Application Support/Firefox/Profiles/*.default-beta*/"
	fireFoxDevProfilePath     = "/Users/*/Library/Application Support/Firefox/Profiles/*.dev-edition-default*/"
	fireFoxNightlyProfilePath = "/Users/*/Library/Application Support/Firefox/Profiles/*.default-nightly*/"
	fireFoxESRProfilePath     = "/Users/*/Library/Application Support/Firefox/Profiles/*.default-esr*/"
	chromeProfilePath         = "/Users/*/Library/Application Support/Google/Chrome/*/"
	chromeBetaProfilePath     = "/Users/*/Library/Application Support/Google/Chrome Beta/*/"
	chromiumProfilePath       = "/Users/*/Library/Application Support/chromium/*/"
	edgeProfilePath           = "/Users/*/Library/Application Support/Microsoft Edge/*/"
	braveProfilePath          = "/Users/*/Library/Application Support/BraveSoftware/Brave-Browser/*/"
	operaProfilePath          = "/Users/*/Library/Application Support/com.operasoftware.Opera/"
	operaGXProfilePath        = "/Users/*/Library/Application Support/com.operasoftware.OperaGX/"
	vivaldiProfilePath        = "/Users/*/Library/Application Support/Vivaldi/*/"
)

const (
	chromeStorageName     = "Chrome"
	chromeBetaStorageName = "Chrome"
	chromiumStorageName   = "chromium"
	edgeStorageName       = "Microsoft Edge"
	braveStorageName      = "Brave"
	operaStorageName      = "Opera"
	vivaldiStorageName    = "Vivaldi"
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
		"chromium": {
			ProfilePath: chromiumProfilePath,
			Name:        chromiumName,
			Storage:     chromiumStorageName,
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
		"opera": {
			ProfilePath: operaProfilePath,
			Name:        operaName,
			Storage:     operaStorageName,
			New:         NewChromium,
		},
		"opera-gx": {
			ProfilePath: operaGXProfilePath,
			Name:        operaGXName,
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
	var (
		cmd            *exec.Cmd
		stdout, stderr bytes.Buffer
	)
	// ➜ security find-generic-password -wa 'Chrome'
	cmd = exec.Command("security", "find-generic-password", "-wa", c.GetStorage())
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	if stderr.Len() > 0 {
		err = errors.New(stderr.String())
		return err
	}
	temp := stdout.Bytes()
	chromeSecret := temp[:len(temp)-1]
	if chromeSecret == nil {
		return ThrowErrorChromeSecretIsEmpty()
	}
	var chromeSalt = []byte("saltysalt")
	// @https://source.chromium.org/chromium/chromium/src/+/master:components/os_crypt/os_crypt_mac.mm;l=157
	key := pbkdf2.Key(chromeSecret, chromeSalt, 1003, 16, sha1.New)
	c.SetSecretKey(key)
	return nil
}
