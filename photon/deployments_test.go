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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware/photon-controller-go-sdk/photon/internal/mocks"
)

var _ = Describe("Deployment", func() {
	var (
		server         *mocks.Server
		client         *Client
		deploymentSpec *DeploymentCreateSpec
	)

	BeforeEach(func() {
		if isIntegrationTest() {
			Skip("Skipping deployment test on integration mode. Need undeployed environment")
		}
		server, client = testSetup()
		deploymentSpec = &DeploymentCreateSpec{
			ImageDatastores:         []string{randomString(10, "go-sdk-deployment-")},
			UseImageDatastoreForVms: true,
			Auth: &AuthInfo{},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("CreateGetAndDeleteDeployment", func() {
		It("Deployment create and delete succeeds", func() {
			if isIntegrationTest() {
				Skip("Skipping deployment test on integration mode. Need undeployed environment")
			}
			mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.Create(deploymentSpec)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("CREATE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockDeployment := Deployment{
				ImageDatastores:         deploymentSpec.ImageDatastores,
				UseImageDatastoreForVms: deploymentSpec.UseImageDatastoreForVms,
				Auth:                 &AuthInfo{},
				NetworkConfiguration: &NetworkConfiguration{Enabled: false},
			}
			server.SetResponseJson(200, mockDeployment)

			mockTask = createMockTask("DELETE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			task, err = client.Deployments.Delete(task.Entity.ID)
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("DELETE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("GetDeployment", func() {
		Context("GetVM needs a vm", func() {
			var (
				tenantID     string
				projID       string
				imageID      string
				flavorSpec   *FlavorCreateSpec
				vmFlavorSpec *FlavorCreateSpec
				vmSpec       *VmCreateSpec
			)

			BeforeEach(func() {
				tenantID = createTenant(server, client)
				projID = createProject(server, client, tenantID)
				imageID = createImage(server, client)
				flavorSpec = &FlavorCreateSpec{
					[]QuotaLineItem{QuotaLineItem{"COUNT", 1, "ephemeral-disk.cost"}},
					"ephemeral-disk",
					randomString(10, "go-sdk-flavor-"),
				}

				_, err := client.Flavors.Create(flavorSpec)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				vmFlavorSpec = &FlavorCreateSpec{
					Name: randomString(10, "go-sdk-flavor-"),
					Kind: "vm",
					Cost: []QuotaLineItem{
						QuotaLineItem{"GB", 2, "vm.memory"},
						QuotaLineItem{"COUNT", 4, "vm.cpu"},
					},
				}
				_, err = client.Flavors.Create(vmFlavorSpec)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				vmSpec = &VmCreateSpec{
					Flavor:        vmFlavorSpec.Name,
					SourceImageID: imageID,
					AttachedDisks: []AttachedDisk{
						AttachedDisk{
							CapacityGB: 1,
							Flavor:     flavorSpec.Name,
							Kind:       "ephemeral-disk",
							Name:       randomString(10),
							State:      "STARTED",
							BootDisk:   true,
						},
					},
					Name: randomString(10, "go-sdk-vm-"),
				}
			})

			AfterEach(func() {
				cleanVMs(client, projID)
				cleanImages(client)
				cleanFlavors(client)
				cleanTenants(client)
			})

			It("Get VMs succeeds", func() {
				mockTask := createMockTask("CREATE_DEPLOYMENT", "COMPLETED")
				server.SetResponseJson(200, mockTask)

				task, err := client.Deployments.Create(deploymentSpec)
				task, err = client.Tasks.Wait(task.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				mockTask = createMockTask("CREATE_VM", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				vmTask, err := client.Projects.CreateVM(projID, vmSpec)
				vmTask, err = client.Tasks.Wait(vmTask.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				server.SetResponseJson(200, createMockVmsPage(VM{Name: vmSpec.Name}))
				vmList, err := client.Deployments.GetVms(task.Entity.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				var found bool
				for _, vm := range vmList.Items {
					if vm.Name == vmSpec.Name && vm.ID == vmTask.Entity.ID {
						found = true
						break
					}
				}
				Expect(found).Should(BeTrue())

				mockTask = createMockTask("DELETE_VM", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				task, err = client.VMs.Delete(vmTask.Entity.ID)
				task, err = client.Tasks.Wait(task.ID)

				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				mockTask = createMockTask("DELETE_DEPLOYMENT", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				task, err = client.Deployments.Delete(task.Entity.ID)
				task, err = client.Tasks.Wait(task.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
			})
		})

	})

	Describe("EnableAndDisableServiceType", func() {
		It("Enable And Disable Service Type", func() {
			serviceType := "KUBERNETES"
			serviceImageId := "testImageId"
			serviceConfigSpec := &ServiceConfigurationSpec{
				Type:    serviceType,
				ImageID: serviceImageId,
			}

			mockTask := createMockTask("CONFIGURE_SERVICE", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			enableTask, err := client.Deployments.EnableServiceType("deploymentId", serviceConfigSpec)
			enableTask, err = client.Tasks.Wait(enableTask.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(enableTask).ShouldNot(BeNil())
			Expect(enableTask.Operation).Should(Equal("CONFIGURE_SERVICE"))
			Expect(enableTask.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("DELETE_SERVICE_CONFIGURATION", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			disableTask, err := client.Deployments.DisableServiceType("deploymentId", serviceConfigSpec)
			disableTask, err = client.Tasks.Wait(disableTask.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(disableTask).ShouldNot(BeNil())
			Expect(disableTask.Operation).Should(Equal("DELETE_SERVICE_CONFIGURATION"))
			Expect(disableTask.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("ConfigureNsx", func() {
		It("Configure NSX", func() {
			nsxAddress := "nsxAddress"
			nsxUsername := "nsxUsername"
			nsxPassword := "nsxPassword"

			nsxConfigSpec := &NsxConfigurationSpec{
				NsxAddress:  nsxAddress,
				NsxUsername: nsxUsername,
				NsxPassword: nsxPassword,
			}

			mockTask := createMockTask("CONFIGURE_NSX", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			enableTask, err := client.Deployments.ConfigureNsx("deploymentId", nsxConfigSpec)
			enableTask, err = client.Tasks.Wait(enableTask.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(enableTask).ShouldNot(BeNil())
			Expect(enableTask.Operation).Should(Equal("CONFIGURE_NSX"))
			Expect(enableTask.State).Should(Equal("COMPLETED"))
		})
	})
})
