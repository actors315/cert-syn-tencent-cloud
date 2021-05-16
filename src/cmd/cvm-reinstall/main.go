package main

import (
	"flag"
	"fmt"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"qcloud-tools/src/cvmVo"
	"qcloud-tools/src/tools"
)

func main() {

	var group string
	var file string

	flag.StringVar(&group, "group", "tiyan", "分组")
	flag.StringVar(&file, "config", "", "配置文件地址")
	flag.Parse()

	config := tools.NewQcloudConfig(file)
	credential,cpf := tools.GetCredential(config.SecretId,config.SecretKey)
	cpf.HttpProfile.Endpoint = "cvm.tencentcloudapi.com"

	cvmItem := config.GetCvmItem(group)

	client, _ := cvm.NewClient(credential, cvmItem.Region, cpf)

	reinstall := new(cvmVo.ReinstallInfo)

	check := reinstall.CheckReinstall(client)

	if !check {
		fmt.Println("不重建")
		return
	}

	reinstall.Reinstall(client)

	cvmVo.ClearExpiredImage(client)
}
