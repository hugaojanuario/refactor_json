// Package cli cuida de toda a interação via terminal do refactor_doc.
package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AskUser solicita ao usuário o tipo de documento (outorgante/outorgado)
// e a frase a ser flexionada. Retorna ambos os valores ou um erro.
func AskUser() (docType string, phrase string, err error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("========================================")
	fmt.Println("  Refactor Doc — Flexão Linguística")
	fmt.Println("========================================")
	fmt.Println()
	fmt.Println("Escolha o tipo de documento:")
	fmt.Println("  1 - Outorgante")
	fmt.Println("  2 - Outorgado")
	fmt.Print("\nOpção: ")

	opt, err := reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("lendo opção: %w", err)
	}

	switch strings.TrimSpace(opt) {
	case "1":
		docType = "outorgante"
	case "2":
		docType = "outorgado"
	default:
		return "", "", fmt.Errorf("opção inválida %q: esperado 1 ou 2", strings.TrimSpace(opt))
	}

	fmt.Printf("\nTipo selecionado: %s\n", docType)
	fmt.Print("\nDigite a frase no masculino singular\n(ex: \"o advogado constitui\"): ")

	phrase, err = reader.ReadString('\n')
	if err != nil {
		return "", "", fmt.Errorf("lendo frase: %w", err)
	}

	phrase = strings.TrimSpace(phrase)
	if phrase == "" {
		return "", "", fmt.Errorf("a frase não pode ser vazia")
	}

	return docType, phrase, nil
}
