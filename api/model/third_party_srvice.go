// Copyright (C) 2014-2018 Goodrain Co., Ltd.
// RAINBOND, Application Management Platform

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

package model

// AddEndpiontsReq is one of the Endpoints in the request to add the endpints.
type AddEndpiontsReq struct {
	Address  string `json:"address" validate:"address|required"`
	IsOnline bool   `json:"is_online" validate:"required"`
}

// UpdEndpiontsReq is one of the Endpoints in the request to update the endpints.
type UpdEndpiontsReq struct {
	EpID     string `json:"ep_id" validate:"required|len:32"`
	Address  string `json:"address"`
	IsOnline bool   `json:"is_online" validate:"required"`
}

// DelEndpiontsReq is one of the Endpoints in the request to update the endpints.
type DelEndpiontsReq struct {
	EpID string `json:"ep_id" validate:"required|len:32"`
}

// EndpointResp is one of the Endpoints list in the response to list, add,
// update or delete the endpints.
type EndpointResp struct {
	EpID     string `json:"ep_id"`
	Address  string `json:"address"`
	Status   string `json:"status"`
	IsOnline bool   `json:"is_online"`
	IsStatic bool   `json:"is_static"`
}

// ThridPartyServiceProbe is the json obejct in the request
// to update or fetch the ThridPartyServiceProbe.
type ThridPartyServiceProbe struct {
	Scheme       string `json:"scheme;"`
	Path         string `json:"path"`
	Port         int    `json:"port"`
	TimeInterval int    `json:"time_interval"`
	MaxErrorNum  int    `json:"max_error_num"`
	Action       string `json:"action"`
}
