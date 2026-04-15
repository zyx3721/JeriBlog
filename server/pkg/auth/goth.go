package auth

import (
	"flec_blog/config"
	"flec_blog/pkg/auth/providers/qq"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/markbates/goth/providers/microsoftonline"
)

// workerProxy Cloudflare Worker 的代理地址
const workerProxy = "https://proxy.flec.top"

// UpdateConfig 根据配置动态更新 OAuth 配置
func UpdateConfig(cfg *config.OAuthConfig) {
	if cfg.SessionSecret != "" {
		store := sessions.NewCookieStore([]byte(cfg.SessionSecret))
		store.MaxAge(86400 * 30)
		store.Options.Path = "/"
		store.Options.HttpOnly = true
		store.Options.Secure = true
		store.Options.SameSite = http.SameSiteNoneMode
		gothic.Store = store
	}

	httpClient := &http.Client{
		Transport: &workerTransport{base: http.DefaultTransport},
		Timeout:   30 * time.Second,
	}

	var providers []goth.Provider
	if cfg.Github.Enabled && cfg.Github.ClientID != "" {
		p := github.New(cfg.Github.ClientID, cfg.Github.ClientSecret, cfg.Github.RedirectURL)
		p.HTTPClient = httpClient
		providers = append(providers, p)
	}
	if cfg.Google.Enabled && cfg.Google.ClientID != "" {
		p := google.New(cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURL)
		p.HTTPClient = httpClient
		providers = append(providers, p)
	}
	// QQ 不走代理，直接使用默认 HTTP client
	if cfg.QQ.Enabled && cfg.QQ.ClientID != "" {
		p := qq.New(cfg.QQ.ClientID, cfg.QQ.ClientSecret, cfg.QQ.RedirectURL)
		providers = append(providers, p)
	}
	if cfg.Microsoft.Enabled && cfg.Microsoft.ClientID != "" {
		p := microsoftonline.New(cfg.Microsoft.ClientID, cfg.Microsoft.ClientSecret, cfg.Microsoft.RedirectURL)
		p.SetName("microsoft")
		p.HTTPClient = httpClient
		providers = append(providers, p)
	}

	if len(providers) > 0 {
		goth.UseProviders(providers...)
	} else {
		goth.ClearProviders()
	}
}

// workerTransport 将 OAuth API 请求转发到 Cloudflare Worker
type workerTransport struct {
	base http.RoundTripper
}

// RoundTrip 将请求转发到 Cloudflare Worker
func (t *workerTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	prefixes := map[string]string{
		"github.com":                "/github",
		"api.github.com":            "/github-api",
		"accounts.google.com":       "/google",
		"www.googleapis.com":        "/google-api",
		"oauth2.googleapis.com":     "/google-oauth2",
		"login.microsoftonline.com": "/microsoft",
		"graph.microsoft.com":       "/microsoft-graph",
	}

	prefix, ok := prefixes[req.URL.Host]
	if !ok {
		return t.base.RoundTrip(req)
	}

	newReq := req.Clone(req.Context())
	workerURL, _ := url.Parse(workerProxy)

	newReq.URL.Scheme = workerURL.Scheme
	newReq.URL.Host = workerURL.Host
	newReq.URL.Path = prefix + req.URL.Path
	newReq.Host = workerURL.Host

	return t.base.RoundTrip(newReq)
}
