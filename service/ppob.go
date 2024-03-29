package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/go-playground/validator"
	"github.com/sirupsen/logrus"
	"github.com/tensuqiuwulu/be-service-bupda-bali/config"
	"github.com/tensuqiuwulu/be-service-bupda-bali/exceptions"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/ppob"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/request"
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/response"
	"github.com/tensuqiuwulu/be-service-bupda-bali/repository"
	"github.com/tensuqiuwulu/be-service-bupda-bali/utilities"
	"gorm.io/gorm"
)

type PpobServiceInterface interface {
	GetPrepaidPulsaPriceList(requestId string, numberPhone string) (priceListResponse []response.GetPrepaidPriceListResponse)
	GetPrepaidDataPriceList(requestId string, numberPhone string) (priceListResponse []response.GetPrepaidPriceListResponse)
	GetPrepaidPlnPriceList(requestId string, idPelanggan string) (priceListResponse []response.GetPrepaidPriceListResponse)
	GetPostpaidPdamProduct(requestId string) (postPaidPadmProductResponse []response.GetPostpaidPdamProductResponse)
	GetPostpaidTelcoProduct(requestId string) (postPaidTelcoProductResponse []response.GetPostpaidTelcoProductResponse)
	InquiryPrepaidPln(requestId string, inquiryPrepaidPlnRequest *request.InquiryPrepaidPlnRequest) (inquiryPrepaidPlnResponse response.InquiryPrepaidPlnResponse)
	InquiryPostpaidPln(requestId string, inquiryPostpaidPlnRequest *request.InquiryPostpaidPlnRequest) (inquiryPostpadPlnResponse response.InquiryPostpaidPlnResponse)
	InquiryPostpaidPdam(requestId string, inquiryPostpaidPdamRequest *request.InquiryPostpaidPdamRequest) (inquiryPostpaidPdamResponse response.InquiryPostpaidPdamResponse)
	InquiryPostpaidTelco(requestId string, inquiryPostpaidTelcoRequest *request.InquiryPostpaidTelcoRequest) (inquiryPostpaidTelcoResponse response.InquiryPostpaidTelcoResponse)
	PrepaidCheckStatusTransaction(requestId, NumberOrder string)
}

type PpobServiceImplementation struct {
	DB                                *gorm.DB
	Validate                          *validator.Validate
	Logger                            *logrus.Logger
	OperatorPrefixRepositoryInterface repository.OperatorPrefixRepositoryInterface
	OrderServiceInterface             OrderServiceInterface
}

func NewPpobService(
	db *gorm.DB,
	validate *validator.Validate,
	logger *logrus.Logger,
	operatorPrefixRepositoryInterface repository.OperatorPrefixRepositoryInterface,
	orderServiceInterface OrderServiceInterface,
) PpobServiceInterface {
	return &PpobServiceImplementation{
		DB:                                db,
		Validate:                          validate,
		Logger:                            logger,
		OperatorPrefixRepositoryInterface: operatorPrefixRepositoryInterface,
		OrderServiceInterface:             orderServiceInterface,
	}
}

func PrefixNumber(phone string) string {
	split := strings.Split(phone, "")
	phoneJoin := split[0] + split[1] + split[2] + split[3]
	return phoneJoin
}

func (service *PpobServiceImplementation) GetPrepaidPlnPriceList(requestId string, idPelanggan string) (priceListResponses []response.GetPrepaidPriceListResponse) {
	prepaidPlnPriceList := service.GetPrepaidPriceList(requestId, idPelanggan, "pln", "pln")
	priceListResponses = response.ToGetPrepaidPriceListResponse(prepaidPlnPriceList.Data.Data)
	return priceListResponses
}

func (service *PpobServiceImplementation) GetPrepaidPulsaPriceList(requestId string, numberPhone string) (priceListResponses []response.GetPrepaidPriceListResponse) {
	phone := PrefixNumber(numberPhone)

	opereratorPrefixResult, err := service.OperatorPrefixRepositoryInterface.FindOperatorPrefixByPhone(service.DB, phone)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(opereratorPrefixResult.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("operator tidak ditemukan"), requestId, []string{"operator not found"}, service.Logger)
	}

	prepaidPulsaPriceList := service.GetPrepaidPriceList(requestId, numberPhone, "pulsa", opereratorPrefixResult.KodeOperator)
	priceListResponses = response.ToGetPrepaidPriceListResponse(prepaidPulsaPriceList.Data.Data)
	return priceListResponses
}

func (service *PpobServiceImplementation) GetPrepaidDataPriceList(requestId string, numberPhone string) (priceListResponses []response.GetPrepaidPriceListResponse) {
	phone := PrefixNumber(numberPhone)

	opereratorPrefixResult, err := service.OperatorPrefixRepositoryInterface.FindOperatorPrefixByPhone(service.DB, phone)
	exceptions.PanicIfError(err, requestId, service.Logger)
	if len(opereratorPrefixResult.Id) == 0 {
		exceptions.PanicIfRecordNotFound(errors.New("operator tidak ditemukan"), requestId, []string{"operator not found"}, service.Logger)
	}

	prepaidDataPriceList := service.GetPrepaidPriceList(requestId, numberPhone, "data", opereratorPrefixResult.KodeOperator)
	priceListResponses = response.ToGetPrepaidDataPriceListResponse(prepaidDataPriceList.Data.Data)
	return priceListResponses
}

func (service *PpobServiceImplementation) GetPrepaidPriceList(requestId string, id, tipe, operator string) *ppob.PrepaidPriceListResponse {

	// Create Request
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "pl"))
	body, _ := json.Marshal(map[string]interface{}{
		"status":   "all",
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PrepaidHost + "/pricelist" + "/" + tipe + "/" + operator
	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	defer resp.Body.Close()

	prepaidPriceList := &ppob.PrepaidPriceListResponse{}
	// fmt.Printf("body: %s\n", prepaidPriceList)

	if err = json.Unmarshal([]byte(data), &prepaidPriceList); err != nil {
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	if prepaidPriceList.Data.Rc != "00" {
		// fmt.Printf("body: %s\n", prepaidPriceList.Data)
		exceptions.PanicIfError(errors.New("error from IAK"), requestId, service.Logger)
	}

	return prepaidPriceList
}

func (service *PpobServiceImplementation) InquiryPrepaidPln(requestId string, inquiryPrepaidPlnRequest *request.InquiryPrepaidPlnRequest) (inquiryPrepaidPlnResponse response.InquiryPrepaidPlnResponse) {
	var err error

	request.ValidateRequest(service.Validate, inquiryPrepaidPlnRequest, requestId, service.Logger)

	// Create Request
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + inquiryPrepaidPlnRequest.CustomerId))
	body, _ := json.Marshal(map[string]interface{}{
		"username":    config.GetConfig().Ppob.Username,
		"customer_id": inquiryPrepaidPlnRequest.CustomerId,
		"sign":        hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PrepaidHost + "/inquiry-pln"

	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	inquiryPrepaidPln := &ppob.InquiryPrepaidPln{}

	if err = json.Unmarshal([]byte(data), &inquiryPrepaidPln); err != nil {
		inquiryPrepaidPlnErrorHandle := &ppob.InquiryPrepaidPlnErrorHandle{}
		if err = json.Unmarshal([]byte(data), &inquiryPrepaidPlnErrorHandle); err != nil {
			exceptions.PanicIfError(err, requestId, service.Logger)
		} else {
			if inquiryPrepaidPlnErrorHandle.Data.Rc == "208" {
				exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
			}
			if inquiryPrepaidPlnErrorHandle.Data.Rc == "14" {
				exceptions.PanicIfBadRequest(errors.New("costumer id not found"), requestId, []string{"Costumer Id Not Found"}, service.Logger)
			}
			exceptions.PanicIfError(err, requestId, service.Logger)
		}
	}

	if inquiryPrepaidPln.Data.Rc != "00" {
		fmt.Printf("body: %s\n", inquiryPrepaidPln.Data)
		exceptions.PanicIfError(errors.New("error from IAK"), requestId, service.Logger)
	}

	inquiryPrepaidPlnResponse = response.ToInquiryPrepaidPlnResponse(inquiryPrepaidPln)

	return inquiryPrepaidPlnResponse
}

func (service *PpobServiceImplementation) InquiryPostpaidPln(requestId string, inquiryPostpaidPlnRequest *request.InquiryPostpaidPlnRequest) (inquiryPostpaidPlnResponse response.InquiryPostpaidPlnResponse) {
	var err error

	request.ValidateRequest(service.Validate, inquiryPostpaidPlnRequest, requestId, service.Logger)

	// generate number order yg akan digunakan sebagai ref id
	refId := utilities.GenerateRefId()

	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + refId))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "inq-pasca",
		"username": config.GetConfig().Ppob.Username,
		"code":     "PLNPOSTPAID",
		"hp":       inquiryPostpaidPlnRequest.CustomerId,
		"ref_id":   refId,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl

	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	inquiryPostpaidPln := &ppob.InquiryPostpaidPln{}
	postpaidErrorResponse := &ppob.PostpaidErrorResponse{}

	if err = json.Unmarshal([]byte(data), &inquiryPostpaidPln); err != nil {
		if err = json.Unmarshal([]byte(data), &postpaidErrorResponse); err != nil {
			exceptions.PanicIfRecordNotFound(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
		}
	}

	if inquiryPostpaidPln.Data.ResponseCode == "01" {
		exceptions.PanicPPOBHandler(errors.New("error 001"), requestId, inquiryPostpaidPln.Data.Message, []string{"001"}, service.Logger)
	}

	if inquiryPostpaidPln.Data.ResponseCode == "14" {
		exceptions.PanicPPOBHandler(errors.New("error 002"), requestId, inquiryPostpaidPln.Data.Message, []string{"002"}, service.Logger)
	}

	inquiryPostpaidPlnResponse = response.ToInquiryPostpaidPlnResponse(inquiryPostpaidPln, inquiryPostpaidPln.Data.Desc.Tagihan.Detail, refId)

	return inquiryPostpaidPlnResponse
}

func (service *PpobServiceImplementation) GetPostpaidPdamProduct(requestId string) (postpaidPdamProductResponse []response.GetPostpaidPdamProductResponse) {
	var err error
	// Create Request
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "pl"))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "pricelist-pasca",
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
		"status":   "all",
		"province": "bali",
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl + "/pdam"

	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	postpaidPriceList := &ppob.PostpaidPriceListResponse{}

	if err = json.Unmarshal([]byte(data), &postpaidPriceList); err != nil {
		exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
	}

	postpaidPdamProductResponse = response.ToGetPostpaidPadmProductResponse(postpaidPriceList.Data.Pasca)

	return postpaidPdamProductResponse
}

func (service *PpobServiceImplementation) GetPostpaidTelcoProduct(requestId string) (postpaidTelcoProductResponse []response.GetPostpaidTelcoProductResponse) {
	var err error
	// Create Request
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + "pl"))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "pricelist-pasca",
		"username": config.GetConfig().Ppob.Username,
		"sign":     hex.EncodeToString(sign[:]),
		"status":   "all",
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl + "/hp"

	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	// fmt.Printf("body: %s\n", data)

	postpaidTelcoPriceList := &ppob.InquiryProductPostpaidPPOB{}

	if err = json.Unmarshal([]byte(data), &postpaidTelcoPriceList); err != nil {
		exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
	}

	postpaidTelcoProductResponse = response.ToGetPostpaidTelcoProductResponse(postpaidTelcoPriceList)

	return postpaidTelcoProductResponse
}

func (service *PpobServiceImplementation) InquiryPostpaidPdam(requestId string, inquiryPostpaidPdamRequest *request.InquiryPostpaidPdamRequest) (inquiryPostpaidPdamResponse response.InquiryPostpaidPdamResponse) {

	var err error

	request.ValidateRequest(service.Validate, inquiryPostpaidPdamRequest, requestId, service.Logger)

	refId := utilities.GenerateRefId()
	// Create Request
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + refId))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "inq-pasca",
		"username": config.GetConfig().Ppob.Username,
		"code":     inquiryPostpaidPdamRequest.Code,
		"hp":       inquiryPostpaidPdamRequest.Hp,
		"ref_id":   refId,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl

	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	inquiryPostpaidPdam := &ppob.InquiryPostpaidPdam{}
	postpaidErrorResponse := &ppob.PostpaidErrorResponse{}
	if err = json.Unmarshal([]byte(data), &inquiryPostpaidPdam); err != nil {
		if err = json.Unmarshal([]byte(data), &postpaidErrorResponse); err != nil {
			exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
		}
	}

	if inquiryPostpaidPdam.Data.ResponseCode == "01" {
		exceptions.PanicPPOBHandler(errors.New("error 001"), requestId, inquiryPostpaidPdam.Data.Message, []string{"001"}, service.Logger)
	}

	if inquiryPostpaidPdam.Data.ResponseCode == "14" {
		exceptions.PanicPPOBHandler(errors.New("error 002"), requestId, inquiryPostpaidPdam.Data.Message, []string{"002"}, service.Logger)
	}

	inquiryPostpaidPdamResponse = response.ToInquiryPostpaidPdamResponse(inquiryPostpaidPdam, inquiryPostpaidPdam.Data.Desc.Bill.Detail, refId)

	return inquiryPostpaidPdamResponse
}

func (service *PpobServiceImplementation) InquiryPostpaidTelco(requestId string, inquiryPostpaidTelcoRequest *request.InquiryPostpaidTelcoRequest) (inquiryPostpaidPdamResponse response.InquiryPostpaidTelcoResponse) {

	var err error

	request.ValidateRequest(service.Validate, inquiryPostpaidTelcoRequest, requestId, service.Logger)

	refId := utilities.GenerateRefId()
	// Create Request
	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + refId))
	body, _ := json.Marshal(map[string]interface{}{
		"commands": "inq-pasca",
		"username": config.GetConfig().Ppob.Username,
		"code":     inquiryPostpaidTelcoRequest.Code,
		"hp":       inquiryPostpaidTelcoRequest.Hp,
		"ref_id":   refId,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PostpaidUrl

	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)

	inquiryPostpaidTelco := &ppob.InquiryPostpaidTelco{}
	postpaidErrorResponse := &ppob.PostpaidErrorResponse{}
	if err = json.Unmarshal([]byte(data), &inquiryPostpaidTelco); err != nil {
		if err = json.Unmarshal([]byte(data), &postpaidErrorResponse); err != nil {
			exceptions.PanicIfBadRequest(errors.New("INVALID DATA"), requestId, []string{"INVALID DATA"}, service.Logger)
		}
	}

	if inquiryPostpaidTelco.Data.ResponseCode == "01" {
		exceptions.PanicPPOBHandler(errors.New("error 001"), requestId, inquiryPostpaidTelco.Data.Message, []string{"001"}, service.Logger)
	}

	if inquiryPostpaidTelco.Data.ResponseCode == "14" {
		exceptions.PanicPPOBHandler(errors.New("error 002"), requestId, inquiryPostpaidTelco.Data.Message, []string{"002"}, service.Logger)
	}

	inquiryPostpaidPdamResponse = response.ToInquiryPostpaidTelcoResponse(inquiryPostpaidTelco, inquiryPostpaidTelco.Data.Desc.Tagihan.Tagihan, refId)

	return inquiryPostpaidPdamResponse
}

func (service *PpobServiceImplementation) PrepaidCheckStatusTransaction(requestId, refId string) {
	var err error

	sign := md5.Sum([]byte(config.GetConfig().Ppob.Username + config.GetConfig().Ppob.PpobKey + refId))
	body, _ := json.Marshal(map[string]interface{}{
		"username": config.GetConfig().Ppob.Username,
		"ref_id":   refId,
		"sign":     hex.EncodeToString(sign[:]),
	})

	reqBody := io.NopCloser(strings.NewReader(string(body)))

	urlString := config.GetConfig().Ppob.PrepaidHost + "/check-status"

	// URL
	url, _ := url.Parse(urlString)

	req := &http.Request{
		Method: "POST",
		URL:    url,
		Header: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Body: reqBody,
	}

	reqDump, _ := httputil.DumpRequestOut(req, true)
	fmt.Printf("REQUEST:\n%s", string(reqDump))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("An Error Occured %v", err)
		exceptions.PanicIfError(err, requestId, service.Logger)
	}

	defer resp.Body.Close()

	// Read response body
	data, _ := io.ReadAll(resp.Body)
	fmt.Printf("body: %s\n", data)
}
