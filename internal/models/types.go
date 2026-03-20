// Package models define os tipos de dados compartilhados entre todos os pacotes internos.
// Centralizar aqui evita imports circulares e deixa o contrato explícito.
package models

// Document representa o envelope JSON que envolve cada arquivo de entrada/saída.
// O campo Content é um ZIP codificado em base64 que contém um DOCX ou ODT.
type Document struct {
	Name         string `json:"name"`
	Descricao    string `json:"descricao"`
	ModeloPadrao bool   `json:"modeloPadrao"`
	Content      string `json:"content"` // base64(ZIP(DOCX|ODT))
}

// Flexions contém as quatro formas flexionadas de uma frase em português.
// A frase de entrada é assumida como masculino singular (MS).
type Flexions struct {
	MS string `json:"MS"` // masculino singular — original
	FS string `json:"FS"` // feminino singular
	MP string `json:"MP"` // masculino plural
	FP string `json:"FP"` // feminino plural
}
