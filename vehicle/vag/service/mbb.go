package service

import (
	"net/url"

	"github.com/evcc-io/evcc/util"
	"github.com/evcc-io/evcc/vehicle/vag"
	"github.com/evcc-io/evcc/vehicle/vag/mbb"
)

// MbbTokenSource creates a refreshing token source for use with the MBB api.
// Once the MBB token expires, it is recreated from the token exchanger (either TokenRefreshService or IDK)
func MbbTokenSource(log *util.Logger, toxValues url.Values, clientID string, q url.Values, user, password string) (vag.TokenSource, error) {
	ts, err := TokenRefreshServiceTokenSource(log, toxValues, q, user, password)
	if err != nil {
		return nil, err
	}

	mbb := mbb.New(log, clientID)

	mts := vag.MetaTokenSource(func() (*vag.Token, error) {
		// get TRS token from refreshing TRS token source
		itoken, err := ts.TokenEx()
		if err != nil {
			return nil, err
		}

		// exchange TRS id_token for MBB token
		mtoken, err := mbb.Exchange(url.Values{"id_token": {itoken.IDToken}})
		if err != nil {
			return nil, err
		}

		return mtoken, err

		// produce tokens from refresh MBB token source
	}, mbb.TokenSource)

	return mts, nil
}
