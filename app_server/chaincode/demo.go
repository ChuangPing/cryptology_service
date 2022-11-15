package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"net/http"
)

// User user info
type User struct {
	Role     string `form:"role" json:"role"`
	UserName string `form:"username" json:"userName"`
	Password string `form:"password" json:"password"`
	Email    string `form:"email" json:"email"`
}

// SymmetricKeyEnc  symmetricKey Enc info
type SymmetricKeyEnc struct {
	C1 string `json:"c1" form:"c1"`
	C2 string `json:"c2" form:"c2"`
	C3 string `json:"c3" form:"c3"`
	C4 string `json:"c4" form:"c4"`
	C5 string `json:"c5" form:"c5"`
}

type DemoChaincode struct {
}

func (dc *DemoChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

// user info id
func constructUserKey(userId string) string {
	return fmt.Sprintf("user_%s", userId)
}

// symmetricKeyEnc info id
func constructSymmetricKeyEncKey(userId string) string {
	return fmt.Sprintf("symmetricKeyEnc_%s", userId)
}

// Invoke 用户开户
func (dc *DemoChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcStr, args := stub.GetFunctionAndParameters()
	switch funcStr {
	case "userRegister":
		return dc.userRegister(stub, args)
	case "login":
		return dc.Login(stub, args)
	case "checkUser":
		return dc.checkUser(stub, args)
	case "storeSymmetricKey":
		return dc.storeSymmetricKey(stub, args)
	default:
		return shim.Error(fmt.Sprintf("unsupported function:%s", funcStr))
	}
}

// userRegister user register
func (dc *DemoChaincode) userRegister(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 套路1：检查参数的个数
	if len(args) != 4 {
		return shim.Error(fmt.Sprintf("param:%s err", args))
	}

	// 套路2：验证参数的正确性
	roleByte := args[0]
	usernameByte := args[1]
	passwordByte := args[2]
	emailByte := args[3]
	if roleByte == "" || usernameByte == "" || passwordByte == "" || emailByte == "" {
		return shim.Error(fmt.Sprintf("invaild args user"))
	}

	// 套路3：验证数据是否存在 应该存在 or 不应该存在 --
	//if userBytes, err := stub.GetState(constructUserKey(id)); err == nil && len(userBytes) != 0 {
	//	return shim.Error("user already exist")
	//}

	// 套路4：写入状态
	registerUser := &User{
		Role:     roleByte,
		UserName: usernameByte,
		Password: passwordByte,
		Email:    emailByte,
	}
	// 序列化对象
	registerUserByte, err := json.Marshal(registerUser)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal user error:%s", err))
	}

	if err := stub.PutState(constructUserKey(registerUser.Email), registerUserByte); err != nil {
		return shim.Error(fmt.Sprintf("register user err:%s", err))

	}
	// 成功返回
	return shim.Success(nil)
}

// checkUser check user email is only
func (dc *DemoChaincode) checkUser(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	res := peer.Response{}
	if len(args) != 1 {
		res.Status = http.StatusInternalServerError
		res.Message = fmt.Sprintf("invalid arg:%s", args[0])
		return res
	}

	email := args[0]
	if email == "" {
		res.Status = http.StatusInternalServerError
		res.Message = fmt.Sprintf("invalid arg:%s", args[0])
		return res
	}

	userBytes, err := stub.GetState(constructUserKey(email))
	if err != nil || len(userBytes) == 0 {
		//res.Status = http.StatusNotFound
		//res.Message = fmt.Sprintf("user emil:%s not used", args[0])
		//return res
		return shim.Success(nil)
	} else {
		//res.Status = http.StatusOK
		//res.Message = fmt.Sprintf("user emil:%s is used", args[0])
		//return res
		return shim.Success(userBytes)
	}
}

// Login ch
func (dc *DemoChaincode) Login(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	res := peer.Response{}

	if len(args) != 1 {
		return shim.Error(fmt.Sprintf("param:%s err", args))
	}

	userEmail := args[0]
	if userEmail == "" {
		return shim.Error("invalid id")
	}

	userBytes, err := stub.GetState(constructUserKey(userEmail))
	if err != nil || len(userBytes) == 0 {
		//	user not found,are you register?
		return shim.Success(nil)
	}

	//Unmarshal
	user := new(User)
	err = json.Unmarshal(userBytes, user)
	if err != nil {
		res.Status = http.StatusInternalServerError
		res.Message = fmt.Sprintf("ivalid arg:%s", args[0])
		return res
		//return shim.Error(fmt.Sprintf("user unmarshal err:%s", err))
	}

	// return login data
	return shim.Success(userBytes)
}

func (dc *DemoChaincode) storeSymmetricKey(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 套路1：检查参数的个数
	if len(args) != 5 {
		return shim.Error(fmt.Sprintf("param:%s err", args))
	}
	email := args[0]
	C1 := args[1]
	C2 := args[2]
	C3 := args[3]
	C4 := args[4]
	C5 := args[5]

	// check param is valid
	if email == "" || C1 == "" || C2 == "" || C3 == "" || C4 == "" || C5 == "" {
		return shim.Error("param is invalid")
	}

	// init store blockchain data
	var symmetricKeyEnc = new(SymmetricKeyEnc)
	symmetricKeyEnc.C1 = C1
	symmetricKeyEnc.C2 = C2
	symmetricKeyEnc.C3 = C3
	symmetricKeyEnc.C4 = C4
	symmetricKeyEnc.C5 = C5
	symmetricKeyEncByte, err := json.Marshal(symmetricKeyEnc)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal symmetricKeyEnc error:%s", err))
	}

	// store blockchain
	if err := stub.PutState(constructSymmetricKeyEncKey(email), symmetricKeyEncByte); err != nil {
		return shim.Error(fmt.Sprintf("store symmetricKeyEnc failed,PutState:err%v", err))
	}

	// TODO log need to store blockChain
	// 成功返回
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(DemoChaincode))
	if err != nil {
		fmt.Printf("启动EducationChaincode时发生错误: %s", err)
	}
}
