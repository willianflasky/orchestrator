package logic

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"orcguard/dao/mysql"
	. "orcguard/mylogger"
	"time"

	"github.com/bitly/go-simplejson"
)

type Info struct {
	RWDomain  string
	RODomain  string
	Oldmaster string
	Newmaster string
}

func NewInfo(old, new string) *Info {
	return &Info{
		Oldmaster: old,
		Newmaster: new,
	}
}

func (self *Info) Run() {
	var err error
	self.RWDomain, self.RODomain, err = mysql.OpertionDB_dao(self.Oldmaster, self.Newmaster)
	if err != nil {
		L.Error("DB change failed: %v", err)
		return
	} else {
		// 1. 修改写域名
		host := self.dnsapi_get(self.RWDomain)
		if len(host) != 1 {
			L.Error("rw a record ! = 1")
		} else {
			if self.Oldmaster == host[0] {
				ret := self.dnsapi_update(self.Oldmaster, self.Newmaster, self.RWDomain)
				if ret == false {
					L.Error("self.dnsapi_update failed.")
				}
			} else {
				L.Error("self.Oldmaster != host[0], please check.")
			}
		}
		// 2. 修改读域名
		hosts := self.dnsapi_get(self.RODomain)
		if len(hosts) > 1 {
			ret := self.dnsapi_del(self.Newmaster, self.RODomain)
			if ret == false {
				L.Error("self.dnsapi_del failed.")
			}
		} else {
			L.Info("%v a record <= 1, so do not to del.", self.RODomain)
		}
	}

}

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
			fmt.Println(mp)
			hosts = append(hosts, mp["content"])
		}
		L.Info("request success. hosts: %v", hosts...)
		return hosts
	} else {
		L.Error("request failed")
		fmt.Println("request failed")
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
		fmt.Println(msg)
		L.Error(msg)
		return false
	} else {
		fmt.Println("request success.")
		L.Info("request success.")
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
		fmt.Println("request success.")
		L.Info("request success.")
	}
	return true
}
