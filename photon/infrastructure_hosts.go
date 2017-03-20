// Copyright (c) 2017 VMware, Inc. All Rights Reserved.
//
// This product is licensed to you under the Apache License, Version 2.0 (the "License").
// You may not use this product except in compliance with the License.
//
// This product may include a number of subcomponents with separate copyright notices and
// license terms. Your use of these subcomponents is subject to the terms and conditions
// of the subcomponent's license, as noted in the LICENSE file.

package photon

import (
	"bytes"
	"encoding/json"
)

// Contains functionality for infrastructure hosts API.
type InfraHostsAPI struct {
	client *Client
}

var InfraHostsUrl string = rootUrl + "/infrastructure/hosts"

// Register a host with photon platform
func (api *InfraHostsAPI) Create(hostSpec *HostCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(hostSpec)
	if err != nil {
		return
	}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+InfraHostsUrl,
		"application/json",
		bytes.NewReader(body),
		api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// List all hosts
func (api *InfraHostsAPI) GetHosts() (result *Hosts, err error) {
	res, err := api.client.restClient.Get(api.client.Endpoint+InfraHostsUrl, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &Hosts{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets a host with the specified ID.
func (api *InfraHostsAPI) Get(id string) (host *Host, err error) {
	res, err := api.client.restClient.Get(api.client.Endpoint+InfraHostsUrl+"/"+id, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Host
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}




// Deletes a host with specified ID.
func (api *InfraHostsAPI) Delete(id string) (task *Task, err error) {
	res, err := api.client.restClient.Delete(api.client.Endpoint+InfraHostsUrl+"/"+id, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Suspend the host with the specified id
func (api *InfraHostsAPI) Suspend(id string) (task *Task, err error) {
	body := []byte{}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+InfraHostsUrl+"/"+id+"/suspend",
		"application/json",
		bytes.NewReader(body),
		api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Resume the host with the specified id
func (api *InfraHostsAPI) Resume(id string) (task *Task, err error) {
	body := []byte{}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+InfraHostsUrl+"/"+id+"/resume",
		"application/json",
		bytes.NewReader(body),
		api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Host with the specified id enter maintenance mode
func (api *InfraHostsAPI) EnterMaintenanceMode(id string) (task *Task, err error) {
	body := []byte{}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+InfraHostsUrl+"/"+id+"/enter-maintenance",
		"application/json",
		bytes.NewReader(body),
		api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Host with the specified id exit maintenance mode
func (api *InfraHostsAPI) ExitMaintenanceMode(id string) (task *Task, err error) {
	body := []byte{}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+InfraHostsUrl+"/"+id+"/exit-maintenance",
		"application/json",
		bytes.NewReader(body),
		api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}




