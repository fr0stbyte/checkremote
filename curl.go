package main

import (
	curl "github.com/andelf/go-curl"
	"net/url"
)

type CurlResponse struct {
	name          string
	dns           float64
	connect       float64
	appConnect    float64
	preTransfer   float64
	startTransfer float64
	total         float64
	redirect      float64
	statusCode    int
	bytes         float64
}

func nullHandler(buf []byte, userdata interface{}) bool {
	return true
}

func checkCurl(provider string, queue chan *CurlResponse) {
	u, _ := url.Parse(provider)
	res := new(CurlResponse)
	res.name = u.Host
	easy := curl.EasyInit()
	defer easy.Cleanup()
	if len(*opts.headers) > 0 {
		easy.Setopt(curl.OPT_HTTPHEADER, *opts.headers)
	}
	easy.Setopt(curl.OPT_PROTOCOLS, curl.OPT_USE_SSL)
	easy.Setopt(curl.OPT_USERAGENT, *opts.agent)
	if *opts.http11 {
		easy.Setopt(curl.OPT_HTTP_VERSION, curl.HTTP_VERSION_1_1)
	} else {
		easy.Setopt(curl.OPT_HTTP_VERSION, curl.HTTP_VERSION_1_0)
	}

	// make curl quiet
	easy.Setopt(curl.OPT_WRITEFUNCTION, nullHandler)

	easy.Setopt(curl.OPT_URL, provider)
	easy.Perform()

	x, err := easy.Getinfo(curl.INFO_NAMELOOKUP_TIME)
	if err != nil {
		res.dns = 0
	} else {
		res.dns = x.(float64) * 1000
	}

	x, err = easy.Getinfo(curl.INFO_CONNECT_TIME)
	if err != nil {
		res.connect = 0
	} else {
		res.connect = x.(float64) * 1000
	}

	x, err = easy.Getinfo(curl.INFO_APPCONNECT_TIME)
	if err != nil {
		res.appConnect = 0
	} else {
		res.appConnect = x.(float64) * 1000
	}

	x, err = easy.Getinfo(curl.INFO_PRETRANSFER_TIME)
	if err != nil {
		res.preTransfer = 0
	} else {
		res.preTransfer = x.(float64) * 1000
	}

	x, err = easy.Getinfo(curl.INFO_STARTTRANSFER_TIME)
	if err != nil {
		res.startTransfer = 0
	} else {
		res.startTransfer = x.(float64) * 1000
	}

	x, err = easy.Getinfo(curl.INFO_REDIRECT_TIME)
	if err != nil {
		res.redirect = 0
	} else {
		res.redirect = x.(float64) * 1000
	}

	x, err = easy.Getinfo(curl.INFO_TOTAL_TIME)
	if err != nil {
		res.total = 0
	} else {
		res.total = x.(float64) * 1000
	}

	y, err := easy.Getinfo(curl.INFO_RESPONSE_CODE)
	if err != nil {
		res.statusCode = 0
	} else {
		res.statusCode = y.(int)
	}

	x, err = easy.Getinfo(curl.INFO_SIZE_DOWNLOAD)
	if err != nil {
		res.bytes = 0
	} else {
		res.bytes = x.(float64)
	}

	processResponse(res)
	queue <- res
}

func processResponse(res *CurlResponse) {
	res.startTransfer -= res.preTransfer
	if res.appConnect > 0 {
		res.preTransfer -= res.appConnect
		res.appConnect -= res.connect
	} else {
		res.preTransfer -= res.connect
	}

	res.connect -= res.dns
}
