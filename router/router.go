package router

import (
	"github.com/gin-gonic/gin"
	"github.com/jicg/AppBg/middleware/jwt"
	"github.com/jicg/AppBg/handle/login"
	"github.com/jicg/AppBg/handle/offorder/order"
	"github.com/jicg/AppBg/handle/user"
	"github.com/jicg/AppBg/handle/pay"
	"github.com/jicg/AppBg/handle/offorder/setting"
)

var (
	JWT_KEY = "jicg"
)

func Router(r *gin.Engine) *gin.Engine {
	jwtauth.SetSignKey(JWT_KEY)
	// 无权限
	v1 := r.Group("/v1")
	no_auth(v1)
	// 权限
	v1_auth := v1.Group("/", jwtauth.JWTAuth())
	auth(v1_auth)
	return r
}

func no_auth(r *gin.RouterGroup) {
	r.POST("/login", login.Login)
	r.GET("/pay/:no", pay.WxPay)
	payRouter := r.Group("offline")
	{
		payRouter.GET("/wxpay/:id", order.WxPay)
		payRouter.GET("/ws", order.WsPayHandler)
		payRouter.GET("/wstext/:id", order.Text)
		xlsOrderRouter := payRouter.Group("xls/order")
		{
			xlsOrderRouter.GET("", order.QueryXlsOrders)
		}
	}

}

func auth(r *gin.RouterGroup) {
	userRouter := r.Group("user")
	{
		userRouter.GET("", user.GetUser)
		userRouter.PATCH("pwdchange", user.ChangePwd)
	}

	offline := r.Group("/offline")
	{
		settingR := offline.Group("setting")
		{
			settingR.GET("", setting.GetSetting)
			settingR.POST("", setting.Update)
		}
		orderRouter := offline.Group("order")
		{
			orderRouter.GET("", order.QueryOrders)
			orderRouter.POST("", order.Add)
			orderRouter.GET("/:id", order.GetById)
			orderRouter.PUT("/:id", order.Update)
			orderRouter.DELETE("/:id", order.Del)
		}

	}
}
