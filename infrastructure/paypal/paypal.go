package paypal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/nimil-jp/gin-utils/errors"

	"go-gin-ddd/config"
)

type IPaypal interface {
	ConnectURL(email string) (string, error)
}

type paypal struct {
	url         string
	clientID    string
	secret      string
	accessToken string
}

func New() IPaypal {
	paypalURL := "https://api-m.paypal.com"
	if config.Env.Paypal.Sandbox {
		paypalURL = "https://api-m.sandbox.paypal.com"
	}

	return &paypal{
		url:      paypalURL,
		clientID: config.Env.Paypal.ClientID,
		secret:   config.Env.Paypal.Secret,
	}
}

func (p *paypal) getAccessToken() error {
	form := url.Values{}
	form.Add("grant_type", "client_credentials")
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/v1/oauth2/token", p.url), strings.NewReader(form.Encode()))
	if err != nil {
		return errors.NewUnexpected(err)
	}

	req.SetBasicAuth(p.clientID, p.secret)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return errors.NewUnexpected(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != 200 {
		return errors.NewExpected(http.StatusInternalServerError, "アクセストークンの取得に失敗しました")
	}

	var body map[string]interface{}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return errors.NewUnexpected(err)
	}

	if accessToken, ok := body["access_token"].(string); ok {
		p.accessToken = accessToken
		return nil
	} else {
		return errors.NewExpected(http.StatusInternalServerError, "アクセストークンを取り出せませんでした")
	}
}

type requestOption struct {
	method       string
	path         string
	body         io.Reader
	status       int
	errorMessage string
	succeeded    func(body []byte) error
}

func (p paypal) request(option requestOption) error {
	if p.accessToken == "" {
		err := p.getAccessToken()
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequest(option.method, fmt.Sprintf("%s%s", p.url, option.path), option.body)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.accessToken)

	client := new(http.Client)
	res, err := client.Do(req)
	if err != nil {
		return errors.NewUnexpected(err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	if res.StatusCode != option.status {
		return errors.NewExpected(http.StatusInternalServerError, option.errorMessage)
	}

	jsonBytes, err := io.ReadAll(res.Body)
	return option.succeeded(jsonBytes)
}

func (p *paypal) ConnectURL(email string) (string, error) {
	var jsonStr = []byte(fmt.Sprintf(`{
	"tracking_id": "%s",
	"email": "%s",
    "operations": [
        {
            "operation": "API_INTEGRATION",
            "api_integration_preference": {
                "rest_api_integration": {
                    "integration_method": "PAYPAL",
                    "integration_type": "THIRD_PARTY",
                    "third_party_details": {
                        "features": [
                            "PAYMENT",
                            "REFUND"
                        ]
                    }
                }
            }
        }
    ],
    "products": [
        "PPCP"
    ],
    "legal_consents": [
        {
            "type": "SHARE_DATA_CONSENT",
            "granted": true
        }
    ],
    "partner_configuration_override": {
        "return_url": "%s/dashboard?paypal=connected"
    }
}`, email, email, config.Env.App.URL))
	var res struct {
		Links []struct {
			Href        string `json:"href"`
			Rel         string `json:"rel"`
			Method      string `json:"method"`
			Description string `json:"description"`
		} `json:"links"`
	}
	err := p.request(requestOption{
		method:       http.MethodPost,
		path:         "/v2/customer/partner-referrals",
		body:         bytes.NewBuffer(jsonStr),
		status:       http.StatusCreated,
		errorMessage: "Paypal接続URLの取得に失敗しました",
		succeeded: func(body []byte) error {
			err := json.Unmarshal(body, &res)
			if err != nil {
				return errors.NewUnexpected(err)
			}
			return nil
		},
	})
	if err != nil {
		return "", err
	}
	return res.Links[1].Href, nil
}
