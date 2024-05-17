

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"

build-w-compress:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -w"
	upx --best --lzma mikdecap


add-dummy-eth:
	ip link add dummy0 type dummy
	ip link set dummy0 mtu 9000 up
