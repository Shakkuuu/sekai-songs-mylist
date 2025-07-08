package googleoauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

func SendVerificationEmail(ctx context.Context, toEmail, token string) error {
	baseURL := os.Getenv("BACK_END_URL")
	verificationURL := fmt.Sprintf(baseURL+"/verify?token=%s&email=%s", token, toEmail)
	resendVerificationURL := fmt.Sprintf(baseURL+"/verify/resend?token=%s&email=%s", token, toEmail)

	subject := "メールアドレスの確認"
	encodedSubject := encodeSubject(subject)

	body := fmt.Sprintf("以下のリンクをクリックしてメールアドレスを確認してください:\n\n%s\n\n期限が切れていた場合はこちら:\n\n%s", verificationURL, resendVerificationURL)

	var message gmail.Message
	fromEmail := os.Getenv("FROM_EMAIL_ADDRESS")
	emailBody := "From: " + fromEmail + "\r\n" +
		"To: " + toEmail + "\r\n" +
		"Subject: " + encodedSubject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=\"UTF-8\"\r\n" +
		"Content-Transfer-Encoding: 7bit\r\n" +
		"\r\n" + body + "\r\n"

	message.Raw = base64.URLEncoding.EncodeToString([]byte(emailBody))

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Gmail client: %w", err)
	}

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return fmt.Errorf("unable to send email: %v", err)
	}

	return nil
}

// 日本語の件名をMIME Base64でエンコード
func encodeSubject(subject string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(subject))
	return fmt.Sprintf("=?UTF-8?B?%s?=", encoded)
}
