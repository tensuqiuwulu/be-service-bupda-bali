package routes

import (
	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/controller"
	authMiddlerware "github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
)

func AuthRoute(e *echo.Echo, authControllerInterface controller.AuthControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/auth/login", authControllerInterface.Login, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/auth/new-token", authControllerInterface.NewToken, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func OtpManagerRoute(e *echo.Echo, otpManagerControllerInterface controller.OtpManagerControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/otp/send/sms", otpManagerControllerInterface.SendOtpBySms, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/otp/verify", otpManagerControllerInterface.VerifyOtp, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func KecamatanRoute(e *echo.Echo, kecamatanControllerInterface controller.KecamatanControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/kecamatan", kecamatanControllerInterface.FindKecamatan, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func KelurahanRoute(e *echo.Echo, kelurahanControllerInterface controller.KelurahanControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/kelurahan", kelurahanControllerInterface.FindKelurahanByIdKeca, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func DesaRoute(e *echo.Echo, desaControllerInterface controller.DesaControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/desa", desaControllerInterface.FindDesaByIdKelu, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func UserRoute(e *echo.Echo, jwt config.Jwt, userControllerInterface controller.UserControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/user/non_surveyed", userControllerInterface.CreateUserNonSuveyed, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user", userControllerInterface.FindUserById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func ProductDesaRoute(e *echo.Echo, jwt config.Jwt, productDesaControllerInterface controller.ProductDesaControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/products", productDesaControllerInterface.FindProductsDesa, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/product", productDesaControllerInterface.FindProductDesaById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/category", productDesaControllerInterface.FindProductsDesaByCategory, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/sub_category", productDesaControllerInterface.FindProductsDesaBySubCategory, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/promo", productDesaControllerInterface.FindProductsDesaByPromo, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func PromoRoute(e *echo.Echo, jwt config.Jwt, promoDesaControllerInterface controller.PromoControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/promo", promoDesaControllerInterface.FindPromo, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func CartRoute(e *echo.Echo, jwt config.Jwt, cartControllerInterface controller.CartControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/cart/add", cartControllerInterface.CreateCart, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.PUT("/cart/update", cartControllerInterface.UpdateCart, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/cart/user", cartControllerInterface.FindCartByUser, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func PointRoute(e *echo.Echo, jwt config.Jwt, pointControllerInterface controller.PointControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/point", pointControllerInterface.FindPointByUser, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func OrderRoute(e *echo.Echo, jwt config.Jwt, orderControllerInterface controller.OrderControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/order/create", orderControllerInterface.CreateOrder, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/orders/user", orderControllerInterface.FindOrderByUser, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/order", orderControllerInterface.FindOrderById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.PUT("/order/cancel", orderControllerInterface.CancelOrderById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.PUT("/order/complete", orderControllerInterface.CancelOrderById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func PaymentChannelRoute(e *echo.Echo, jwt config.Jwt, paymentChannelControllerInterface controller.PaymentChannelControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/payment_channel", paymentChannelControllerInterface.FindPaymentChannel, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func SettingRoute(e *echo.Echo, jwt config.Jwt, settingControllerInterface controller.SettingControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/setting/shippingcost", settingControllerInterface.FindSettingShippingCost, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func UserShippingAddressRoute(e *echo.Echo, jwt config.Jwt, userShippingAddress controller.UserShippingAddressControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/shipping_address/create", userShippingAddress.CreateUserShippingAddress, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/shipping_address", userShippingAddress.FindUserShippingAddress, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/shipping_address/delete", userShippingAddress.DeleteUserShippingAddress, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func PpobRoute(e *echo.Echo, jwt config.Jwt, ppob controller.PpobControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/pricelist/pulsa", ppob.GetPulsaPriceList, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}
