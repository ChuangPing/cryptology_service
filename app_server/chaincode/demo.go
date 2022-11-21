package main

import (
	"crypto/sha256"
	bptree "education/app_server/chaincode/BPlusTree"
	"education/app_server/chaincode/tool"
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric-chaincode-go/shim"
	"github.com/hyperledger/fabric-protos-go/peer"
	"log"
	"net/http"
	"time"
)

// AuthorityType
const (
	PRECISE_AUTHORITY     string = "precise"
	FUZZY_AUTHORITY       string = "fuzzy"
	ALLORCANCEL_AUTHORITY string = "all_cancel"
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

// KeywordIndexInfo keywordIndex info
type KeywordIndexInfo struct {
	CurrentIndex    int64              `json:"currentIndex" form:"currentIndex"`
	KeywordIndexMap []map[string]int64 `json:"keywordIndexMap" form:"keywordIndexMap"`
}

// KeywordEnc info
type KeywordEnc struct {
	A  string `form:"a" json:"a"`
	B  string `form:"b" json:"b"`
	Ci string `form:"ci" json:"ci"`
}

// Metadata info
type Metadata struct {
	AuthorityType string   `form:"authorityType" json:"authorityType"`
	KeywordIndex  int64    `form:"keywordIndex" json:"keywordIndex"`
	AuthorityHash [32]byte `form:"authorityHash" json:"authorityHash"`
	KeywordEnc    `form:"keyword" json:"keywordEnc"`
}

//	MetadataLog when store metadata to blockchain
type MetadataLog struct {
	Username     string `form:"username" json:"username"`
	Email        string `form:"email" json:"email"`
	Role         string `form:"role" json:"role"`
	KeywordIndex int64  `form:"keywordIndex" json:"keywordIndex"`
	Time         string `json:"time"`
	Operation    string `form:"operation" json:"operation"`
	Result       string `form:"result" json:"result"`
}
type DemoChaincode struct {
}

func (dc *DemoChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println(" ==== Init ====")

	return shim.Success(nil)
}

// keyword index
var (
	keywordIndex int64
	keywordMap   map[string]int64
)

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
	case "setKeywordIndex":
		return dc.setKeywordIndex(stub, args)
	case "storeSymmetricKeyEnc":
		return dc.storeSymmetricKeyEnc(stub, args)
	case "setMetadata":
		return dc.storeMetadata(stub, args)
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

	if err := stub.PutState(tool.ConstructUserKey(registerUser.Email), registerUserByte); err != nil {
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

	userBytes, err := stub.GetState(tool.ConstructUserKey(email))
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

	userBytes, err := stub.GetState(tool.ConstructUserKey(userEmail))
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

func (dc *DemoChaincode) storeSymmetricKeyEnc(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 套路1：检查参数的个数
	if len(args) != 6 {
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
	if err := stub.PutState(tool.ConstructSymmetricKeyEncKey(email), symmetricKeyEncByte); err != nil {
		return shim.Error(fmt.Sprintf("store symmetricKeyEnc failed,PutState:err%v", err))
	}

	// TODO log need to store blockChain
	// 成功返回
	return shim.Success(nil)
}

// store metadata to blockchain
func (dc *DemoChaincode) storeMetadata(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 套路1：检查参数的个数
	if len(args) != 6 {
		return shim.Error(fmt.Sprintf("param:%s err", args))
	}
	email := args[0]
	authority := args[1]
	keywordIndexStr := args[2]
	A := args[3]
	B := args[4]
	Ci := args[5]

	// check param is valid
	if email == "" || authority == "" || keywordIndexStr == "" || A == "" || Ci == "" || B == "" {
		return shim.Error("param is invalid")
	}
	keywordIndexByte := []byte(keywordIndexStr)
	index := tool.BytesToInt64(keywordIndexByte)

	// hash authority
	authorityHash := sha256.Sum256([]byte(authority))

	// use email get user info
	userBytes, err := stub.GetState(tool.ConstructUserKey(email))
	if err != nil || len(userBytes) == 0 {
		//	user not found,are you register?
		return shim.Success(nil)
	}

	//Unmarshal
	user := new(User)
	err = json.Unmarshal(userBytes, user)
	if err != nil {
		return shim.Error(fmt.Sprintf("user unmarshal err:%s", err))
	}

	// check user
	if user.Role != "0" {
		//	failed log
		metadataLog := new(MetadataLog)
		metadataLog.Username = user.UserName
		metadataLog.Email = user.Email
		metadataLog.Role = user.Role
		metadataLog.KeywordIndex = index
		metadataLog.Time = time.Now().Format("2006-01-02 15:04:05")
		metadataLog.Operation = "store metadata to blockchain"
		metadataLog.Operation = "failed"
		// store log to blockchain
		metadataLogByte, err := json.Marshal(metadataLog)
		if err != nil {
			log.Fatal("metadata marshal failed,err:", err)
		}
		err = stub.PutState(tool.ConstructMetadataLogKey(), metadataLogByte)
		if err != nil {
			return shim.Error(fmt.Sprintf("store metadata to blockchain failed,err:%v", err))
		}
		//	stop
		return shim.Error(fmt.Sprintf("user:%s role:%s not have authority", user.UserName, user.Role))
	}

	// init metadata
	metadata := Metadata{
		// precise authority
		AuthorityType: PRECISE_AUTHORITY,
		KeywordIndex:  index,
		AuthorityHash: authorityHash,
		KeywordEnc: KeywordEnc{
			A:  A,
			B:  B,
			Ci: Ci,
		},
	}
	//metadata := new(Metadata)
	//keywordEnc := KeywordEnc{
	//	A:  A,
	//	B:  B,
	//	Ci: Ci,
	//}
	//metadata.KeywordEnc = keywordEnc
	//metadata.AuthorityHash = authorityHash

	// Marshal metadata
	metadataByte, err := json.Marshal(metadata)
	if err != nil {
		log.Fatal("metadata marshal failed,err:", err)
	}

	//store metadata to bPlusTree
	// init bPTree
	var bPTree = bptree.NewBPTree(5)
	bPTree.Set(metadata.KeywordIndex, metadataByte)
	err = stub.PutState(tool.ConstructMetadataKey(email+keywordIndexStr), metadataByte)
	if err != nil {
		return shim.Error(fmt.Sprintf("store metadata to blockchain failed,err:%v", err))
	}

	// init success log
	metadataLog := new(MetadataLog)
	metadataLog.Username = user.UserName
	metadataLog.Email = user.Email
	metadataLog.Role = user.Role
	metadataLog.KeywordIndex = index
	metadataLog.Time = time.Now().Format("2006-01-02 15:04:05")
	metadataLog.Operation = "store metadata to blockchain"
	metadataLog.Operation = "success"

	// store log to blockchain
	metadataLogByte, err := json.Marshal(metadataLog)
	if err != nil {
		log.Fatal("metadata marshal failed,err:", err)
	}
	err = stub.PutState(tool.ConstructMetadataLogKey(), metadataLogByte)
	if err != nil {
		return shim.Error(fmt.Sprintf("store metadataLog to blockchain failed,err:%v", err))
	}
	return shim.Success(nil)
}

// setKeywordIndex store keyword to blockchain return index
func (dc *DemoChaincode) setKeywordIndex(stub shim.ChaincodeStubInterface, args []string) peer.Response {
	// 套路1：检查参数的个数
	if len(args) != 3 {
		return shim.Error(fmt.Sprintf("param:%s err", args))
	}
	email := args[0]
	keyword := args[1]
	filename := args[2]

	if email == "" || keyword == "" || filename == "" {
		return shim.Error("param is invalid")
	}
	mapKey := email + keyword + filename
	if keywordIndex != 0 {
		// get blockchain keyword info
		keywordIndexInfoBytes, err := stub.GetState(tool.ConstructKeywordIndexKey())
		if err != nil || len(keywordIndexInfoBytes) == 0 {
			return shim.Error("get blockchain keyword info failed")
		}
		var keywordIndexInfo = new(KeywordIndexInfo)
		err = json.Unmarshal(keywordIndexInfoBytes, keywordIndexInfo)
		if err != nil {
			return shim.Error("get blockchain keyword info success but Unmarshal failed")
		}
		keywordIndex = keywordIndexInfo.CurrentIndex
		// keyword index
		keywordIndex++
		keywordMap := make(map[string]int64)
		keywordMap[mapKey] = keywordIndex
		keywordIndexInfo.KeywordIndexMap = append(keywordIndexInfo.KeywordIndexMap, keywordMap)
		keywordIndexInfo.CurrentIndex = keywordIndex

		// store to blockchain
		storeKeywordIndexInfoBytes, err := json.Marshal(keywordIndexInfo)
		if err != nil {
			return shim.Error("get blockchain keyword info success but marshal failed")
		}
		err = stub.PutState("keywordIndexInfo", storeKeywordIndexInfoBytes)
		if err != nil {
			return shim.Error("store blockchain keywordIndex info failed")
		}

		// return index
		indexByte := tool.Int64ToBytes(keywordIndexInfo.CurrentIndex)
		return shim.Success(indexByte)
	} else {
		// first store keyword to blockchain
		keywordIndexInfo := KeywordIndexInfo{}
		keywordIndex++
		keywordMap := make(map[string]int64)
		keywordMap[mapKey] = keywordIndex

		keywordIndexInfo.KeywordIndexMap = append(keywordIndexInfo.KeywordIndexMap, keywordMap)
		keywordIndexInfo.CurrentIndex = keywordIndex

		// store to blockchain
		storeKeywordIndexInfoBytes, err := json.Marshal(keywordIndexInfo)
		if err != nil {
			return shim.Error("get blockchain keyword info success but marshal failed")
		}
		err = stub.PutState("keywordIndexInfo", storeKeywordIndexInfoBytes)
		if err != nil {
			return shim.Error("store blockchain keywordIndex info failed")
		}

		// return index
		indexByte := tool.Int64ToBytes(keywordIndexInfo.CurrentIndex)
		return shim.Success(indexByte)
	}

}

func main() {
	err := shim.Start(new(DemoChaincode))
	if err != nil {
		fmt.Printf("启动EducationChaincode时发生错误: %s", err)
	}
}
