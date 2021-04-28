package gohttpspeedtest

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Provider is a struct to hold the data need to speed test, download and upload Urls
type Provider struct {
	DownloadURL string
	UploadURL   string
}

const (
	// magic number to fit the tests under 30 seconds
	numMeasures = 5
)

// MeasureDownloadAnUpload receive a provider with a downloadURL and an uploadURL.
// Both downloadURL and uploadURL are used numMeasures times, defined by the constant.
// After each download or upload, the Mbits transfered and time consumed are saved to calculate the value of speed in Mbps.
func MeasureDownloadAndUpload(provider *Provider) (download float64, upload float64, err error) {
	download, err = measureDownload(provider.DownloadURL, numMeasures)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("error on measure download: %s", err)
	}
	upload, err = measureUpload(provider.UploadURL, numMeasures)
	if err != nil {
		return 0.0, 0.0, fmt.Errorf("error on measure upload: %s", err)
	}
	return download, upload, nil
}

func measureDownload(downloadUrl string, num int) (float64, error) {
	var totalMBits float64 // In Mbit
	var totalTime float64  // In seconds

	for i := 0; i < num; i++ {
		size, elapsed, err := singleDownload(http.DefaultClient, downloadUrl)
		if err != nil {
			return 0.0, fmt.Errorf("error on download: %s", err)
		}

		// size is in bytes, should convert to Mbits
		totalMBits += 8 * float64(size) / (1024.0 * 1024.0)
		totalTime += elapsed.Seconds()

	}
	// Mbits/seconds = Mbps
	return totalMBits / totalTime, nil
}

func measureUpload(uploadUrl string, num int) (float64, error) {
	var totalMBits float64 // In Mbit
	var totalTime float64  // In seconds

	for i := 0; i < num; i++ {
		size, elapsed, err := singleUpload(http.DefaultClient, uploadUrl)
		if err != nil {
			return 0.0, fmt.Errorf("error on upload: %s", err)
		}

		// size is in bytes, should convert to Mbits
		totalMBits += 8 * float64(size) / (1024.0 * 1024.0)
		totalTime += elapsed.Seconds()

	}
	// Mbits/seconds = Mbps
	return totalMBits / totalTime, nil
}

// singleDownload get the resource in downloadUrl and measure the size of the resource and elapsed time.
func singleDownload(client *http.Client, downloadUrl string) (size int, elapsed time.Duration, err error) {
	start := time.Now()
	res, err := client.Get(downloadUrl)
	if err != nil || res.StatusCode != http.StatusOK {
		return 0, time.Duration(0), fmt.Errorf("error on get request to %s: %s", downloadUrl, err)
	}
	defer res.Body.Close()

	// read the body to get the size of resource
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, time.Duration(0), fmt.Errorf("error on read body: %s", err)
	}
	elapsed = time.Since(start)

	return len(body), elapsed, nil
}

// singleDownload send 25MB to uploadUrl and measure the elapsed time.
func singleUpload(client *http.Client, uploadUrl string) (size int, elapsed time.Duration, err error) {
	v := url.Values{}
	// generate a 25MB object by repeasing 1M times 25 chars.
	v.Add("content", strings.Repeat("0123456789012345678901234", 1024*1024))

	start := time.Now()
	res, err := client.PostForm(uploadUrl, v)
	if err != nil || res.StatusCode != http.StatusOK {
		return 0, time.Duration(0), fmt.Errorf("error on get request to %s: %s", uploadUrl, err)
	}
	defer res.Body.Close()
	elapsed = time.Since(start)

	return len(v.Get("content")), elapsed, nil
}
