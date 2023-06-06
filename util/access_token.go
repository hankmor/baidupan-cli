package util

import openapi "baidupan-cli/openxpanapi"

func MockAccessToken() *openapi.OauthTokenDeviceTokenResponse {
	token := "126.cf86d7e7d994ce1dedc27301aade8472.YmeTyvNpEf1sQlTESsrLT235EFkyhTB26ZNukrS.u-pTdA"
	return &openapi.OauthTokenDeviceTokenResponse{
		AccessToken: &token,
	}
}
