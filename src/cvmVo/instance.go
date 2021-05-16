package cvmVo

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"time"
)

type Item struct {
	InstanceId string
}

type ReinstallInfo struct {
	sourceItem *Item
	targetItem *Item
}

func (reinstall *ReinstallInfo) CheckReinstall(client *cvm.Client) bool {
	request := cvm.NewDescribeInstancesRequest()
	response, err := client.DescribeInstances(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		msg := fmt.Sprintf("An API error has returned: %s", err)
		panic(msg)
		return false
	}
	if err != nil {
		return false
	}

	if *response.Response.TotalCount < 2 {
		return false
	}

	for _, instance := range response.Response.InstanceSet {
		cvmItem := new(Item)
		if "EXPIRED" == *instance.RestrictState {
			cvmItem.InstanceId = *instance.InstanceId
			reinstall.sourceItem = cvmItem
		} else if "NORMAL" == *instance.RestrictState && "RUNNING" == *instance.InstanceState {
			if "CentOS 8.0 64位" != *instance.OsName {
				cvmItem.InstanceId = *instance.InstanceId
				reinstall.targetItem = cvmItem
			}
		}
	}

	if reinstall.sourceItem == nil || reinstall.targetItem == nil {
		return false
	}

	return true
}

func (reinstall *ReinstallInfo) Reinstall(client *cvm.Client) {

	imageName := "bak-" + time.Now().Format("2006-01-02")
	image := GetImageInfo(client,"")
	imageId := image.ImageId
	if image.ImageName != imageName {
		imageId = CreateImage(client, reinstall.sourceItem.InstanceId)
	}
	if "" == imageId {
		fmt.Println("创建源镜像失败")
		return
	}

	ticker := time.NewTicker(60 * time.Second)
	ch := make(chan int)
	go func() {
		var x int
		for x < 10 {
			select {
			case <-ticker.C:
				image = GetImageInfo(client, imageId)
				if image.ImageState == "NORMAL" {
					x = 10
				}
				x++
			}
		}
		ticker.Stop()
		ch <- 0
	}()
	<-ch

	CloseInstance(client, reinstall.sourceItem.InstanceId)

	request := cvm.NewResetInstanceRequest()

	request.InstanceId = common.StringPtr(reinstall.targetItem.InstanceId)
	request.ImageId = common.StringPtr(imageId)
	request.LoginSettings = &cvm.LoginSettings {
		KeepImageLogin: common.StringPtr("TRUE"),
	}

	response, _ := client.ResetInstance(request)

	fmt.Printf("%s", response.ToJsonString())
}

func CloseInstance(client *cvm.Client, instanceId string) {

	request := cvm.NewStopInstancesRequest()
	request.InstanceIds = common.StringPtrs([]string{instanceId})

	response, _ := client.StopInstances(request)

	fmt.Printf("%s", response.ToJsonString())
}
