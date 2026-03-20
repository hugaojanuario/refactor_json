// Package client implementa o cliente HTTP que chama o serviço Python de flexão.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hugaojanuario/refactor_doc/internal/models"
)

const defaultFlexServiceURL = "http://localhost:8000/flex"

// httpClient configurado com timeout para não travar o pipeline.
var httpClient = &http.Client{Timeout: 10 * time.Second}

type flexRequest struct {
	Text string `json:"text"`
}

// GetFlexions chama o serviço Python e retorna as quatro formas flexionadas.
// Retorna erro descritivo se o serviço não estiver rodando.
func GetFlexions(phrase string) (models.Flexions, error) {
	body, err := json.Marshal(flexRequest{Text: phrase})
	if err != nil {
		return models.Flexions{}, fmt.Errorf("serializando request: %w", err)
	}

	resp, err := httpClient.Post(defaultFlexServiceURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return models.Flexions{}, fmt.Errorf(
			"chamando serviço de flexão em %s: %w\n→ Inicie o serviço Python: cd flex-service && uvicorn app.main:app",
			defaultFlexServiceURL, err,
		)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.Flexions{}, fmt.Errorf("serviço de flexão retornou status %d", resp.StatusCode)
	}

	var flexions models.Flexions
	if err := json.NewDecoder(resp.Body).Decode(&flexions); err != nil {
		return models.Flexions{}, fmt.Errorf("decodificando resposta: %w", err)
	}

	return flexions, nil
}
