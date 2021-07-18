package certificate

import (
	"fmt"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ecdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ecdn/v20191012"
	"time"
)

type ISync interface {
	UpdateCredential() bool
	GetCredential() (*common.Credential, *profile.ClientProfile)
	GetCertRequestParam() (params string)
}

type Sync struct {
	SecretId       string
	SecretKey      string
	Domain         string
	PrivateKeyData string
	PublicKeyData  string
}

func (sync Sync) GetCredential() (*common.Credential, *profile.ClientProfile) {
	credential := common.NewCredential(
		sync.SecretId,
		sync.SecretKey,
	)

	cpf := profile.NewClientProfile()

	return credential, cpf
}

func (sync Sync) GetCertRequestParam() (params string) {

	params = "{\"Domain\":\"%s\",\"Https\":{\"Switch\":\"on\",\"Http2\":\"on\",\"CertInfo\":{\"Certificate\":\"%s\",\"PrivateKey\":\"%s\",\"Message\":\"%s\"}}}"
	params = fmt.Sprintf(params, sync.Domain, sync.PublicKeyData, sync.PrivateKeyData, time.Now().Format("2006-01-02"))

	return
}

type CdnSync struct {
	Sync
}

func (sync CdnSync) UpdateCredential() bool {
	credential, cpf := sync.GetCredential()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"

	client, _ := cdn.NewClient(credential, "", cpf)
	request := cdn.NewUpdateDomainConfigRequest()

	params := sync.GetCertRequestParam()

	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}

	request.ForceRedirect = &cdn.ForceRedirect{
		Switch:             common.StringPtr("on"),
		RedirectType:       common.StringPtr("https"),
		RedirectStatusCode: common.Int64Ptr(301),
	}

	response, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
	}
	if err != nil {
		fmt.Printf("UpdateDomainConfig returned: %s", err)
		return false
	}
	fmt.Printf("%s \n", response.ToJsonString())

	return true
}

type EcdnSync struct {
	Sync
}

func (sync EcdnSync) UpdateCredential() bool {
	credential, cpf := sync.GetCredential()
	cpf.HttpProfile.Endpoint = "ecdn.tencentcloudapi.com"

	client, _ := ecdn.NewClient(credential, "", cpf)
	request := ecdn.NewUpdateDomainConfigRequest()

	params := sync.GetCertRequestParam()

	err := request.FromJsonString(params)
	if err != nil {
		panic(err)
	}

	request.ForceRedirect = &ecdn.ForceRedirect{
		Switch:             common.StringPtr("on"),
		RedirectType:       common.StringPtr("https"),
		RedirectStatusCode: common.Int64Ptr(301),
	}

	response, err := client.UpdateDomainConfig(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Println("An API error has returned: ", err)
	}
	if err != nil {
		fmt.Println("An API error has returned: ", err)
		return false
	}
	fmt.Println(response.ToJsonString())

	return true
}