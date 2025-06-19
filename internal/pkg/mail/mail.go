package mail

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/labstack/gommon/log"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var client *http.Client

func Init() error {
	ctx := context.Background()

	// credentials.json は Google Cloud Console から取得したクライアントIDなど
	b, err := os.ReadFile(os.Getenv("MAIL_CREDENTIALS_FILE"))
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "Unable to read client secret file"))
	}

	// Gmail API 用の認可スコープ
	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		return errors.WithStack(errors.WithMessage(err, "Unable to parse client secret file to config"))
	}

	client, err = getClient(ctx, config)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	mailTokenFile := os.Getenv("MAIL_TOKEN_FILE")
	tok, err := tokenFromFile(mailTokenFile)
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "トークン読み込み失敗"))
	}

	// 自動リフレッシュ付きのトークンソースを作成
	tokenSource := config.TokenSource(ctx, tok)

	// 新しいアクセストークンが取得された場合、ファイルを更新
	refreshedToken, err := tokenSource.Token()
	if err != nil {
		return nil, errors.WithStack(errors.WithMessage(err, "アクセストークン取得失敗"))
	}
	if err := saveToken(mailTokenFile, refreshedToken); err != nil {
		return nil, errors.WithStack(err)
	}

	// 自動リフレッシュ付きクライアントを返す
	return oauth2.NewClient(ctx, tokenSource), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("📦 トークンを %s に保存中...\n", path)
	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return errors.WithStack(err)
	}
	fmt.Println("📦 トークン保存完了！")
	return nil
}

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
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	_, err = srv.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Fatalf("Unable to send email: %v", err)
	}

	return nil
}

// 日本語の件名をMIME Base64でエンコード
func encodeSubject(subject string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(subject))
	return fmt.Sprintf("=?UTF-8?B?%s?=", encoded)
}
