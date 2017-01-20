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
		server *mocks.Server
		client *Client
	)

	BeforeEach(func() {
		if isIntegrationTest() {
			Skip("Skipping deployment test on integration mode. Need undeployed environment")
		}
		server, client = testSetup()
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("GetDeployment", func() {
		It("Get Hosts succeeds", func() {
			mockTask := createMockTask("CREATE_HOST", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			hostSpec := &HostCreateSpec{
				Username: randomString(10),
				Password: randomString(10),
				Address:  randomAddress(),
				Tags:     []string{"CLOUD"},
				Metadata: map[string]string{"test": "go-sdk-host"},
			}

			mockTask = createMockTask("CREATE_HOST", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			hostTask, err := client.Hosts.Create(hostSpec, "deployment-ID")
			hostTask, err = client.Tasks.Wait(hostTask.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())

			server.SetResponseJson(200, createMockHostsPage(Host{Tags: hostSpec.Tags, ID: hostTask.Entity.ID}))
			hostList, err := client.Deployments.GetHosts("deployment-ID")

			var found bool
			for _, host := range hostList.Items {
				if host.ID == hostTask.Entity.ID {
					found = true
					break
				}
			}
			Expect(found).Should(BeTrue())
		})

		Context("GetVM needs a vm", func() {
			var (
				tenantID     string
				resName      string
				projID       string
				imageID      string
				flavorSpec   *FlavorCreateSpec
				vmFlavorSpec *FlavorCreateSpec
				vmSpec       *VmCreateSpec
			)

			BeforeEach(func() {
				tenantID = createTenant(server, client)
				resName = createResTicket(server, client, tenantID)
				projID = createProject(server, client, tenantID, resName)
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
				mockTask := createMockTask("CREATE_VM", "COMPLETED")
				server.SetResponseJson(200, mockTask)
				vmTask, err := client.Projects.CreateVM(projID, vmSpec)
				vmTask, err = client.Tasks.Wait(vmTask.ID)
				GinkgoT().Log(err)
				Expect(err).Should(BeNil())

				server.SetResponseJson(200, createMockVmsPage(VM{Name: vmSpec.Name}))
				vmList, err := client.Deployments.GetVms("deployment-id")
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
				task, err := client.VMs.Delete(vmTask.Entity.ID)
				task, err = client.Tasks.Wait(task.ID)

				GinkgoT().Log(err)
				Expect(err).Should(BeNil())
			})
		})

	})

	Describe("InitializeAndFinalizeMigrateDeployment", func() {
		It("Initialize and Finalize deployment migration succeeds", func() {
			if isIntegrationTest() {
				Skip("Skipping deployment test on integration mode. Need undeployed environment")
			}
			mockTask := createMockTask("INITIALIZE_MIGRATE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			initializeMigrationOperation := &InitializeMigrationOperation{SourceNodeGroupReference: "sourceAddr"}
			task, err := client.Deployments.InitializeDeploymentMigration(initializeMigrationOperation, "deployment-ID")
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("INITIALIZE_MIGRATE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("FINALIZE_MIGRATE_DEPLOYMENT", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			finalizeMigrationOperation := &FinalizeMigrationOperation{SourceNodeGroupReference: "sourceAddr"}
			task, err = client.Deployments.FinalizeDeploymentMigration(finalizeMigrationOperation, "deployment-ID")
			task, err = client.Tasks.Wait(task.ID)
			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("FINALIZE_MIGRATE_DEPLOYMENT"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("SetSecurityGroups", func() {
		It("sets security groups for a project", func() {
			mockTask := createMockTask("SET_SECURITY_GROUPS", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			// Set security groups for the project
			expected := &Deployment{
				Auth: &AuthInfo{
					SecurityGroups: []string{
						randomString(10),
						randomString(10),
					},
				},
			}

			payload := SecurityGroupsSpec{
				Items: expected.Auth.SecurityGroups,
			}
			updateTask, err := client.Deployments.SetSecurityGroups("deployment-ID", &payload)
			updateTask, err = client.Tasks.Wait(updateTask.ID)
			Expect(err).Should(BeNil())

			// Get the security groups for the project
			server.SetResponseJson(200, expected)
			deployment, err := client.Deployments.Get("deployment-ID")
			Expect(err).Should(BeNil())
			Expect(deployment.Auth.SecurityGroups).To(ContainElement(payload.Items[0]))
			Expect(deployment.Auth.SecurityGroups).To(ContainElement(payload.Items[1]))
		})
	})

	Describe("SetImageDatastores", func() {
		It("Succeeds", func() {
			mockTask := createMockTask("UPDATE_IMAGE_DATASTORES", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			imageDatastores := &ImageDatastores{
				[]string{"imageDatastore1", "imageDatastore2"},
			}
			createdTask, err := client.Deployments.SetImageDatastores("deploymentId", imageDatastores)
			createdTask, err = client.Tasks.Wait(createdTask.ID)

			Expect(err).Should(BeNil())
			Expect(createdTask.Operation).Should(Equal("UPDATE_IMAGE_DATASTORES"))
			Expect(createdTask.State).Should(Equal("COMPLETED"))
		})

		It("Fails", func() {
			mockApiError := createMockApiError("INVALID_IMAGE_DATASTORES", "Not a super set", 400)
			server.SetResponseJson(400, mockApiError)

			imageDatastores := &ImageDatastores{
				[]string{"imageDatastore1", "imageDatastore2"},
			}
			createdTask, err := client.Deployments.SetImageDatastores("deploymentId", imageDatastores)

			Expect(err).Should(Equal(*mockApiError))
			Expect(createdTask).Should(BeNil())
		})
	})

	Describe("SyncHostsConfig", func() {
		It("Sync Hosts Config succeeds", func() {
			mockTask := createMockTask("SYNC_HOSTS_CONFIG", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.SyncHostsConfig("deploymentId")
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("SYNC_HOSTS_CONFIG"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("PauseSystemAndPauseBackgroundTasks", func() {
		It("Pause System and Resume System succeeds", func() {
			mockTask := createMockTask("PAUSE_SYSTEM", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.PauseSystem("deploymentId")
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("PAUSE_SYSTEM"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("RESUME_SYSTEM", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err = client.Deployments.PauseBackgroundTasks("deploymentId")
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("RESUME_SYSTEM"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})

		It("Pause Background Tasks and Resume System succeeds", func() {
			mockTask := createMockTask("PAUSE_BACKGROUND_TASKS", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err := client.Deployments.PauseBackgroundTasks("deploymentId")
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("PAUSE_BACKGROUND_TASKS"))
			Expect(task.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("RESUME_SYSTEM", "COMPLETED")
			server.SetResponseJson(200, mockTask)

			task, err = client.Deployments.PauseBackgroundTasks("deploymentId")
			task, err = client.Tasks.Wait(task.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(task).ShouldNot(BeNil())
			Expect(task.Operation).Should(Equal("RESUME_SYSTEM"))
			Expect(task.State).Should(Equal("COMPLETED"))
		})
	})

	Describe("EnableAndDisableClusterType", func() {
		It("Enable And Disable Cluster Type", func() {
			clusterType := "SWARM"
			clusterImageId := "testImageId"
			clusterConfigSpec := &ClusterConfigurationSpec{
				Type:    clusterType,
				ImageID: clusterImageId,
			}

			mockTask := createMockTask("CONFIGURE_CLUSTER", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			enableTask, err := client.Deployments.EnableClusterType("deploymentId", clusterConfigSpec)
			enableTask, err = client.Tasks.Wait(enableTask.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(enableTask).ShouldNot(BeNil())
			Expect(enableTask.Operation).Should(Equal("CONFIGURE_CLUSTER"))
			Expect(enableTask.State).Should(Equal("COMPLETED"))

			mockTask = createMockTask("DELETE_CLUSTER_CONFIGURATION", "COMPLETED")
			server.SetResponseJson(200, mockTask)
			disableTask, err := client.Deployments.DisableClusterType("deploymentId", clusterConfigSpec)
			disableTask, err = client.Tasks.Wait(disableTask.ID)

			GinkgoT().Log(err)
			Expect(err).Should(BeNil())
			Expect(disableTask).ShouldNot(BeNil())
			Expect(disableTask.Operation).Should(Equal("DELETE_CLUSTER_CONFIGURATION"))
			Expect(disableTask.State).Should(Equal("COMPLETED"))
		})
	})
})
