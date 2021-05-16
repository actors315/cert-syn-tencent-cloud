package tools

import (
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
)

func GetCredential(secretId string, secretKey string) (*common.Credential, *profile.ClientProfile) {
	credential := common.NewCredential(
		secretId,
		secretKey,
	)

	cpf := profile.NewClientProfile()

	return credential, cpf
}
