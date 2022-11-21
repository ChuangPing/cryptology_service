package main

import (
	"crypto/sha256"
	"fmt"
)

func main() {
	str := "chuanping"
	//hash := sha256.Sum256([]byte(str))
	//h := sha256.New()
	//h.Write([]byte(str))
	////hash := h.Sum()
	hash := sha256.Sum256([]byte(str))
	//sliceByte1 := make([]byte, len(hash))
	//sliceByte1 = append(sliceByte1, hash)

	fmt.Printf("hash:%x\n", hash)
	str1 := "chuanping"
	hash1 := sha256.Sum256([]byte(str1))

	for i := 0; i < len(hash); i++ {
		if hash[i] != hash1[i] {
			fmt.Println(0)
			return
		}
	}
	fmt.Println(1)

    //= time.Now().Format("2006-01-02 15:04:05")

	//var num int64 = 122
	////mystr := strconv.FormatInt(num, 10)
	////fmt.Printf("str:%s\n", mystr)
	////
	//strByte := Int64ToBytes(num)
	//fmt.Printf("strbyte:%v\n", strByte)
	//
	//nums := BytesToInt64(strByte)
	//fmt.Printf("nums:%d\n", nums)
}
//func BytesToInt64(buf []byte) int64 {
//	return int64(binary.BigEndian.Uint64(buf))
//}
//
//func Int64ToBytes(i int64) []byte {
//	var buf = make([]byte, 8)
//	binary.BigEndian.PutUint64(buf, uint64(i))
//	return buf
//}
