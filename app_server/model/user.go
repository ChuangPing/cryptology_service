package model

import (
	"education/app_server/service"
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	"github.com/sirupsen/logrus"
	"net/http"
)

// CheckUser check user email only
func CheckUser(email string, serverInfo service.InvokeServer) (respone channel.Response, err error) {
	req := channel.Request{
		ChaincodeID: serverInfo.ChaincodeID,
		Fcn:         "checkUser",
		Args: [][]byte{
			[]byte(email),
		},
	}
	respone, err = serverInfo.Client.Execute(req)
	fmt.Printf("client execute respone:%+v\n", respone)
	if err != nil && respone.ChaincodeStatus != http.StatusNotFound {
		logrus.Error("client execute  checkUser failed,err:", err)
		return
	}
	return
}

// Register user register
func Register(user User, serverInfo service.InvokeServer) (respone channel.Response, err error) {
	// 将user对象序列化成为字节数组
	//userByte, err := json.Marshal(user)
	//if err != nil {
	//	logrus.Error("json marshal failed, err:", err)
	//	return
	//}

	req := channel.Request{
		ChaincodeID: serverInfo.ChaincodeID,
		Fcn:         "userRegister",
		Args: [][]byte{
			[]byte(user.Role),
			[]byte(user.UserName),
			[]byte(user.Password),
			[]byte(user.Email),
		},
	}

	respone, err = serverInfo.Client.Execute(req)
	if err != nil {
		logrus.Error("client execute failed,err:", err)
		return
	}
	return
}

// Login user login
func Login(userEmail string, serverInfo service.InvokeServer) (respone channel.Response, err error) {
	req := channel.Request{
		ChaincodeID: serverInfo.ChaincodeID,
		Fcn:         "login",
		Args: [][]byte{
			[]byte(userEmail),
		},
	}

	respone, err = serverInfo.Client.Execute(req)
	if err != nil {
		logrus.Error("client execute failed,err:", err)
		return
	}
	return
}
