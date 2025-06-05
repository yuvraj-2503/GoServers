package mail

import (
	"context"
	"testing"
)

const (
	content = "<!DOCTYPE html>\n<html>\n<head>\n    <meta charset=\"UTF-8\">\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n    <title>Vibely OTP</title>\n    <style>\n        body {\n            font-family: Arial, sans-serif;\n            background-color: #f4f4f4;\n            padding: 20px;\n            display: flex;\n            justify-content: center;\n            align-items: center;\n        }\n        .container {\n            background: #ffffff;\n            padding: 20px;\n            border-radius: 10px;\n            box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);\n            max-width: 400px;\n            text-align: center;\n        }\n        .logo {\n            width: 100px;\n            margin-bottom: 20px;\n        }\n        .otp {\n            font-size: 24px;\n            font-weight: bold;\n            color: #ff6600;\n        }\n        .footer {\n            margin-top: 20px;\n            font-size: 12px;\n            color: #666;\n        }\n    </style>\n</head>\n<body>\n    <div class=\"container\">\n        <img class=\"logo\" src=\"https://media-hosting.imagekit.io/13b9224618c14f28/ChatGPT%20Image%20Apr%201,%202025,%2003_20_09%20PM.png?Expires=1838114676&Key-Pair-Id=K2ZIVPTIP2VGHC&Signature=JH5bohs85~9dsQULc9gUc~rG7Zu-Axsn46SA-HD7P0bY2RJd27BaFsUORcXXyWR4Tf8HfauxkDtsZpZpmK0EBwa~wcdeHxdr5LysyQKij3PBopqweDHykIhnSqEIzJSx34886aX58kqtleFoxokQhNRDvFtZ4jU04vLtsRZlcwOozH9F8HBDEDocggodpzgucFJ-iUqBuSb8zTuLh1lJHNP~rEX8QS6lpednS7Ae~xle4pkgMRzJTHfFPeYrRFmCcSejmP90yx-4L6XHvKl58EH3NmM23LOy6iiYo8kZ98qnlLhmnDXvnz2X9czufe3jBOpcL8frKzxEIms-2wGH6A__\" alt=\"Vibely Logo\">\n        <h2>Welcome to Vibely!</h2>\n        <p>Your One-Time Password (OTP) for Vibely is:</p>\n        <p class=\"otp\">123456</p>\n        <p>Please use this OTP to complete your action within <strong>10 minutes</strong>.</p>\n        <p>⚠ <strong>Do not share this OTP with anyone</strong> for security reasons.</p>\n        <div class=\"footer\">Stay social, stay Vibely!<br>The Vibely Team</div>\n    </div>\n</body>\n</html>"
)

func TestMailSenderImpl_Send(t *testing.T) {
	ctx := context.Background()
	type fields struct {
		apiKey string
	}
	type args struct {
		ctx  *context.Context
		mail *Mail
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "success",
			fields: fields{
				apiKey: "mlsn.8043f79cd78114bce0655d37efd072eb9f0b7de0ddee1a1a3a2e42b9d27f0210",
			},
			args: args{
				ctx: &ctx,
				mail: &Mail{
					From:    "noreply.vibely@trial-zkq340ezw30gd796.mlsender.net",
					To:      []string{"singh.yuvraj1047@gmail.com"},
					Subject: "OTP For SignUp to Vibely",
					Content: content,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MailSenderImpl{
				apiKey: tt.fields.apiKey,
			}
			if err := m.Send(tt.args.ctx, tt.args.mail); (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
