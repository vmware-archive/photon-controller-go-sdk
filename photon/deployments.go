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

// Contains functionality for deployments API.
type DeploymentsAPI struct {
	client *Client
}

var deploymentUrl string = rootUrl + "/deployments"

// Creates a deployment
func (api *DeploymentsAPI) Create(deploymentSpec *DeploymentCreateSpec) (task *Task, err error) {
	body, err := json.Marshal(deploymentSpec)
	if err != nil {
		return
	}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+deploymentUrl,
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

// Deletes a deployment with specified ID.
func (api *DeploymentsAPI) Delete(id string) (task *Task, err error) {
	res, err := api.client.restClient.Delete(api.getEntityUrl(id), api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Returns all deployments.
func (api *DeploymentsAPI) GetAll() (result *Deployments, err error) {
	res, err := api.client.restClient.Get(api.client.Endpoint+deploymentUrl, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	result = &Deployments{}
	err = json.NewDecoder(res.Body).Decode(result)
	return
}

// Gets all the vms with the specified deployment ID.
func (api *DeploymentsAPI) GetVms(id string) (result *VMs, err error) {
	uri := api.getEntityUrl(id) + "/vms"
	res, err := api.client.restClient.GetList(api.client.Endpoint, uri, api.client.options.TokenOptions)
	if err != nil {
		return
	}

	result = &VMs{}
	err = json.Unmarshal(res, result)
	return
}

//  Enable service type with specified deployment ID.
func (api *DeploymentsAPI) EnableServiceType(id string, serviceConfigSpec *ServiceConfigurationSpec) (task *Task, err error) {
	body, err := json.Marshal(serviceConfigSpec)
	if err != nil {
		return
	}
	res, err := api.client.restClient.Post(
		api.getEntityUrl(id)+"/enable_service_type",
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

//  Disable service type with specified deployment ID.
func (api *DeploymentsAPI) DisableServiceType(id string, serviceConfigSpec *ServiceConfigurationSpec) (task *Task, err error) {
	body, err := json.Marshal(serviceConfigSpec)
	if err != nil {
		return
	}
	res, err := api.client.restClient.Post(
		api.getEntityUrl(id)+"/disable_service_type",
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

// Configure NSX.
func (api *DeploymentsAPI) ConfigureNsx(id string, nsxConfigSpec *NsxConfigurationSpec) (task *Task, err error) {
	body, err := json.Marshal(nsxConfigSpec)
	if err != nil {
		return
	}

	res, err := api.client.restClient.Post(
		api.getEntityUrl(id)+"/configure_nsx",
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

// Configure NSX.
func (api *DeploymentsAPI) ConfigureNsxCni(id string, nsxCniConfigSpec *NsxCniConfigurationSpec) (task *Task, err error) {
	body, err := json.Marshal(nsxCniConfigSpec)
	if err != nil {
		return
	}

	res, err := api.client.restClient.Post(
		api.getEntityUrl(id)+"/configure_nsx_cni",
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

func (api *DeploymentsAPI) getEntityUrl(id string) (url string) {
	return api.client.Endpoint + deploymentUrl + "/" + id
}
