package googleoauth

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

func UploadFile(ctx context.Context, fileName, folderID string, f io.Reader) (string, error) {
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return "", fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	file := &drive.File{
		Name:    fileName,
		Parents: []string{folderID},
	}
	uploadedFile, err := srv.Files.Create(file).Media(f).Do()
	if err != nil {
		return "", fmt.Errorf("could not upload file: %v", err)
	}

	perm := &drive.Permission{
		Type: "anyone", // すべてのユーザー
		Role: "reader", // 読み取り専用
	}
	_, err = srv.Permissions.Create(uploadedFile.Id, perm).Do()
	if err != nil {
		return "", fmt.Errorf("unable to set permission: %v", err)
	}

	return uploadedFile.Id, nil
}

func DownloadFile(ctx context.Context, fileID string) (*drive.File, *http.Response, error) {
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	meta, err := srv.Files.Get(fileID).Fields("name", "mimeType").Do()
	if err != nil {
		return nil, nil, err
	}

	resp, err := srv.Files.Get(fileID).Download()
	if err != nil {
		return nil, nil, err
	}

	return meta, resp, nil
}

func DeleteFile(ctx context.Context, fileID string) error {
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	if err := srv.Files.Delete(fileID).Do(); err != nil {
		return fmt.Errorf("failed to delete file: %v", err)
	}

	return nil
}
