/*
项目名称：JeriBlog
文件名称：markdown.go
创建时间：2026-04-26 11:15:30

系统用户：Jerion
作　　者：Jerion
联系邮箱：416685476@qq.com
功能描述：Markdown 内容解析工具，提取文件引用
*/

package utils

import (
	"regexp"
	"strings"
)

// FileType 文件类型
type FileType string

const (
	FileTypeImage FileType = "文章配图"
	FileTypeVideo FileType = "文章视频"
	FileTypeAudio FileType = "文章音频"
	FileTypeFile  FileType = "文章附件"
)

// FileReference 文件引用信息
type FileReference struct {
	URL  string
	Type FileType
}

// ExtractFileReferencesFromMarkdown 从 Markdown 内容中提取所有文件引用及其类型
func ExtractFileReferencesFromMarkdown(content string) []FileReference {
	if content == "" {
		return []FileReference{}
	}

	urlMap := make(map[string]FileType) // URL -> 文件类型
	var references []FileReference

	// 1. 提取标准 Markdown 图片：![alt](url) -> 文章配图
	markdownImageRegex := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
	matches := markdownImageRegex.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 && match[1] != "" {
			url := strings.TrimSpace(match[1])
			if _, exists := urlMap[url]; !exists {
				urlMap[url] = FileTypeImage
			}
		}
	}

	// 2. 提取照片墙中的图片：:::photo ... :::endphoto -> 文章配图
	photoWallRegex := regexp.MustCompile(`(?s):::photo\s+(.*?)\s+:::endphoto`)
	photoMatches := photoWallRegex.FindAllStringSubmatch(content, -1)
	for _, photoMatch := range photoMatches {
		if len(photoMatch) > 1 {
			photoContent := photoMatch[1]

			// 提取照片墙内的 Markdown 图片：![](url)
			photoImageRegex := regexp.MustCompile(`!\[.*?\]\((.*?)\)`)
			photoImageMatches := photoImageRegex.FindAllStringSubmatch(photoContent, -1)
			for _, imgMatch := range photoImageMatches {
				if len(imgMatch) > 1 && imgMatch[1] != "" {
					url := strings.TrimSpace(imgMatch[1])
					if _, exists := urlMap[url]; !exists {
						urlMap[url] = FileTypeImage
					}
				}
			}

			// 提取照片墙内的直接 URL
			lines := strings.Split(photoContent, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line == ":::n" || line == "" {
					continue
				}

				parts := strings.Fields(line)
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, "http://") ||
						strings.HasPrefix(part, "https://") ||
						strings.HasPrefix(part, "/") {
						if !strings.Contains(part, "](") {
							if _, exists := urlMap[part]; !exists {
								urlMap[part] = FileTypeImage
							}
						}
					}
				}
			}
		}
	}

	// 3. 提取视频 URL：:::video url ::: -> 文章视频
	// 修复：使用 [ \t] 替代 \s，避免匹配换行符导致连续视频块被错误合并
	videoRegex := regexp.MustCompile(`:::video[ \t]+([^\s]+)(?:[ \t]+([^\s]+))?[ \t]*:::`)
	videoMatches := videoRegex.FindAllStringSubmatch(content, -1)
	for _, videoMatch := range videoMatches {
		if len(videoMatch) > 1 {
			firstParam := strings.TrimSpace(videoMatch[1])
			if (strings.HasPrefix(firstParam, "http://") ||
				strings.HasPrefix(firstParam, "https://") ||
				strings.HasPrefix(firstParam, "/")) &&
				firstParam != "bilibili" &&
				firstParam != "youtube" {
				// 视频类型优先级高于图片
				urlMap[firstParam] = FileTypeVideo
			}
		}
	}

	// 4. 提取音频 URL：:::audio 标题 url ::: -> 文章音频
	audioRegex := regexp.MustCompile(`:::audio\s+[^\s]+\s+([^\s]+)\s+:::`)
	audioMatches := audioRegex.FindAllStringSubmatch(content, -1)
	for _, audioMatch := range audioMatches {
		if len(audioMatch) > 1 {
			audioURL := strings.TrimSpace(audioMatch[1])
			if strings.HasPrefix(audioURL, "http://") ||
				strings.HasPrefix(audioURL, "https://") ||
				strings.HasPrefix(audioURL, "/") {
				// 音频类型优先级高于图片
				urlMap[audioURL] = FileTypeAudio
			}
		}
	}

	// 5. 提取链接卡片中的本地文件链接：:::link 标题 url 描述 ::: -> 文章附件
	linkCardRegex := regexp.MustCompile(`:::link\s+[^\s]+\s+([^\s]+)(?:\s+.*?)?\s+:::`)
	linkMatches := linkCardRegex.FindAllStringSubmatch(content, -1)
	for _, linkMatch := range linkMatches {
		if len(linkMatch) > 1 {
			linkURL := strings.TrimSpace(linkMatch[1])
			if strings.HasPrefix(linkURL, "/") && strings.Contains(linkURL, ".") {
				parts := strings.Split(linkURL, "/")
				if len(parts) > 0 {
					lastPart := parts[len(parts)-1]
					if strings.Contains(lastPart, ".") {
						if _, exists := urlMap[linkURL]; !exists {
							urlMap[linkURL] = FileTypeFile
						}
					}
				}
			}
		}
	}

	// 转换为数组
	for url, fileType := range urlMap {
		references = append(references, FileReference{
			URL:  url,
			Type: fileType,
		})
	}

	return references
}

// ExtractFileURLsFromMarkdown 从 Markdown 内容中提取所有文件 URL（兼容旧接口）
func ExtractFileURLsFromMarkdown(content string) []string {
	references := ExtractFileReferencesFromMarkdown(content)
	urls := make([]string, 0, len(references))
	for _, ref := range references {
		urls = append(urls, ref.URL)
	}
	return urls
}
