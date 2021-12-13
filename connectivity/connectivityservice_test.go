package connectivity

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_NewHttpConnectivityService(t *testing.T) {
	type args struct {
		httpClient HttpClient
		testUrl    string
	}
	tests := []struct {
		name    string
		args    args
		want    *HttpConnectivityService
		wantErr bool
	}{{
		name:    "Bad Url",
		args:    args{http.DefaultClient, ":://badurl"},
		want:    nil,
		wantErr: true,
	},
		{
			name: "valid client",
			args: args{httpClient: http.DefaultClient, testUrl: "https://google.com"},
			want: &HttpConnectivityService{httpClient: http.DefaultClient, testUrl: "https://google.com"},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewHttpConnectivityService(tt.args.httpClient, tt.args.testUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHttpConnectivityService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHttpConnectivityService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_IsConnectedToInternet_happyPath(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockHttpClient := NewMockHttpClient(ctrl)
	// This mock will always return without an error
	mockHttpClient.EXPECT().Do(gomock.Any()).Return(&http.Response{}, nil)

	s, err := NewHttpConnectivityService(mockHttpClient, "http://example.com")

	statusChan := make(chan bool)
	s.NotifyInternetChange(func(status bool) {
		fmt.Println("sending")
		statusChan <- status
	})

	actualStatus, err := s.IsConnectedToInternet(context.TODO())

	assert.Equal(t, actualStatus, true, "internet status")
	assert.Nil(t, err)

	ticker := time.NewTicker(time.Second)
	select {
	case <-statusChan:
		// callback was called. This test passed
	case <-ticker.C:
		t.Error("Callback was never called")
	}
}
