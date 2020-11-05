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

package probe

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/goodrain/rainbond/node/nodem/service"
	v1 "github.com/goodrain/rainbond/util/prober/types/v1"
)

// HTTPProbe probes through the http protocol
type HTTPProbe struct {
	Name          string
	Address       string
	ResultsChan   chan *v1.HealthStatus
	Ctx           context.Context
	Cancel        context.CancelFunc
	TimeInterval  int
	MaxErrorsNum  int
	TimeoutSecond int
}

//Check starts http probe.
func (h *HTTPProbe) Check() {
	go h.HTTPCheck()
}

//Stop stops http probe.
func (h *HTTPProbe) Stop() {
	h.Cancel()
}

//HTTPCheck http check
func (h *HTTPProbe) HTTPCheck() {
	if h.TimeInterval == 0 {
		h.TimeInterval = 5
	}
	timer := time.NewTimer(time.Second * time.Duration(h.TimeInterval))
	defer timer.Stop()
	for {
		HealthMap := h.GetHTTPHealth()
		result := &v1.HealthStatus{
			Name:   h.Name,
			Status: HealthMap["status"],
			Info:   HealthMap["info"],
		}
		h.ResultsChan <- result
		timer.Reset(time.Second * time.Duration(h.TimeInterval))
		select {
		case <-h.Ctx.Done():
			return
		case <-timer.C:
		}
	}
}

// Return true if the underlying error indicates a http.Client timeout.
//
// Use for errors returned from http.Client methods (Get, Post).
func isClientTimeout(err error) bool {
	if uerr, ok := err.(*url.Error); ok {
		if nerr, ok := uerr.Err.(net.Error); ok && nerr.Timeout() {
			return true
		}
	}
	return false
}

//GetHTTPHealth get http health
func (h *HTTPProbe) GetHTTPHealth() map[string]string {
	address := h.Address
	c := &http.Client{
		Timeout: time.Duration(h.TimeoutSecond) * time.Second,
	}
	if strings.HasPrefix(address, "https://") {
		c.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		logrus.Warnf("address %s do not has scheme, auto add http scheme", address)
		address = "http://" + address
	}
	addr, err := url.Parse(address)
	if err != nil {
		logrus.Errorf("%s is invalid %s", address, err.Error())
		return map[string]string{"status": service.Stat_healthy, "info": "check url is invalid"}
	}
	if addr.Scheme == "" {
		addr.Scheme = "http"
	}
	logrus.Debugf("http probe check address; %s", address)
	resp, err := c.Get(addr.String())
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		if isClientTimeout(err) {
			return map[string]string{"status": service.Stat_death, "info": "Request service timeout"}
		}
		logrus.Debugf("http probe request error %s", err.Error())
		return map[string]string{"status": service.Stat_unhealthy, "info": err.Error()}
	}
	if resp.StatusCode >= 400 {
		logrus.Debugf("http probe check address %s return code %d", address, resp.StatusCode)
		return map[string]string{"status": service.Stat_unhealthy, "info": "Service unhealthy"}
	}
	return map[string]string{"status": service.Stat_healthy, "info": "service health"}
}
