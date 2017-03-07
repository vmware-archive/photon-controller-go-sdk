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

// Contains functionality for services API.
type ServicesAPI struct {
	client *Client
}

var serviceUrl string = "/services/"

const ExtendedPropertyDNS string = "dns"
const ExtendedPropertyGateway string = "gateway"
const ExtendedPropertyNetMask string = "netmask"
const ExtendedPropertyLoadBalancerIP string = "load_balancer_ip"
const ExtendedPropertyMasterCount string = "master_count"
const ExtendedPropertyMasterIPs string = "master_ips"
const ExtendedPropertyMasterIP string = "master_ip"
const ExtendedPropertyMasterIP2 string = "master_ip2"
const ExtendedPropertyContainerNetwork string = "container_network"
const ExtendedPropertyZookeeperIP1 string = "zookeeper_ip1"
const ExtendedPropertyZookeeperIP2 string = "zookeeper_ip2"
const ExtendedPropertyZookeeperIP3 string = "zookeeper_ip3"
const ExtendedPropertyETCDCount string = "etcd_count"
const ExtendedPropertyETCDIP1 string = "etcd_ip1"
const ExtendedPropertyETCDIP2 string = "etcd_ip2"
const ExtendedPropertyETCDIP3 string = "etcd_ip3"
const ExtendedPropertySSHKey string = "ssh_key"
const ExtendedPropertyRegistryCACert string = "registry_ca_cert"
const ExtendedPropertyAdminPassword string = "admin_password"

// Deletes a service with specified ID.
func (api *ServicesAPI) Delete(id string) (task *Task, err error) {
	res, err := api.client.restClient.Delete(api.client.Endpoint+serviceUrl+id, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	task, err = getTask(getError(res))
	return
}

// Gets a service with the specified ID.
func (api *ServicesAPI) Get(id string) (service *Service, err error) {
	res, err := api.client.restClient.Get(api.client.Endpoint+serviceUrl+id, api.client.options.TokenOptions)
	if err != nil {
		return
	}
	defer res.Body.Close()
	res, err = getError(res)
	if err != nil {
		return
	}
	var result Service
	err = json.NewDecoder(res.Body).Decode(&result)
	return &result, nil
}

// Gets vms for service with the specified ID.
func (api *ServicesAPI) GetVMs(id string) (result *VMs, err error) {
	uri := api.client.Endpoint + serviceUrl + id + "/vms"
	res, err := api.client.restClient.GetList(api.client.Endpoint, uri, api.client.options.TokenOptions)
	if err != nil {
		return
	}

	result = &VMs{}
	err = json.Unmarshal(res, result)
	return
}

// Resize a service to specified count.
func (api *ServicesAPI) Resize(id string, resize *ServiceResizeOperation) (task *Task, err error) {
	body, err := json.Marshal(resize)
	if err != nil {
		return
	}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+serviceUrl+id+"/resize",
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

// Start a background process to recreate failed VMs in a service with the specified ID.
func (api *ServicesAPI) TriggerMaintenance(id string) (task *Task, err error) {
	body := []byte{}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+serviceUrl+id+"/trigger_maintenance",
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

// Change a service version to the specified image by destroying and recreating the VMs.
func (api *ServicesAPI) ChangeVersion(id string, changeVersion *ServiceChangeVersionOperation) (task *Task, err error) {
	body, err := json.Marshal(changeVersion)
	if err != nil {
		return
	}
	res, err := api.client.restClient.Post(
		api.client.Endpoint+serviceUrl+id+"/change_version",
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
