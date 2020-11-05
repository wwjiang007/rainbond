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

package region

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/goodrain/rainbond/node/api/model"
	"github.com/goodrain/rainbond/node/nodem/client"
	"github.com/pquerna/ffjson/ffjson"

	//"github.com/goodrain/rainbond/grctl/cmd"

	"github.com/goodrain/rainbond/api/util"
	utilhttp "github.com/goodrain/rainbond/util/http"
)

func (r *regionImpl) Nodes() NodeInterface {
	return &node{regionImpl: *r, prefix: "/v2/nodes"}
}
func (r *regionImpl) Configs() ConfigsInterface {
	return &configs{regionImpl: *r, prefix: "/v2/configs"}
}

type task struct {
	regionImpl
	prefix string
}
type node struct {
	regionImpl
	prefix string
}

func (n *node) Get(node string) (*client.HostNode, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var gc client.HostNode
	res.Bean = &gc
	code, err := n.DoRequest(n.prefix+"/"+node, "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get node detail code %d", code))
	}
	return &gc, nil
}

func (n *node) GetNodeResource(node string) (*client.NodePodResource, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var gc client.NodePodResource
	res.Bean = &gc
	code, err := n.DoRequest(n.prefix+"/"+node+"/resource", "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get node resource code %d", code))
	}
	return &gc, nil
}

func (n *node) GetNodeByRule(rule string) ([]*client.HostNode, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var gc []*client.HostNode
	res.List = &gc
	code, err := n.DoRequest(n.prefix+"/rule/"+rule, "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get node rule code %d", code))
	}
	return gc, nil
}
func (n *node) UpdateNodeStatus(nid, status string) (*client.HostNode, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var gc client.HostNode
	res.Bean = &gc
	req := fmt.Sprintf(`{"status":"%s"}`, status)
	code, err := n.DoRequest(n.prefix+"/"+nid+"/status", "PUT", bytes.NewBuffer([]byte(req)), &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Update node status failure code %d", code))
	}
	return &gc, nil
}
func (n *node) List() ([]*client.HostNode, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var gc []*client.HostNode
	res.List = &gc
	code, err := n.DoRequest(n.prefix, "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get node list code %d", code))
	}
	return gc, nil
}

func (n *node) GetAllNodeHealth() (map[string][]map[string]string, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var gc map[string][]map[string]string
	res.Bean = &gc
	code, err := n.DoRequest(n.prefix+"/all_node_health", "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	if code != 200 {
		return nil, util.CreateAPIHandleError(code, fmt.Errorf("Get node health code %d", code))
	}
	return gc, nil
}

func (n *node) Add(node *client.APIHostNode) (*client.HostNode, *util.APIHandleError) {
	body, err := json.Marshal(node)
	if err != nil {
		return nil, util.CreateAPIHandleError(400, err)
	}
	var res utilhttp.ResponseBody
	var renode client.HostNode
	res.Bean = &renode
	code, err := n.DoRequest(n.prefix, "POST", bytes.NewBuffer(body), &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	return &renode, handleAPIResult(code, res)
}
func (n *node) Label(nid string) NodeLabelInterface {
	return &nodeLabelImpl{nodeImpl: n, NodeID: nid}
}
func (n *node) Condition(nid string) NodeConditionInterface {
	return &nodeConditionImpl{nodeImpl: n, NodeID: nid}
}

type nodeLabelImpl struct {
	nodeImpl *node
	NodeID   string
}

func (nl *nodeLabelImpl) List() (map[string]string, *util.APIHandleError) {
	var decode map[string]string
	var res utilhttp.ResponseBody
	res.Bean = &decode
	code, err := nl.nodeImpl.DoRequest(nl.nodeImpl.prefix+"/"+nl.NodeID+"/labels", "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	return decode, handleAPIResult(code, res)
}
func (nl *nodeLabelImpl) Delete(k string) *util.APIHandleError {
	var decode map[string]string
	var res utilhttp.ResponseBody
	res.Bean = &decode
	code, err := nl.nodeImpl.DoRequest(nl.nodeImpl.prefix+"/"+nl.NodeID+"/labels", "GET", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	if handleAPIResult(code, res) != nil {
		return handleAPIResult(code, res)
	}
	delete(decode, k)
	body, err := json.Marshal(decode)
	if err != nil {
		return util.CreateAPIHandleError(400, err)
	}
	code, err = nl.nodeImpl.DoRequest(nl.nodeImpl.prefix+"/"+nl.NodeID+"/labels", "PUT", bytes.NewBuffer(body), &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}
func (nl *nodeLabelImpl) Add(k, v string) *util.APIHandleError {
	var decode map[string]string
	var res utilhttp.ResponseBody
	res.Bean = &decode
	code, err := nl.nodeImpl.DoRequest(nl.nodeImpl.prefix+"/"+nl.NodeID+"/labels", "GET", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	if handleAPIResult(code, res) != nil {
		return handleAPIResult(code, res)
	}
	decode[k] = v
	body, err := json.Marshal(decode)
	if err != nil {
		return util.CreateAPIHandleError(400, err)
	}
	code, err = nl.nodeImpl.DoRequest(nl.nodeImpl.prefix+"/"+nl.NodeID+"/labels", "PUT", bytes.NewBuffer(body), &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}

type nodeConditionImpl struct {
	nodeImpl *node
	NodeID   string
}

func (nl *nodeConditionImpl) List() ([]client.NodeCondition, *util.APIHandleError) {
	var decode []client.NodeCondition
	var res utilhttp.ResponseBody
	res.List = &decode
	code, err := nl.nodeImpl.DoRequest(nl.nodeImpl.prefix+"/"+nl.NodeID+"/conditions", "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	return decode, handleAPIResult(code, res)
}
func (nl *nodeConditionImpl) Delete(k client.NodeConditionType) ([]client.NodeCondition, *util.APIHandleError) {
	var decode []client.NodeCondition
	var res utilhttp.ResponseBody
	res.List = &decode
	code, err := nl.nodeImpl.DoRequest(nl.nodeImpl.prefix+"/"+nl.NodeID+"/conditions/"+string(k), "DELETE", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	return decode, handleAPIResult(code, res)
}

func (n *node) Delete(nid string) *util.APIHandleError {
	var res utilhttp.ResponseBody
	code, err := n.DoRequest(n.prefix+"/"+nid, "DELETE", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}

func (n *node) Up(nid string) *util.APIHandleError {
	var res utilhttp.ResponseBody
	code, err := n.DoRequest(n.prefix+"/"+nid+"/up", "POST", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}
func (n *node) Down(nid string) *util.APIHandleError {
	var res utilhttp.ResponseBody
	code, err := n.DoRequest(n.prefix+"/"+nid+"/down", "POST", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}
func (n *node) UnSchedulable(nid string) *util.APIHandleError {
	var res utilhttp.ResponseBody
	code, err := n.DoRequest(n.prefix+"/"+nid+"/unschedulable", "PUT", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}
func (n *node) ReSchedulable(nid string) *util.APIHandleError {
	var res utilhttp.ResponseBody
	code, err := n.DoRequest(n.prefix+"/"+nid+"/reschedulable", "PUT", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}
func (n *node) Install(nid string) *util.APIHandleError {
	var res utilhttp.ResponseBody
	code, err := n.DoRequest(n.prefix+"/"+nid+"/install", "POST", nil, &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}

type configs struct {
	regionImpl
	prefix string
}

func (c *configs) Get() (*model.GlobalConfig, *util.APIHandleError) {
	var res utilhttp.ResponseBody
	var gc = model.GlobalConfig{
		Configs: make(map[string]*model.ConfigUnit),
	}
	res.Bean = &gc
	code, err := c.DoRequest(c.prefix+"/datacenter", "GET", nil, &res)
	if err != nil {
		return nil, util.CreateAPIHandleError(code, err)
	}
	return &gc, handleAPIResult(code, res)
}

func (c *configs) Put(gc *model.GlobalConfig) *util.APIHandleError {
	rebody, err := ffjson.Marshal(gc)
	if err != nil {
		return util.CreateAPIHandleError(400, err)
	}
	var res utilhttp.ResponseBody
	code, err := c.DoRequest(c.prefix+"/datacenter", "PUT", bytes.NewBuffer(rebody), &res)
	if err != nil {
		return util.CreateAPIHandleError(code, err)
	}
	return handleAPIResult(code, res)
}

//TaskInterface task api
type TaskInterface interface {
	Get(name string) (*model.Task, *util.APIHandleError)
	GetTaskStatus(task string) (map[string]*model.TaskStatus, *util.APIHandleError)
	Add(task *model.Task) *util.APIHandleError
	AddGroup(group *model.TaskGroup) *util.APIHandleError
	Exec(name string, nodes []string) *util.APIHandleError
	List() ([]*model.Task, *util.APIHandleError)
}

//NodeInterface node api
type NodeInterface interface {
	GetNodeByRule(rule string) ([]*client.HostNode, *util.APIHandleError)
	Get(node string) (*client.HostNode, *util.APIHandleError)
	GetNodeResource(node string) (*client.NodePodResource, *util.APIHandleError)
	List() ([]*client.HostNode, *util.APIHandleError)
	GetAllNodeHealth() (map[string][]map[string]string, *util.APIHandleError)
	Add(node *client.APIHostNode) (*client.HostNode, *util.APIHandleError)
	Up(nid string) *util.APIHandleError
	Down(nid string) *util.APIHandleError
	UnSchedulable(nid string) *util.APIHandleError
	ReSchedulable(nid string) *util.APIHandleError
	Delete(nid string) *util.APIHandleError
	Label(nid string) NodeLabelInterface
	Condition(nid string) NodeConditionInterface
	Install(nid string) *util.APIHandleError
	UpdateNodeStatus(nid, status string) (*client.HostNode, *util.APIHandleError)
}

//NodeLabelInterface node label interface
type NodeLabelInterface interface {
	Add(k, v string) *util.APIHandleError
	Delete(k string) *util.APIHandleError
	List() (map[string]string, *util.APIHandleError)
}

//NodeConditionInterface node condition manage api
type NodeConditionInterface interface {
	List() ([]client.NodeCondition, *util.APIHandleError)
	Delete(conditionType client.NodeConditionType) ([]client.NodeCondition, *util.APIHandleError)
}

//ConfigsInterface 数据中心配置API
type ConfigsInterface interface {
	Get() (*model.GlobalConfig, *util.APIHandleError)
	Put(*model.GlobalConfig) *util.APIHandleError
}
