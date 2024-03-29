package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type SettingControllerInterface interface {
	FindSettingShippingCost(c echo.Context) error
	FindNewVersion(c echo.Context) error
}

type SettingControllerImplementation struct {
	SettingServiceInterface service.SettingServiceInterface
}

func NewSettingController(
	settingServiceInterface service.SettingServiceInterface) SettingControllerInterface {
	return &SettingControllerImplementation{
		SettingServiceInterface: settingServiceInterface,
	}
}

func (controller *SettingControllerImplementation) FindNewVersion(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	OS, _ := strconv.Atoi(c.QueryParam("os"))
	settingResponse := controller.SettingServiceInterface.FindNewVersion(requestId, OS)
	responses := response.Response{Code: 200, Mssg: "success", Data: settingResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *SettingControllerImplementation) FindSettingShippingCost(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idDesa := middleware.TokenClaimsIdDesa(c)
	settingResponse := controller.SettingServiceInterface.FindSettingShippingCost(requestId, idDesa)
	responses := response.Response{Code: 200, Mssg: "success", Data: settingResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}
