package model

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// RegisterUser register info
type RegisterUser struct {
	Role           string `form:"role" json:"role"`
	UserName       string `form:"username" json:"username"`
	Password       string `form:"password" json:"password"`
	VerifyPassword string `form:"verifyPassword" json:"verifyPassword"`
	Email          string `form:"email" json:"email"`
}

// User user info
type User struct {
	Role     string `form:"role"`
	UserName string `form:"username"`
	Password string `form:"password"`
	Email    string `form:"email" json:"email"`
}

// LoginUser user login info
type LoginUser struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
}

// AesFormData user crypto info
type AesFormData struct {
	Username string `json:"username" form:"username"` //用户名
	// user email TODO :form need change
	Email string `form:"email" json:"email"`
	//File         string `json:"file" form:"f"`         //需要加密的文件
	//Time         time.Time `form:"time" json:"time"`                 //发布时间
	IsShare      bool   `form:"isShare" json:"isShare"`           //是否上链
	SymmetricKey string `form:"symmetricKey" json:"symmetricKey"` // 对称密钥  -- 不存储数据库
	KeyWords     string `form:"keyWords" json:"keyWords"`         // 数据类型
}

//	用户加密信息表  -- 与用户表的关系是一对多的关系：一个用户有多个加密信息
type CryptInfo struct {
	Email          string
	Username       string    // 加密文件上传用户
	CreateDate     time.Time // 创建时间
	KeyWords       string    // 数据关键字
	IsShare        bool      // 是否分享
	AesEncFileName string    // aes 加密文件名
	SymmetricKey   string    `form:"symmetricKey" json:"symmetricKey"` // 对称密钥
	gorm.Model
}

// DataInfo metadataAndCm info
type DataInfo struct {
	AesEncFileName string `form:"aesEncFileName" json:"aesEncFileName"`
	//SymmetricKey string `form:"symmetricKey" json:"symmetricKey"` // 对称密钥  -- 不存储数据库
	FileHash string `form:"fileHash" json:"fileHash"`
	Email    string `form:"email" json:"email"`
}

// Metadata info
type Metadata struct {
	Index             int  `form:"index" json:"index"`
	Authority         byte `form:"authority" json:"authority"`
	KeywordCiphertext byte `form:"KeywordCiphertext" json:"keywordCiphertext"`
}

func NewGorm() *gorm.DB {
	dsn := "root:root@tcp(127.0.0.1:3306)/cryptology?charset=utf8mb4&parseTime=True&loc=Local"

	//全局模式配置logger -- 打印sql语句 每个abi背后的sql语句才是重点
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             time.Second, // 慢 SQL 阈值
			LogLevel:                  logger.Info, // 日志级别
			IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  true,        // 禁用彩色打印
		},
	)
	// 全局模式配置logger
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		fmt.Println("err:", err)
	}
	return db
}
