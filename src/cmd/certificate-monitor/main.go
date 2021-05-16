package main

import (
	"fmt"
	"net/http"
	"qcloud-tools/src/certificate"
	"qcloud-tools/src/web"
)

func main() {

	// 开启一个定时器
	go certificate.TickerSchedule()

	http.HandleFunc("/add-domain", AddDomain)
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

func AddDomain(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Content-Type", "text/html")
	fmt.Fprint(writer, "<div>AddDomain ~~</div>")
}