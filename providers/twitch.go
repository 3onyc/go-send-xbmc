package providers

import (
	"errors"
	"github.com/3onyc/go-send-xbmc"
	"net/url"
	"strings"
)

type TwitchProvider struct {
}

func (p TwitchProvider) Name() string {
	return "twitch"
}

func (p TwitchProvider) CanResolve(url *url.URL) bool {
	return strings.Contains(url.Host, "twitch.tv")
}

func (p TwitchProvider) GetUrl(url *url.URL) (string, error) {
	sName, err := p.parseName(url)
	if err != nil {
		return "", err
	}

	return "plugin://plugin.video.twitch/playLive/" + sName, nil
}

func (p TwitchProvider) parseName(url *url.URL) (string, error) {
	urlParts := strings.Split(url.Path, "/")
	if len(urlParts) < 2 {
		return "", errors.New("Stream name missing from path")
	}

	return urlParts[1], nil
}

func init() {
	sendxbmc.Register(TwitchProvider{})
}
