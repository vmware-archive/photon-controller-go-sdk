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
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Auth", func() {
	var (
		server *testServer
		client *Client
	)

	BeforeEach(func() {
		server, client = testSetup()
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("GetAuth", func() {
		It("returns auth info", func() {
			expected := &AuthInfo{
				Enabled: false,
				Port:    0,
			}
			server.SetResponseJson(200, expected)
			info, err := client.Auth.Get()
			fmt.Fprintf(GinkgoWriter, "Got auth info: %+v\n", info)
			Expect(info).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
		})
	})
})

var _ = Describe("Tokens", func() {
	var (
		server *testServer
		client *Client
	)

	BeforeEach(func() {
		server, client = testSetup()
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("GetTokensByPassword", func() {
		It("returns tokens", func() {
			expected := &TokenOptions{
				AccessToken:  "fake_access_token",
				ExpiresIn:    36000,
				RefreshToken: "fake_refresh_token",
				IdToken:      "fake_id_token",
				TokenType:    "Bearer",
			}
			server.SetResponseJson(200, expected)
			info, err := client.Auth.GetTokensByPassword("username", "password")
			fmt.Fprintf(GinkgoWriter, "Got tokens: %+v\n", info)
			Expect(info).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
		})
	})

	Describe("GetTokensByRefreshToken", func() {
		It("returns tokens", func() {
			expected := &TokenOptions{
				AccessToken: "fake_access_token",
				ExpiresIn:   36000,
				IdToken:     "fake_id_token",
				TokenType:   "Bearer",
			}
			server.SetResponseJson(200, expected)
			info, err := client.Auth.GetTokensByRefreshToken("refresh_token")
			fmt.Fprintf(GinkgoWriter, "Got tokens: %+v\n", info)
			Expect(info).ShouldNot(BeNil())
			Expect(err).Should(BeNil())
		})
	})
})
