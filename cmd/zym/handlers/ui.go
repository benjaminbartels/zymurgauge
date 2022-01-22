package handlers

import (
	"context"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/benjaminbartels/zymurgauge/web"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

type UIHandler struct {
	FileReader web.FileReader
}

func (h *UIHandler) Get(ctx context.Context, w http.ResponseWriter, r *http.Request, p httprouter.Params) error {
	var filePath string

	if strings.HasPrefix(r.URL.Path, base+"/static") {
		filePath = uiDir + strings.TrimPrefix(r.URL.Path, base)
	} else {
		filePath = uiDir + "/index.html"
	}

	b, err := h.FileReader.ReadFile(filePath)
	if err != nil {
		return errors.Wrap(err, "could not read file")
	}

	if contentType := mime.TypeByExtension(filepath.Ext(r.URL.Path)); len(contentType) > 0 {
		h := w.Header()
		h.Add("Content-Type", contentType)
	}

	if _, err := w.Write(b); err != nil {
		return errors.Wrap(err, "could not write response")
	}

	return nil
}
