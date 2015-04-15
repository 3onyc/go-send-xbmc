package main

import (
	"flag"
	"github.com/3onyc/go-send-xbmc"
	_ "github.com/3onyc/go-send-xbmc/providers"
	"github.com/atotto/clipboard"
	"log"
)

var (
	XBMC_HOST = flag.String("host", "127.0.0.1", "XBMC host")
	XBMC_PORT = flag.Uint("port", 80, "XBMC API port")

	WEB        = flag.Bool("web", false, "Run a webserver that URLs can be sent to")
	WEB_LISTEN = flag.String("web-listen", ":8080", "IP:Port for the webserver to listen on")

	CLIP = flag.Bool("clip", false, "Grab URL from clipboard")
)

func syntax() {
	log.Println("Syntax: send-xbmc <url>")
	log.Println("        send-xbmc -web")
}

func getInputUrl() (string, error) {
	if *CLIP {
		return clipboard.ReadAll()
	} else {
		return flag.Arg(0), nil
	}
}

func main() {
	flag.Parse()
	api := sendxbmc.NewXbmcApi(*XBMC_HOST, *XBMC_PORT)

	if *WEB {
		s := sendxbmc.NewWebserver(api)
		s.Listen(*WEB_LISTEN)

		return
	}

	inputUrl, err := getInputUrl()
	if err != nil {
		api.SendErrorNotification(err)
		return
	}

	if inputUrl == "" {
		syntax()
		return
	}

	url, err := sendxbmc.Resolve(inputUrl)
	if err != nil {
		api.SendErrorNotification(err)
		return
	}

	if err := api.SendNotification("Opening", url); err != nil {
		log.Fatal(err)
	}

	if err := api.Play(url); err != nil {
		api.SendErrorNotification(err)
	}
}
