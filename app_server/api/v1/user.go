package v1

import (
	"education/app_server/model"
	"education/app_server/service"
	"education/app_server/tool"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
	"strings"
)

func RegisterHandler(ctx *gin.Context) {
	userInfo := new(model.RegisterUser)
	err := ctx.ShouldBind(userInfo)
	if err != nil {
		//向前端返回错误
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "服务端处理数据出错",
		})
		fmt.Println("err:\n", err)
		return
	}
	fmt.Printf("user info %v\n", userInfo)
	//服务端验证前端传递的数据
	if strings.TrimSpace(userInfo.UserName) == "" || strings.TrimSpace(userInfo.Password) == "" || strings.TrimSpace(userInfo.Email) == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "注册时：注册的相关信息不能为空",
		})
		//打印日志
		fmt.Println("注册时：注册的相关信息不能为空", userInfo)
		// 阻止代码继续运行
		return
	}

	if userInfo.VerifyPassword != userInfo.Password {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "两次密码不一致，请检查",
		})
		//打印日志
		fmt.Println("注册时：两次密码不一致，请检查")
		// 阻止代码继续运行
		return
	}

	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(userInfo.Email) //匹配成功返回比配后的字符串，不成功返回空
	if res == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "邮箱格式错误",
		})
		fmt.Println("注册时：邮箱格式不正确，请重新输入")
		return
	}

	// 将校验通过的数据存储在数据库
	var user model.User

	user.Role = userInfo.Role
	user.UserName = userInfo.UserName
	// 将用户密码进行加密
	user.Password = tool.ScryptyPw(userInfo.Password)
	user.Email = userInfo.Email

	// 校验email -- only user email
	response, err := model.CheckUser(user.Email, *service.Server)
	if err != nil {
		logrus.Error("checkUser failed, err:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "inner err",
		})
		return
	}
	fmt.Printf("see return response:%v", response)
	// check chaincode is err
	if response.Payload != nil {
		// email is used
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "用户邮箱已被注册",
		})
		return
	}
	// 验证通过插入数据库
	response, err = model.Register(user, *service.Server)
	if err != nil {
		logrus.Error("checkUser failed, err:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "inner err",
		})
		return
	}
	// check chaincode is err
	if response.ChaincodeStatus != 200 {
		// email is used
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "chaincode err",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"msg":  "注册成功",
	})
}

func Login(ctx *gin.Context) {
	var loginInfo model.LoginUser
	err := ctx.ShouldBind(&loginInfo)
	if err != nil {
		logrus.Error("get login info err:", err)
		// 向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "数据处理错误",
		})
		return
	}
	fmt.Printf("loginInfo:", loginInfo)
	// 验证邮箱合法性
	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(loginInfo.Email) //匹配成功返回比配后的字符串，不成功返回空
	if res == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "邮箱格式错误",
		})
		fmt.Println("登录时：邮箱格式不正确，请重新输入")
		return
	}

	// 根据登录信息查询数据库
	response, err := model.Login(loginInfo.Email, *service.Server)
	if err != nil {
		logrus.Error("Login failed, err:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "inner err",
		})
		return
	}
	// check user is register
	if response.Payload == nil {
		// email is not register
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "账号邮箱还未注册",
		})
		return
	}
	userByte := response.Payload
	fmt.Printf("userByte:%+v\n", userByte)
	//unmarshal
	var user model.User
	err = json.Unmarshal(userByte, &user)
	if err != nil {
		logrus.Error("login info unmarshal failed,err:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "数据处理错误",
		})
		return
	}

	// 比对密码
	if user.Password == tool.ScryptyPw(loginInfo.Password) {

		ctx.JSON(http.StatusOK, gin.H{
			"code":     http.StatusOK,
			"msg":      "登录成功",
			"role":     user.Role,
			"username": user.UserName,
			"email":    user.Email,
		})
	} else {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusUnauthorized,
			"msg":  "密码或账户错误",
			"role": user.Role,
		})
	}
}
