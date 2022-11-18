package service

import (
	"log"
	"time"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	invelirepository "github.com/tensuqiuwulu/be-service-bupda-bali/repository/inveli_repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gorm.io/gorm"
)

type PaymentChannelServiceInterface interface {
	FindPaymentChannel(requestId string, idUser string, requestPayChan *request.GetPaymentChannelRequest) (paymentChannelResponses []response.FindPaymentChannelResponse)
}

type PaymentChannelServiceImplementation struct {
	DB                                *gorm.DB
	Validate                          *validator.Validate
	Logger                            *logrus.Logger
	PaymentChannelRepositoryInterface repository.PaymentChannelRepositoryInterface
	InveliAPIRepositoryInterface      invelirepository.InveliAPIRepositoryInterface
	UserRepositoryInterface           repository.UserRepositoryInterface
	OrderRepositoryInterface          repository.OrderRepositoryInterface
}

func NewPaymentChannelService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	paymentChannelRepositoryInterface repository.PaymentChannelRepositoryInterface,
	inveliAPIRepositoryInterface invelirepository.InveliAPIRepositoryInterface,
	userRepositoryInterface repository.UserRepositoryInterface,
	orderRepositoryInterface repository.OrderRepositoryInterface,
) PaymentChannelServiceInterface {
	return &PaymentChannelServiceImplementation{
		DB:                                db,
		Validate:                          validate,
		Logger:                            logger,
		PaymentChannelRepositoryInterface: paymentChannelRepositoryInterface,
		InveliAPIRepositoryInterface:      inveliAPIRepositoryInterface,
		UserRepositoryInterface:           userRepositoryInterface,
		OrderRepositoryInterface:          orderRepositoryInterface,
	}
}

func (service *PaymentChannelServiceImplementation) FindPaymentChannel(requestId string, idUser string, requestPayChan *request.GetPaymentChannelRequest) (paymentChannelResponses []response.FindPaymentChannelResponse) {
	var err error
	paymentChannelResponse, _ := service.PaymentChannelRepositoryInterface.FindPaymentChannel(service.DB)
	if paymentChannelResponse == nil {
		exceptions.PanicIfRecordNotFound(err, requestId, []string{"Data paymanet channel not found"}, service.Logger)
	}

	user, _ := service.UserRepositoryInterface.FindUserById(service.DB, idUser)

	var jmlOrder float64

	var biayaTanggungRenteng float64
	jmlOrderPayLate, err := service.OrderRepositoryInterface.FindOrderPayLaterById(service.DB, idUser)
	if err != nil {
		log.Println(err.Error())
	}
	jmlOrder = 0
	for _, v := range jmlOrderPayLate {
		jmlOrder = jmlOrder + v.TotalBill
	}

	jmlOrder += requestPayChan.TotalBill

	userPaylaterFlag, _ := service.UserRepositoryInterface.GetUserPayLaterFlagThisMonth(service.DB, idUser)
	if len(userPaylaterFlag.Id) == 0 {
		err := service.UserRepositoryInterface.CreateUserPayLaterFlag(service.DB, &entity.UsersPaylaterFlag{
			Id:                  utilities.RandomUUID(),
			IdUser:              idUser,
			TanggungRentengFlag: 1,
			PaylaterDate:        time.Now(),
			CreatedAt:           time.Now(),
		})

		if err != nil {
			log.Println(err.Error())
		}

		biayaTanggungRenteng = 2500
	} else {
		if len(jmlOrderPayLate) == 0 {
			biayaTanggungRenteng = 2500
		} else {
			if int(jmlOrder) > (userPaylaterFlag.TanggungRentengFlag * 1000000) {
				biayaTanggungRenteng = 2500
			} else {
				biayaTanggungRenteng = 0
			}
		}

	}

	paymentChannelResponses = response.ToFindPaymentChannelResponse(paymentChannelResponse, user.User.StatusPaylater, biayaTanggungRenteng, user.User.IsPaylater)
	return paymentChannelResponses
}
