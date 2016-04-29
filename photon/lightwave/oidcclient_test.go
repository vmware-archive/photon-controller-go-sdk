// Copyright (c) 2016 VMware, Inc. All Rights Reserved.
//
// This product is licensed to you under the Apache License, Version 2.0 (the "License").
// You may not use this product except in compliance with the License.
//
// This product may include a number of subcomponents with separate copyright notices and
// license terms. Your use of these subcomponents is subject to the terms and conditions
// of the subcomponent's license, as noted in the LICENSE file.

package lightwave

import (
	"fmt"
	"strings"

	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/vmware/photon-controller-go-sdk/photon/internal/mocks"
)

var _ = Describe("OIDCClient", func() {
	var (
		client *OIDCClient
		server *mocks.Server
	)

	AfterEach(func() {
		if server != nil {
			server.Close()
			server = nil
		}
	})

	Describe("NewOIDCClient", func() {
		It("Trims trailing '/' from endpoint", func() {
			endpointList := []string{
				"http://10.146.1.0/",
				"http://10.146.1.0",
			}

			for index, endpoint := range endpointList {
				client := NewOIDCClient(endpoint, nil, nil)
				Expect(client.Endpoint).To(
					Equal(strings.TrimRight(endpoint, "/")),
					fmt.Sprintf("Test data index: %v", index))
			}
		})
	})

	Describe("GetRootCerts", func() {
		Context("with fake server", func() {
			BeforeEach(func() {
				client, server = testSetupFakeServer()
			})

			Context("when server responds with valid certificate", func() {
				BeforeEach(func() {
					template := &x509.Certificate{
						IsCA: true,
						BasicConstraintsValid: true,
						SubjectKeyId:          []byte{1, 2, 3},
						SerialNumber:          big.NewInt(1234),
						Subject: pkix.Name{
							Country:      []string{"Earth"},
							Organization: []string{"Mother Nature"},
						},
						NotBefore: time.Now(),
						NotAfter:  time.Now().AddDate(5, 5, 5),
						// see http://golang.org/pkg/crypto/x509/#KeyUsage
						ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
						KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
					}

					// generate private key
					privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
					Expect(err).To(BeNil())

					cert, err := x509.CreateCertificate(rand.Reader, template, template, &privatekey.PublicKey, privatekey)
					Expect(err).To(BeNil())

					certOut := new(bytes.Buffer)
					err = pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: cert})
					Expect(err).To(BeNil())

					body := []lightWaveCert{
						lightWaveCert{Value: certOut.String()},
					}
					server.SetResponseJsonForPath(certDownloadPath, 200, body)
				})

				It("retrieves certificates", func() {
					certList, err := client.GetRootCerts()
					Expect(err).To(BeNil())
					Expect(certList).ToNot(BeNil())
					Expect(len(certList)).To(BeNumerically(">", 0))
				})
			})

			Context("when server responds with unsupported format", func() {
				BeforeEach(func() {
					server.SetResponseForPath(certDownloadPath, 200, "text")
				})

				It("returns an error", func() {
					certList, err := client.GetRootCerts()
					Expect(certList).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("invalid character 'e' in literal true (expecting 'r')"))
				})
			})

			Context("when server responds with unparasble cert data", func() {
				BeforeEach(func() {
					body := []lightWaveCert{
						lightWaveCert{Value: "text"},
					}

					server.SetResponseJson(200, body)
				})

				It("returns an error", func() {
					certList, err := client.GetRootCerts()
					Expect(certList).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("Unexpected response format: &[{text}]"))
				})
			})

			Context("when server responds with error", func() {
				BeforeEach(func() {
					server.SetResponseForPath(certDownloadPath, 400, "")
				})

				It("returns an error", func() {
					certList, err := client.GetRootCerts()
					Expect(certList).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("Unexpected error retrieving auth server certs: 400 400 Bad Request"))
				})
			})
		})

		Context("with real server", func() {
			BeforeEach(func() {
				client, _, _ = testSetupRealServer(nil)
				if client == nil {
					Skip("LIGHTWAVE_ENDPOINT must be set for this test to run")
				}
			})

			It("retrieves certificates", func() {
				certList, err := client.GetRootCerts()
				Expect(err).To(BeNil())
				Expect(certList).ToNot(BeNil())
				Expect(len(certList)).To(BeNumerically(">", 0))
			})
		})
	})

	Describe("GetTokensPasswordGrant", func() {
		Context("with fake server", func() {
			BeforeEach(func() {
				client, server = testSetupFakeServer()
			})

			Context("when server responds with valid data", func() {
				var (
					expected *OIDCTokenResponse
				)

				BeforeEach(func() {
					expected = &OIDCTokenResponse{
						AccessToken:  "fake_access_token",
						ExpiresIn:    36000,
						RefreshToken: "fake_refresh_token",
						IdToken:      "fake_id_token",
						TokenType:    "Bearer",
					}
					server.SetResponseJsonForPath(tokenPath, 200, expected)
				})

				It("retrieves tokens", func() {
					resp, err := client.GetTokenByPasswordGrant("u", "p")
					Expect(err).To(BeNil())
					Expect(resp).To(BeEquivalentTo(expected))
				})
			})

			Context("when server responds with unsupported format", func() {
				BeforeEach(func() {
					server.SetResponseForPath(tokenPath, 200, "text")
				})

				It("returns an error", func() {
					resp, err := client.GetTokenByPasswordGrant("u", "p")
					Expect(resp).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("invalid character 'e' in literal true (expecting 'r')"))
				})
			})

			Context("when server responds with error", func() {
				BeforeEach(func() {
					server.SetResponseForPath(tokenPath, 400, "Error")
				})

				It("returns an error", func() {
					resp, err := client.GetTokenByPasswordGrant("u", "p")
					Expect(resp).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("Status: 400 Bad Request, Body: Error\n [<nil>]"))
				})
			})
		})

		Context("with real server", func() {
			BeforeEach(func() {
				options := &OIDCClientOptions{
					IgnoreCertificate: true,
				}

				client, _, _ = testSetupRealServer(options)
				if client == nil {
					Skip("LIGHTWAVE_ENDPOINT must be set for this test to run")
				}
			})

			Context("when username/password is wrong", func() {
				It("retrieves certificates", func() {
					resp, err := client.GetTokenByPasswordGrant("u", "p")
					Expect(resp).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.(OIDCError).Code).To(BeEquivalentTo("invalid_grant"))
					Expect(err.(OIDCError).Message).ToNot(BeNil())
				})
			})
		})
	})

	Describe("GetTokensPasswordGrant", func() {
		Context("with fake server", func() {
			BeforeEach(func() {
				client, server = testSetupFakeServer()
			})

			Context("when server responds with valid data", func() {
				var (
					expected *OIDCTokenResponse
				)

				BeforeEach(func() {
					expected = &OIDCTokenResponse{
						AccessToken:  "fake_access_token",
						ExpiresIn:    36000,
						RefreshToken: "fake_refresh_token",
						IdToken:      "fake_id_token",
						TokenType:    "Bearer",
					}
					server.SetResponseJsonForPath(tokenPath, 200, expected)
				})

				It("retrieves tokens", func() {
					resp, err := client.GetTokenByRefreshTokenGrant("rf")
					Expect(err).To(BeNil())
					Expect(resp).To(BeEquivalentTo(expected))
				})
			})

			Context("when server responds with unsupported format", func() {
				BeforeEach(func() {
					server.SetResponseForPath(tokenPath, 200, "text")
				})

				It("returns an error", func() {
					resp, err := client.GetTokenByRefreshTokenGrant("rf")
					Expect(resp).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("invalid character 'e' in literal true (expecting 'r')"))
				})
			})

			Context("when server responds with error", func() {
				BeforeEach(func() {
					server.SetResponseForPath(tokenPath, 400, "Error")
				})

				It("returns an error", func() {
					resp, err := client.GetTokenByRefreshTokenGrant("rf")
					Expect(resp).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err).To(MatchError("Status: 400 Bad Request, Body: Error\n [<nil>]"))
				})
			})
		})

		Context("with real server", func() {
			BeforeEach(func() {
				options := &OIDCClientOptions{
					IgnoreCertificate: true,
				}

				client, _, _ = testSetupRealServer(options)
				if client == nil {
					Skip("LIGHTWAVE_ENDPOINT must be set for this test to run")
				}
			})

			Context("when username/password is wrong", func() {
				It("retrieves certificates", func() {
					resp, err := client.GetTokenByRefreshTokenGrant("rt")
					Expect(resp).To(BeNil())
					Expect(err).ToNot(BeNil())
					Expect(err.(OIDCError).Code).To(BeEquivalentTo("invalid_grant"))
					Expect(err.(OIDCError).Message).ToNot(BeNil())
				})
			})
		})
	})

	Describe("Token Retrieval flow with Real Server", func() {
		var (
			username string
			password string
		)

		BeforeEach(func() {
			options := &OIDCClientOptions{
				IgnoreCertificate: true,
			}

			client, username, password = testSetupRealServer(options)
			if client == nil || len(username) == 0 || len(password) == 0 {
				Skip("LIGHTWAVE_ENDPOINT, LIGHTWAVE_USERNAME, LIGHTWAVE_PASSWORD must be set for this test to run")
			}
		})

		It("retrieves tokens", func() {
			resp, err := client.GetTokenByPasswordGrant(username, password)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())

			client.logger.Printf("TokenType: %v\n", resp.TokenType)
			client.logger.Printf("AcceessToken: %v\n", resp.AccessToken)
			client.logger.Printf("ExpiresIn: %v\n", resp.ExpiresIn)
			client.logger.Printf("IdToken: %v\n", resp.IdToken)
			client.logger.Printf("RefreshToken: %v\n", resp.RefreshToken)

			resp, err = client.GetTokenByRefreshTokenGrant(resp.RefreshToken)
			Expect(err).To(BeNil())
			Expect(resp).ToNot(BeNil())

			client.logger.Printf("TokenType: %v\n", resp.TokenType)
			client.logger.Printf("AcceessToken: %v\n", resp.AccessToken)
			client.logger.Printf("ExpiresIn: %v\n", resp.ExpiresIn)
			client.logger.Printf("IdToken: %v\n", resp.IdToken)
			client.logger.Printf("RefreshToken: %v\n", resp.RefreshToken)
		})
	})

	Describe("ParseTokenDetails", func() {
		Context("with fake server", func() {
			BeforeEach(func() {
				client, server = testSetupFakeServer()
			})

			Context("parsing token details", func() {
				var (
					expected *JWTToken
				)

				BeforeEach(func() {
					server.SetResponseJsonForPath(tokenPath, 200, expected)
				})

				It("parses tokens", func() {
					expected := &JWTToken{
						TokenId:    "CfPby7BAlaOI3Uj_TEq_UJOJmYXJiVOYuCYAXPw2l2U",
						Algorithm:  "RS256",
						Subject:    "ec-admin@esxcloud",
						Audience:   []string{"ec-admin@esxcloud", "rs_esxcloud"},
						Groups:     []string{"esxcloud\\ESXCloudAdmins", "esxcloud\\Everyone"},
						Issuer:     "https://10.146.64.238/openidconnect/esxcloud",
						IssuedAt:   1461795927,
						Expires:    1461803127,
						Scope:      "openid offline_access rs_esxcloud at_groups",
						TokenType:  "Bearer",
						TokenClass: "access_token",
						Tenant:     "esxcloud",
					}
					resp := client.ParseTokenDetails(
						"eyJhbGciOiJSUzI1NiJ9.eyJzdWIiOiJlYy1hZG1pbkBlc3hjbG91ZCIsImF1ZCI6WyJlYy1hZG1pbkB" +
							"lc3hjbG91ZCIsInJzX2VzeGNsb3VkIl0sInNjb3BlIjoib3BlbmlkIG9mZmxpbmVfYWNjZXNzIHJzX2VzeGNsb3VkIGF0X" +
							"2dyb3VwcyIsImlzcyI6Imh0dHBzOlwvXC8xMC4xNDYuNjQuMjM4XC9vcGVuaWRjb25uZWN0XC9lc3hjbG91ZCIsImdyb3V" +
							"wcyI6WyJlc3hjbG91ZFxcRVNYQ2xvdWRBZG1pbnMiLCJlc3hjbG91ZFxcRXZlcnlvbmUiXSwidG9rZW5fY2xhc3MiOiJhY" +
							"2Nlc3NfdG9rZW4iLCJ0b2tlbl90eXBlIjoiQmVhcmVyIiwiZXhwIjoxNDYxODAzMTI3LCJpYXQiOjE0NjE3OTU5MjcsInR" +
							"lbmFudCI6ImVzeGNsb3VkIiwianRpIjoiQ2ZQYnk3QkFsYU9JM1VqX1RFcV9VSk9KbVlYSmlWT1l1Q1lBWFB3MmwyVSJ9." +
							"QOpb-8L8if1kEHPEQvsGe_Z8v_gdlPDpjWcu8LxMnAxZELQx6YBn7UM2MO83Qgo-0bqu2ysbcSpjz0mP4pf48z_DyKlMCa" +
							"B6ViStwavIx7lM1TENrt5nURpjqxlzQY0CxjyYIWxoYQIUbn7c5MXe-vt-OTXAg8bGkwphltj7xUak90mQlZGSBrHFCT_Q" +
							"PGwxRTNsRwWq45tF7LgKr49L4z5PnkLQ3LpC8jI7x1SUFBiYcJgi76pGNlD4qihpmKhGJK0WpspEAvXhtsGwBVavGxeXzL" +
							"-PBTYz7Zs1EjD4Isar-91pq-HeTVfhd_KBBqktaQq0WO48Vu0KtHHRv_Us90-Qs53gsY0CnrxHV8qyNR27LyaIMWhG24hq" +
							"TyBsZVgT-gzs9_-QdLqtkXNgr4Oiqoy9Gi8LAmARGFCgTXOS7uPqZ6_ut71WPhwwoUIuXVUG8vvuRD6_UIIGXyPjBM0sfg" +
							"X5rMeo45bYO51mNjqAysz7FBwMetkZUqKg6pxWmTmO_xnH5D55I1P2zd_VBo5be-hr7jjTqqDAGkGMU0PM8IajpnWe24wu" +
							"lPzQqRr5-HlQx50B0nwhYFJVCd_3KW6qCw-MmfGB-1aX-GVG2wa_vUKzc4gDDn65-z0rP_gYtrB9q8oNR-hPY4v18DQEdY" +
							"bsuoJoqriXk1A0zkeoX13kFXY")
					Expect(resp).To(BeEquivalentTo(expected))
				})
			})
		})
	})
})
