package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Config struct {
	BaseURL            string
	Username           string
	Password           string
	AuthType           string
	InsecureSkipVerify bool
}

type Client struct {
	httpClient *http.Client
	baseURL    string
	token      string
	username   string
	password   string
	authType   string
}

func NewClient(ctx context.Context, cfg Config) (*Client, error) {
	httpClient := &http.Client{
		Timeout:   60 * time.Second,
		Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.InsecureSkipVerify}},
	}

	c := &Client{
		httpClient: httpClient,
		baseURL:    strings.TrimRight(cfg.BaseURL, "/"),
		username:   cfg.Username,
		password:   cfg.Password,
		authType:   cfg.AuthType,
	}

	if err := c.authenticate(ctx, cfg.Username, cfg.Password, cfg.AuthType); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) authenticate(ctx context.Context, username, password, authType string) error {
	payload := map[string]string{
		"userName": username,
		"password": password,
		"authType": authType,
	}

	var resp authResponse
	if err := c.doRequest(ctx, http.MethodPost, "/admin/v1/authentication", payload, &resp, false); err != nil {
		return err
	}

	c.token = resp.AuthToken
	return nil
}

func (c *Client) GetServer(ctx context.Context) (*Server, error) {
	var resp serverResponse
	if err := c.doRequest(ctx, http.MethodGet, "/admin/v2/server", nil, &resp, true); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) UpdateServerSMTP(ctx context.Context, smtp SMTPSettings) (*Server, error) {
	req := serverPatchRequest{
		Data: serverPatchData{
			Type: "server",
			Attributes: serverPatchAttributes{
				SMTP: smtp,
			},
		},
	}

	var resp serverResponse
	if err := c.doRequest(ctx, http.MethodPatch, "/admin/v2/server", req, &resp, true); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) ListSites(ctx context.Context) ([]Site, error) {
	var resp sitesResponse
	if err := c.doRequest(ctx, http.MethodGet, "/admin/v2/sites", nil, &resp, true); err != nil {
		return nil, err
	}

	return resp.Data, nil
}

func (c *Client) GetSiteUser(ctx context.Context, siteID, userID string) (*User, error) {
	var resp userResponse
	path := fmt.Sprintf("/admin/v2/sites/%s/users/%s", siteID, userID)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp, true); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) CreateSiteUser(ctx context.Context, siteID string, attrs UserAttributes) (*User, error) {
	req := userRequest{
		Data: userData{
			Type:       "user",
			Attributes: attrs,
		},
	}

	var resp userResponse
	path := fmt.Sprintf("/admin/v2/sites/%s/users", siteID)
	if err := c.doRequest(ctx, http.MethodPost, path, req, &resp, true); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) UpdateSiteUser(ctx context.Context, siteID, userID string, attrs UserAttributes) (*User, error) {
	req := userRequest{
		Data: userData{
			Type:       "user",
			Attributes: attrs,
		},
	}

	var resp userResponse
	path := fmt.Sprintf("/admin/v2/sites/%s/users/%s", siteID, userID)
	if err := c.doRequest(ctx, http.MethodPatch, path, req, &resp, true); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

func (c *Client) DeleteSiteUser(ctx context.Context, siteID, userID string) error {
	path := fmt.Sprintf("/admin/v2/sites/%s/users/%s", siteID, userID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil, true)
}

func (c *Client) GetEventRule(ctx context.Context, siteID, ruleID string) (*EventRule, error) {
	var resp eventRuleResponse
	path := fmt.Sprintf("/admin/v2/sites/%s/event-rules/%s", siteID, ruleID)
	if err := c.doRequest(ctx, http.MethodGet, path, nil, &resp, true); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) CreateEventRule(ctx context.Context, siteID string, data EventRuleRequestData) (*EventRule, error) {
	req := eventRuleRequest{Data: data}
	var resp eventRuleResponse
	path := fmt.Sprintf("/admin/v2/sites/%s/event-rules", siteID)
	if err := c.doRequest(ctx, http.MethodPost, path, req, &resp, true); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) UpdateEventRule(ctx context.Context, siteID, ruleID string, data EventRuleRequestData) (*EventRule, error) {
	req := eventRuleRequest{Data: data}
	var resp eventRuleResponse
	path := fmt.Sprintf("/admin/v2/sites/%s/event-rules/%s", siteID, ruleID)
	if err := c.doRequest(ctx, http.MethodPatch, path, req, &resp, true); err != nil {
		return nil, err
	}
	return &resp.Data, nil
}

func (c *Client) DeleteEventRule(ctx context.Context, siteID, ruleID string) error {
	path := fmt.Sprintf("/admin/v2/sites/%s/event-rules/%s", siteID, ruleID)
	return c.doRequest(ctx, http.MethodDelete, path, nil, nil, true)
}

type authResponse struct {
	AuthToken string `json:"authToken"`
}

type serverPatchRequest struct {
	Data serverPatchData `json:"data"`
}

type serverPatchData struct {
	Type       string                `json:"type"`
	Attributes serverPatchAttributes `json:"attributes"`
}

type serverPatchAttributes struct {
	SMTP SMTPSettings `json:"smtp"`
}

func (c *Client) doRequest(ctx context.Context, method, path string, body any, dest any, includeAuth bool) error {
	var bodyBytes []byte
	var err error
	if body != nil {
		buf := &bytes.Buffer{}
		if err = json.NewEncoder(buf).Encode(body); err != nil {
			return err
		}
		bodyBytes = buf.Bytes()
	}

	makeRequest := func() (*http.Response, error) {
		var bodyReader io.Reader
		if len(bodyBytes) > 0 {
			bodyReader = bytes.NewReader(bodyBytes)
		}

		url := c.baseURL + "/" + strings.TrimLeft(path, "/")
		req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
		if err != nil {
			return nil, err
		}

		if body != nil {
			req.Header.Set("Content-Type", "application/json")
		}

		if includeAuth {
			req.Header.Set("Authorization", fmt.Sprintf("EFTAdminAuthToken %s", c.token))
		}

		return c.httpClient.Do(req)
	}

	resp, err := makeRequest()
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusUnauthorized && includeAuth {
		resp.Body.Close()
		if err := c.authenticate(ctx, c.username, c.password, c.authType); err != nil {
			return err
		}
		resp, err = makeRequest()
		if err != nil {
			return err
		}
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("globalscape EFT API %s %s failed: %s", method, path, strings.TrimSpace(string(raw)))
	}

	if dest == nil {
		io.Copy(io.Discard, resp.Body)
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(dest)
}

// Server models provide just the fields that are surfaced to Terraform.
type Server struct {
	Type       string           `json:"type"`
	ID         string           `json:"id"`
	Attributes ServerAttributes `json:"attributes"`
}

type ServerAttributes struct {
	Version          string           `json:"version"`
	General          ServerGeneral    `json:"general"`
	ListenerSettings ListenerSettings `json:"listenerSettings"`
	SMTP             SMTPSettings     `json:"smtp"`
}

type ServerGeneral struct {
	ConfigFilePath       string `json:"configFilePath"`
	EnableUtcInListings  bool   `json:"enableUtcInListings"`
	LastModifiedBy       string `json:"lastModifiedBy"`
	LastModifiedUnixTime int64  `json:"lastModifiedTime"`
}

type ListenerSettings struct {
	AdminPort                  int64    `json:"adminPort"`
	EnableRemoteAdministration bool     `json:"enableRemoteAdministration"`
	ListenIPs                  []string `json:"listenIps"`
}

type SMTPSettings struct {
	Login             string `json:"login"`
	Password          string `json:"password"`
	Port              int64  `json:"port"`
	SenderAddress     string `json:"senderAddr"`
	SenderName        string `json:"senderName"`
	Server            string `json:"server"`
	UseAuthentication bool   `json:"useAuthentication"`
	UseImplicitTLS    bool   `json:"useImplicitTLS"`
}

type serverResponse struct {
	Data Server `json:"data"`
}

type sitesResponse struct {
	Data []Site `json:"data"`
}

type Site struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes SiteAttributes `json:"attributes"`
}

type SiteAttributes struct {
	Name string `json:"name"`
}

type userResponse struct {
	Data User `json:"data"`
}

type userRequest struct {
	Data userData `json:"data"`
}

type userData struct {
	Type       string         `json:"type"`
	Attributes UserAttributes `json:"attributes"`
}

type User struct {
	Type       string         `json:"type"`
	ID         string         `json:"id"`
	Attributes UserAttributes `json:"attributes"`
}

type UserAttributes struct {
	LoginName            string             `json:"loginName"`
	AccountEnabled       string             `json:"accountEnabled,omitempty"`
	Password             *UserPassword      `json:"password,omitempty"`
	Personal             *UserPersonal      `json:"personal,omitempty"`
	HomeFolder           *UserHomeFolder    `json:"homeFolder,omitempty"`
	HasHomeFolderAsRoot  string             `json:"hasHomeFolderAsRoot,omitempty"`
	AgreementToTerms     string             `json:"agreementToTermsOfService,omitempty"`
	ConsentToPrivacy     string             `json:"consentToPrivacyPolicy,omitempty"`
	IsEUDataSubject      string             `json:"isEuDataSubject,omitempty"`
	ExternalAuth         string             `json:"externalAuthentication,omitempty"`
	ChangePasswordPolicy *ChangePasswordSet `json:"changePassword,omitempty"`
}

type UserPassword struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type UserPersonal struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Email       string `json:"email,omitempty"`
}

type UserHomeFolder struct {
	Enabled string               `json:"enabled,omitempty"`
	Value   *UserHomeFolderValue `json:"value,omitempty"`
}

type UserHomeFolderValue struct {
	Path string `json:"path,omitempty"`
}

type ChangePasswordSet struct {
	Enabled string                   `json:"enabled,omitempty"`
	Value   *ChangePasswordSetValues `json:"value,omitempty"`
}

type ChangePasswordSetValues struct {
	MustChangePassword bool `json:"mustChangePassword,omitempty"`
}

type EventRule struct {
	Type          string          `json:"type"`
	ID            string          `json:"id"`
	Attributes    json.RawMessage `json:"attributes"`
	Relationships json.RawMessage `json:"relationships,omitempty"`
}

type eventRuleResponse struct {
	Data EventRule `json:"data"`
}

type eventRuleRequest struct {
	Data EventRuleRequestData `json:"data"`
}

type EventRuleRequestData struct {
	Type          string          `json:"type"`
	ID            string          `json:"id,omitempty"`
	Attributes    json.RawMessage `json:"attributes"`
	Relationships json.RawMessage `json:"relationships,omitempty"`
}
