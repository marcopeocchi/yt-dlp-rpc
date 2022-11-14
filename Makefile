default:
	go build -o yt-dlp-rpc *.go
	mkdir -p build
	mv yt-dlp-rpc* build

clean:
	rm -r build

multiarch:
	GOOS=linux GOARCH=arm go build -o yt-dlp-rpc_linux-arm *.go
	GOOS=linux GOARCH=arm64 go build -o yt-dlp-rpc_linux-arm64 *.go
	GOOS=linux GOARCH=amd64 go build -o yt-dlp-rpc_linux-amd64 *.go

	mkdir -p build
	mv yt-dlp-rpc* build