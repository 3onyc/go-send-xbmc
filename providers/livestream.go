package providers

import (
	"errors"
	"github.com/3onyc/go-send-xbmc"
	"net/url"
	"strings"
)

type LivestreamProvider struct {
}

func (p LivestreamProvider) Name() string {
	return "DI.fm"
}

func (p LivestreamProvider) CanResolve(url *url.URL) bool {
	return strings.Contains(url.Host, "di.fm")
}

func (p LivestreamProvider) GetUrl(url *url.URL) (string, error) {
	sName, err := p.parseName(url)
	if err != nil {
		return "", err
	}

	return "rtsp://mobilestr2.livestream.com/livestreamiphone/" + sName, nil
}

func (p LivestreamProvider) parseName(url *url.URL) (string, error) {
	urlParts := strings.Split(url.Path, "/")
	if len(urlParts) < 2 {
		return "", errors.New("Stream name missing from path")
	}

	return urlParts[1], nil
}

func init() {
	sendxbmc.Register(LivestreamProvider{})
}
