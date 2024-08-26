package core

import (
	"log"
	"net/url"
)

func ConvertToEndpoint(rawUrl string, weight uint64) (ServerEndpoint, error) {
	parsedUrl, err := url.Parse(rawUrl)

	if err != nil {
		return ServerEndpoint{}, err
	}

	return ServerEndpoint{
		URL:    *parsedUrl,
		Weight: weight,
	}, nil
}

func MustConvertToEndpoint(rawUrl string, weight uint64) ServerEndpoint {
	endpoint, err := ConvertToEndpoint(rawUrl, weight)

	if err != nil {
		log.Panic(err)
	}

	return endpoint
}
