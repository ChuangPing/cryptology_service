package v1

import (
	"education/app_server/model"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func InterestKeywordHandle(ctx *gin.Context) {
	// get param from client
	interestKeyword := ctx.Query("interestKeyword")
	if interestKeyword == "" {
		logrus.Error("get client param is invalid,interestKeyword:", interestKeyword)
		ctx.JSON(http.StatusOK, gin.H{
			"code": http.StatusBadRequest,
			"msg":  "参数无效",
		})
	}

	// use interestKeyword select mysql
	db := model.NewGorm()
	var cryptInfo model.CryptInfo
	result := db.Where("key_words LIKE ?", interestKeyword+"%").Find(&cryptInfo)
	if result.Error != nil {
		logrus.Warnf("DU user interestKeyword not store system, intersetKeyword:%s\n", interestKeyword)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": http.StatusInternalServerError,
			"msg":  "系统中不存在此类关键字，请您重新选择！",
		})
		return
	}
}
