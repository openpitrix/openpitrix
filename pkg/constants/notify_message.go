package constants

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
)

const (
	en            = "en"
	zhCN          = "zh_cn"
	defaultLocale = zhCN
)

const EmailNotifyName = "email"

type EmailNotifyContent struct {
	Content []string
}

const EmailNotifyTemplate = `
<!DOCTYPE html>
<html>
    <head>
        <meta charset="UTF-8">
        <title>OpenPitrix Notification</title>
    </head>
    <body>
		{{range .Content}}{{ . }}<br />{{end}}
	</body>
</html>
`

type NotifyMessage struct {
	en   string
	zhCN string
}

type NotifyTitle struct {
	NotifyMessage
}

type NotifyContent struct {
	NotifyMessage
}

func (n *NotifyMessage) GetMessage(locale string, params ...interface{}) string {
	switch locale {
	case en:
		return fmt.Sprintf(n.en, params...)
	case zhCN:
		return fmt.Sprintf(n.zhCN, params...)
	default:
		return fmt.Sprintf(n.zhCN, params...)
	}
}

func (n *NotifyTitle) GetDefaultMessage(params ...interface{}) string {
	return n.GetMessage(defaultLocale, params...)
}

func (n *NotifyContent) GetDefaultMessage(params ...interface{}) string {
	t, _ := template.New(EmailNotifyName).Parse(EmailNotifyTemplate)
	b := bytes.NewBuffer([]byte{})
	emailContent := &EmailNotifyContent{
		Content: strings.Split(n.GetMessage(defaultLocale, params...), "\n"),
	}
	t.Execute(b, emailContent)
	return b.String()
}

var (
	AdminInviteIsvNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】邀请您成为平台服务商",
		},
	}
	AdminInviteIsvNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
邀请您入驻应用市场，成为优质服务商，为平台用户提供企业解决方案、产品和集成服务，共享快速收益。 
用户：%s 
密码：%s 
首次登陆后请修改密码。
`,
		},
	}

	AdminInviteUserNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】邀请您成为平台用户",
		},
	}
	AdminInviteUserNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
邀请成为平台用户。
用户：%s 
密码：%s 
首次登陆后请修改密码。
`,
		},
	}

	IsvInviteMemberNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【%s】邀请您加入 %s 平台",
		},
	}
	IsvInviteMemberNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
【%s】邀请您加入 %s 平台协同工作。 
用户：%s 
密码：%s 
首次登陆后请修改密码。
`,
		},
	}

	SubmitVendorNotifyAdminTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用服务商资质申请",
		},
	}
	SubmitVendorNotifyAdminContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
收到 %s 应用服务商资质申请，请尽快完成审核。
`,
		},
	}

	SubmitVendorNotifyIsvTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】已收到您的应用服务商资质申请",
		},
	}
	SubmitVendorNotifyIsvContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
已收到您的应用服务商资质申请，我们会在3个工作日内完成审核，请您耐心等待。
`,
		},
	}

	PassVendorNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】您的 %s 应用服务商资质申请已通过",
		},
	}
	PassVendorNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
恭喜您，应用服务商资质申请通过审核，正式成为 %s 应用服务商。
`,
		},
	}

	RejectVendorNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】已拒绝您的 %s 应用服务商资质申请",
		},
	}
	RejectVendorNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
您提交的 %s 应用服务商资质申请信息有误，请核对相关内容，完善申请后重新提交。
`,
		},
	}

	SubmitAppVersionNotifyReviewerTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本审核申请",
		},
	}
	SubmitAppVersionNotifyReviewerContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
收到 %s 应用 %s 版本审核申请，请尽快完成审核。
`,
		},
	}

	SubmitAppVersionNotifySubmitterTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】已收到您的 %s 应用 %s 版本审核申请",
		},
	}
	SubmitAppVersionNotifySubmitterContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
已收到您的 %s 应用 %s 版本审核申请，请您耐心等待。
`,
		},
	}

	PassAppVersionInfoNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本通过应用信息审核",
		},
	}
	PassAppVersionInfoNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
恭喜您，%s 应用 %s 版本已通过应用信息审核，等待平台商务审核。
`,
		},
	}

	PassAppVersionBusinessNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本通过平台商务审核",
		},
	}
	PassAppVersionBusinessNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
恭喜您，%s 应用 %s 版本已通过平台商务审核，等待平台技术审核。
`,
		},
	}

	PassAppVersionTechnicalNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本通过平台技术审核",
		},
	}
	PassAppVersionTechnicalNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好：
恭喜您，%s 应用 %s 版本已通过平台技术审核，请尽快完成应用版本上架。 
`,
		},
	}

	RejectAppVersionInfoNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本未通过应用信息审核",
		},
	}
	RejectAppVersionInfoNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
您提交的 %s 应用 %s 版本未通过应用信息审核，请核对相关内容，完善后重新提交。
`,
		},
	}

	RejectAppVersionBusinessNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本未通过平台商务审核",
		},
	}
	RejectAppVersionBusinessNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
您提交的 %s 应用 %s 版本未通过平台商务审核，请核对相关内容，完善后重新提交。
`,
		},
	}

	RejectAppVersionTechnicalNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本未通过平台技术审核",
		},
	}
	RejectAppVersionTechnicalNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
您提交的 %s 应用 %s 版本未通过平台技术审核，请核对相关内容，完善后重新提交。
`,
		},
	}

	ReleaseAppVersionNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本已上架",
		},
	}
	ReleaseAppVersionNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
%s 应用 %s 版本已上架到应用市场。
`,
		},
	}

	SuspendAppVersionNotifyTitle = NotifyTitle{
		NotifyMessage: NotifyMessage{
			zhCN: "【OpenPitrix】%s 应用 %s 版本已下架",
		},
	}
	SuspendAppVersionNotifyContent = NotifyContent{
		NotifyMessage: NotifyMessage{
			zhCN: `
%s 您好： 
%s 应用 %s 版本已从应用市场下架。
`,
		},
	}
)
