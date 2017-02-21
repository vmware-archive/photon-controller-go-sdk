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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware/photon-controller-go-sdk/photon/internal/mocks"
)

var _ = Describe("Router", func() {
	var (
		server *mocks.Server
		client *Client
	)

	BeforeEach(func() {
		server, client = testSetup()
	})

	AfterEach(func() {
		cleanImages(client)
		server.Close()
	})

	Describe("UpdateRouter", func() {
		It("update router's name", func() {
			mockTask := createMockTask("UPDATE_ROUTER", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			routerSpec := &RouterUpdateSpec{RouterName: "router-1"}
			task, err := client.Routers.UpdateRouter("router-Id", routerSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("UPDATE_ROUTER"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("DeleteRouter", func() {
		It("Delete router succeeds", func() {
			mockTask := createMockTask("DELETE_ROUTER", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Routers.Delete("routerId")
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)

			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_ROUTER"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})
})
