package providers

import (
	"errors"
	"github.com/3onyc/go-send-xbmc"
	"net/url"
	"strings"
)

type YoutubeProvider struct {
}

func (p YoutubeProvider) Name() string {
	return "youtube"
}

func (p YoutubeProvider) CanResolve(url *url.URL) bool {
	return strings.Contains(url.Host, "youtu")
}

func (p YoutubeProvider) GetUrl(url *url.URL) (string, error) {
	vId, err := p.parseId(url)
	if err != nil {
		return "", err
	}

	return "plugin://plugin.video.youtube/?action=play_video&videoid=" + vId, nil
}

func (p YoutubeProvider) parseId(url *url.URL) (string, error) {
	if strings.HasSuffix(url.Host, "youtu.be") {
		return p.parseShortId(url)
	}

	if _, ok := url.Query()["v"]; !ok {
		return "", errors.New("No video ID found in YouTube URL")
	}

	return url.Query()["v"][0], nil
}

// Parses a YouTu.be URL
func (p YoutubeProvider) parseShortId(url *url.URL) (string, error) {
	if len(url.Path) < 2 {
		return "", errors.New("No video ID found in YouTu.be URL")
	}

	if id := strings.Split(url.Path[1:], "/")[0]; id != "" {
		return id, nil
	} else {
		return "", errors.New("No video ID found in YouTu.be URL")
	}
}

func init() {
	sendxbmc.Register(YoutubeProvider{})
}
