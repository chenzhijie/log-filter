package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

var filter *Filter

type HTTPResp struct {
	Result string `json:"result"`
	Data   string `json:"data"`
}

func StartHTTPServer() {
	filter = NewFilter()
	filter.StartCleanCache()

	fmt.Println("HTTP server start listen 3001")
	http.Handle("/", http.FileServer(http.Dir("./dist")))
	http.HandleFunc("/v0/log", filterHandler)

	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func setHeader(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	w.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型
	w.Header().Set("content-type", "application/json")
}

func filterHandler(w http.ResponseWriter, req *http.Request) {
	setHeader(w)
	i := 1
	var err error
	if len(req.FormValue("node")) > 0 {
		i, err = strconv.Atoi(req.FormValue("node"))
		if err != nil {
			ret, _ := newHTTPResp("failed", err)
			io.WriteString(w, ret)
			return
		}
	}

	logType := LOG_ALL
	if len(req.FormValue("type")) > 0 {
		logType = (int)(filter.LogTypeFromString(req.FormValue("type")))
	}
	more := 0
	if len(req.FormValue("more")) > 0 {
		more, err = strconv.Atoi(req.FormValue("more"))
		if err != nil {
			ret, _ := newHTTPResp("failed", err)
			io.WriteString(w, ret)
			return
		}
	}

	fmt.Printf("GET node:%d, type: %d, more:%d\n", i, logType, more)
	data, err := filter.GetFileContString(uint8(i), uint8(logType), more)
	if err != nil {
		ret, _ := newHTTPResp("failed", err)
		io.WriteString(w, ret)
		return
	}
	ret, _ := newHTTPResp("success", data)
	io.WriteString(w, ret)
}

func newHTTPResp(result string, data interface{}) (string, error) {
	resp := &HTTPResp{
		Result: result,
		Data:   fmt.Sprintf("%s", data),
	}
	buf, err := json.Marshal(resp)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}
