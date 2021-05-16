cert-monitor:
	go build -o bin/cert-monitor src/cmd/certificate-monitor/main.go
	chmod +x bin/cert-monitor

cert-sync:
	go build -o bin/cert-sync src/cmd/certificate-sync/main.go
	chmod +x bin/cert-sync

cvm-renew:
	go build -o bin/cvm-renew src/cmd/cvm-reinstall/main.go
	chmod +x bin/cvm-renew

clean:
	rm -rf bin/*
	go clean
