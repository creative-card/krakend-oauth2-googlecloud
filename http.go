package oauth2googleclient

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/oauth2/jwt"

	"github.com/luraproject/lura/v2/config"
	"github.com/luraproject/lura/v2/transport/http/client"
)

// Namespace is the key to use to store and access the custom config data
const Namespace = "github.com/creative-card/krakend-oauth2-googlecloud"

// NewHTTPClient creates a HTTPClientFactory with an http client configured for dealing
// with all the logic related to the oauth2 Google credentials grant
func NewHTTPClient(cfg *config.Backend) client.HTTPClientFactory {
	oauth, ok := configGetter(cfg.ExtraConfig).(Config)
	if !ok || oauth.IsDisabled {
		return client.NewHTTPClient
	}
	c := &jwt.Config{
		Email:      	  oauth.Email,
		PrivateKey:     []byte(oauth.PrivateKey),
		TokenURL:       oauth.TokenURL,
		Scopes:         strings.Split(oauth.Scopes, ","),
	}
	cli := c.Client(context.Background())
	return func(_ context.Context) *http.Client {
		return cli
	}
}

// Config is the custom config struct containing the params for the auth/googlecloud package
type Config struct {
	IsDisabled     bool
	Email          string
	PrivateKey     string
	TokenURL       string
	Scopes         string
}

// ZeroCfg is the zero value for the Config struct
var ZeroCfg = Config{}

func configGetter(e config.ExtraConfig) interface{} {
	v, ok := e[Namespace]
	if !ok {
		return nil
	}
	tmp, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	cfg := Config{}
	if v, ok := tmp["is_disabled"]; ok {
		cfg.IsDisabled = v.(bool)
	}
	if v, ok := tmp["email"]; ok {
		cfg.Email = v.(string)
	}
	if v, ok := tmp["private_key"]; ok {
		cfg.PrivateKey = v.(string)
	}
	if v, ok := tmp["token_url"]; ok {
		cfg.TokenURL = v.(string)
	}
	if v, ok := tmp["scopes"]; ok {
		cfg.Scopes = v.(string)
	}
	return cfg
}
