package main

import (
	"net/http"
	oauth2 "github.com/kanosaki/pixiv_oauth2" // modified for pixiv
	"golang.org/x/net/context"
	"io/ioutil"
	"encoding/json"
	"net/http/cookiejar"
	log "github.com/Sirupsen/logrus"
)

type Config struct {
	ClientID     string  `json:"client_id"`
	ClientSecret string  `json:"client_secret"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
}

func CreateOAuthConfig(clientId string, clientSecret string) *oauth2.Config {
	return &oauth2.Config{
		ClientID: clientId,
		ClientSecret: clientSecret,
		Scopes: []string{},
		Endpoint: oauth2.Endpoint{
			AuthURL: "", // never used
			TokenURL: "https://oauth.secure.pixiv.net/auth/token",
		},
	}
}

type Pixiv struct {
	AuthConnection *http.Client
	Token          *oauth2.Token
	Config         *Config
}

func NewFromConfigFile(path string) (*Pixiv, error) {
	configJson, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config *Config
	err = json.Unmarshal(configJson, &config)
	if err != nil {
		return nil, err
	}
	return New(config), nil
}

func New(config *Config) *Pixiv {
	return &Pixiv{
		AuthConnection: nil,
		Token: nil,
		Config: config,
	}
}

func (px *Pixiv) FetchToken(ctx context.Context) error {
	log.Debugf("Login: %s", px.Config.Username)
	config := px.Config
	oauthConfig := CreateOAuthConfig(config.ClientID, config.ClientSecret)
	jar, err := cookiejar.New(nil)
	if err != nil {
		return err
	}
	oauthClient := &http.Client{
		Jar: jar,
	}
	oauthCtx := context.WithValue(ctx, oauth2.HTTPClient, oauthClient)
	token, err := oauthConfig.PasswordCredentialsToken(oauthCtx, config.Username, config.Password)
	if err != nil {
		return err
	}
	px.AuthConnection = oauthConfig.Client(context.Background(), token)
	px.AuthConnection.Jar = oauthClient.Jar
	px.Token = token
	return nil
}

// IsAuthorized != oauth2 authorize completed, IsAuthorized == oauth2 token ready.
func (px *Pixiv) IsAuthorized() bool {
	return px.Token != nil && px.Token.Valid() && px.AuthConnection != nil
}

func (px *Pixiv) AuthClient() (*http.Client, error) {
	if ! px.IsAuthorized() {
		err := px.FetchToken(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return px.AuthConnection, nil
}

func (px *Pixiv) PlainClient() (*http.Client, error) {
	if ! px.IsAuthorized() {
		err := px.FetchToken(context.Background())
		if err != nil {
			return nil, err
		}
	}
	return &http.Client{
		Jar: px.AuthConnection.Jar,
	}, nil
}
