package utils

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

func HandlePostImageUpload(r *http.Request, maxUploadSize int64, formName string) (sql.NullString, error) {
	fmt.Println("weiss")
	file, header, err := r.FormFile(formName)
	if err != nil {
		fmt.Println(err.Error())
		return sql.NullString{}, nil
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		return sql.NullString{}, errors.New("File too large (max 20MB)")
	}

	buff := make([]byte, 512)
	if _, err := file.Read(buff); err != nil {
		return sql.NullString{}, errors.New("Invalid file")
	}

	if _, err := file.Seek(0, 0); err != nil {
		return sql.NullString{}, errors.New("File error")
	}

	filetype := http.DetectContentType(buff)
	if filetype != "image/jpeg" && filetype != "image/png" && filetype != "image/gif" {
		return sql.NullString{}, errors.New("Only JPEG, PNG and GIF images are allowed")
	}

	ext := filepath.Ext(header.Filename)
	filename := uuid.New().String() + ext
	filePath := filepath.Join("../frontend/public/uploads", "posts", filename)

	if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
		return sql.NullString{}, errors.New("Unable to create upload directory")
	}

	dst, err := os.Create(filePath)
	if err != nil {
		return sql.NullString{}, errors.New("Failed to create file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return sql.NullString{}, errors.New("Failed to save file")
	}

	relativePath := strings.TrimPrefix(filePath, "../frontend/public")
	fmt.Println("jejehje",relativePath)
	return sql.NullString{String: relativePath, Valid: true}, nil
}
