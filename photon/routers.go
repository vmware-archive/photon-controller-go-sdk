// Copyright (c) 2016 VMware, Inc. All Rights Reserved.
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

// Contains functionality for routers API.
type RoutersAPI struct {
	client *Client
}

var routerUrl string = "/routers/"

// Gets a router with the specified ID.
func (api *RoutersAPI) Get(id string) (router *Router, err error) {
	res, err := api.client.restClient.Get(api.client.Endpoint+routerUrl+id, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Router
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

// Sets router's name.
func (api *RoutersAPI) SetName(id string, routerName *RouterSetNameOperation) (task *Task, err error) {
	body, err := json.Marshal(routerName)
	if err != nil {
		return
	}

	res, err := api.client.restClient.Post(
		api.client.Endpoint+routerUrl+id+"/set_router_name",
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