// +build windows

package lightwave

import (
	"encoding/base64"
	"fmt"
	"github.com/vmware/photon-controller-go-sdk/SSPI"
	"math/rand"
	"net"
	"net/url"
	"strings"
	"time"
)

const gssTicketGrantFormatString = "grant_type=urn:vmware:grant_type:gss_ticket&gss_ticket=%s&context_id=%s&scope=%s"

func (client *OIDCClient) GetTokensFromWindowsLogInContext() (tokens *OIDCTokenResponse, err error) {
	spn, err := client.buildSPN()
	if err != nil {
		return nil, err
	}

	auth, _ := SSPI.GetAuth("", "", spn, "")

	userContext, err := auth.InitialBytes()
	if err != nil {
		return nil, err
	}

	contextId := client.generateRandomString()
	body := fmt.Sprintf(gssTicketGrantFormatString, url.QueryEscape(base64.StdEncoding.EncodeToString(userContext)), contextId, client.Options.TokenScope)
	tokens, err = client.getToken(body)

	for {
		if err == nil {
			break
		}

		// In case of error the response will be in format: invalid_grant: gss_continue_needed:'context id':'token from server'
		gssToken := client.validateAndExtractGSSResponse(err, contextId)
		if gssToken == "" {
			return nil, err
		}

		data, err := base64.StdEncoding.DecodeString(gssToken)
		if err != nil {
			return nil, err
		}

		userContext, err := auth.NextBytes(data)
		body := fmt.Sprintf(gssTicketGrantFormatString, url.QueryEscape(base64.StdEncoding.EncodeToString(userContext)), contextId, client.Options.TokenScope)
		tokens, err = client.getToken(body)
	}

	return tokens, err
}

// Gets the SPN (Service Principal Name) in the format host/FQDN of lightwave
func (client *OIDCClient) buildSPN() (spn string, err error) {
	u, err := url.Parse(client.Endpoint)
	if err != nil {
		return "", err
	}

	host, _, err := net.SplitHostPort(u.Host)
	if err != nil {
		return "", err
	}

	addr, err := net.LookupAddr(host)
	if err != nil {
		return "", err
	}

	var s = strings.TrimSuffix(addr[0], ".")
	return "host/" + s, nil
}

func (client *OIDCClient) validateAndExtractGSSResponse(err error, contextId string) string {
	parts := strings.Split(err.Error(), ":")
	if !(len(parts) == 4 && strings.TrimSpace(parts[1]) == "gss_continue_needed" && parts[2] == contextId) {
		return ""
	} else {
		return parts[3]
	}
}

func (client *OIDCClient) generateRandomString() string {
	const length = 10
	rand.Seed(time.Now().UTC().UnixNano())
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}
