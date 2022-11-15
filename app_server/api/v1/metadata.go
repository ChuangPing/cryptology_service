package v1

import (
	gRPCClient "education/app_server/gRPC/client"
	"education/app_server/model"
	"education/app_server/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"regexp"
)

func MetadataAndCm(ctx *gin.Context) {
	dataInfo := new(model.DataInfo)
	err := ctx.ShouldBind(dataInfo)
	if err != nil {
		logrus.Error("dataInfo err:", err)
		//	向前端提示错误
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "数据处理错误",
		})
		return
	}
	fmt.Printf("dataInfo:%+v\n", dataInfo)

	// use Email check user is authority
	reg, _ := regexp.Compile("^[A-Za-z0-9\u4e00-\u9fa5]+@[a-zA-Z0-9_-]+(\\.[a-zA-Z0-9_-]+)+$")
	res := reg.FindString(dataInfo.Email) //匹配成功返回比配后的字符串，不成功返回空
	if res == "" {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "邮箱格式错误",
		})
		fmt.Println("注册时：邮箱格式不正确，请重新输入")
		return
	}

	// 校验email -- only user email
	response, err := model.CheckUser(dataInfo.Email, *service.Server)
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
		// email is used
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "用户未注册",
		})
		return
	}

	// check user authority
	userByte := response.Payload
	fmt.Printf("userByte:%+v\n", userByte)
	//unmarshal
	var user model.User
	err = json.Unmarshal(userByte, &user)
	if err != nil {
		logrus.Error("user info unmarshal failed,err:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "数据处理错误",
		})
		return
	}
	if user.Role != "0" {
		logrus.Warnf("username:%v,role:%v not authority to store data to blockchain\n", user.UserName, user.Role)
		// user not authority to store data to blockchain
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "用户未注册",
		})
		return
	}

	// store data to blockchain

	// 1.encryption SymmetricKey
	C1, C2, C3, C4, C5, err := gRPCClient.SymmetricKeyEnc("asdfghjklkjhgfds")
	if err != nil {
		logrus.Error("gRPCClient:SymmetricKeyEnc failed,err", err)
		return
	}
	response, err = model.StoreSymmetricKeyEnc(dataInfo.Email, C1, C2, C3, C4, C5, *service.Server)
	if err != nil {
		logrus.Error("StoreSymmetricKeyEnc failed, err:", err)
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
		"msg":  "对称密钥加密成功",
	})
	return
	// 2.encryption keywords TODO

	// 3.set Metadata TODO

	// 4.store blockchain TODO

}
