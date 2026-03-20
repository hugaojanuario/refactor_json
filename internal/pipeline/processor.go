package pipeline

import (
	"fmt"

	"github.com/hugaojanuario/refactor_doc/internal/docx"
	"github.com/hugaojanuario/refactor_doc/internal/models"
)

// ProcessFile é a função central do sistema. Orquestra o pipeline completo
// para um único arquivo JSON (que pode conter um ou mais documentos):
//
//  1. Lê o JSON (suporta objeto único ou array)
//  2. Para cada documento: decodifica base64 → bytes do ODT/DOCX
//  3. Detecta o tipo (ODT ou DOCX)
//  4. Substitui os placeholders *MS*/*FS*/*MP*/*FP* no XML interno
//  5. Recodifica → base64
//  6. Salva o JSON de saída preservando o formato original
func ProcessFile(inputPath, outputPath string, flexions models.Flexions) error {
	// Passo 1: lê o JSON (array ou objeto)
	fc, err := ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("passo 1 (leitura): %w", err)
	}

	// Passo 2-5: processa cada documento dentro do arquivo
	for i, doc := range fc.Documents {
		newContent, err := processDocument(doc.Content, flexions)
		if err != nil {
			return fmt.Errorf("documento[%d] %q: %w", i, doc.Name, err)
		}
		doc.Content = newContent
	}

	// Passo 6: salva o JSON de saída
	if err := WriteFile(outputPath, fc); err != nil {
		return fmt.Errorf("passo 6 (escrita): %w", err)
	}

	return nil
}

// processDocument processa o campo content de um único documento:
// decodifica base64, detecta tipo, substitui placeholders, recodifica.
func processDocument(content string, flexions models.Flexions) (string, error) {
	// Decode base64 → bytes brutos do ODT/DOCX
	docData, err := DecodeBase64(content)
	if err != nil {
		return "", fmt.Errorf("decode base64: %w", err)
	}

	// Detecta se é ODT ou DOCX
	ext, err := DetectDocType(docData)
	if err != nil {
		return "", fmt.Errorf("detectando tipo: %w", err)
	}

	// Processa o XML interno (substitui placeholders)
	newDocData, err := docx.ProcessDocument(docData, ext, flexions)
	if err != nil {
		return "", fmt.Errorf("processando %s: %w", ext, err)
	}

	// Recodifica → base64
	return EncodeBase64(newDocData), nil
}
