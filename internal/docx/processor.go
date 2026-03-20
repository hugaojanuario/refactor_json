// Package docx manipula arquivos DOCX e ODT (ambos são ZIPs contendo XML).
// Responsabilidade única: abrir o ZIP do documento, localizar o XML principal,
// substituir os placeholders e reempacotar.
package docx

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"github.com/hugaojanuario/refactor_doc/internal/models"
)

// xmlTargetFile retorna o caminho do XML principal dentro do documento.
// DOCX usa word/document.xml; ODT usa content.xml.
func xmlTargetFile(ext string) string {
	if ext == ".odt" {
		return "content.xml"
	}
	return "word/document.xml"
}

// ProcessDocument abre um DOCX ou ODT (ambos são ZIPs), substitui os placeholders
// *MS/*FS/*MP/*FP no XML principal e retorna os bytes do documento atualizado.
func ProcessDocument(docData []byte, ext string, flexions models.Flexions) ([]byte, error) {
	r, err := zip.NewReader(bytes.NewReader(docData), int64(len(docData)))
	if err != nil {
		return nil, fmt.Errorf("abrindo ZIP do documento: %w", err)
	}

	target := xmlTargetFile(ext)
	xmlFound := false

	var outBuf bytes.Buffer
	w := zip.NewWriter(&outBuf)

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf("abrindo %s: %w", f.Name, err)
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			return nil, fmt.Errorf("lendo %s: %w", f.Name, err)
		}

		if f.Name == target {
			xmlFound = true
			content = replacePlaceholders(content, flexions)
		}

		// CreateHeader preserva método de compressão e metadados originais.
		// Isso é importante para manter a integridade do DOCX/ODT.
		fw, err := w.CreateHeader(&f.FileHeader)
		if err != nil {
			return nil, fmt.Errorf("criando header para %s: %w", f.Name, err)
		}
		if _, err := fw.Write(content); err != nil {
			return nil, fmt.Errorf("escrevendo %s: %w", f.Name, err)
		}
	}

	if !xmlFound {
		return nil, fmt.Errorf("XML principal %q não encontrado no documento", target)
	}

	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("fechando ZIP do documento: %w", err)
	}

	return outBuf.Bytes(), nil
}

// replacePlaceholders substitui *MS, *FS, *MP, *FP no conteúdo XML.
//
// Decisão de design: usamos string replacement simples em vez de parsear o XML.
// Razão: os placeholders são marcadores artificiais inseridos pelo usuário e
// dificilmente estarão divididos entre elementos XML. Parsear e reconstruir
// o XML completo arriscaria alterar atributos de namespace e formatação que
// o Word/LibreOffice verifica na abertura do arquivo.
//
// Os valores de substituição são escapados para XML (& → &amp; etc.).
func replacePlaceholders(xmlData []byte, f models.Flexions) []byte {
	s := string(xmlData)

	// Placeholders no formato *XX* (com asterisco no início e no fim)
	s = strings.ReplaceAll(s, "*MS*", escapeXML(f.MS))
	s = strings.ReplaceAll(s, "*FS*", escapeXML(f.FS))
	s = strings.ReplaceAll(s, "*MP*", escapeXML(f.MP))
	s = strings.ReplaceAll(s, "*FP*", escapeXML(f.FP))

	return []byte(s)
}

// escapeXML escapa uma string para uso seguro dentro de conteúdo XML.
// Exemplo: "advogado & sócio" → "advogado &amp; sócio"
func escapeXML(s string) string {
	var buf strings.Builder
	if err := xml.EscapeText(&buf, []byte(s)); err != nil {
		// xml.EscapeText só falha com UTF-8 inválido; nesse caso retorna o original
		return s
	}
	return buf.String()
}
