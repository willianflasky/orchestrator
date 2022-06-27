package logic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	. "orcguard/mylogger"
	"time"

	"github.com/bitly/go-simplejson"
)

func (self *Info) dnsapi_get(domain string) (hosts []interface{}) {
	apiurl := "https://oap-lab.eeo-inc.com/api/v2/opsTool/getDomain"
	data := url.Values{}
	data.Set("domain", domain)
	u, err := url.ParseRequestURI(apiurl)
	if err != nil {
		fmt.Printf("parse url requestUrl failed, err:%v\n", err)
	}
	u.RawQuery = data.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		fmt.Printf("post failed, err:%v\n", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read from resp.Body failed, err:%v\n", err)
		return
	}
	js, _ := simplejson.NewJson(body)
	code := js.Get("Code").MustInt()
	if code == 200 {
		list, _ := js.Get("Response").Array()
		for _, item := range list {
			mp := item.(map[string]interface{})
			hosts = append(hosts, mp["content"])
		}
		L.Info("dnsapi_get success. hosts: ")
		for _, host := range hosts {
			L.Info("[%v]", host)
		}

		return hosts
	} else {
		L.Error("dnsapi_get failed")
		fmt.Println("dnsapi_get failed")
	}
	return hosts
}

func (self *Info) dnsapi_update(oldmaster, newmaster, domain string) bool {
	url := "https://oap-lab.eeo-inc.com/api/v2/opsTool/updateExistDomainResolution"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("sourceip", oldmaster)
	_ = writer.WriteField("changeip", newmaster)
	_ = writer.WriteField("domain", domain)
	err := writer.Close()
	if err != nil {
		L.Error("%v", err)
		return false
	}

	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		L.Error("%v", err)
		return false
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		L.Error("%v", err)
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		L.Error("read from resp.Body failed, err:%v\n", err)
		return false
	}
	// fmt.Print(string(body))
	js, _ := simplejson.NewJson(body)
	code := js.Get("Code").MustInt()
	if code != 200 {
		msg := js.Get("Response").Get("Error").Get("Message").MustString()
		L.Error(msg)
		return false
	} else {
		fmt.Println("dnsapi_update success.")
		L.Info("dnsapi_update success. domain: [%v] old: [%v] new: [%v]", domain, oldmaster, newmaster)
	}
	return true
}

func (self *Info) dnsapi_del(newmaster, domain string) bool {
	url := "https://oap-lab.eeo-inc.com/api/v2/opsTool/delExistDomainResolution"
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("ip", newmaster)
	_ = writer.WriteField("domain", domain)
	err := writer.Close()
	if err != nil {
		L.Error("%v", err)
		return false
	}
	client := &http.Client{Timeout: 5 * time.Second}
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		L.Error("%v", err)
		return false
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := client.Do(req)
	if err != nil {
		L.Error("%v", err)
		return false
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		L.Error("read from resp.Body failed, err:%v\n", err)
		return false
	}
	// fmt.Print(string(body))
	js, _ := simplejson.NewJson(body)
	code := js.Get("Code").MustInt()
	if code != 200 {
		msg := js.Get("Response").Get("Error").Get("Message").MustString()
		fmt.Println(msg)
		L.Error(msg)
		return false
	} else {
		fmt.Println("dnsapi_del success.")
		L.Info("dnsapi_del success.")
	}
	return true
}
