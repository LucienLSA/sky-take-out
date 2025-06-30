package middlewares

import (
	"net/http"
	"skytakeout/common"
	"skytakeout/common/e"
	"skytakeout/common/enum"
	"skytakeout/common/retcode"
	"skytakeout/common/utils"
	"skytakeout/global"
	"skytakeout/internal/cache"
	"skytakeout/logger"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

// 验证JWT管理员
func VerifyJWTAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("sky-take-out")
		ctx2, span := tracer.Start(c, "VerifyJWTAdmin")
		defer span.End()
		code := e.SUCCESS
		token := c.Request.Header.Get(global.Config.Jwt.Admin.Name)
		if token == "" {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("c.Request.Header.Get Error")
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, global.Config.Jwt.Admin.Secret)
		if err != nil {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("utils.ParseToken Error", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		rToken, err := cache.GetJwtToken(c, payLoad.UserName)
		// redis获取token失败，分别为内部错误和未登录
		if err == retcode.NewError(e.RedisERR, "rdb.Get failed") {
			code = e.RedisERR
			logger.Logger(ctx2).Error("RedisERR", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.Result{Code: code})
			c.Abort()
			return
		}
		if err == retcode.NewError(e.ErrorUserNotLogin, "token is empty") {
			code = e.ErrorUserNotLogin
			logger.Logger(ctx2).Error("ErrorUserNotLogin", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 如果无错误，而是token不一致，则说明在另一端登录
		if token != rToken {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("Repeated login", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 在上下文设置载荷信息
		c.Set(enum.CurrentId, payLoad.UserId)
		c.Set(enum.CurrentName, payLoad.UserName)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
	}
}

// 验证JWT用户
func VerifyJWTUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		token := c.Request.Header.Get(global.Config.Jwt.User.Name)
		if token == "" {
			code = e.UNKNOW_IDENTITY
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 解析获取用户载荷信息
		payLoad, err := utils.ParseToken(token, global.Config.Jwt.User.Secret)
		if err != nil {
			code = e.UNKNOW_IDENTITY
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 在上下文设置载荷信息
		c.Set(enum.CurrentId, payLoad.UserId)
		c.Set(enum.CurrentName, payLoad.UserName)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
	}
}
