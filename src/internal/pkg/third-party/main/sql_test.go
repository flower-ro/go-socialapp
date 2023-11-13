package main

import (
	"crypto/sha256"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"testing"
)

func Test_Md5(t *testing.T) {

	storeContainer, err := sqlstore.New("sqlite3",
		fmt.Sprintf("file:%s?_foreign_keys=off", "E:\\software\\sqlite3\\db\\8613027979536_back.db"), nil)
	//storeContainer, err := sqlstore.New(*dbDialect, *dbAddress, dbLog)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return
	}

	d, err := storeContainer.GetFirstDevice()
	//aa := sha256.New()
	ss := sha256.Sum256(d.NoiseKey.Priv[:])
	result := fmt.Sprintf("%s", ss)
	spew.Dump(result)
	//md5.Sum()
	//spew.Dump("---------------=", (d.NoiseKey.Priv).(*[]byte))
	//
	//spew.Dump(storeContainer.GetFirstDevice())

	//now := time.Now()
	//fmt.Println(now.Unix())               // 1565084298 秒
	//fmt.Println(now.UnixNano())           // 1565084298178502600 纳秒
	//fmt.Sprintf("%d", now.UnixNano()/1e6) // 1565084298178 毫秒
	//
	//timess := fmt.Sprintf("%d", now.Unix())
	//fmt.Println(timess)
	//str := "8ar4tc9v1" + timess
	//data := []byte(str)
	//has := md5.Sum(data)
	//md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制
	//
	//fmt.Println(md5str1)
}
