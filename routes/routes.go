package routes

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/controller"
	authMiddlerware "github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
)

func AuthRoute(e *echo.Echo, authControllerInterface controller.AuthControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/auth/login", authControllerInterface.Login, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/auth/ubah-password/inveli", authControllerInterface.FirstTimeUbahPasswordInveli, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
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

func ListPinjamanRoute(e *echo.Echo, jwt config.Jwt, listPinjamanControllerInterface controller.ListPinjamanControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/paylater/list-pinjaman", listPinjamanControllerInterface.FindListPinjamanByUser, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/paylater/list-pinjaman-detail", listPinjamanControllerInterface.FindListPinjamanById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func InfoDesaRoute(e *echo.Echo, jwt config.Jwt, infoDesaControllerInterface controller.InfoDesaControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/infodesa", infoDesaControllerInterface.FindInfoDesaByIdDesa, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func UserRoute(e *echo.Echo, jwt config.Jwt, userControllerInterface controller.UserControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/user/non_surveyed", userControllerInterface.CreateUserNonSuveyed, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/user/surveyed", userControllerInterface.CreateUserSuveyed, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/user/bigis", userControllerInterface.FindUserFromBigis, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user", userControllerInterface.FindUserById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.DELETE("/user/delete", userControllerInterface.DeleteUserById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.PUT("/user/update/password", userControllerInterface.UpdateUserPassword, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.PUT("/user/update/forgotpassword", userControllerInterface.UpdateUserForgotPassword, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.PUT("/user/update/profile", userControllerInterface.UpdateUserProfile, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.PUT("/user/update/phone", userControllerInterface.UpdateUserPhone, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user/get-balance-paylater", userControllerInterface.GetSimpananKhususBalance, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user/get-balance-bima", userControllerInterface.GetUserAccountBimaByID, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/user/aktivasi-inveli", userControllerInterface.AktivasiAkunInveli, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user/get-tunggakan-paylater", userControllerInterface.GetTunggakanPaylater, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user/get-limit-paylater", userControllerInterface.GetLimitPayLater, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user/get-va-tab-bima", userControllerInterface.GetVATabBimaNasabah, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/user/no-rekening", userControllerInterface.GetNoRekening, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func InveliTestRoutes(e *echo.Echo, inveliTestControllerInterface controller.InveliTestingController) {
	group := e.Group("api/v1")
	group.POST("/inveli/test", inveliTestControllerInterface.GetAccountInfo, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/inveli/test-akun-info", inveliTestControllerInterface.GetStatusAkun, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/inveli/balance", inveliTestControllerInterface.GetBalanceAccount, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/inveli/get-riwayat-pinjaman", inveliTestControllerInterface.GetRiwayatPinjaman, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/inveli/get-saldo-bupda", inveliTestControllerInterface.GetSaldoBupda, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func ProductDesaRoute(e *echo.Echo, jwt config.Jwt, productDesaControllerInterface controller.ProductDesaControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/products", productDesaControllerInterface.FindProductsDesa, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/product", productDesaControllerInterface.FindProductDesaById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/category", productDesaControllerInterface.FindProductsDesaByCategory, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/sub_category", productDesaControllerInterface.FindProductsDesaBySubCategory, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/promo", productDesaControllerInterface.FindProductsDesaByPromo, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/notoken", productDesaControllerInterface.FindProductsDesaNotoken, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/category/notoken", productDesaControllerInterface.FindProductsDesaByCategoryNotoken, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/products/sub_category/notoken", productDesaControllerInterface.FindProductsDesaBySubCategoryNotoken, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
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

func MerchantRoute(e *echo.Echo, jwt config.Jwt, merchantControllerInterface controller.MerchantControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/merchant/request_approve", merchantControllerInterface.CreateMerchantApproveList, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/merchant/status_approve", merchantControllerInterface.FindMerchantStatusApproveByUserResponse, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
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
	group.PUT("/order/complete", orderControllerInterface.CompleteOrderById, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/order/update/payment", orderControllerInterface.UpdateOrderPaymentStatus, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/order/callback/ppob", orderControllerInterface.CallbackPpobTransaction, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/order/paylater", orderControllerInterface.FindOrderPayLater, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func PaymentChannelRoute(e *echo.Echo, jwt config.Jwt, paymentChannelControllerInterface controller.PaymentChannelControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/payment_channel", paymentChannelControllerInterface.FindPaymentChannel, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func SettingRoute(e *echo.Echo, jwt config.Jwt, settingControllerInterface controller.SettingControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/setting/shippingcost", settingControllerInterface.FindSettingShippingCost, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/version", settingControllerInterface.FindNewVersion, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func BannerRoute(e *echo.Echo, jwt config.Jwt, bannerControllerInterface controller.BannerControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/banner", bannerControllerInterface.FindBannerByDesa, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/banner/no_token", bannerControllerInterface.FindBannerAll, authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func UserShippingAddressRoute(e *echo.Echo, jwt config.Jwt, userShippingAddress controller.UserShippingAddressControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/shipping_address/create", userShippingAddress.CreateUserShippingAddress, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/shipping_address", userShippingAddress.FindUserShippingAddress, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/shipping_address/delete", userShippingAddress.DeleteUserShippingAddress, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func PpobRoute(e *echo.Echo, jwt config.Jwt, ppob controller.PpobControllerInterface) {
	group := e.Group("api/v1")
	group.GET("/prepaid/pricelist/pulsa", ppob.GetPrepaidPulsaPriceList, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/prepaid/pricelist/data", ppob.GetPrepaidDataPriceList, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/prepaid/pricelist/pln", ppob.GetPrepaidPlnPriceList, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/prepaid/inquiry/pln", ppob.InquiryPrepaidPln, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/postpaid/inquiry/pln", ppob.InquiryPostpaidPln, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/postpaid/list/pdam", ppob.GetPostpaidPdamProduct, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/postpaid/list/telco", ppob.GetPostpaidTelcoProduct, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/postpaid/inquiry/pdam", ppob.InquiryPostpaidPdam, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/postpaid/inquiry/telco", ppob.InquiryPostpaidTelco, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

func PaylaterRoute(e *echo.Echo, jwt config.Jwt, paylaterControllerInterface controller.PaylaterControllerInterface) {
	group := e.Group("api/v1")
	group.POST("/paylater/pay", paylaterControllerInterface.PayPaylater, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/paylater/get-tagihan", paylaterControllerInterface.GetTagihanPaylater, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/order/paylater-monthly", paylaterControllerInterface.GetRiwayatPaylaterPerBulan, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/order/paylater-month", paylaterControllerInterface.GetOrderPaylaterByMonth, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/paylater/paid-transaction", paylaterControllerInterface.GetPembayaranTransaksiByIdUser, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/bima/mutation", paylaterControllerInterface.GetTabunganBimaMutation, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.GET("/paylater/tagihan-pelunasan", paylaterControllerInterface.GetTagihanPelunasan, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
	group.POST("/paylater/debet-per-transaksi", paylaterControllerInterface.DebetPerTransaksi, authMiddlerware.Authentication(jwt), authMiddlerware.RateLimit(), authMiddlerware.Timeout())
}

// Main Route
func MainRoute(e *echo.Echo) {
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}})
	})
}
