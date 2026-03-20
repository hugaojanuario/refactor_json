package pipeline

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/hugaojanuario/refactor_doc/internal/docx"
	"github.com/hugaojanuario/refactor_doc/internal/models"
)

// DecodeBase64 decodifica o campo content (base64 → bytes do ZIP).
// Tenta primeiro o encoding padrão, depois o URL-safe.
func DecodeBase64(encoded string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		data, err = base64.URLEncoding.DecodeString(encoded)
		if err != nil {
			return nil, fmt.Errorf("base64 decode falhou: %w", err)
		}
	}
	return data, nil
}

// EncodeBase64 codifica bytes para string base64.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ProcessZIP abre o ZIP externo, encontra o(s) arquivo(s) DOCX/ODT,
// processa cada um (substituindo os placeholders no XML interno),
// e retorna os bytes do novo ZIP externo com os documentos atualizados.
func ProcessZIP(outerZIPData []byte, flexions models.Flexions) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(outerZIPData), int64(len(outerZIPData)))
	if err != nil {
		return nil, fmt.Errorf("abrindo ZIP externo: %w", err)
	}

	var outBuf bytes.Buffer
	w := zip.NewWriter(&outBuf)

	found := false
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("abrindo entrada %s: %w", f.Name, err)
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("lendo entrada %s: %w", f.Name, err)
		}

		// Detecta DOCX ou ODT pela extensão (path.Ext usa '/' de ZIPs corretamente)
		ext := strings.ToLower(path.Ext(f.Name))
		if ext == ".docx" || ext == ".odt" {
			found = true
			content, err = docx.ProcessDocument(content, ext, flexions)
			if err != nil {
				return nil, fmt.Errorf("processando documento %s: %w", f.Name, err)
			}
		}

		// Preserva metadados originais (compressão, timestamps) via CreateHeader
		fw, err := w.CreateHeader(&f.FileHeader)
		if err != nil {
			return nil, fmt.Errorf("criando header para %s: %w", f.Name, err)
		}
		if _, err := fw.Write(content); err != nil {
			return nil, fmt.Errorf("escrevendo %s no ZIP de saída: %w", f.Name, err)
		}
	}

	if !found {
		return nil, fmt.Errorf("nenhum arquivo .docx ou .odt encontrado no ZIP")
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("fechando ZIP de saída: %w", err)
	}

	return outBuf.Bytes(), nil
}
