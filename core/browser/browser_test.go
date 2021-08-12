// Package browser
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-12
package browser

import (
	"fmt"
	"reflect"
	"testing"
)

func TestListBrowser(t *testing.T) {
	browsers := ListBrowser()
	for index, value := range browsers {
		fmt.Printf("browsers[%d] : %s\n", index, reflect.ValueOf(value).String())
	}
}

func TestListItem(t *testing.T) {
	items := ListItem()
	for index, item := range items {
		fmt.Printf("items[%d] : %s\n", index, reflect.ValueOf(item).String())
	}
}

func TestPickBrowser(t *testing.T) {
	name := "chrome"
	browsers, err := PickBrowser(name)
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}
	chrome := browsers[0]

	fmt.Printf("browser: %s\n", chrome.GetName())
	item, err := chrome.GetItem("cookie")
	if err != nil {
		fmt.Printf("err: %s\n", err)
	}

	fmt.Printf("item: %s\n", reflect.ValueOf(item).String())
	key := chrome.GetSecretKey()
	fmt.Printf("key: %s\n", reflect.ValueOf(key).String())

	err = item.ChromeParse(key)
	if err != nil {
		fmt.Printf("err: %s\n", err)
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
	//	{"string with format specifiers: %v", errors.New("string with format specifiers: %v")},
	//}
	//
	//for _, tt := range tests {
	//	got := New(tt.err)
	//	if got.Error() != tt.want.Error() {
	//		t.Errorf("New.Error(): got: %q, want %q", got, tt.want)
	//	}
	//}
}