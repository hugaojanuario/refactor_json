package pipeline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hugaojanuario/refactor_doc/internal/models"
)

// FileContent guarda os documentos lidos e o formato original (array ou objeto).
// Isso garante que a saída preserve o mesmo formato do arquivo de entrada.
type FileContent struct {
	Documents []*models.Document
	IsArray   bool // true se o JSON de entrada era um array [...]
}

// ReadFile lê um arquivo JSON que pode ser um objeto único {...} ou um array [{...}].
func ReadFile(path string) (*FileContent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lendo %s: %w", path, err)
	}

	fc := &FileContent{}

	// Tenta deserializar como array primeiro
	var docs []*models.Document
	if err := json.Unmarshal(data, &docs); err == nil {
		fc.Documents = docs
		fc.IsArray = true
		return fc, nil
	}

	// Tenta como objeto único
	var doc models.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parseando JSON de %s: %w", path, err)
	}

	fc.Documents = []*models.Document{&doc}
	fc.IsArray = false
	return fc, nil
}

// WriteFile serializa e salva o FileContent preservando o formato original.
func WriteFile(path string, fc *FileContent) error {
	var (
		data []byte
		err  error
	)

	if fc.IsArray {
		data, err = json.MarshalIndent(fc.Documents, "", "  ")
	} else {
		data, err = json.MarshalIndent(fc.Documents[0], "", "  ")
	}

	if err != nil {
		return fmt.Errorf("serializando JSON: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("escrevendo %s: %w", path, err)
	}

	return nil
}
