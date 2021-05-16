package cvmVo

import (
	"fmt"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ecdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ecdn/v20191012"
	"qcloud-tools/src/tools"
)

func UpdateCdnCredential(credential *common.Credential, cpf *profile.ClientProfile, params string) {

	client, _ := cdn.NewClient(credential, "", cpf)
	request := cdn.NewUpdateDomainConfigRequest()

	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	fmt.Printf("%s \n", response.ToJsonString())
}

func UpdateEcdnCredential(credential *common.Credential, cpf *profile.ClientProfile, params string) {

	client, _ := ecdn.NewClient(credential, "", cpf)
	request := ecdn.NewUpdateDomainConfigRequest()

	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}
	response, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	fmt.Printf("%s \n", response.ToJsonString())
}

func UpdateCredential(file string, group string) {
	config := tools.NewQcloudConfig(file)

	certItem, found := config.Certificates[group]
	if !found {
		panic("配置不存在")
	}

	credential, cpf := tools.GetCredential(config.SecretId, config.SecretKey)

	params := config.GetCertRequestParam(group)

	switch certItem.Product {
	case "ecdn":
		cpf.HttpProfile.Endpoint = "ecdn.tencentcloudapi.com"
		UpdateEcdnCredential(credential, cpf, params)
	default:
		cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
		UpdateCdnCredential(credential, cpf, params)
	}

}
