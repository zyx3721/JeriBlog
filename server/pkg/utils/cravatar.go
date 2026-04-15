package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// GetEmailHash 计算邮箱的 MD5 哈希
func GetEmailHash(email string) string {
	email = strings.TrimSpace(strings.ToLower(email))
	hash := md5.Sum([]byte(email))
	return hex.EncodeToString(hash[:])
}

// DownloadCravatarAvatar 下载 Cravatar 头像
func DownloadCravatarAvatar(email string) (io.Reader, error) {
	emailHash := GetEmailHash(email)
	url := fmt.Sprintf("https://cravatar.cn/avatar/%s?s=200&d=robohash", emailHash)
	return DownloadRemoteImage(url)
}

// DownloadRemoteImage 下载远程图片
func DownloadRemoteImage(url string) (io.Reader, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(data), nil
}
