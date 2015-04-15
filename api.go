package sendxbmc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type XbmcApi struct {
	Host string
	Port uint
}

func NewXbmcApi(h string, p uint) *XbmcApi {
	return &XbmcApi{
		Host: h,
		Port: p,
	}
}

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

func (a XbmcApi) SendXbmc(r *XbmcRequest) error {
	xbmcUrl := fmt.Sprintf("http://%s:%d/jsonrpc", a.Host, a.Port)
	log.Println(xbmcUrl)

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

func (a XbmcApi) SendNotification(title, msg string) error {
	return a.SendXbmc(NewXbmcRequest("GUI.ShowNotification", NewNotification(title, msg)))
}

func (a XbmcApi) SendErrorNotification(err error) {
	log.Println(err)

	if err := a.SendNotification("send-xbmc error", err.Error()); err != nil {
		log.Fatal(err)
	}
}

func (a XbmcApi) Play(url string) error {
	return a.SendXbmc(NewXbmcRequest("Player.Open", NewPlayerOpenParams(url)))
}
