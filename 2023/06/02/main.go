package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

func main() {
	Barrier([]string{"https://www.baidu.com", "https://www.weibo.com", "https://www.shirdon.com"}...)
}

type BarrierResponse struct {
	Err    error
	Resp   string
	Status int
}

func doRequest(out chan<- BarrierResponse, url string) {
	res := BarrierResponse{}
	client := http.Client{
		Timeout: time.Duration(20 * time.Second),
	}

	resp, err := client.Get(url)
	if resp != nil {
		res.Status = resp.StatusCode
	}
	if err != nil {
		res.Err = err
		out <- res
		return
	}

	byt, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		res.Err = err
		out <- res
		return
	}

	res.Resp = string(byt)
	out <- res
}

func Barrier(urls ...string) {
	requestNumber := len(urls)

	in := make(chan BarrierResponse, requestNumber)
	response := make([]BarrierResponse, requestNumber)

	defer close(in)

	for _, urls := range urls {
		go doRequest(in, urls)
	}

	var hasError bool
	for i := 0; i < requestNumber; i++ {
		res := <-in
		if res.Err != nil {
			zap.L().Error("Error in request", zap.Error(res.Err))
			hasError = true
		}
		fmt.Println(res.Status)
		response[i] = res
	}

	if !hasError {
		for _, resp := range response {
			zap.L().Info("Response", zap.String("response", resp.Resp))
		}
	}
}
