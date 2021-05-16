package cvmVo

import (
	"fmt"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	cvm "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	"time"
)

type Image struct {
	ImageId    string
	ImageState string
	ImageName string
}

func GetImageInfo(client *cvm.Client, imageId string) *Image {
	image := new(Image)

	request := cvm.NewDescribeImagesRequest()

	if imageId != "" {
		request.Filters = []*cvm.Filter{
			&cvm.Filter{
				Name:   common.StringPtr("image-id"),
				Values: common.StringPtrs([]string{imageId}),
			},
		}
	} else {
		request.Filters = []*cvm.Filter{
			&cvm.Filter{
				Name:   common.StringPtr("image-type"),
				Values: common.StringPtrs([]string{"PRIVATE_IMAGE"}),
			},
		}
	}

	response, err := client.DescribeImages(request)
	if err != nil {
		return image
	}

	image.ImageId = imageId
	image.ImageState = *response.Response.ImageSet[0].ImageState
	image.ImageName = *response.Response.ImageSet[0].ImageName

	return image
}

func ClearExpiredImage(client *cvm.Client) {

	request := cvm.NewDescribeImagesRequest()

	request.Filters = []*cvm.Filter{
		&cvm.Filter{
			Name:   common.StringPtr("image-type"),
			Values: common.StringPtrs([]string{"PRIVATE_IMAGE"}),
		},
	}

	response, err := client.DescribeImages(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}

	if *response.Response.TotalCount < 5 {
		return
	}

	var deleteImageIds []string
	var i int64 = 5
	for ; i < *response.Response.TotalCount; i++ {
		image := response.Response.ImageSet[i]
		if *image.ImageState == "USING" || *image.ImageState == "CREATING" {
			continue
		}
		deleteImageIds = append(deleteImageIds, *image.ImageId)
	}

	if len(deleteImageIds) > 0 {
		DeleteImages(client, deleteImageIds)
	}
}

func CreateImage(client *cvm.Client, instanceId string) string {

	request := cvm.NewCreateImageRequest()

	request.InstanceId = common.StringPtr(instanceId)
	request.ImageName = common.StringPtr("bak-" + time.Now().Format("2006-01-02"))

	response, err := client.CreateImage(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		msg := fmt.Sprintf("An API error has returned: %s", err)
		panic(msg)
	}
	if err != nil {
		panic(err)
	}

	return *response.Response.ImageId
}

func DeleteImages(client *cvm.Client, imageIds []string) {

	request := cvm.NewDeleteImagesRequest()

	request.ImageIds = common.StringPtrs(imageIds)

	response, err := client.DeleteImages(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s", response.ToJsonString())
}
