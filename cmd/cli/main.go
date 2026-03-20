// Ponto de entrada da CLI refactor_doc.
// Orquestra: input do usuário → flexão via Python → processamento de todos os JSONs.
package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hugaojanuario/refactor_doc/internal/cli"
	"github.com/hugaojanuario/refactor_doc/internal/client"
	"github.com/hugaojanuario/refactor_doc/internal/pipeline"
)

func main() {
	// 1. Interação com o usuário
	docType, phrase, err := cli.AskUser()
	if err != nil {
		log.Fatalf("erro na entrada: %v", err)
	}

	// 2. Chama o serviço Python para obter as quatro formas flexionadas
	fmt.Printf("\nConsultando serviço de flexão para: %q\n", phrase)
	flexions, err := client.GetFlexions(phrase)
	if err != nil {
		log.Fatalf("erro no serviço de flexão: %v", err)
	}

	fmt.Println("\nFlexões geradas:")
	fmt.Printf("  MS: %s\n", flexions.MS)
	fmt.Printf("  FS: %s\n", flexions.FS)
	fmt.Printf("  MP: %s\n", flexions.MP)
	fmt.Printf("  FP: %s\n\n", flexions.FP)

	// 3. Define os diretórios de entrada e saída baseados no tipo escolhido
	inputDir := filepath.Join("input", docType)
	outputDir := filepath.Join("output", docType)

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		log.Fatalf("criando diretório de saída %s: %v", outputDir, err)
	}

	// 4. Lê todos os arquivos JSON do diretório de entrada
	entries, err := os.ReadDir(inputDir)
	if err != nil {
		log.Fatalf("lendo diretório de entrada %s: %v\n→ Crie o diretório e adicione arquivos JSON.", inputDir, err)
	}

	// 5. Processa cada arquivo JSON
	processed, failed := 0, 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		inputPath := filepath.Join(inputDir, entry.Name())
		outputPath := filepath.Join(outputDir, entry.Name())

		fmt.Printf("Processando: %s\n", entry.Name())

		if err := pipeline.ProcessFile(inputPath, outputPath, flexions); err != nil {
			fmt.Fprintf(os.Stderr, "  [ERRO] %s: %v\n", entry.Name(), err)
			failed++
			continue
		}

		fmt.Printf("  → Salvo em: %s\n", outputPath)
		processed++
	}

	fmt.Printf("\nConcluído: %d processado(s), %d com erro(s).\n", processed, failed)
	if failed > 0 {
		os.Exit(1)
	}
}
