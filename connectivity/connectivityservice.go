package connectivity

//go:generate mockgen -destination connectivityservice_mock.go -package connectivity . ConnectivityService
//go:generate mockgen -destination httpclient_mock.go -package connectivity . HttpClient

import (
	"context"
	"net/http"
	"net/url"
)

type ConnectivityService interface {
	IsConnectedToInternet(ctx context.Context) (bool, error)
	NotifyInternetChange(callback NotifyCallback)
}

type NotifyCallback func(bool)

type HttpConnectivityService struct {
	httpClient HttpClient
	callback   NotifyCallback
	testUrl    string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

var _ ConnectivityService = &HttpConnectivityService{}

func NewHttpConnectivityService(httpClient HttpClient, testUrl string) (*HttpConnectivityService, error) {
	// validate url
	_, err := url.Parse(testUrl)
	if err != nil {
		return nil, err
	}

	return &HttpConnectivityService{
		httpClient: httpClient,
		testUrl:    testUrl,
	}, nil
}

func (s *HttpConnectivityService) IsConnectedToInternet(ctx context.Context) (bool, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", s.testUrl, nil)
	if err != nil {
		return false, err
	}

	_, err = s.httpClient.Do(req)
	connected := err == nil

	// callback
	if s.callback != nil {
		go s.callback(connected)
	}

	return connected, nil
}

func (s *HttpConnectivityService) NotifyInternetChange(callback NotifyCallback) {
	s.callback = callback
}
