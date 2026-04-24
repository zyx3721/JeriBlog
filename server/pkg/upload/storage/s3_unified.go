/*
项目名称：JeriBlog
文件名称：s3_unified.go
创建时间：2026-04-16 14:59:17

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：S3 兼容存储实现
*/

package storage

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"jeri_blog/config"
)

// S3UnifiedStorage 基于MinIO SDK的统一S3存储
type S3UnifiedStorage struct {
	cfg         config.UploadConfig
	storageType string
	client      *minio.Client
	bucketName  string
	baseURL     string
	mu          sync.Mutex
}

// NewS3UnifiedStorage 创建统一S3存储实例
func NewS3UnifiedStorage(cfg config.UploadConfig, storageType string) (*S3UnifiedStorage, error) {
	// 根据存储类型获取对应的 bucket
	var bucketName string
	switch storageType {
	case "s3":
		bucketName = cfg.S3.Bucket
	case "cos":
		bucketName = cfg.COS.Bucket
	case "oss":
		bucketName = cfg.OSS.Bucket
	case "kodo":
		bucketName = cfg.Kodo.Bucket
	case "r2":
		bucketName = cfg.R2.Bucket
	case "minio":
		bucketName = cfg.MinIO.Bucket
	}

	storage := &S3UnifiedStorage{
		cfg:         cfg,
		storageType: storageType,
		bucketName:  bucketName,
	}
	return storage, nil
}

// ensureClient 确保客户端已初始化
func (s *S3UnifiedStorage) ensureClient() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.client != nil {
		return nil
	}

	// 根据存储类型获取对应的配置
	var accessKey, secretKey, region, endpoint, domain string
	var useSSL bool

	switch s.storageType {
	case "s3":
		accessKey = s.cfg.S3.AccessKey
		secretKey = s.cfg.S3.SecretKey
		region = s.cfg.S3.Region
		endpoint = s.cfg.S3.Endpoint
		domain = s.cfg.S3.Domain
		useSSL = true
		if region == "" {
			return fmt.Errorf("AWS S3 需要配置 region")
		}
		if endpoint == "" {
			endpoint = fmt.Sprintf("s3.%s.amazonaws.com", region)
		}
	case "cos":
		accessKey = s.cfg.COS.SecretID
		secretKey = s.cfg.COS.SecretKey
		region = s.cfg.COS.Region
		domain = s.cfg.COS.Domain
		useSSL = true
		if region == "" {
			return fmt.Errorf("腾讯云 COS 需要配置 region")
		}
		endpoint = fmt.Sprintf("cos.%s.myqcloud.com", region)
	case "oss":
		accessKey = s.cfg.OSS.AccessKey
		secretKey = s.cfg.OSS.SecretKey
		region = s.cfg.OSS.Region
		domain = s.cfg.OSS.Domain
		useSSL = true
		if region == "" {
			return fmt.Errorf("阿里云 OSS 需要配置 region")
		}
		endpoint = fmt.Sprintf("oss-%s.aliyuncs.com", region)
	case "kodo":
		accessKey = s.cfg.Kodo.AccessKey
		secretKey = s.cfg.Kodo.SecretKey
		region = s.cfg.Kodo.Region
		domain = s.cfg.Kodo.Domain
		useSSL = true
		if region == "" {
			return fmt.Errorf("七牛云 Kodo 需要配置 region")
		}
		endpoint = fmt.Sprintf("s3.%s.qiniucs.com", region)
	case "r2":
		accessKey = s.cfg.R2.AccessKey
		secretKey = s.cfg.R2.SecretKey
		endpoint = s.cfg.R2.Endpoint
		domain = s.cfg.R2.Domain
		useSSL = s.cfg.R2.UseSSL
		region = "auto"
		if endpoint == "" {
			return fmt.Errorf("Cloudflare R2 需要配置 endpoint")
		}
	case "minio":
		accessKey = s.cfg.MinIO.AccessKey
		secretKey = s.cfg.MinIO.SecretKey
		region = s.cfg.MinIO.Region
		endpoint = s.cfg.MinIO.Endpoint
		domain = s.cfg.MinIO.Domain
		useSSL = s.cfg.MinIO.UseSSL
		if region == "" {
			region = "us-east-1"
		}
		if endpoint == "" {
			return fmt.Errorf("MinIO 需要配置 endpoint")
		}
	default:
		return fmt.Errorf("不支持的存储类型: %s", s.storageType)
	}

	// 验证必需参数
	if accessKey == "" || secretKey == "" {
		return fmt.Errorf("存储类型 %s 需要配置访问密钥", s.storageType)
	}
	if s.bucketName == "" {
		return fmt.Errorf("存储类型 %s 需要配置存储桶名称", s.storageType)
	}

	// 设置 Bucket 查找方式
	bucketLookup := minio.BucketLookupAuto
	if s.storageType == "cos" || s.storageType == "oss" || s.storageType == "kodo" {
		bucketLookup = minio.BucketLookupDNS
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:        credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure:       useSSL,
		Region:       region,
		BucketLookup: bucketLookup,
	})
	if err != nil {
		return fmt.Errorf("创建存储实例失败: %w", err)
	}

	s.client = client
	s.baseURL = buildBaseURL(endpoint, s.bucketName, useSSL, domain, s.storageType)
	return nil
}

// buildObjectKey 构建完整的对象键
func (s *S3UnifiedStorage) buildObjectKey(filePath string) string {
	return strings.TrimPrefix(filePath, "/")
}

// buildBaseURL 构建基础访问URL
func buildBaseURL(endpoint, bucket string, useSSL bool, customDomain, storageType string) string {
	if customDomain != "" {
		return strings.TrimSuffix(customDomain, "/")
	}

	scheme := "http"
	if useSSL {
		scheme = "https"
	}

	// COS 使用 virtual-hosted style：bucket.cos.region.myqcloud.com
	if storageType == "cos" {
		return fmt.Sprintf("%s://%s.%s", scheme, bucket, endpoint)
	}

	return fmt.Sprintf("%s://%s/%s", scheme, endpoint, bucket)
}

// Save 实现 Storage 接口 - 保存文件
func (s *S3UnifiedStorage) Save(reader io.Reader, filePath string, size int64) error {
	if err := s.ensureClient(); err != nil {
		return err
	}
	objectKey := s.buildObjectKey(filePath)
	_, err := s.client.PutObject(context.Background(), s.bucketName, objectKey, reader, size, minio.PutObjectOptions{
		ContentType: getContentType(filePath),
	})
	return err
}

// Delete 实现 Storage 接口 - 删除文件
func (s *S3UnifiedStorage) Delete(filePath string) error {
	if err := s.ensureClient(); err != nil {
		return err
	}
	objectKey := s.buildObjectKey(filePath)
	return s.client.RemoveObject(context.Background(), s.bucketName, objectKey, minio.RemoveObjectOptions{})
}

// GetURL 实现 Storage 接口 - 获取文件访问 URL
func (s *S3UnifiedStorage) GetURL(filePath string, _ string) string {
	objectKey := s.buildObjectKey(filePath)

	if s.baseURL != "" {
		return fmt.Sprintf("%s/%s", s.baseURL, objectKey)
	}

	endpointURL := s.client.EndpointURL()
	return fmt.Sprintf("%s://%s/%s", endpointURL.Scheme, endpointURL.Host, objectKey)
}

// GetPreSignedURL 获取预签名URL（用于临时访问）
func (s *S3UnifiedStorage) GetPreSignedURL(filePath string, expires time.Duration) (string, error) {
	objectKey := s.buildObjectKey(filePath)
	presignedURL, err := s.client.PresignedGetObject(context.Background(), s.bucketName, objectKey, expires, nil)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

// BatchDelete 批量删除文件
func (s *S3UnifiedStorage) BatchDelete(filePaths []string) error {
	if len(filePaths) == 0 {
		return nil
	}

	ctx := context.Background()
	objectsCh := make(chan minio.ObjectInfo)

	go func() {
		defer close(objectsCh)
		for _, filePath := range filePaths {
			objectKey := s.buildObjectKey(filePath)
			objectsCh <- minio.ObjectInfo{Key: objectKey}
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}
	errorCh := s.client.RemoveObjects(ctx, s.bucketName, objectsCh, opts)

	var errors []string
	for err := range errorCh {
		if err.Err != nil {
			errors = append(errors, fmt.Sprintf("删除 %s 失败: %v", err.ObjectName, err.Err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("批量删除部分失败: %s", strings.Join(errors, "; "))
	}

	return nil
}

// Exists 实现 Storage 接口 - 检查文件是否存在
func (s *S3UnifiedStorage) Exists(filePath string) bool {
	objectKey := s.buildObjectKey(filePath)
	_, err := s.client.StatObject(context.Background(), s.bucketName, objectKey, minio.StatObjectOptions{})
	return err == nil
}

// HealthCheck 检查存储可用性
func (s *S3UnifiedStorage) HealthCheck() error {
	if err := s.ensureClient(); err != nil {
		return fmt.Errorf("初始化存储客户端失败: %w", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err := s.client.BucketExists(ctx, s.bucketName)
	if err != nil {
		return fmt.Errorf("存储桶不可访问: %w", err)
	}
	return nil
}

// GetObjectInfo 获取对象信息
func (s *S3UnifiedStorage) GetObjectInfo(filePath string) (*ObjectInfo, error) {
	objectKey := s.buildObjectKey(filePath)
	objInfo, err := s.client.StatObject(context.Background(), s.bucketName, objectKey, minio.StatObjectOptions{})
	if err != nil {
		return nil, err
	}

	return &ObjectInfo{
		Key:          objInfo.Key,
		Size:         objInfo.Size,
		LastModified: objInfo.LastModified,
		ContentType:  objInfo.ContentType,
		ETag:         objInfo.ETag,
	}, nil
}

// ObjectInfo 对象信息
type ObjectInfo struct {
	Key          string    `json:"key"`
	Size         int64     `json:"size"`
	LastModified time.Time `json:"last_modified"`
	ContentType  string    `json:"content_type"`
	ETag         string    `json:"etag"`
}

// getContentType 根据文件扩展名获取Content-Type
func getContentType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))

	contentTypes := map[string]string{
		".jpg":  "image/jpeg",
		".jpeg": "image/jpeg",
		".png":  "image/png",
		".gif":  "image/gif",
		".webp": "image/webp",
		".svg":  "image/svg+xml",
		".mp4":  "video/mp4",
		".webm": "video/webm",
		".pdf":  "application/pdf",
		".txt":  "text/plain",
		".json": "application/json",
		".xml":  "application/xml",
		".css":  "text/css",
		".js":   "application/javascript",
		".html": "text/html",
	}

	if contentType, exists := contentTypes[ext]; exists {
		return contentType
	}

	return "application/octet-stream"
}
