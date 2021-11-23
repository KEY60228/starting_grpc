package handler

import (
	"bytes"
	"image/upload/gen/api"
	"io"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type ImageUploadHandler struct {
	sync.Mutex
	files map[string][]byte
}

func NewImageUploadHandler() *ImageUploadHandler {
	return &ImageUploadHandler{
		files: make(map[string][]byte),
	}
}

// 画像アップロード処理
func (h *ImageUploadHandler) Upload(stream api.ImageUploadService_UploadServer) error {
	// 最初のリクエストを受け取る
	req, err := stream.Recv()
	if err != nil {
		return err
	}

	// 初回は必ずメタ情報が送られる
	meta := req.GetFileMeta()
	filename := meta.Filename

	// UUIDの生成
	u, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	uuid := u.String()

	// 画像データ格納用バッファ
	buf := &bytes.Buffer{}

	// 塊ごとにアップロードされたバイナリをループしながら全て受け取る
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		chunk := r.GetData()
		_, err = buf.Write(chunk)
		if err != nil {
			return err
		}
	}

	// バッファから画像バイナリにする
	data := buf.Bytes()
	// 画像バイナリからファイル形式を検出する
	mimeType := http.DetectContentType(data)

	// h.uploaded.files[filename] = data
	h.files[uuid] = data

	err = stream.SendAndClose(&api.ImageUploadResponse{
		Uuid:        uuid,
		Size:        int32(len(data)),
		Filename:    filename,
		ContentType: mimeType,
	})

	return err
}
