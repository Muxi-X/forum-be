package audit

import (
	"context"
	"fmt"

	"github.com/muxi-Infra/auditor-Backend/sdk/v2/api/request"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/client"
	"github.com/muxi-Infra/auditor-Backend/sdk/v2/dto"
)

type auditClient struct {
	client     *client.Client
	webHookURL string
}

var AuditClient *auditClient

func InitAuditClient(apiKey, webHookURL, RegionURL string, timeout int) {
	c, err := client.NewClient(client.Config{
		ApiKey:         apiKey,
		ConnectTimeout: timeout,
		Region:         RegionURL,
	})

	if err != nil {
		panic(err)
	}

	AuditClient = &auditClient{
		client:     c,
		webHookURL: webHookURL,
	}
}

func (a *auditClient) SubmitToAudit(id uint, author, title, content string, pictures []string) error {
	auditContent := dto.NewContents(
		dto.WithTopicText(title, content),
		dto.WithTopicPictures(pictures),
	)

	req, err := request.NewUploadReq(
		a.webHookURL,
		id,
		auditContent,
		request.WithUploadAuthor(author),
	)
	if err != nil {
		return fmt.Errorf("构建审核请求失败: %w", err)
	}

	resp, err := a.client.UploadItem(context.Background(), req)
	if err != nil {
		return fmt.Errorf("送审失败: %w", err)
	}
	if resp.Basic.Errorx != nil {
		return fmt.Errorf("审核服务错误: %v", resp.Basic.Errorx)
	}

	return nil
}
