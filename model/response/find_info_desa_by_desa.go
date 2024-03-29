package response

import (
	"github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"
)

type FindInfoDesaByIdDesaResponse struct {
	Id          string `json:"id"`
	Attachments string `json:"attachments"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Url         string `json:"url"`
	CreatedDate string `json:"date"`
}

func ToFindInfoDesaByIdDesaResponse(infoDesas []entity.InfoDesa) (infoDesasResponses []FindInfoDesaByIdDesaResponse) {
	for _, infoDesa := range infoDesas {
		var infoDesaResponse FindInfoDesaByIdDesaResponse
		infoDesaResponse.Id = infoDesa.Id
		infoDesaResponse.Attachments = infoDesa.Attachments
		infoDesaResponse.Title = infoDesa.Title
		infoDesaResponse.Content = infoDesa.Content
		infoDesaResponse.Url = infoDesa.Url
		infoDesaResponse.CreatedDate = infoDesa.CreatedDate.Time.Format("2006-01-02 15:04:05")
		infoDesasResponses = append(infoDesasResponses, infoDesaResponse)
	}
	return infoDesasResponses
}
