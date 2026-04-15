package wechatmp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	defaultAPIBase  = "https://api.weixin.qq.com"
	accessTokenPath = "/cgi-bin/token"
	uploadImagePath = "/cgi-bin/media/uploadimg"
	addMaterialPath = "/cgi-bin/material/add_material"
	createDraftPath = "/cgi-bin/draft/add"
	tokenGuardAhead = 60 * time.Second
)

// Config 微信公众号客户端配置
type Config struct {
	AppID      string
	AppSecret  string
	BaseURL    string
	HTTPClient *http.Client
}

// Client 微信公众号 API 客户端，封装草稿箱和素材库相关接口
type Client struct {
	appID      string
	appSecret  string
	baseURL    string
	httpClient *http.Client

	mu          sync.Mutex
	accessToken string
	tokenExpiry time.Time
}

// NewClient 根据配置创建微信公众号客户端实例
func NewClient(cfg Config) (*Client, error) {
	if strings.TrimSpace(cfg.AppID) == "" {
		return nil, fmt.Errorf("微信 app_id 不能为空")
	}
	if strings.TrimSpace(cfg.AppSecret) == "" {
		return nil, fmt.Errorf("微信 app_secret 不能为空")
	}

	httpClient := cfg.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}

	baseURL := strings.TrimRight(cfg.BaseURL, "/")
	if baseURL == "" {
		baseURL = defaultAPIBase
	}

	return &Client{
		appID:      cfg.AppID,
		appSecret:  cfg.AppSecret,
		baseURL:    baseURL,
		httpClient: httpClient,
	}, nil
}

// DraftArticle 微信草稿文章结构
type DraftArticle struct {
	Title            string
	Author           string
	Content          string
	Digest           string
	ContentSourceURL string
	ThumbMediaID     string
	NeedOpenComment  int
}

// UploadImageResult 上传图片接口返回结果
type UploadImageResult struct {
	URL string `json:"url"`
}

// MaterialResult 新增素材接口返回结果
type MaterialResult struct {
	MediaID string `json:"media_id"`
	URL     string `json:"url"`
}

// DraftResult 创建草稿接口返回结果
type DraftResult struct {
	MediaID string `json:"media_id"`
}

// UploadImage 上传文章内图片，返回微信 CDN 地址
func (c *Client) UploadImage(ctx context.Context, filename string, data []byte) (*UploadImageResult, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	values := url.Values{"access_token": []string{token}}
	endpoint := c.buildURL(uploadImagePath, values)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return nil, fmt.Errorf("创建 multipart 表单失败: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("写入 multipart 数据失败: %w", err)
	}
	_ = writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp UploadImageResult
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// AddThumbMaterial 上传封面图到素材库
func (c *Client) AddThumbMaterial(ctx context.Context, filename string, data []byte) (*MaterialResult, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	values := url.Values{
		"access_token": []string{token},
		"type":         []string{"image"},
	}
	endpoint := c.buildURL(addMaterialPath, values)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("media", filename)
	if err != nil {
		return nil, fmt.Errorf("创建 multipart 表单失败: %w", err)
	}
	if _, err := part.Write(data); err != nil {
		return nil, fmt.Errorf("写入 multipart 数据失败: %w", err)
	}
	_ = writer.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	var resp MaterialResult
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}

	// 验证返回的 media_id 是否有效
	if resp.MediaID == "" {
		return nil, fmt.Errorf("微信返回的封面素材 media_id 为空")
	}

	return &resp, nil
}

// CreateDraft 将一篇或多篇文章推送到微信草稿箱
func (c *Client) CreateDraft(ctx context.Context, articles []DraftArticle) (*DraftResult, error) {
	if len(articles) == 0 {
		return nil, fmt.Errorf("至少需要一篇文章")
	}

	token, err := c.getAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"articles": convertDraftArticles(articles),
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("序列化草稿数据失败: %w", err)
	}

	endpoint := c.buildURL(createDraftPath, url.Values{"access_token": []string{token}})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	var resp DraftResult
	if err := c.doRequest(req, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// convertDraftArticles 将草稿文章转换为微信 API 格式
func convertDraftArticles(items []DraftArticle) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(items))
	for _, item := range items {
		result = append(result, map[string]interface{}{
			"title":              item.Title,
			"author":             item.Author,
			"content":            item.Content,
			"digest":             item.Digest,
			"content_source_url": item.ContentSourceURL,
			"thumb_media_id":     item.ThumbMediaID,
			"need_open_comment":  item.NeedOpenComment,
		})
	}
	return result
}

// buildURL 构建完整的 API 请求地址
func (c *Client) buildURL(path string, values url.Values) string {
	base := c.baseURL + path
	if len(values) == 0 {
		return base
	}
	return base + "?" + values.Encode()
}

// doRequest 执行 HTTP 请求并解析响应
func (c *Client) doRequest(req *http.Request, out interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("调用微信 API 失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("读取微信响应失败: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("微信 API 状态码 %d: %s", resp.StatusCode, string(body))
	}

	var apiErr apiError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.ErrCode != 0 {
		return apiErr
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("解析微信响应失败: %w", err)
	}
	return nil
}

// getAccessToken 获取或刷新 access_token
func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.accessToken != "" && time.Now().Add(tokenGuardAhead).Before(c.tokenExpiry) {
		return c.accessToken, nil
	}

	values := url.Values{
		"grant_type": []string{"client_credential"},
		"appid":      []string{c.appID},
		"secret":     []string{c.appSecret},
	}

	endpoint := c.buildURL(accessTokenPath, values)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return "", err
	}

	var resp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := c.doRequest(req, &resp); err != nil {
		return "", err
	}
	if resp.AccessToken == "" {
		return "", fmt.Errorf("微信返回的 access_token 为空")
	}

	c.accessToken = resp.AccessToken
	if resp.ExpiresIn > 0 {
		c.tokenExpiry = time.Now().Add(time.Duration(resp.ExpiresIn) * time.Second)
	} else {
		c.tokenExpiry = time.Now().Add(2 * time.Hour)
	}
	return c.accessToken, nil
}

// apiError 微信 API 错误响应
type apiError struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func (e apiError) Error() string {
	return fmt.Sprintf("微信 API 错误 %d: %s", e.ErrCode, e.ErrMsg)
}
