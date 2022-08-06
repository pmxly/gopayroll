package middleware

import (
	"gopayroll/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"time"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var basicParam common.BasicParam
		if err := c.ShouldBindBodyWith(&basicParam, binding.JSON); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//log.Println(basicParam)
		secretToken := basicParam.SecretToken
		dateStr, err := common.Decrypt(secretToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		loc, _ := time.LoadLocation(common.StdLocation)
		secretDate, _ := time.ParseInLocation(common.TimeLayout, dateStr, loc)
		curDate := common.CurStdDate()
		if secretDate.Before(curDate) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token is expired"})
			return
		}
		c.Next()
	}
}