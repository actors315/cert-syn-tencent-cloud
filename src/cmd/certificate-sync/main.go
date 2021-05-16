package main

import (
	"flag"
	"qcloud-tools/src/cvmVo"
)

func main() {

	var group string
	var file string

	flag.StringVar(&group, "group", "profile", "分组")
	flag.StringVar(&file, "config", "", "配置文件地址")
	flag.Parse()

	cvmVo.UpdateCredential(file,group)

}
