// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package browser

import (
	"database/sql"
	"fmt"
	"github.com/teocci/go-chrome-cookies/core/data"
	"github.com/teocci/go-chrome-cookies/logger"
	"reflect"
	"testing"
)

func TestListBrowser(t *testing.T) {
	browsers := ListBrowser()
	for index, value := range browsers {
		fmt.Printf("browsers[%d] : %s\n", index, reflect.ValueOf(value).String())
	}
}

func TestPickBrowser(t *testing.T) {
	name := "chrome"
	browsers, err := PickBrowser(name)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	var chrome Browser
	chrome = browsers[0]
	err = chrome.InitSecretKey()
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Printf("browser: %s\n", chrome.GetName())

	item, err := chrome.GetItem("cookie")
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Printf("item: %s\n", reflect.ValueOf(item).String())
	key := chrome.GetSecretKey()
	str := string(key)
	fmt.Println("key:", str)

	fmt.Printf("key: %s\n", reflect.ValueOf(key).String())
	err = item.CopyDB()
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	err = item.ChromeParse(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	err = item.OutPut(data.GetFormat(data.FormatNameConsole), "D:/Temp/", "chrome")
	if err != nil {
		return
	}

	//fmt.Printf("item: %s\n", item)
}

func TestNew(t *testing.T) {
	//tests := []struct {
	//	err  string
	//	want error
	//}{
	//	{"", fmt.Errorf("")},
	//	{"foo", fmt.Errorf("foo")},
	//	{"foo", New("foo")},
	//	{"string with format specifiers: %v", throw.New("string with format specifiers: %v")},
	//}
	//
	//for _, tt := range tests {
	//	got := New(tt.err)
	//	if got.Error() != tt.want.Error() {
	//		t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
	//	}
	//}
}

func TestDB(t *testing.T) {
	cookieDB, err := sql.Open("sqlite3", data.ChromeCookieFile)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	defer func() {
		if err := cookieDB.Close(); err != nil {
			logger.Debug(err)
		}
	}()

	rows, err := cookieDB.Query("SELECT name FROM sqlite_master WHERE type='table';")
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			logger.Debug(err)
		}
	}()

	fmt.Println("rows:" )
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			fmt.Printf("err: %s\n", err)
		}
		fmt.Println("name:", name )
	}
}
