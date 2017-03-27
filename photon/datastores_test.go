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
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/vmware/photon-controller-go-sdk/photon/internal/mocks"
)

var _ = ginkgo.Describe("Datastores", func() {
	var (
		server *mocks.Server
		client *Client
	)

	ginkgo.BeforeEach(func() {
		server, client = testSetup()
	})

	ginkgo.AfterEach(func() {
		server.Close()
	})

	ginkgo.Describe("Get", func() {
		ginkgo.It("Get a single datastore successfully", func() {
			kind := "datastore"
			datastoreType := "LOCAL_VMFS"
			datastoreID := "1234"
			server.SetResponseJson(200,
				Datastore{
					Kind:     kind,
					Type:     datastoreType,
					ID:       datastoreID,
					SelfLink: "https://192.0.0.2/datastores/1234",
				})

			datastore, err := client.Datastores.Get("1234")
			ginkgo.GinkgoT().Log(err)

			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(datastore).ShouldNot(gomega.BeNil())
			gomega.Expect(datastore.Kind).Should(gomega.Equal(kind))
			gomega.Expect(datastore.Type).Should(gomega.Equal(datastoreType))
			gomega.Expect(datastore.ID).Should(gomega.Equal(datastoreID))
		})
	})
	ginkgo.Describe("GetAll", func() {
		ginkgo.It("Get all datastores successfully", func() {
			kind := "datastore"
			datastoreType := "LOCAL_VMFS"
			datastoreID := "1234"
			datastore := Datastore{
				Kind:     kind,
				Type:     datastoreType,
				ID:       datastoreID,
				SelfLink: "https://192.0.0.2/datastores/1234",
			}
			datastoresExpected := Datastores{
				Items: []Datastore{datastore},
			}
			server.SetResponseJson(200, datastoresExpected)

			datastores, err := client.Datastores.GetAll()
			ginkgo.GinkgoT().Log(err)

			gomega.Expect(err).Should(gomega.BeNil())
			gomega.Expect(datastores).ShouldNot(gomega.BeNil())
			gomega.Expect(datastores.Items[0].Kind).Should(gomega.Equal(kind))
			gomega.Expect(datastores.Items[0].Type).Should(gomega.Equal(datastoreType))
			gomega.Expect(datastores.Items[0].ID).Should(gomega.Equal(datastoreID))
		})
	})
})
