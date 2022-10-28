package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	modelService "github.com/tensuqiuwulu/be-service-bupda-bali/model/service"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthServiceInterface interface {
	Login(requestId string, loginRequest *request.LoginRequest) (loginResponse interface{})
	FirstTimeLoginInveli(requestId string, loginInveliRequest *request.LoginInveliRequest) (loginResponse response.LoginInveliResponse)
	FirstTimeUbahPasswordInveli(requestId string, ubahPasswordInveliRequest *request.UpdateUserPasswordInveliRequest) error
	InveliUbahPin(requestId string, ubahPinRequest *request.LoginInveliRequest) error
	NewToken(requestId string, refreshToken string) (token string)
	GenerateToken(user modelService.User) (token string, err error)
	GenerateRefreshToken(user modelService.User) (token string, err error)
}

type AuthServiceImplementation struct {
	DB                            *gorm.DB
	ConfigJwt                     config.Jwt
	Validate                      *validator.Validate
	Logger                        *logrus.Logger
	UserRepositoryInterface       repository.UserRepositoryInterface
	InveliAPIRespositoryInterface invelirepository.InveliAPIRepositoryInterface
}

func NewAuthService(
	db *gorm.DB,
	configJwt config.Jwt,
	validate *validator.Validate,
	logger *logrus.Logger,
	userRepositoryInterface repository.UserRepositoryInterface,
	inveliAPIRespositoryInterface invelirepository.InveliAPIRepositoryInterface,
) AuthServiceInterface {
	return &AuthServiceImplementation{
		DB:                            db,
		ConfigJwt:                     configJwt,
		Validate:                      validate,
		Logger:                        logger,
		UserRepositoryInterface:       userRepositoryInterface,
		InveliAPIRespositoryInterface: inveliAPIRespositoryInterface,
	}
}

func (service *AuthServiceImplementation) Login(requestId string, loginRequest *request.LoginRequest) (loginResponse interface{}) {
	var userModelService modelService.User

	request.ValidateRequest(service.Validate, loginRequest, requestId, service.Logger)

	// jika username tidak ditemukan
	user, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, loginRequest.Phone)
	if user.Id == "" {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"not found"}, service.Logger)
	}

	if user.IsDelete == 1 {
		exceptions.PanicIfRecordNotFound(errors.New("user not found"), requestId, []string{"not found"}, service.Logger)
	}

	if user.IsActive == 1 {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
		exceptions.PanicIfBadRequest(err, requestId, []string{"Invalid Credentials"}, service.Logger)

		userModelService.Id = user.Id
		userModelService.IdDesa = user.IdDesa
		userModelService.AccountType = user.AccountType

		token, err := service.GenerateToken(userModelService)
		exceptions.PanicIfError(err, requestId, service.Logger)

		refreshToken, err := service.GenerateRefreshToken(userModelService)
		exceptions.PanicIfError(err, requestId, service.Logger)

		_, err = service.UserRepositoryInterface.SaveUserRefreshToken(service.DB, userModelService.Id, refreshToken)
		exceptions.PanicIfError(err, requestId, service.Logger)

		loginResponse = response.ToLoginResponse(token, refreshToken)

		return loginResponse
	} else {
		exceptions.PanicIfUnauthorized(errors.New("account is not active"), requestId, []string{"not active"}, service.Logger)
		return nil
	}
}

func (service *AuthServiceImplementation) FirstTimeLoginInveli(requestId string, loginInveliRequest *request.LoginInveliRequest) (loginResponse response.LoginInveliResponse) {

	request.ValidateRequest(service.Validate, loginInveliRequest, requestId, service.Logger)

	loginResult := service.InveliAPIRespositoryInterface.InveliLogin(loginInveliRequest.Phone, loginInveliRequest.Pin)

	fmt.Println("inveli login : ", loginResult)

	if len(loginResult.AccessToken) == 0 {
		exceptions.PanicIfBadRequest(errors.New("invalid credentials"), requestId, []string{"Invalid Credentials Inveli Login"}, service.Logger)
	}

	user := &entity.User{
		InveliAccessToken: loginResult.AccessToken,
		InveliIDMember:    loginResult.UserID,
	}

	userResult, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, loginInveliRequest.Phone)
	if len(userResult.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("invalid credentials"), requestId, []string{"User Not Found"}, service.Logger)
	}

	service.UserRepositoryInterface.SaveUserInveliToken(service.DB, userResult.Id, user)

	loginResponse = response.ToLoginInveliResponse(loginResult.AccessToken, loginResult.UserID)

	return loginResponse
}

func (service *AuthServiceImplementation) FirstTimeUbahPasswordInveli(requestId string, ubahPasswordInveliRequest *request.UpdateUserPasswordInveliRequest) error {

	request.ValidateRequest(service.Validate, ubahPasswordInveliRequest, requestId, service.Logger)

	userResult, _ := service.UserRepositoryInterface.FindUserByPhone(service.DB, ubahPasswordInveliRequest.Phone)
	if len(userResult.Id) == 0 {
		exceptions.PanicIfBadRequest(errors.New("invalid credentials"), requestId, []string{"User Not Found"}, service.Logger)
	}

	err := service.InveliAPIRespositoryInterface.InveliUbahPassword(userResult.InveliIDMember, ubahPasswordInveliRequest.NewPassword, userResult.InveliAccessToken)

	if err != nil {
		exceptions.PanicIfBadRequest(errors.New("invalid credentials"), requestId, []string{"cant change password"}, service.Logger)
	}

	return nil

}
func (service *AuthServiceImplementation) InveliUbahPin(requestId string, ubahPinRequest *request.LoginInveliRequest) error {

	request.ValidateRequest(service.Validate, ubahPinRequest, requestId, service.Logger)

	ubahPinResult := service.InveliAPIRespositoryInterface.InveliUbahPin(ubahPinRequest.Phone, ubahPinRequest.Pin)

	fmt.Println(ubahPinResult)
	return nil
}

func (service *AuthServiceImplementation) NewToken(requestId string, refreshToken string) (token string) {
	tokenParse, err := jwt.ParseWithClaims(refreshToken, &modelService.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(service.ConfigJwt.Key), nil
	})

	if !tokenParse.Valid {
		exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			exceptions.PanicIfUnauthorized(err, requestId, []string{"invalid token"}, service.Logger)
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			exceptions.PanicIfUnauthorized(err, requestId, []string{"expired token"}, service.Logger)
		} else {
			exceptions.PanicIfError(err, requestId, service.Logger)
		}
	}

	if claims, ok := tokenParse.Claims.(*modelService.TokenClaims); ok && tokenParse.Valid {
		user, err := service.UserRepositoryInterface.FindUserByIdAndRefreshToken(service.DB, claims.Id, refreshToken)
		exceptions.PanicIfRecordNotFound(err, requestId, []string{"User tidak ada"}, service.Logger)

		var userModelService modelService.User
		userModelService.Id = user.Id
		userModelService.IdDesa = user.IdDesa
		userModelService.AccountType = user.AccountType
		token, err := service.GenerateRefreshToken(userModelService)
		exceptions.PanicIfError(err, requestId, service.Logger)
		return token
	} else {
		err := errors.New("no claims")
		exceptions.PanicIfBadRequest(err, requestId, []string{"no claims"}, service.Logger)
		return ""
	}
}

func (service *AuthServiceImplementation) GenerateToken(user modelService.User) (token string, err error) {
	// Create the Claims
	claims := modelService.TokenClaims{
		Id:          user.Id,
		IdDesa:      user.IdDesa,
		AccountType: user.AccountType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(service.ConfigJwt.Tokenexpiredtime)).Unix(),
			Issuer:    "cyrilia",
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenWithClaims.SignedString([]byte(service.ConfigJwt.Key))
	if err != nil {
		return "", err
	}
	return token, err
}

func (service *AuthServiceImplementation) GenerateRefreshToken(user modelService.User) (token string, err error) {
	// Create the Claims
	claims := modelService.TokenClaims{
		Id:          user.Id,
		IdDesa:      user.IdDesa,
		AccountType: user.AccountType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().AddDate(0, 0, int(service.ConfigJwt.Refreshtokenexpiredtime)).Unix(),
			Issuer:    "cyrilia",
		},
	}

	tokenWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenWithClaims.SignedString([]byte(service.ConfigJwt.Key))
	if err != nil {
		return "", err
	}
	return token, err
}
