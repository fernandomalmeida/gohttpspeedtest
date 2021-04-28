package gohttpspeedtest

import (
	"fmt"
	"strings"

	"github.com/showwin/speedtest-go/speedtest"
	"gopkg.in/ddo/go-fast.v0"
)

func FastProvider() (*Provider, error) {

	fastCom := fast.New()
	fastCom.Init()

	urls, err := fastCom.GetUrls()
	if err != nil {
		return nil, fmt.Errorf("error on get Urls from fast.com: %s", err)
	}

	return &Provider{
		DownloadURL: urls[0],
		UploadURL:   urls[0],
	}, nil
}

func OoklaProvider() (*Provider, error) {
	user, err := speedtest.FetchUserInfo()
	if err != nil {
		return nil, fmt.Errorf("error on fetch user info: %s", err)
	}

	serverList, err := speedtest.FetchServerList(user)
	if err != nil {
		return nil, fmt.Errorf("error on fetch server list: %s", err)
	}

	targets, err := serverList.FindServer([]int{})
	if err != nil || len(targets) < 1 {
		return nil, fmt.Errorf("error on find servers: %s", err)
	}

	baseURL := strings.Split(targets[0].URL, "/upload.php")[0]

	return &Provider{
		DownloadURL: fmt.Sprintf("%s/random%dx%d.jpg", baseURL, 4000, 4000),
		UploadURL:   fmt.Sprintf("%s/upload.php", baseURL),
	}, nil
}
