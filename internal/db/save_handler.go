package db

import "fmt"

type SaveHandler struct {
}

func NewSaveHandler() *SaveHandler {
	return &SaveHandler{}
}

func (h *SaveHandler) UploadFile() error {
	return fmt.Errorf("error")
}
