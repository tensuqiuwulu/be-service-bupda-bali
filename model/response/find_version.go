package response

import "github.com/tensuqiuwulu/be-service-bupda-bali/model/entity"

type FindVersionResponse struct {
	OSName  string `json:"os"`
	Current string `json:"current"`
	New     string `json:"new"`
}

func ToFindNewVersionResponse(setting []entity.Setting, os int) (settingResponse FindVersionResponse) {
	if os == 1 {
		settingResponse.OSName = "Android"
	} else {
		settingResponse.OSName = "iOS"
	}
	settingResponse.Current = setting[0].SettingName
	settingResponse.New = setting[1].SettingName
	return settingResponse
}