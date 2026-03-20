// Package pipeline orquestra o pipeline completo de processamento de documentos:
// leitura JSON → decode base64 → descompactar ZIP → processar XML → reempacotar → JSON.
package pipeline

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/hugaojanuario/refactor_doc/internal/models"
)

// ReadDocument lê e parseia um arquivo JSON de documento do disco.
func ReadDocument(path string) (*models.Document, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("lendo %s: %w", path, err)
	}

	var doc models.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parseando JSON de %s: %w", path, err)
	}

	if doc.Content == "" {
		return nil, fmt.Errorf("campo 'content' está vazio em %s", path)
	}

	return &doc, nil
}

// WriteDocument serializa o documento (com Content atualizado) e salva no disco.
func WriteDocument(path string, doc *models.Document, newContent string) error {
	doc.Content = newContent

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return fmt.Errorf("serializando documento: %w", err)
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return fmt.Errorf("escrevendo %s: %w", path, err)
	}

	return nil
}
