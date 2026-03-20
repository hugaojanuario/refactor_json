package pipeline

import (
	"fmt"

	"github.com/hugaojanuario/refactor_doc/internal/models"
)

// ProcessFile é a função central do sistema. Orquestra o pipeline completo
// para um único arquivo JSON:
//
//  1. Lê o JSON do disco
//  2. Decodifica base64 → bytes do ZIP externo
//  3. Dentro do ZIP, encontra o DOCX/ODT e substitui os placeholders no XML
//  4. Reempacota: novo DOCX/ODT → novo ZIP → base64
//  5. Salva o JSON de saída com o novo content
func ProcessFile(inputPath, outputPath string, flexions models.Flexions) error {
	// Passo 1: lê o envelope JSON
	doc, err := ReadDocument(inputPath)
	if err != nil {
		return fmt.Errorf("passo 1 (leitura): %w", err)
	}

	// Passo 2: decodifica base64 → bytes brutos do ZIP
	zipData, err := DecodeBase64(doc.Content)
	if err != nil {
		return fmt.Errorf("passo 2 (base64): %w", err)
	}

	// Passo 3: processa o ZIP (encontra DOCX/ODT, substitui placeholders)
	newZipData, err := ProcessZIP(zipData, flexions)
	if err != nil {
		return fmt.Errorf("passo 3 (processamento): %w", err)
	}

	// Passo 4: recodifica ZIP → base64
	newContent := EncodeBase64(newZipData)

	// Passo 5: salva JSON de saída com content atualizado
	if err := WriteDocument(outputPath, doc, newContent); err != nil {
		return fmt.Errorf("passo 5 (escrita): %w", err)
	}

	return nil
}
