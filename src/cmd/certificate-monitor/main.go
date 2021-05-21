package main

import (
	"context"
	"qcloud-tools/src/certificate"
	"qcloud-tools/src/core"
	"time"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	go core.SignalHandler(cancel)

	// 开启一个定时器
	certificate.TickerSchedule(ctx)

	time.Sleep(time.Second * 5)
}
