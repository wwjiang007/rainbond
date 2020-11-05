// RAINBOND, Application Management Platform
// Copyright (C) 2014-2017 Goodrain Co., Ltd.

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version. For any non-GPL usage of Rainbond,
// one or multiple Commercial Licenses authorized by Goodrain Co., Ltd.
// must be obtained first.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
	api_model "github.com/goodrain/rainbond/api/model"
	"github.com/goodrain/rainbond/cmd"
	envoyv2 "github.com/goodrain/rainbond/node/core/envoy/v2"
	"github.com/goodrain/rainbond/util"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		cmd.ShowVersion("sidecar")
	}
	loggerFile, _ := os.Create("/var/log/sidecar.log")
	if loggerFile != nil {
		defer loggerFile.Close()
		logrus.SetOutput(loggerFile)
	}
	if os.Getenv("DEBUG") == "true" {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if err := Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

//Run run
func Run() error {
	// start run first
	run()
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	//step finally: listen Signal
	term := make(chan os.Signal)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	for {
		select {
		case <-term:
			logrus.Warn("Received SIGTERM, exiting gracefully...")
			return nil
		case <-ticker.C:
			run()
		}
	}
}

var oldHosts = make(map[string]string)

func run() {
	configs := discoverConfig()
	if configs != nil {
		if hosts := getHosts(configs); hosts != nil {
			if err := writeHosts(hosts); err != nil {
				logrus.Errorf("write hosts failure %s", err.Error())
			} else {
				logrus.Debugf("rewrite hosts file success, %+v", hosts)
				oldHosts = hosts
			}
		}
	}
}

func haveChange(hosts, oldHosts map[string]string) bool {
	if len(hosts) != len(oldHosts) {
		return true
	}
	for k, v := range hosts {
		if ov, exist := oldHosts[k]; !exist || v != ov {
			return true
		}
	}
	for k, v := range oldHosts {
		if ov, exist := hosts[k]; !exist || v != ov {
			return true
		}
	}
	return false
}

func discoverConfig() *api_model.ResourceSpec {
	discoverURL := fmt.Sprintf("http://%s:6100/v1/resources/%s/%s/%s", os.Getenv("XDS_HOST_IP"), os.Getenv("TENANT_ID"), os.Getenv("SERVICE_NAME"), os.Getenv("PLUGIN_ID"))
	http.DefaultClient.Timeout = time.Second * 5
	res, err := http.Get(discoverURL)
	if err != nil {
		logrus.Errorf("get config failure %s", err.Error())
	}
	if res != nil && res.Body != nil {
		defer res.Body.Close()
		var rs api_model.ResourceSpec
		if err := json.NewDecoder(res.Body).Decode(&rs); err != nil {
			logrus.Errorf("parse config body failure %s", err.Error())
		} else {
			return &rs
		}
	}
	return nil
}

func getHosts(configs *api_model.ResourceSpec) map[string]string {
	hosts := make(map[string]string)
	for _, service := range configs.BaseServices {
		options := envoyv2.GetOptionValues(service.Options)
		for _, domain := range options.Domains {
			if domain != "" && domain != "*" {
				hosts[domain] = "127.0.0.1"
			}
		}
	}
	if len(hosts) == 0 {
		return nil
	}
	return hosts
}

func writeHosts(ipnames map[string]string) error {
	hostFilePath := os.Getenv("HOST_FILE_PATH")
	if hostFilePath == "" {
		hostFilePath = "/etc/hosts"
	}
	hosts, err := util.NewHosts(hostFilePath)
	if err != nil {
		return err
	}
	for name, ip := range ipnames {
		hosts.Add(ip, name)
	}
	return hosts.Flush()
}
