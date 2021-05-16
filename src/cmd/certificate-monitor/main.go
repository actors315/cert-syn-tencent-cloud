package main

import (
	"fmt"
	"net/http"
	"qcloud-tools/src/certificate"
	"qcloud-tools/src/tools"
	"qcloud-tools/src/web"
)

func main() {

	// 开启一个定时器
	go certificate.TickerSchedule()

	rootPath := tools.GetRootPath()
	staticPath := fmt.Sprintf("%s/web/static", rootPath)

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticPath))))
	http.HandleFunc("/add-domain", web.AddDomain)
	http.HandleFunc("/list", web.GetList)
	http.HandleFunc("/", Welcome)

	if err := http.ListenAndServe(":80", nil); err != nil {
		fmt.Println(err)
	}
}

func Welcome(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	fmt.Fprint(writer, "<div>Welcome ~~</div>")
}
