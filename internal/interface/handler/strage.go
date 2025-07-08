package handler

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/Shakkuuu/sekai-songs-mylist/internal/pkg/googleoauth"
	"github.com/cockroachdb/errors"
)

type StorageHandler struct {
}

func NewStorageHandler() *StorageHandler {
	return &StorageHandler{}
}

const maxUploadSize = 2 << 30 // 2GB

func (h *StorageHandler) GetImageHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		err := errors.New("missing id parameter")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	meta, data, err := googleoauth.DownloadFile(ctx, fileID)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := data.Body.Close(); err != nil {
			log.Printf("failed to close resp.Body: %v", err)
		}
	}()

	w.Header().Set("Content-Type", meta.MimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf(`inline; filename="%s"`, meta.Name))

	if _, err := io.Copy(w, data.Body); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "failed to stream image: "+err.Error(), http.StatusInternalServerError)
	}
}

func (h *StorageHandler) UploadThumbnailHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "File size over", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	fileID, err := googleoauth.UploadFile(ctx, handler.Filename, os.Getenv("THUMBNAIL_FOLDER_ID"), file)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "file upload error", http.StatusInternalServerError)
		return
	}

	publicURL := fmt.Sprintf(os.Getenv("BACK_END_URL")+"/image?id=%s", fileID)
	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintf(w, `{"url":"%s"}`, publicURL); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (h *StorageHandler) UploadAttachmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "File size over", http.StatusBadRequest)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("failed to close file: %v", err)
		}
	}()

	fileID, err := googleoauth.UploadFile(ctx, handler.Filename, os.Getenv("ATTACHMENT_FOLDER_ID"), file)
	if err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "file upload error", http.StatusInternalServerError)
		return
	}

	publicURL := fmt.Sprintf(os.Getenv("BACK_END_URL")+"/image?id=%s", fileID)
	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintf(w, `{"url":"%s"}`, publicURL); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

func (h *StorageHandler) DeleteAttachmentHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID := r.URL.Query().Get("id")
	if fileID == "" {
		err := errors.New("missing id parameter")
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := googleoauth.DeleteFile(ctx, fileID); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := fmt.Fprintf(w, `{"deleted":"%s"}`, fileID); err != nil {
		cerr := errors.WithStack(err)
		log.Printf("%+v\n", cerr)
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}
