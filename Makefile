all: win32 win64 linux32 linux64 linux-armv7 linux-armv6

build:
	mkdir build

win32: build
	GOOS=windows GOARCH=386 go build -o build/sendxbmc-i386.exe -ldflags '-H=windowsgui' github.com/3onyc/go-send-xbmc/send-xbmc

win64: build
	GOOS=windows GOARCH=amd64 go build -o build/sendxbmc-amd64.exe -ldflags '-H=windowsgui' github.com/3onyc/go-send-xbmc/send-xbmc

linux32: build
	GOOS=linux GOARCH=386 go build -o build/sendxbmc-i386 github.com/3onyc/go-send-xbmc/send-xbmc

linux64: build
	GOOS=linux GOARCH=amd64 go build -o build/sendxbmc-amd64 github.com/3onyc/go-send-xbmc/send-xbmc

linux-armv7: build
	GOOS=linux GOARCH=arm GOARM=7 go build -o build/sendxbmc-armv7 github.com/3onyc/go-send-xbmc/send-xbmc

linux-armv6: build
	GOOS=linux GOARCH=arm GOARM=6 go build -o build/sendxbmc-armv6 github.com/3onyc/go-send-xbmc/send-xbmc

.PHONY: all win32 win64 linux32 linux64 linux-armv7 linux-armv6
