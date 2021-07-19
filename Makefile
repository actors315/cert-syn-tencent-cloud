cert-monitor:
	go build -o bin/cert-monitor src/cmd/certificate-monitor/main.go
	chmod +x bin/cert-monitor

cert-http:
	go build -o bin/cert-http src/cmd/http/main.go
	chmod +x bin/cert-http

clean:
	rm -rf bin/*
	go clean
