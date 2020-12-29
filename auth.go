package tradestation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const stateName = "tradestationAuth"
const securityAuthorizePath = "/security/authorize"
const authorizePath = "/authorize"

// sample from doc, note that response when refreshing does not include a refresh token
//{
//    "refresh_token": "eGlhc2xvTTVJaEdXMWs4VjhraWx4bk5QMHJMaA==",
//    "expires_in": 1200,
//    "access_token": "eGlhc2xvozT2IxWnVITmdwGVFPQ==",
//    "token_type": "AccessToken",
//    "userid": "testUser"
//}
type accessTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Userid       string `json:"userid"`
}

func (r *tradestation) AccessTokenStatus() string {
	if r.AccessToken == "" {
		err := r.PersistentState.Retrieve(stateName, &r.authState)
		if err != nil {
			return "No access token or error reading state - " + err.Error()
		}
		if r.AccessToken == "" {
			return "State read but still no access token, auth flow must be incomplete"
		}
	}
	if time.Now().UTC().After(r.Expires) {
		return "Access token expired"
	}
	return "Access token valid"
}

func (r *tradestation) AccessTokenOK() error {
	if r.AccessToken == "" {
		r.PersistentState.Retrieve(stateName, &r.authState)
		if r.AccessToken == "" {
			return ErrNoAccessToken
		}
	}
	if time.Now().UTC().Before(r.Expires) {
		return nil
	}

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("client_id", r.clientID)
	data.Set("redirect_uri", r.RedredirectURI)
	data.Set("client_secret", r.clientSecret)
	data.Set("refresh_token", r.RefreshToken)
	data.Set("response_type", "token")
	requestData := data.Encode()

	req, err := http.NewRequest(
		http.MethodPost,
		r.baseURL+securityAuthorizePath,
		strings.NewReader(requestData),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	// should set req.ContentLength instead but for strings.Reader this is done automatically
	//req.Header.Add("Content-Length", strconv.Itoa(len(requestData)))
	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrStatusNotOK
	}

	decoder := json.NewDecoder(resp.Body)
	decodedResponse := &accessTokenResponse{}
	err = decoder.Decode(decodedResponse)
	if err != nil {
		return err
	}
	r.AccessToken = decodedResponse.AccessToken
	r.Expires = time.Now().UTC().Add(time.Second * time.Duration(decodedResponse.ExpiresIn-20))

	return r.PersistentState.Save(stateName, &r.authState)
}

func (r *tradestation) StartAuth(redirectURI string) (string, error) {
	redirectURI = strings.TrimSpace(redirectURI)
	if redirectURI == "" {
		return "", fmt.Errorf("Argument empty: %v", "redirectURI")
	}

	r.RefreshToken = ""
	r.AccessToken = ""
	r.Expires = time.Now().UTC()
	r.RedredirectURI = redirectURI
	err := r.PersistentState.Save(stateName, &r.authState)
	if err != nil {
		return "", err
	}

	data := url.Values{}
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", r.clientID)
	data.Set("response_type", "code")

	return r.baseURL + authorizePath + "?" + data.Encode(), nil
}

func (r *tradestation) FinishAuth(code string) error {
	code = strings.TrimSpace(code)
	if code == "" {
		return fmt.Errorf("Argument empty: %v", "code")
	}

	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("client_id", r.clientID)
	data.Set("redirect_uri", r.RedredirectURI)
	data.Set("client_secret", r.clientSecret)
	data.Set("code", code)
	data.Set("response_type", "token")
	requestData := data.Encode()

	req, err := http.NewRequest(
		http.MethodPost,
		r.baseURL+securityAuthorizePath,
		strings.NewReader(requestData),
	)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := r.RoundTripper.RoundTrip(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ErrStatusNotOK
	}

	decoder := json.NewDecoder(resp.Body)
	decodedResponse := &accessTokenResponse{}
	err = decoder.Decode(decodedResponse)
	if err != nil {
		return err
	}
	r.AccessToken = decodedResponse.AccessToken
	r.RefreshToken = decodedResponse.RefreshToken
	r.Expires = time.Now().UTC().Add(time.Second * time.Duration(decodedResponse.ExpiresIn-20))

	err = r.PersistentState.Save(stateName, &r.authState)
	if err != nil {
		return err
	}

	return nil
}
