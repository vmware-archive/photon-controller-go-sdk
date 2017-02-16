// +build !windows

package lightwave

import "errors"

func (client *OIDCClient) GetTokensFromWindowsLogInContext() (tokens *OIDCTokenResponse, err error) {
	return nil, errors.New("Not supported on this OS")
}