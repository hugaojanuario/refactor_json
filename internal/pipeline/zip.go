package pipeline

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"fmt"
)

// DecodeBase64 decodifica o campo content (base64 → bytes do documento).
// Tenta encoding padrão e depois URL-safe.
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

// EncodeBase64 codifica bytes para string base64 padrão.
func EncodeBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// DetectDocType analisa os bytes de um documento ZIP (DOCX ou ODT)
// e retorna ".odt" ou ".docx". Retorna erro se não conseguir identificar.
func DetectDocType(data []byte) (string, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", fmt.Errorf("abrindo documento como ZIP: %w", err)
	}

	for _, f := range r.File {
		switch f.Name {
		case "mimetype":
			// ODT sempre tem um arquivo 'mimetype'
			return ".odt", nil
		case "word/document.xml":
			// DOCX sempre tem word/document.xml
			return ".docx", nil
		}
	}

	return "", fmt.Errorf("não foi possível identificar o tipo do documento (nem ODT nem DOCX)")
}
