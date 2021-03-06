package main

import (
	"fmt"
	"net/http"
	"qcloud-tools/src/config"
	"qcloud-tools/src/tools"
	"qcloud-tools/src/web"
)

func main() {

	rootPath := tools.GetRootPath()
	staticPath := fmt.Sprintf("%s/web/static", rootPath)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	http.HandleFunc("/login", web.CheckLogin)
	http.HandleFunc("/sync/add", web.AddSync)
	http.HandleFunc("/info/add", web.AddDomain)
	http.HandleFunc("/",  web.GetList)

	addr := fmt.Sprintf(":%d", config.QcloudTool.Http.Port)

	if err := http.ListenAndServe(addr, nil); err != nil {
		fmt.Println(err)
	}
}