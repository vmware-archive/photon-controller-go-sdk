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
	"strings"
	"fmt"
)

var _ = Describe("Client", func() {
	Describe("NewClient", func() {
		It("Trims trailing '/' from endpoint", func() {
			endpointList := [] string {
				"http://10.146.1.0/",
				"http://10.146.1.0",
			}

			for index, endpoint := range endpointList {
				client := NewClient(endpoint, "", nil)
				Expect(client.Endpoint).To(
					Equal(strings.TrimRight(endpoint, "/")),
					fmt.Sprintf("Test data index: %v", index))
			}
		})

		It("Trims trailing '/' from authEndpoint", func() {
			endpoint := "http://10.146.1.0/"
			authEndpointList := [] string {
				"http://10.146.1.1/",
				"http://10.146.1.0",
			}

			for index, authEndpoint := range authEndpointList {
				client := NewClient(endpoint, authEndpoint, nil)
				Expect(client.AuthEndpoint).To(
					Equal(strings.TrimRight(authEndpoint, "/")),
					fmt.Sprintf("Test data index: %v", index))
			}
		})
	})
})
