package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"

	"github.com/Azure/go-autorest/autorest/adal"
)

const (
	msiEndpoint             = "http://localhost:50342/oauth2/token"
	activeDirectoryEndpoint = "https://login.microsoftonline.com/"
	armResource             = "https://management.core.windows.net/"
	clientAppID             = "04b07795-8ddb-461a-bbee-02f9e1bf7b46"
)

type responseJSON struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Resource     string `json:"resource"`
	TokenType    string `json:"token_type"`
}

func defaultTokenCachePath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	defaultTokenPath := usr.HomeDir + "/.adal/accessToken.json"
	return defaultTokenPath
}

func acquireTokenDeviceCodeFlow(oauthConfig adal.OAuthConfig,
	applicationID string,
	resource string,
	callbacks ...adal.TokenRefreshCallback) (*adal.ServicePrincipalToken, error) {

	oauthClient := &http.Client{}
	deviceCode, err := adal.InitiateDeviceAuth(
		oauthClient,
		oauthConfig,
		applicationID,
		resource)
	if err != nil {
		return nil, fmt.Errorf("Failed to start device auth flow: %s", err)
	}

	fmt.Println(*deviceCode.Message)

	token, err := adal.WaitForUserCompletion(oauthClient, deviceCode)
	if err != nil {
		return nil, fmt.Errorf("Failed to finish device auth flow: %s", err)
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(
		oauthConfig,
		applicationID,
		resource,
		*token,
		callbacks...)
	return spt, err
}

func refreshToken(oauthConfig adal.OAuthConfig,
	applicationID string,
	resource string,
	tokenCachePath string,
	callbacks ...adal.TokenRefreshCallback) (*adal.ServicePrincipalToken, error) {

	token, err := adal.LoadToken(tokenCachePath)
	if err != nil {
		return nil, fmt.Errorf("failed to load token from cache: %v", err)
	}

	spt, err := adal.NewServicePrincipalTokenFromManualToken(
		oauthConfig,
		applicationID,
		resource,
		*token,
		callbacks...)
	if err != nil {
		return nil, err
	}
	return spt, spt.Refresh()
}

func saveToken(spt adal.Token) error {
	err := adal.SaveToken(defaultTokenCachePath(), 0600, spt)
	if err != nil {
		return err
	}

	log.Printf("Acquired token was saved in '%s' file\n", defaultTokenCachePath())
	return nil
}

func acquireAuthTokenDeviceFlow() (string, error) {
	oauthConfig, err := adal.NewOAuthConfig(activeDirectoryEndpoint, "common")
	if err != nil {
		panic(err)
	}

	callback := func(token adal.Token) error {
		return saveToken(token)
	}

	if _, err := os.Stat(defaultTokenCachePath()); err == nil {
		token, err := adal.LoadToken(defaultTokenCachePath())
		if err != nil {
			return "", err
		}

		var spt *adal.ServicePrincipalToken
		if token.IsExpired() {
			spt, err = refreshToken(*oauthConfig, clientAppID, armResource, defaultTokenCachePath(), callback)
			if err != nil {
				return "", err
			}

			return fmt.Sprintf("%s %s", spt.Type, spt.AccessToken), nil
		}

		return fmt.Sprintf("%s %s", token.Type, token.AccessToken), nil
	}

	var spt *adal.ServicePrincipalToken
	spt, err = acquireTokenDeviceCodeFlow(
		*oauthConfig,
		clientAppID,
		armResource,
		callback)
	if err == nil {
		err = saveToken(spt.Token)
	}

	return fmt.Sprintf("%s %s", spt.Type, spt.AccessToken), nil
}

func acquireAuthTokenMSI() (string, error) {
	msiendpoint, _ := url.Parse(msiEndpoint)

	parameters := url.Values{}
	parameters.Add("resource", armResource)

	msiendpoint.RawQuery = parameters.Encode()

	req, err := http.NewRequest("GET", msiendpoint.String(), nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Metadata", "true")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	responseBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}

	var r responseJSON
	err = json.Unmarshal(responseBytes, &r)
	if err != nil {
		return "", err
	}

	return r.TokenType + " " + r.AccessToken, nil
}

func acquireAuthToken() (string, error) {
	_, isCloudShell := os.LookupEnv("ACC_CLOUD")

	if !isCloudShell {
		token, err := acquireAuthTokenDeviceFlow()
		if err != nil {
			return "", err
		}

		return token, nil
	}

	token, err := acquireAuthTokenMSI()
	if err != nil {
		return "", err
	}

	return token, nil
}
