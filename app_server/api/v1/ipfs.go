package v1

import (
	"education/app_server/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetFileInfo get store to ipfs file info
func GetFileInfo(ctx *gin.Context) {
	username := ctx.Query("username")
	fileName := ctx.Query("fileName")
	db := model.NewGorm()
	var cryptInfo model.CryptInfo
	result := db.Where("username=? AND aes_enc_file_name LIKE ?", username, "%"+fileName+"%").Find(&cryptInfo)
	if result.Error != nil {
		fmt.Println("获取加密文件信息失败", result.Error)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "获取加密文件信息失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"msg":      "获取加密文件信息成功",
		"username": cryptInfo.Username,
		"isShare":  cryptInfo.IsShare,
	})
}
