package tool

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
)

type Aes struct {
	symmetricKey string // 对称密钥  -- 其它
	fileName     string // 文件路径
}

// 初始化结构体
func NewAes(key, fileName string) *Aes {
	return &Aes{
		symmetricKey: key,
		fileName:     fileName,
	}
}

//加密过程：
//  1、处理数据，对数据进行填充，采用PKCS7（当密钥长度不够时，缺几位补几个几）的方式。
//  2、对数据进行加密，采用AES加密方法中CBC加密模式
//  3、对得到的加密数据，进行base64加密，得到字符串
// 解密过程相反

//16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
//key不能泄露
//var PwdKey = []byte("ABCDABCDABCDABCD")

//pkcs7Padding 填充
func pkcs7Padding(data []byte, blockSize int) []byte {
	//判断缺少几位长度。最少1，最多 blockSize
	padding := blockSize - len(data)%blockSize
	//补足位数。把切片[]byte{byte(padding)}复制padding个
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

//pkcs7UnPadding 填充的反向操作
func pkcs7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("加密字符串错误！")
	}
	//	获取填充的个数
	unPadding := int(data[length-1])
	return data[:(length - unPadding)], nil
}

//AESEncrypt 加密
func aesEncrypt(data []byte, key []byte) ([]byte, error) {
	//	创建加密实例
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//	判断加密快的大小
	blockSize := block.BlockSize()
	//	填充
	encryptBytes := pkcs7Padding(data, blockSize)
	//	初始化加密数据接受切片
	crypted := make([]byte, len(encryptBytes))
	//	使用cbc加密模式
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	//	执行加密
	blockMode.CryptBlocks(crypted, encryptBytes)
	return crypted, nil
}

//AESDecrypt 解密
func aesDecrypt(data []byte, key []byte) ([]byte, error) {
	//	创建实例
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//	获取快的大小
	blockSize := block.BlockSize()
	//	使用cbc
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	//	初始化解密数据接收切片
	crypted := make([]byte, len(data))
	//	执行解密
	blockMode.CryptBlocks(crypted, data)
	//	去除填充
	crypted, err = pkcs7UnPadding(crypted)
	if err != nil {
		return nil, err
	}
	return crypted, nil
}

//EncryptByAes 加密，使用base64编码
func encryptByAes(data []byte, key string) (string, error) {
	res, err := aesEncrypt(data, []byte(key))
	if err != nil {
		return "", nil
	}
	return base64.StdEncoding.EncodeToString(res), nil
}

//使用DecryptAes 解密并使用base64进行解码
func decryptByAes(data string, key string) ([]byte, error) {
	dataByte, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	return aesDecrypt(dataByte, []byte(key))
}

//文件的加解密函数  --对外暴露的接口

//	EncryptFile 文件加密， filePath 需要加密的文件路径，fName 加密后的文件名
func (ae *Aes) EncryptFile() (err error) {
	fName := ae.fileName
	//加密文件路径
	filePath := "/root/go/src/education/app_server/upload/" + fName
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("未找到文件")
		return
	}
	defer f.Close()

	fInfo, _ := f.Stat()
	fLen := fInfo.Size()
	fmt.Println("待处理文件大小:", fLen)
	maxLen := 1024 * 1024 * 100 //100mb  每 100mb 进行加密一次
	var forNum int64 = 0
	getLen := fLen

	if fLen > int64(maxLen) {
		getLen = int64(maxLen)
		forNum = fLen / int64(maxLen)
		fmt.Println("需要加密次数：", forNum+1)
	}
	//加密后存储的文件
	ff, err := os.OpenFile("/root/go/src/education/app_server/upload/encrypt/encryptFile_"+fName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件写入错误")
		return err
	}
	defer ff.Close()
	//循环加密，并写入文件
	for i := 0; i < int(forNum+1); i++ {
		a := make([]byte, getLen)
		n, err := f.Read(a)
		if err != nil {
			fmt.Println("文件读取错误")
			return err
		}
		getByte, err := encryptByAes(a[:n], ae.symmetricKey)
		if err != nil {
			fmt.Println("加密错误")
			return err
		}
		getBytes := append([]byte(getByte), []byte("\n")...)
		//写入
		buf := bufio.NewWriter(ff)
		buf.WriteString(string(getBytes[:]))
		buf.Flush()
	}
	ffInfo, _ := ff.Stat()
	fmt.Printf("文件加密成功，生成文件名为：%s，文件大小为：%v Byte \n", ffInfo.Name(), ffInfo.Size())
	return nil
}

//DecryptFile 文件解密
func (ae *Aes) DecryptFile() (err error) {
	fName := ae.fileName
	filePath := "/root/go/src/education/app_server/upload/encrypt/encryptFile_" + fName
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println("未找到文件")
		return
	}
	defer f.Close()
	fInfo, _ := f.Stat()
	fmt.Println("待处理文件大小:", fInfo.Size())

	br := bufio.NewReader(f)
	ff, err := os.OpenFile("/root/go/src/education/app_server/upload/decrypt/decryptFile_"+fName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件写入错误")
		return err
	}
	defer ff.Close()
	num := 0
	//逐行读取密文，进行解密，写入文件
	for {
		num = num + 1
		a, err := br.ReadString('\n')
		fmt.Println("a", a)
		if err != nil {
			break
		}
		getByte, err := decryptByAes(a, ae.symmetricKey)
		if err != nil {
			fmt.Println("解密错误")
			return err
		}

		buf := bufio.NewWriter(ff)
		buf.Write(getByte)
		buf.Flush()
	}
	fmt.Println("解密次数：", num)
	ffInfo, _ := ff.Stat()
	fmt.Printf("文件解密成功，生成文件名为：%s，文件大小为：%v Byte \n", ffInfo.Name(), ffInfo.Size())
	return
}
