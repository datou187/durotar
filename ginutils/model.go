package ginutils

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"

	"talent.com/server/durotar.git/user"
)

const (
	//通用错误码
	OK          = 0
	ServerError = 1
	ParamError  = 2
	DataError   = 3
	NotAllowed  = 4
	NotFound    = 5
	OverLimit   = 6
)

type Resp struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Data   interface{} `json:"data"`
}

//通过header验证用户身份
func TokenChecker(db *gorm.DB, allowAnonymous bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		userStr := c.Request.Header.Get("X-AUTH-USER")
		token := c.Request.Header.Get("X-AUTH-TOKEN")
		userId, _ := strconv.ParseUint(userStr, 10, 64)
		if checked, err := user.CheckToken(db, userId, token); !checked || err != nil {
			userId = 0
			return
		}
		if userId == 0 {
			if allowAnonymous {
				c.Set("UserID", 0)
			} else {
				c.String(http.StatusForbidden, "token invalid")
				return
			}
		}
		c.Set("UserID", userId)
		c.Next()
	}
}

//不使用gin自带的logger，允许所有跨域
func GetLeafLoggerRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = append(config.AllowHeaders, "X-AUTH-USER", "X-AUTH-TOKEN", "X-SIGN", "X-AUTH-SOURCE")
	router.Use(cors.New(config))
	router.Use(gin.Recovery())
	return router
}
