package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindProductsDesaResponse struct {
	Id              string  `json:"id"`
	IdBrand         int     `json:"id_brand"`
	IdCategory      int     `json:"id_category"`
	IdSubCategory   int     `json:"id_sub_category"`
	IdType          int     `json:"id_type"`
	IdUnit          int     `json:"id_unit"`
	NoSku           string  `json:"no_sku"`
	ProductName     string  `json:"product_name"`
	Price           float64 `json:"price_normal"`
	PricePromo      float64 `json:"price_promo"`
	PromoPercentage float64 `json:"promo_percentage"`
	FlagPromo       int     `json:"flag_promo"`
	Description     string  `json:"description"`
	PictureUrl      string  `json:"picture_url"`
	Thumbnail       string  `json:"thumbnail"`
	StockOpname     int     `json:"stock_opname"`
	PriceInfo       string  `json:"price_info"`
	AccountType     string  `json:"account_type"`
}

func ToFindProductsDesaResponse(productsDesas []entity.ProductsDesa, AccountType int) (productsDesaResponses []FindProductsDesaResponse) {
	for _, productDesa := range productsDesas {
		productsDesaResponse := FindProductsDesaResponse{}
		productsDesaResponse.Id = productDesa.Id
		productsDesaResponse.IdBrand = productDesa.ProductsMaster.IdBrand
		productsDesaResponse.IdType = productDesa.IdType
		productsDesaResponse.IdCategory = productDesa.ProductsMaster.IdCategory
		productsDesaResponse.IdSubCategory = productDesa.ProductsMaster.IdSubCategory
		productsDesaResponse.IdUnit = productDesa.ProductsMaster.IdUnit
		productsDesaResponse.NoSku = productDesa.ProductsMaster.NoSku
		productsDesaResponse.ProductName = productDesa.ProductsMaster.ProductName
		productsDesaResponse.FlagPromo = productDesa.IsPromo
		if AccountType == 1 {
			if productDesa.IsPromo == 1 {
				productsDesaResponse.Price = productDesa.Price
				productsDesaResponse.PricePromo = productDesa.PricePromo
				productsDesaResponse.PromoPercentage = productDesa.PercentagePromo
			} else {
				productsDesaResponse.Price = productDesa.Price
				productsDesaResponse.PricePromo = 0
				productsDesaResponse.PromoPercentage = 0
			}
			productsDesaResponse.AccountType = "User Biasa"
			productsDesaResponse.PriceInfo = "Krama Harga Normal"
		} else if AccountType == 2 {
			productsDesaResponse.AccountType = "User Merchant"
			productsDesaResponse.PriceInfo = "Krama Harga Grosir"
			productsDesaResponse.Price = productDesa.PriceGrosir
		}
		productsDesaResponse.Description = productDesa.ProductsMaster.Description
		productsDesaResponse.PictureUrl = productDesa.ProductsMaster.PictureUrl
		productsDesaResponse.Thumbnail = productDesa.ProductsMaster.Thumbnail
		productsDesaResponse.Thumbnail = productDesa.ProductsMaster.Thumbnail
		productsDesaResponse.StockOpname = productDesa.StockOpname
		productsDesaResponses = append(productsDesaResponses, productsDesaResponse)
	}
	return productsDesaResponses
}
