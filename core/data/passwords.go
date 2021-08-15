// Package data
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-14
package data

import (
	"bytes"
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/teocci/go-chrome-cookies/core/decrypt"
	"github.com/teocci/go-chrome-cookies/filemgmt"
	"github.com/teocci/go-chrome-cookies/logger"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type loginData struct {
	UserName    string
	encryptPass []byte
	encryptUser []byte
	Password    string
	LoginUrl    string
	CreateDate  time.Time
}

type passwords struct {
	mainPath string
	subPath  string
	logins   []loginData
}

func NewFPasswords(main, sub string) Item {
	return &passwords{mainPath: main, subPath: sub}
}

func NewCPasswords(main, sub string) Item {
	return &passwords{mainPath: main}
}

func (p *passwords) ChromeParse(key []byte) error {
	loginDB, err := sql.Open("sqlite3", ChromePasswordFile)
	if err != nil {
		return err
	}
	defer func() {
		if err := loginDB.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	rows, err := loginDB.Query(QueryChromiumLogin)
	if err != nil {
		return err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	for rows.Next() {
		var (
			url, username string
			pwd, password []byte
			create        int64
		)
		err = rows.Scan(&url, &username, &pwd, &create)
		if err != nil {
			logger.Error(err)
		}
		login := loginData{
			UserName:    username,
			encryptPass: pwd,
			LoginUrl:    url,
		}
		if key == nil {
			password, err = decrypt.DPApi(pwd)
		} else {
			password, err = decrypt.ChromePass(key, pwd)
		}
		if err != nil {
			logger.Debugf("%s have empty password %s", login.LoginUrl, err.Error())
		}
		if create > time.Now().Unix() {
			login.CreateDate = filemgmt.TimeEpochFormat(create)
		} else {
			login.CreateDate = filemgmt.TimeStampFormat(create)
		}
		login.Password = string(password)
		p.logins = append(p.logins, login)
	}
	return nil
}

func (p *passwords) FirefoxParse() error {
	globalSalt, metaBytes, nssA11, nssA102, err := getFirefoxDecryptKey()
	if err != nil {
		return err
	}
	keyLin := []byte{248, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	metaPBE, err := decrypt.NewASN1PBE(metaBytes)
	if err != nil {
		logger.Error("decrypt meta data failed", err)
		return err
	}
	// default master password is empty
	var masterPwd []byte
	k, err := metaPBE.Decrypt(globalSalt, masterPwd)
	if err != nil {
		logger.Error("decrypt firefox meta bytes failed", err)
		return err
	}
	if bytes.Contains(k, []byte("password-check")) {
		logger.Debug("password-check success")
		m := bytes.Compare(nssA102, keyLin)
		if m == 0 {
			nssPBE, err := decrypt.NewASN1PBE(nssA11)
			if err != nil {
				logger.Error("decode firefox nssA11 bytes failed", err)
				return err
			}
			finallyKey, err := nssPBE.Decrypt(globalSalt, masterPwd)
			finallyKey = finallyKey[:24]
			if err != nil {
				logger.Error("get firefox finally key failed")
				return err
			}
			allLogins, err := getFirefoxLoginData()
			if err != nil {
				return err
			}
			for _, v := range allLogins {
				userPBE, err := decrypt.NewASN1PBE(v.encryptUser)
				if err != nil {
					logger.Error("decode firefox user bytes failed", err)
				}
				pwdPBE, err := decrypt.NewASN1PBE(v.encryptPass)
				if err != nil {
					logger.Error("decode firefox password bytes failed", err)
				}
				user, err := userPBE.Decrypt(finallyKey, masterPwd)
				if err != nil {
					logger.Error(err)
				}
				pwd, err := pwdPBE.Decrypt(finallyKey, masterPwd)
				if err != nil {
					logger.Error(err)
				}
				logger.Debug("decrypt firefox success")
				p.logins = append(p.logins, loginData{
					LoginUrl:   v.LoginUrl,
					UserName:   string(decrypt.PKCS5UnPadding(user)),
					Password:   string(decrypt.PKCS5UnPadding(pwd)),
					CreateDate: v.CreateDate,
				})
			}
		}
	}
	return nil
}

func (p *passwords) CopyDB() error {
	err := CopyToLocalPath(p.mainPath, filepath.Base(p.mainPath))
	if err != nil {
		logger.Error(err)
	}
	if p.subPath != "" {
		err = CopyToLocalPath(p.subPath, filepath.Base(p.subPath))
	}
	return err
}

func (p *passwords) Release() error {
	err := os.Remove(filepath.Base(p.mainPath))
	if err != nil {
		logger.Error(err)
	}
	if p.subPath != "" {
		err = os.Remove(filepath.Base(p.subPath))
	}
	return err
}

func (p *passwords) OutPut(format OutputFormat, browser, dir string) error {
	sort.Sort(p)
	switch format {
	case formatCSV:
		err := p.outPutCsv(browser, dir)
		return err
	case formatConsole:
		p.outPutConsole()
		return nil
	default:
		err := p.outPutJson(browser, dir)
		return err
	}
}

func (p *passwords) outPutCsv(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNamePassword, GetFormatName(formatJson))
	if err := WriteToCsv(filename, p.logins); err != nil {
		return err
	}
	fmt.Printf("%s Get %d passwords, filename is %s \n", filemgmt.Prefix, len(p.logins), filename)
	return nil
}

func (p *passwords) outPutJson(browser, dir string) error {
	filename := filemgmt.FormatFileName(dir, browser, ItemNamePassword, GetFormatName(formatCSV))
	err := WriteToJson(filename, p.logins)
	if err != nil {
		return err
	}
	fmt.Printf("%s Get %d passwords, filename is %s \n", filemgmt.Prefix, len(p.logins), filename)
	return nil
}

func (p *passwords) outPutConsole() {
	for _, v := range p.logins {
		fmt.Printf("%+v\n", v)
	}
}

func (p passwords) Len() int {
	return len(p.logins)
}

func (p passwords) Less(i, j int) bool {
	return p.logins[i].CreateDate.After(p.logins[j].CreateDate)
}

func (p passwords) Swap(i, j int) {
	p.logins[i], p.logins[j] = p.logins[j], p.logins[i]
}

// getFirefoxDecryptKey get value from key4.db
func getFirefoxDecryptKey() (item1, item2, a11, a102 []byte, err error) {
	var (
		keyDB   *sql.DB
		pwdRows *sql.Rows
		nssRows *sql.Rows
	)
	keyDB, err = sql.Open("sqlite3", FirefoxKey4File)
	if err != nil {
		logger.Error(err)
		return nil, nil, nil, nil, err
	}
	defer func() {
		if err := keyDB.Close(); err != nil {
			logger.Error(err)
		}
	}()

	pwdRows, err = keyDB.Query(QueryMetaData)
	defer func() {
		if err := pwdRows.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	for pwdRows.Next() {
		if err := pwdRows.Scan(&item1, &item2); err != nil {
			logger.Error(err)
			continue
		}
	}
	if err != nil {
		logger.Error(err)
	}
	nssRows, err = keyDB.Query(QueryNssPrivate)
	defer func() {
		if err := nssRows.Close(); err != nil {
			logger.Debug(err)
		}
	}()
	for nssRows.Next() {
		if err := nssRows.Scan(&a11, &a102); err != nil {
			logger.Debug(err)
		}
	}
	return item1, item2, a11, a102, nil
}

// getFirefoxLoginData used to get firefox
func getFirefoxLoginData() (l []loginData, err error) {
	s, err := ioutil.ReadFile(FirefoxLoginFile)
	if err != nil {
		return nil, err
	}
	h := gjson.GetBytes(s, "logins")
	if h.Exists() {
		for _, v := range h.Array() {
			var (
				m loginData
				u []byte
				p []byte
			)
			m.LoginUrl = v.Get("formSubmitURL").String()
			u, err = base64.StdEncoding.DecodeString(v.Get("encryptedUsername").String())
			m.encryptUser = u
			if err != nil {
				logger.Debug(err)
			}
			p, err = base64.StdEncoding.DecodeString(v.Get("encryptedPassword").String())
			m.encryptPass = p
			m.CreateDate = filemgmt.TimeStampFormat(v.Get("timeCreated").Int() / 1000)
			l = append(l, m)
		}
	}
	return
}
