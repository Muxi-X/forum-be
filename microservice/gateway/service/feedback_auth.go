package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	feedbacksdk "github.com/luyb177/feedback-sdk/feedback"
	"github.com/spf13/viper"
)

// FeedbackTokenResponse 保留 forum-be 原有的返回类型名称，实际结构由 SDK 提供。
type FeedbackTokenResponse = feedbacksdk.Token

// ExchangeFeedbackToken 使用反馈 SDK 生成 HMAC 签名并兑换反馈中台 Token。
func ExchangeFeedbackToken(ctx context.Context, studentID string) (*FeedbackTokenResponse, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(viper.GetString("feedback_service.base_url")), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("feedback_service.base_url 未配置")
	}

	projectID := strings.TrimSpace(viper.GetString("feedback_service.project_id"))
	if projectID == "" {
		return nil, fmt.Errorf("feedback_service.project_id 未配置")
	}

	keyID := strings.TrimSpace(viper.GetString("feedback_service.key_id"))
	if keyID == "" {
		return nil, fmt.Errorf("feedback_service.key_id 未配置")
	}

	studentID = strings.TrimSpace(studentID)
	if studentID == "" {
		return nil, fmt.Errorf("student_id 不能为空")
	}

	apiKey := strings.TrimSpace(viper.GetString("feedback_service.api_key"))
	if apiKey == "" {
		return nil, fmt.Errorf("feedback_service.api_key 未配置")
	}
	sdkClient, err := feedbacksdk.NewClient(feedbacksdk.Config{
		Endpoint:   baseURL,
		ProjectID:  projectID,
		KeyID:      keyID,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: feedbackServiceTimeout},
	})
	if err != nil {
		return nil, fmt.Errorf("初始化反馈 SDK 失败: %w", err)
	}
	token, err := sdkClient.Exchange(ctx, feedbacksdk.Identity{
		StudentID: studentID,
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}
