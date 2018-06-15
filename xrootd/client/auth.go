// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package client // import "go-hep.org/x/hep/xrootd/client"

import (
	"bytes"
	"context"

	"github.com/pkg/errors"
)

func (client *Client) auth(ctx context.Context, securityInformation []byte) error {
	securityInformation = bytes.TrimLeft(securityInformation, "&")
	providerInfos := bytes.Split(securityInformation, []byte{'&'})

	var errs []error
	for _, providerInfo := range providerInfos {
		providerInfo = bytes.TrimLeft(providerInfo, "P=")[:]
		paramsData := bytes.Split(providerInfo, []byte{','})
		params := make([]string, len(paramsData))
		for i := range paramsData {
			params[i] = string(paramsData[i])
		}
		provider := params[0]
		params = params[1:]

		auther, ok := client.auths[provider]
		if !ok {
			errs = append(errs, errors.Errorf("xrootd: could not authorize using %s: provider was not found", provider))
			continue
		}
		r, err := auther.Request(params)
		if err != nil {
			errs = append(errs, errors.Errorf("xrootd: could not authorize using %s: %v", provider, err))
			continue
		}
		_, err = client.call(ctx, r)
		if err != nil {
			errs = append(errs, errors.Errorf("xrootd: could not authorize using %s: %v", provider, err))
			continue
		}
		return nil
	}

	return errors.Errorf("xrootd: could not authorize:\n%v", errs)
}
