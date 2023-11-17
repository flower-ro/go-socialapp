package main

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/util/keys"
	"io"
	"testing"
)

var keyGUID = []byte("258EAFA5-E914-47DA-95CA-C5AB0DC85B11")

func computeAcceptKey(challengeKey string) string {
	h := sha1.New()
	h.Write([]byte(challengeKey))
	h.Write(keyGUID)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func decodeAcceptKey() {
	de, err := base64.StdEncoding.DecodeString("U9wTSK1cN8+gqJ7Vgb9eLA4g6yrLYz74D/8ohMuwMlo=")
	spew.Dump(err)
	noisekey := keys.NewKeyPairFromPrivateKey(*(*[32]byte)(de))
	spew.Dump(noisekey)
}

func generateChallengeKey() (string, error) {
	p := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, p); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(p), nil
}

func Test_decode(t *testing.T) {
	decodeAcceptKey()

}

func Test_random(t *testing.T) {
	tmp, _ := generateChallengeKey()
	spew.Dump(tmp, computeAcceptKey(tmp))
}

func Test_Md5(t *testing.T) {

	storeContainer, err := sqlstore.New("sqlite3",
		fmt.Sprintf("file:%s?_foreign_keys=off", "E:\\software\\sqlite3\\db\\86xxxx_back.db"), nil)
	//storeContainer, err := sqlstore.New(*dbDialect, *dbAddress, dbLog)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		return
	}

	d, err := storeContainer.GetFirstDevice()

	bytes2 := sha256.Sum256(d.NoiseKey.Priv[:]) //计算哈希值，返回一个长度为32的数组
	hashCode2 := hex.EncodeToString(bytes2[:])  //将数组转换成切片，转换成16进制，返回字符串

	spew.Dump(hashCode2)

	//aa := sha256.New(
	//has := md5.Sum(d.NoiseKey.Priv[:])
	//spew.Dump(has)
	//md5str1 := fmt.Sprintf("%x", has) //将[]byte转成16进制

	//ss := sha256.Sum256(d.NoiseKey.Priv[:])
	//result := fmt.Sprintf("%s", ss)
	//spew.Dump(result)
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

// go判断字符串是否是gbk:
func Test_isGBK(t *testing.T) {

	//storeContainer, err := sqlstore.New("sqlite3",
	//	fmt.Sprintf("file:%s?_foreign_keys=off", "E:\\software\\sqlite3\\db\\8613027979536_back.db"), nil)
	////storeContainer, err := sqlstore.New(*dbDialect, *dbAddress, dbLog)
	//if err != nil {
	//	log.Errorf("Failed to connect to database: %v", err)
	//	return
	//}

	str := "Njl4MzEyOTc0MDY3NiOaxnajjJceHi43lHprq1D9dK0FnA=="

	mm, err := base64.StdEncoding.DecodeString(str)

	spew.Dump(err)
	spew.Dump(string(mm))

	spew.Dump(fmt.Sprintf("%s", md5.Sum([]byte(str))))

}

// go判断字符串是否是gbk:
func isGBK(data []byte) bool {
	length := len(data)
	var i int = 0
	for i < length {
		if data[i] <= 0x7f {
			//编码0~127,只有⼀个字节的编码，兼容ASCII码
			i++
			continue
		} else {
			//⼤于127的使⽤双字节编码，落在gbk编码范围内的字符
			if data[i] >= 0x81 &&
				data[i] <= 0xfe &&
				data[i+1] >= 0x40 &&
				data[i+1] <= 0xfe &&
				data[i+1] != 0xf7 {
				i += 2
				continue
			} else {
				return false
			}
		}
	}
	return true
}

// go判断字符串是否是utf-8:
func preNUm(data byte) int {
	var mask byte = 0x80
	var num int = 0
	//8bit中⾸个0bit前有多少个1bits
	for i := 0; i < 8; i++ {
		if (data & mask) == mask {
			num++
			mask = mask >> 1
		} else {
			break
		}
	}
	return num
}
func isUtf8(data []byte) bool {
	i := 0
	for i < len(data) {
		if (data[i] & 0x80) == 0x00 {
			// 0XXX_XXXX
			i++
			continue
		} else if num := preNUm(data[i]); num > 2 {
			// 110X_XXXX 10XX_XXXX
			// 1110_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_0XXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_10XX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// 1111_110X 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX 10XX_XXXX
			// preNUm() 返回⾸个字节的8个bits中⾸个0bit前⾯1bit的个数，该数量也是该字符所使⽤的字节数   i++
			for j := 0; j < num-1; j++ {
				//判断后⾯的 num - 1 个字节是不是都是10开头
				if (data[i] & 0xc0) != 0x80 {
					return false
				}
				i++
			}
		} else {
			//其他情况说明不是utf-8
			return false
		}
	}
	return true
}

//go gbk与utf-8相互转换:
//gbk与utf-8相互转换可以使⽤官⽅的golang.org/x/text 包来进⾏转换
//需要先安装这个第三⽅包，使⽤以下命令进⾏安装
//go get golang.org/x/text
//import"golang.org/x/text/encoding/simplifiedchinese"
//simplifiedchinese.GBK.NewEncoder().Bytes()//utf-8 转 gbk
//simplifiedchinese.GBK.NewDecoder().Bytes()//gbk 转 utf-8
//⽤例：
//const(
//	GBK string="GBK"
//	UTF8 string="UTF8"
//	UNKNOWN string="UNKNOWN"
//)
////需要说明的是，isGBK()是通过双字节是否落在gbk的编码范围内实现的，
////⽽utf-8编码格式的每个字节都是落在gbk的编码范围内，
////所以只有先调⽤isUtf8()先判断不是utf-8编码，再调⽤isGBK()才有意义
//func GetStrCoding(data []byte)string{
//	if isUtf8(data)==true{
//		return UTF8
//	}else if isGBK(data)==true{
//		return GBK
//	}else{
//		return UNKNOWN
//	}
//}
//func main(){
//	str :="⽉⾊真美，风也温柔，233333333，~！@#"//go字符串编码为utf-8
//	fmt.Println("before
//
//
//	展开
//
//	你想对它做什么：
//	总结概要内容
//
//	解释一下
//
//	中英互译
//
//	你已选取文本：
//
//
//	func isGBK(data []byte)bool{
//		length :=len(data)
//		var i int=0
//		for i < length {
//		if data[i]<=0x7f{
//		//编码0~127,只有⼀个字节的编码，兼容ASCII码
//		i++
//		continue
//	}else{
//		//⼤于127的使⽤双字节编码，落在gbk编码范围内的字符
//		if  data[i]>=0x81&&
//		data[i]<=0xfe&&
//		data[i +1]>=0x40&&
//		data[i +1]<=0xfe&&
//		data[i +1]!=0xf7{
//		i +=2
//		continue
//	}else{
//		return false
//	}
//	}
//	}
//		return true
//	}
