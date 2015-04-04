package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/3onyc/go-send-xbmc"
	_ "github.com/3onyc/go-send-xbmc/providers"
	"github.com/atotto/clipboard"
	"log"
	"math/rand"
	"net/http"
)

var (
	XBMC_HOST = flag.String("host", "127.0.0.1", "XBMC host")
	XBMC_PORT = flag.Uint("port", 80, "XBMC API port")

	WEB  = flag.Bool("web", false, "Run a webserver that URLs can be sent to")
	CLIP = flag.Bool("clip", false, "Grab URL from clipboard")
)

type XbmcRequest struct {
	Method  string      `json:"method"`
	Id      uint32      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Params  interface{} `json:"params"`
}

func NewXbmcRequest(method string, params interface{}) *XbmcRequest {
	return &XbmcRequest{
		Method:  method,
		Id:      rand.Uint32(),
		JsonRpc: "2.0",
		Params:  params,
	}
}

type PlayerOpenParams struct {
	Item PlayerOpenParamsItem `json:"item"`
}

type PlayerOpenParamsItem struct {
	File string `json:"file"`
}

func NewPlayerOpenParams(url string) *PlayerOpenParams {
	return &PlayerOpenParams{
		Item: PlayerOpenParamsItem{
			File: url,
		},
	}
}

type Notification struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewNotification(title, msg string) *Notification {
	return &Notification{
		Title:   title,
		Message: msg,
	}
}

func sendXbmc(r *XbmcRequest) error {
	xbmcUrl := fmt.Sprintf("http://%s:%d/jsonrpc", *XBMC_HOST, *XBMC_PORT)

	enc, err := json.Marshal(r)
	if err != nil {
		return err
		log.Printf("%s\n", enc)
	}

	resp, err := http.Post(xbmcUrl, "application/json", bytes.NewReader(enc))
	if err != nil {
		return err
	}
	log.Printf("%v\n", resp)

	return nil
}

func sendNotification(title, msg string) error {
	return sendXbmc(NewXbmcRequest("GUI.ShowNotification", NewNotification(title, msg)))
}

func sendErrorNotification(err error) {
	log.Println(err)

	if err := sendNotification("send-xbmc error", err.Error()); err != nil {
		log.Fatal(err)
	}
}

func play(url string) error {
	return sendXbmc(NewXbmcRequest("Player.Open", NewPlayerOpenParams(url)))
}

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
	inputUrl, err := getInputUrl()
	if err != nil {
		sendErrorNotification(err)
		return
	}

	if *WEB {
		sendxbmc.Webserver()
		return
	}

	if inputUrl == "" {
		syntax()
		return
	}

	url, err := sendxbmc.Resolve(inputUrl)
	if err != nil {
		sendErrorNotification(err)
		return
	}

	if err := sendNotification("Opening", url); err != nil {
		log.Fatal(err)
	}

	if err := play(url); err != nil {
		sendErrorNotification(err)
	}
}
