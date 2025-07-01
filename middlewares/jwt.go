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
		token := c.Request.Header.Get(global.Config.Jwt.Admin.AccessToken)
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
		rToken, err := cache.GetUserAToken(c, payLoad.UserName)
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

// 验证JWT管理员
func VerifyJWTAdminV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer("sky-take-out")
		ctx2, span := tracer.Start(c, "VerifyJWTAdminV1")
		defer span.End()
		code := e.SUCCESS
		access_token := c.GetHeader(global.Config.Jwt.Admin.AccessToken)
		refresh_token := c.GetHeader(global.Config.Jwt.Admin.RefreshToken)
		if access_token == "" {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("c.GetHeader Error")
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 刷新access_token和refresh_token
		newAccessToken, newRefreshToken, err := utils.ParseRefreshToken(access_token, refresh_token, global.Config.Jwt.Admin.Secret)
		if err != nil {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("utils.ParseRefreshToken Error", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 解析获取用户载荷信息
		claims, err := utils.ParseToken(newAccessToken, global.Config.Jwt.Admin.Secret)
		if err != nil {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("utils.ParseToken Error", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// === 新增：写入Redis ===
		err = cache.StoreUserAToken(c, newAccessToken, claims.UserName)
		if err != nil {
			code = e.RedisERR
			logger.Logger(ctx2).Error("StoreUserAToken Error", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.Result{Code: code})
			c.Abort()
			return
		}
		err = cache.StoreUserRToken(c, newRefreshToken, claims.UserName)
		if err != nil {
			code = e.RedisERR
			logger.Logger(ctx2).Error("StoreUserRToken Error", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.Result{Code: code})
			c.Abort()
			return
		}
		// redis获取token
		rToken, err := cache.GetUserAToken(c, claims.UserName)
		// 如果失败，分别为内部错误和未登录
		if err == retcode.NewError(e.RedisERR, "rdb.Get failed") {
			code = e.RedisERR
			logger.Logger(ctx2).Error("RedisERR", zap.Error(err))
			c.JSON(http.StatusBadGateway, common.Result{Code: code})
			c.Abort()
			return
		}
		if err == retcode.NewError(e.ErrorUserNotLogin, "token is invalid") {
			code = e.ErrorUserNotLogin
			logger.Logger(ctx2).Error("ErrorUserNotLogin", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		// 如果无错误，而是token不一致，则说明在另一端登录
		if access_token != rToken {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("Repeated login", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		SetToken(c, newAccessToken, newRefreshToken)
		// 在上下文设置载荷信息
		c.Set(enum.CurrentId, claims.UserId)
		c.Set(enum.CurrentName, claims.UserName)
		// 这里是否要通知客户端重新保存新的Token
		c.Next()
	}
}

func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(global.Config.Jwt.Admin.AccessToken, accessToken)
	c.Header(global.Config.Jwt.Admin.RefreshToken, refreshToken)
	c.SetCookie(global.Config.Jwt.Admin.AccessToken, accessToken, global.Config.Jwt.Cookie.MaxAge, "/", "", secure, true)
	c.SetCookie(global.Config.Jwt.Admin.RefreshToken, refreshToken, global.Config.Jwt.Cookie.MaxAge, "/", "", secure, true)
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(global.Config.Jwt.Https.HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}

// 验证JWT用户
func VerifyJWTUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := e.SUCCESS
		token := c.Request.Header.Get(global.Config.Jwt.User.AccessToken)
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
