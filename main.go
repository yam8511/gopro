package main

import (
	"fmt"

	"github.com/ipipdotnet/ipdb-go"
)

func main() {

	var s string
	fmt.Scanln(&s)
}

var dbIPv4 *ipdb.City
var dbIPv6 *ipdb.City

// doInit 初始化資料庫
func doInit() error {
	path := "./ipipdotnet/"
	pathOfIPv4db := path + "mydata4vipday2_20200409.ipdb"
	pathOfIPv6db := path + "mydata6vipday2_20200409.ipdb"

	// 初始化資料庫
	f := func(ipv4 bool, path string) (*ipdb.City, error) {
		db, err := ipdb.NewCity(path)
		if err != nil {
			return nil, fmt.Errorf("ipipdonet 資料庫初始化失敗(path:%s): %w", path, err)
		}

		if len(db.Languages()) < 1 {
			return nil, fmt.Errorf("ipipdonet 資料庫初始化失敗(path:%s): db.Languages() < 1", path)
		}

		if ipv4 && !db.IsIPv4() {
			return nil, fmt.Errorf("ipipdonet 資料庫(ipv4)初始化失敗(path:%s)", path)
		} else if !ipv4 && !db.IsIPv6() {
			return nil, fmt.Errorf("ipipdonet 資料庫(ipv6)初始化失敗(path:%s)", path)
		}

		return db, nil
	}
	var err error
	dbIPv4, err = f(true, pathOfIPv4db) // 初始化IPv4資料庫
	if err != nil {
		return err
	}
	dbIPv6, err = f(false, pathOfIPv6db) // 初始化IPv6資料庫
	if err != nil {
		return err
	}

	return nil
}
