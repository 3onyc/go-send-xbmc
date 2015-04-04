package providers

import (
	"errors"
	"github.com/3onyc/go-send-xbmc"
	"net/url"
	"strings"
)

type DiFmProvider struct {
}

func (p DiFmProvider) Name() string {
	return "DI.fm"
}

func (p DiFmProvider) CanResolve(url *url.URL) bool {
	return strings.Contains(url.Host, "di.fm")
}

func (p DiFmProvider) GetUrl(url *url.URL) (string, error) {
	sName, err := p.parseName(url)
	if err != nil {
		return "", err
	}

	return "http://pub6.di.fm:80/di_" + sName + "_aacplus", nil
}

func (p DiFmProvider) parseName(url *url.URL) (string, error) {
	urlParts := strings.Split(url.Path, "/")
	if len(urlParts) < 2 {
		return "", errors.New("Stream name missing from path")
	}

	return urlParts[1], nil
}

func init() {
	sendxbmc.Register(DiFmProvider{})
}
