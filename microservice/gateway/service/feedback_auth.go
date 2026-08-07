package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	feedbacksdk "github.com/luyb177/feedback-sdk/feedback"
	"github.com/spf13/viper"
)

const (
	defaultFeedbackAudience = "feedback-center"
)

// FeedbackTokenResponse 保留 forum-be 原有的返回类型名称，实际结构由 SDK 提供。
type FeedbackTokenResponse = feedbacksdk.Token

// ExchangeFeedbackToken 使用反馈 SDK 生成后端身份断言并兑换反馈中台 Token。
func ExchangeFeedbackToken(ctx context.Context, studentID, tableIdentity string) (*FeedbackTokenResponse, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(viper.GetString("feedback_service.base_url")), "/")
	if baseURL == "" {
		return nil, fmt.Errorf("feedback_service.base_url 未配置")
	}

	projectID := strings.TrimSpace(viper.GetString("feedback_service.project_id"))
	if projectID == "" {
		return nil, fmt.Errorf("feedback_service.project_id 未配置")
	}

	keyID := strings.TrimSpace(viper.GetString("feedback_service.assertion_key_id"))
	if keyID == "" {
		return nil, fmt.Errorf("feedback_service.assertion_key_id 未配置")
	}

	tableIdentity = strings.TrimSpace(tableIdentity)
	if !allowedFeedbackTable(tableIdentity) {
		return nil, fmt.Errorf("feedback_service.table_identity 不允许")
	}

	privateKey, err := loadFeedbackAssertionPrivateKeyBytes()
	if err != nil {
		return nil, err
	}

	studentID = strings.TrimSpace(studentID)
	if studentID == "" {
		return nil, fmt.Errorf("student_id 不能为空")
	}

	issuer := strings.TrimSpace(viper.GetString("feedback_service.assertion_issuer"))
	if issuer == "" {
		issuer = projectID
	}
	audience := strings.TrimSpace(viper.GetString("feedback_service.assertion_audience"))
	if audience == "" {
		audience = defaultFeedbackAudience
	}
	sdkClient, err := feedbacksdk.NewClient(feedbacksdk.Config{
		Endpoint:     baseURL,
		ProjectID:    projectID,
		KeyID:        keyID,
		Issuer:       issuer,
		PrivateKey:   privateKey,
		Audience:     audience,
		HTTPClient:   &http.Client{Timeout: feedbackServiceTimeout},
		AssertionTTL: viper.GetDuration("feedback_service.assertion_ttl"),
	})
	if err != nil {
		return nil, fmt.Errorf("初始化反馈 SDK 失败: %w", err)
	}
	token, err := sdkClient.Exchange(ctx, feedbacksdk.Identity{
		StudentID:     studentID,
		TableIdentity: tableIdentity,
	})
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func allowedFeedbackTable(tableIdentity string) bool {
	for _, allowed := range viper.GetStringSlice("feedback_service.table_identities") {
		if strings.TrimSpace(allowed) == tableIdentity {
			return true
		}
	}
	legacy := strings.TrimSpace(viper.GetString("feedback_service.table_identify"))
	return legacy != "" && legacy == tableIdentity
}

func loadFeedbackAssertionPrivateKeyBytes() ([]byte, error) {
	keyData := []byte(strings.TrimSpace(viper.GetString("feedback_service.assertion_private_key")))
	if len(keyData) == 0 {
		keyFile := strings.TrimSpace(viper.GetString("feedback_service.assertion_private_key_file"))
		if keyFile == "" {
			return nil, fmt.Errorf("feedback_service.assertion_private_key_file 未配置")
		}

		var err error
		keyData, err = os.ReadFile(keyFile)
		if err != nil {
			return nil, fmt.Errorf("读取反馈身份私钥失败: %w", err)
		}
	}

	return keyData, nil
}
