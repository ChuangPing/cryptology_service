package v1

import (
	"education/app_server/model"
	"education/app_server/service"
	"education/app_server/tool"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func AesHandler(ctx *gin.Context) {
	var aesFormData model.AesFormData
	err := ctx.ShouldBind(&aesFormData)
	if err != nil {
		logrus.Error("aesFormData err:", err)
		//	向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "数据处理错误",
		})
		return
	}
	fmt.Printf("aesformdata:%+v\n", aesFormData)

	// check email
	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(aesFormData.Email) //匹配成功返回比配后的字符串，不成功返回空
	if res == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "邮箱格式错误",
		})
		fmt.Println("注册时：邮箱格式不正确，请重新输入")
		return
	}

	//check email is register
	response, err := model.CheckUser(aesFormData.Email, *service.Server)
	if err != nil {
		logrus.Error("checkUser failed, err:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "inner err",
		})
		return
	}
	// check chaincode is err
	if response.Payload == nil {
		// email is not register
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "账号邮箱还未注册",
		})
		return
	}

	// check symmetricKey
	if len(aesFormData.SymmetricKey) != 16 || len(aesFormData.SymmetricKey) != 24 || len(aesFormData.SymmetricKey) != 32 {
		logrus.Warnf("symmetricKey invalid,err:%s", aesFormData.SymmetricKey)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "对称密钥非法",
		})
		return
	}

	// 检查对称密钥
	if strings.TrimSpace(aesFormData.SymmetricKey) == "" {
		//	向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "对称密钥不能为空",
		})
		return
	}

	// receive file
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println(err)
		//	向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "文件读取错误",
		})
		return
	}
	// 获取当前时间戳
	timeInt := time.Now().Unix()
	timeStr := strconv.FormatInt(timeInt, 10)
	// 为了防止上传同一文件导致覆盖，给文件名加上时间戳   E:/LearningCode/Gin_project/Gin_demo/upload/ 斜杠结束表示存储的路径-那个文件夹
	dist := "/root/go/src/education/app_server/upload/" + timeStr + file.Filename
	fmt.Printf("file:%+v\n", file)
	err = ctx.SaveUploadedFile(file, dist)
	if err != nil {
		fmt.Println("文件上传失败", err)
		//	向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "文件上传失败",
		})
		return
	}

	//初始化aes加密工具函数
	aes := tool.NewAes(aesFormData.SymmetricKey, timeStr+file.Filename)
	//调用aes将文件进行加密
	err = aes.EncryptFile()
	// 测试解密
	err = aes.DecryptFile()
	if err != nil {
		fmt.Println("文件加密失败", err)
		//	向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "文件加密失败",
		})
		return
	}

	//创建用户加密信息表
	db := model.NewGorm()
	// 创建用户上传文件加密信息表
	err = db.AutoMigrate(&model.CryptInfo{})
	if err != nil {
		fmt.Println("创建加密信息表失败", err)
		//	向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "文件加密失败",
		})
		return
	}

	// 验证通过插入数据库
	var cryptInfo model.CryptInfo
	cryptInfo = model.CryptInfo{
		Username: aesFormData.Username,
		Email:    aesFormData.Email,
		//加密文件上传时间
		CreateDate:     time.Now(),
		KeyWords:       aesFormData.KeyWords,
		IsShare:        aesFormData.IsShare,
		AesEncFileName: "encryptFile_" + timeStr + file.Filename, // 用户上传文件名为 当前时间戳加文件名
		SymmetricKey:   aesFormData.SymmetricKey,
	}
	//创建用户加密信息表
	result := db.Create(&cryptInfo)
	if result.Error != nil {
		fmt.Println(result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError, // 400
			"msg":  "服务器内部错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "加密成功",
	})
	//fmt.Printf("aesformdata:%+v\n", aesFormData)
}
