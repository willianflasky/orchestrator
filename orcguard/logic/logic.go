package logic

import (
	"orcguard/dao/mysql"
	. "orcguard/mylogger"
)

type Info struct {
	RWDomain  string
	RODomain  string
	Oldmaster string
	Newmaster string
	Port      int
}

func NewInfo(old, new string, port int) *Info {
	return &Info{
		Oldmaster: old,
		Newmaster: new,
		Port:      port,
	}
}

func (self *Info) Run() {
	var err error
	self.RWDomain, self.RODomain, err = mysql.OpertionDB_dao(self.Oldmaster, self.Newmaster)
	if err != nil {
		L.Error("DB change failed: %v", err)
		return
	} else {
		// 记录info信息到日志
		L.Info("=====【data】=====")
		L.Info("RWDomain: %v", self.RWDomain)
		L.Info("RODomain: %v", self.RODomain)
		L.Info("OldMaster: %v", self.Oldmaster)
		L.Info("NewMaster: %v", self.Newmaster)
		L.Info("Port: %v\n", self.Port)

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
