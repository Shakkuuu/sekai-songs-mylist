package googleoauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

var client *http.Client

func Init() error {
	ctx := context.Background()

	// credentials.json は Google Cloud Console から取得したクライアントIDなど
	b, err := os.ReadFile(os.Getenv("MAIL_CREDENTIALS_FILE"))
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %v", err)
	}

	scopes := []string{
		gmail.GmailSendScope,
		"https://www.googleapis.com/auth/drive.file",
	}

	// Gmail API 用の認可スコープ
	config, err := google.ConfigFromJSON(b, scopes...)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %v", err)
	}

	client, err = getClient(ctx, config)
	if err != nil {
		return err
	}

	return nil
}

func getClient(ctx context.Context, config *oauth2.Config) (*http.Client, error) {
	mailTokenFile := os.Getenv("MAIL_TOKEN_FILE")
	tok, err := tokenFromFile(mailTokenFile)
	if err != nil {
		return nil, fmt.Errorf("トークン読み込み失敗: %v", err)
	}

	// 自動リフレッシュ付きのトークンソースを作成
	tokenSource := config.TokenSource(ctx, tok)

	// 新しいアクセストークンが取得された場合、ファイルを更新
	refreshedToken, err := tokenSource.Token()
	if err != nil {
		return nil, fmt.Errorf("アクセストークン取得失敗: %v", err)
	}
	if err := saveToken(mailTokenFile, refreshedToken); err != nil {
		return nil, err
	}

	// 自動リフレッシュ付きクライアントを返す
	return oauth2.NewClient(ctx, tokenSource), nil
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
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
		return err
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()
	if err := json.NewEncoder(f).Encode(token); err != nil {
		return err
	}
	fmt.Println("📦 トークン保存完了！")
	return nil
}
