package gohttpspeedtest

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSingleDownload(t *testing.T) {
	t.Parallel()

	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// generate a payload with 25MB
		payload := strings.Repeat("0123456789012345678901234", 1024*1024)

		w.Write([]byte(payload))
	}))

	size, elapsed, err := singleDownload(mockHttpServer.Client(), mockHttpServer.URL)
	if err != nil {
		t.Errorf("error should be nil: %s", err)
	}
	t.Logf("Size: %d, Elapsed time: %s\n", size, elapsed)
}

func TestFailedReadBodySingleDownload(t *testing.T) {
	t.Parallel()

	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Lie about content-length to raise an error on response body read
		w.Header().Set("Content-Length", "1")

	}))

	_, _, err := singleDownload(mockHttpServer.Client(), mockHttpServer.URL)
	if !strings.Contains(err.Error(), "error on read body:") {
		t.Errorf("should receive a error on read body: %s", err)
	}
}

func TestFailedNot200StatusCodeSingleDownload(t *testing.T) {
	t.Parallel()

	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, _, err := singleDownload(mockHttpServer.Client(), mockHttpServer.URL)
	if !strings.Contains(err.Error(), "error on get request to") {
		t.Errorf("should receive a error on request: %s", err)
	}
}

func TestSingleUpload(t *testing.T) {
	t.Parallel()

	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))

	size, elapsed, err := singleUpload(mockHttpServer.Client(), mockHttpServer.URL)
	if err != nil {
		t.Errorf("error should be nil: %s", err)
	}
	t.Logf("Size: %d, Elapsed time: %s\n", size, elapsed)
}

func TestFailedNot200StatusCodeSingleUpload(t *testing.T) {
	t.Parallel()

	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, _, err := singleUpload(mockHttpServer.Client(), mockHttpServer.URL)
	if !strings.Contains(err.Error(), "error on get request to") {
		t.Errorf("should receive a error on request: %s", err)
	}
}

func TestMeasureDownload(t *testing.T) {
	t.Parallel()

	countURLCalls := 0
	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countURLCalls++

		// generate a payload with 25MB
		payload := strings.Repeat("0123456789012345678901234", 1024*1024)
		w.Write([]byte(payload))

	}))

	measureDownload(mockHttpServer.URL, 5)

	if countURLCalls != 5 {
		t.Errorf("MeasureDownloads should make 5 calls to URL")
	}
}

func TestFailedMeasureDownload(t *testing.T) {
	t.Parallel()

	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, err := measureDownload(mockHttpServer.URL, 5)
	if !strings.Contains(err.Error(), "error on download") {
		t.Errorf("should have an error on download")
	}
}

func TestMeasureUpload(t *testing.T) {
	t.Parallel()

	countURLCalls := 0
	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		countURLCalls++
	}))

	measureUpload(mockHttpServer.URL, 5)

	if countURLCalls != 5 {
		t.Errorf("MeasureDownloads should make 5 calls to URL")
	}
}

func TestFailedMeasureUpload(t *testing.T) {
	t.Parallel()

	mockHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	_, err := measureUpload(mockHttpServer.URL, 5)
	if !strings.Contains(err.Error(), "error on upload") {
		t.Errorf("should have an error on download")
	}
}

func TestMeasureDownloadAndUpload(t *testing.T) {
	t.Parallel()

	okServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))

	tests := []struct {
		name        string
		provider    *Provider
		wantErr     bool
		errContains string
	}{
		{
			name: "download and upload ok",
			provider: &Provider{
				DownloadURL: okServer.URL,
				UploadURL:   okServer.URL,
			},
			wantErr: false,
		},
		{
			name: "error on download",
			provider: &Provider{
				DownloadURL: errorServer.URL,
				UploadURL:   okServer.URL,
			},
			wantErr:     true,
			errContains: "error on measure download",
		},
		{
			name: "error on upload",
			provider: &Provider{
				DownloadURL: okServer.URL,
				UploadURL:   errorServer.URL,
			},
			wantErr:     true,
			errContains: "error on measure upload",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := MeasureDownloadAndUpload(tt.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("MeasureDownloadAndUpload() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestFastProvider(t *testing.T) {
	t.Parallel()

	provider, err := FastProvider()
	if provider.DownloadURL == "" {
		t.Errorf("download url must not be empty")
	}
	if provider.UploadURL == "" {
		t.Errorf("upload url must not be empty")
	}
	if err != nil {
		t.Errorf("error should be nil")
	}

}

func TestOoklaProvider(t *testing.T) {
	t.Parallel()

	provider, err := OoklaProvider()
	if provider.DownloadURL == "" {
		t.Errorf("download url must not be empty")
	}
	if provider.UploadURL == "" {
		t.Errorf("upload url must not be empty")
	}
	if err != nil {
		t.Errorf("error should be nil")
	}

}

func BenchmarkFastProvider(b *testing.B) {
	provider, _ := FastProvider()

	for n := 0; n < b.N; n++ {
		MeasureDownloadAndUpload(provider)
	}

}

func BenchmarkOoklaProvider(b *testing.B) {
	provider, _ := OoklaProvider()

	for n := 0; n < b.N; n++ {
		MeasureDownloadAndUpload(provider)
	}
}
