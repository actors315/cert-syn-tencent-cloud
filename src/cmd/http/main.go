package main

import (
	"fmt"
	"net/http"
	"qcloud-tools/src/config"
	"qcloud-tools/src/tools"
	"qcloud-tools/src/web"
)

func main() {

	if !config.QcloudTool.Switch.OpenHttp {
		return
	}

	rootPath := tools.GetRootPath()
	staticPath := fmt.Sprintf("%s/web/static", rootPath)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	http.HandleFunc("/add-domain", web.AddDomain)
	http.HandleFunc("/",  web.GetList)

	addr := fmt.Sprintf(":%d", config.QcloudTool.Http.Port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}