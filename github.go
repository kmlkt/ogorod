package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func GetToken() {
	resp, _ := http.PostForm("https://github.com/login/device/code", url.Values{"client_id": {"Ov23lixITeS56XKopwWu"}, "scope": {"repo"}})
	body, _ := io.ReadAll(resp.Body)
	values, _ := url.ParseQuery(string(body))
	fmt.Println(values.Get("user_code"))
	device_code := values.Get("device_code")
	resp, _ = http.PostForm("https://github.com/login/oauth/access_token", url.Values{"client_id": {"Ov23lixITeS56XKopwWu"}, "device_code": {device_code}, "grant_type": {"urn:ietf:params:oauth:grant-type:device_code"}})
	body, _ = io.ReadAll(resp.Body)
	values, _ = url.ParseQuery(string(body))
	for values.Get("error") == "authorization_pending" {
		time.Sleep(time.Second * 5)
		resp, _ = http.PostForm("https://github.com/login/oauth/access_token", url.Values{"client_id": {"Ov23lixITeS56XKopwWu"}, "device_code": {device_code}, "grant_type": {"urn:ietf:params:oauth:grant-type:device_code"}})
		body, _ = io.ReadAll(resp.Body)
		values, _ = url.ParseQuery(string(body))
		fmt.Println("Токен не получен")
	}
	fmt.Println(values.Get("access_token"))
}
