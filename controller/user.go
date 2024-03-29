package controller

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/middleware"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/service"
)

type UserControllerInterface interface {
	CreateUserNonSuveyed(c echo.Context) error
	CreateUserSuveyed(c echo.Context) error
	FindUserById(c echo.Context) error
	DeleteUserById(c echo.Context) error
	UpdateUserPassword(c echo.Context) error
	UpdateUserForgotPassword(c echo.Context) error
	UpdateUserProfile(c echo.Context) error
	UpdateUserPhone(c echo.Context) error
	FindUserFromBigis(c echo.Context) error
	GetSimpananKhususBalance(c echo.Context) error
	GetUserAccountBimaByID(c echo.Context) error
	AktivasiAkunInveli(c echo.Context) error
	GetTunggakanPaylater(c echo.Context) error
	GetLimitPayLater(c echo.Context) error
	GetVATabBimaNasabah(c echo.Context) error
	// GetTagihanPaylater(c echo.Context) error
	GetNoRekening(c echo.Context) error
}

type UserControllerImplementation struct {
	Logger               *logrus.Logger
	UserServiceInterface service.UserServiceInterface
}

func NewUserController(
	logger *logrus.Logger,
	userServiceInterface service.UserServiceInterface,
) UserControllerInterface {
	return &UserControllerImplementation{
		UserServiceInterface: userServiceInterface,
	}
}

func (controller *UserControllerImplementation) GetNoRekening(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.GetNoRekeningNasabah(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: map[string]interface{}{
		"no_rekening": userResponse,
	}, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) GetVATabBimaNasabah(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.GetVANasabah(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: map[string]interface{}{
		"virtual_account": userResponse,
	}, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) GetLimitPayLater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.GetLimitPayLater(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: userResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) GetTunggakanPaylater(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.GetTunggakanPaylater(requestId, idUser)
	var responses response.Response
	if userResponse == nil {
		responses = response.Response{Code: 200, Mssg: "success", Data: "Tidak ada tunggakan", Error: []string{}}
	} else {
		responses = response.Response{Code: 200, Mssg: "success", Data: userResponse, Error: []string{}}
	}

	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) AktivasiAkunInveli(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	controller.UserServiceInterface.AktivasiAkunInveli(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: "", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) GetSimpananKhususBalance(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.GetSimpananKhususBalance(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: userResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) GetUserAccountBimaByID(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.GetUserAccountBimaByID(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: userResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) FindUserFromBigis(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromFindBigisResponsesRequestBody(c, requestId, controller.Logger)
	userResponse := controller.UserServiceInterface.FindUserFromBigis(requestId, request)
	responses := response.Response{Code: 200, Mssg: "success", Data: userResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) CreateUserSuveyed(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromCreateUserSurveyedRequestBody(c, requestId, controller.Logger)
	controller.UserServiceInterface.CreateUserSuveyed(requestId, request)
	responses := response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) CreateUserNonSuveyed(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	request := request.ReadFromCreateUserRequestBody(c, requestId, controller.Logger)
	controller.UserServiceInterface.CreateUserNonSuveyed(requestId, request)
	responses := response.Response{Code: 200, Mssg: "success", Data: nil, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) FindUserById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	userResponse := controller.UserServiceInterface.FindUserById(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: userResponse, Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) DeleteUserById(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	controller.UserServiceInterface.DeleteUserById(requestId, idUser)
	responses := response.Response{Code: 200, Mssg: "success", Data: "delete user success", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) UpdateUserPassword(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	updateUserPasswordRequest := request.ReadFromUpdateUserPasswordRequestBody(c, requestId, controller.Logger)
	controller.UserServiceInterface.UpdateUserPassword(requestId, idUser, updateUserPasswordRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: "update password success", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) UpdateUserForgotPassword(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	updateUserForgotPasswordRequest := request.ReadFromUpdateUserForgotPasswordRequestBody(c, requestId, controller.Logger)
	controller.UserServiceInterface.UpdateUserForgotPassword(requestId, updateUserForgotPasswordRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: "update password success", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) UpdateUserProfile(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	updateUserProfileRequest := request.ReadFromUpdateUserProfileRequestBody(c, requestId, controller.Logger)
	controller.UserServiceInterface.UpdateUserProfile(requestId, idUser, updateUserProfileRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: "update user profile success", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}

func (controller *UserControllerImplementation) UpdateUserPhone(c echo.Context) error {
	requestId := c.Response().Header().Get(echo.HeaderXRequestID)
	idUser := middleware.TokenClaimsIdUser(c)
	updateUserPhoneRequest := request.ReadFromUpdateUserPhoneRequestBody(c, requestId, controller.Logger)
	controller.UserServiceInterface.UpdateUserPhone(requestId, idUser, updateUserPhoneRequest)
	responses := response.Response{Code: 200, Mssg: "success", Data: "update user phone success", Error: []string{}}
	return c.JSON(http.StatusOK, responses)
}
