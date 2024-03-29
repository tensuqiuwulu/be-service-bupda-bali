package controller

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type PaylaterControllerInterface interface {
	PayPaylater(c echo.Context) error
	DebetPerTransaksi(c echo.Context) error
	GetTagihanPaylater(c echo.Context) error
	GetRiwayatPaylaterPerBulan(c echo.Context) error
	GetOrderPaylaterByMonth(c echo.Context) error
	GetPembayaranTransaksiByIdUser(c echo.Context) error
	GetTabunganBimaMutation(c echo.Context) error
	GetTagihanPelunasan(c echo.Context) error
}

type PaylaterControllerImplementation struct {
	logger                   *logrus.Logger
	PaylaterServiceInterface service.PaylaterServiceInterface
	PaymentServiceInterface  service.PaymentServiceInterface
}

func NewPaylaterController(
	logger *logrus.Logger,
	paylaterServiceInterface service.PaylaterServiceInterface,
	paymentServiceInterface service.PaymentServiceInterface,
) PaylaterControllerInterface {
	return &PaylaterControllerImplementation{
		logger:                   logger,
		PaylaterServiceInterface: paylaterServiceInterface,
		PaymentServiceInterface:  paymentServiceInterface,
	}
}

func (controller *PaylaterControllerImplementation) DebetPerTransaksi(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	debetRequest := request.ReadFromDebetPerTransaksiRequestBody(c, requestId, controller.logger)
	controller.PaymentServiceInterface.DebetPerTransaksi(requestId, IdUser, debetRequest.LoanId)
	responses := response.Response{Code: 201, Mssg: "success", Data: "success pay paylater", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) GetTabunganBimaMutation(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	startDate := c.QueryParam("start_date")
	endDate := c.QueryParam("end_date")
	mutations := controller.PaymentServiceInterface.GetTabunganBimaMutation(requestId, IdUser, startDate, endDate)
	responses := response.Response{Code: 200, Mssg: "success", Data: mutations, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) GetPembayaranTransaksiByIdUser(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	indexDate := c.QueryParam("index_date")
	riwayatResponse := controller.PaylaterServiceInterface.GetPembayaranTransaksiByIdUser(requestId, IdUser, indexDate)
	responses := response.Response{Code: 200, Mssg: "success", Data: riwayatResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) GetOrderPaylaterByMonth(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	month, _ := strconv.Atoi(c.QueryParam("month"))
	riwayatResponse := controller.PaylaterServiceInterface.GetOrderPaylaterByMonth(requestId, IdUser, month)
	responses := response.Response{Code: 200, Mssg: "success", Data: riwayatResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) GetRiwayatPaylaterPerBulan(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	// month, _ := strconv.Atoi(c.QueryParam("month"))
	riwayatResponse := controller.PaylaterServiceInterface.GetOrderPaylaterPerBulan(requestId, IdUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: riwayatResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) GetTagihanPaylater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	tagihanResponse := controller.PaylaterServiceInterface.GetTagihanPaylater(requestId, IdUser)
	if condition := len(tagihanResponse) == 0; condition {
		exceptions.PanicIfRecordNotFound(errors.New("tagihan paylater not found"), requestId, []string{"tagihan paylater not found"}, controller.logger)
	}

	responses := response.Response{Code: 200, Mssg: "success", Data: tagihanResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) PayPaylater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	controller.PaymentServiceInterface.PayPaylaterBill(requestId, IdUser)
	responses := response.Response{Code: 201, Mssg: "success", Data: "success pay paylater", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *PaylaterControllerImplementation) GetTagihanPelunasan(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	IdUser := middleware.TokenClaimsIdUser(c)
	tagihanResponse := controller.PaymentServiceInterface.GetTagihanPelunasan(requestId, IdUser)
	if condition := len(tagihanResponse) == 0; condition {
		exceptions.PanicIfRecordNotFound(errors.New("tagihan pelunasan not found"), requestId, []string{"tagihan pelunasan not found"}, controller.logger)
	}

	responses := response.Response{Code: 200, Mssg: "success", Data: tagihanResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}
