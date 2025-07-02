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

// 验证JWT管理员 单token
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
		// 只做token一致性校验
		if token != rToken {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("Token不一致，可能存在重复登录",
				zap.String("requestToken", token[:10]+"..."),
				zap.String("redisToken", rToken[:10]+"..."))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		c.Set(enum.CurrentId, payLoad.UserId)
		c.Set(enum.CurrentName, payLoad.UserName)
		c.Next()
	}
}

// 验证JWT管理员 双token
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
		claims, err := utils.ParseToken(newAccessToken, global.Config.Jwt.Admin.Secret)
		if err != nil {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("utils.ParseToken Error", zap.Error(err))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
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
		rToken, err := cache.GetUserAToken(c, claims.UserName)
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
		// 只做token一致性校验
		if newAccessToken != rToken {
			code = e.UNKNOW_IDENTITY
			logger.Logger(ctx2).Error("Token不一致，可能存在重复登录",
				zap.String("newToken", newAccessToken[:10]+"..."),
				zap.String("redisToken", rToken[:10]+"..."))
			c.JSON(http.StatusUnauthorized, common.Result{Code: code})
			c.Abort()
			return
		}
		SetToken(c, newAccessToken, newRefreshToken)
		c.Set(enum.CurrentId, claims.UserId)
		c.Set(enum.CurrentName, claims.UserName)
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
