package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"go.uber.org/zap"

	"gitlab.com/slon/shad-go/distbuild/pkg/build"
)

type BuildClient struct {
	l        *zap.Logger
	endpoint string
}

func NewBuildClient(l *zap.Logger, endpoint string) *BuildClient {
	return &BuildClient{
		l:        l,
		endpoint: endpoint,
	}
}

type statusReader struct {
	r   io.ReadCloser
	dec *json.Decoder
}

func (r *statusReader) Close() error {
	return r.r.Close()
}

func (r *statusReader) Next() (*StatusUpdate, error) {
	var u StatusUpdate
	if err := r.dec.Decode(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (c *BuildClient) StartBuild(ctx context.Context, request *BuildRequest) (*BuildStarted, StatusReader, error) {
	reqJSON, err := json.Marshal(request)
	if err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/build", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("content-type", "application/json")

	rsp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, nil, err
	}
	defer func() {
		if rsp.Body != nil {
			_ = rsp.Body.Close()
		}
	}()

	if rsp.StatusCode != 200 {
		bodyStr, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			return nil, nil, fmt.Errorf("build request failed: %v", err)
		}

		return nil, nil, fmt.Errorf("build failed: %s", bodyStr)
	}

	dec := json.NewDecoder(rsp.Body)
	var started BuildStarted
	if err := dec.Decode(&started); err != nil {
		return nil, nil, err
	}

	r := &statusReader{r: rsp.Body, dec: dec}
	rsp.Body = nil
	return &started, r, nil
}

func (c *BuildClient) SignalBuild(ctx context.Context, buildID build.ID, signal *SignalRequest) (*SignalResponse, error) {
	signalJSON, err := json.Marshal(signal)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.endpoint+"/signal?build_id="+buildID.String(), bytes.NewBuffer(signalJSON))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json")

	rsp, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	rspBody, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("signal request failed: %v", err)
	}

	if rsp.StatusCode != 200 {
		return nil, fmt.Errorf("signal failed: %s", rspBody)
	}

	var signalRsp SignalResponse
	if err = json.Unmarshal(rspBody, &rsp); err != nil {
		return nil, err
	}

	return &signalRsp, err
}
