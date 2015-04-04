package sendxbmc

import (
	"net/url"
	"strings"
)

var Providers = []Provider{}

type Provider interface {
	Name() string
	CanResolve(url *url.URL) bool
	GetUrl(url *url.URL) (string, error)
}

func Register(p Provider) {
	Providers = append(Providers, p)
}

func Resolve(u string) (string, error) {
	if !strings.Contains(u, "://") {
		u = "nil://" + u
	}

	url, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	for _, p := range Providers {
		if p.CanResolve(url) {
			return p.GetUrl(url)
		}
	}

	return u, nil
}
