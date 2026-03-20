.PHONY: build run-go run-python install-python fixtures clean test-flex

# ─── Go ───────────────────────────────────────────────────────────────────────

## Compila o binário da CLI
build:
	go build -o bin/refactor_doc ./cmd/cli

## Compila e executa a CLI interativamente
run-go: build
	./bin/refactor_doc

## Verifica se o código compila sem gerar binário
check:
	go vet ./...

# ─── Python ───────────────────────────────────────────────────────────────────

## Instala dependências Python e baixa o modelo spaCy
install-python:
	cd flex-service && pip install -r requirements.txt
	python -m spacy download pt_core_news_sm

## Inicia o serviço de flexão em modo desenvolvimento (hot-reload)
run-python:
	cd flex-service && uvicorn app.main:app --reload --port 8000

# ─── Fixtures ─────────────────────────────────────────────────────────────────

## Gera os arquivos JSON de entrada para teste
fixtures:
	python scripts/generate_fixtures.py

## Testa o serviço de flexão diretamente (requer curl + serviço rodando)
test-flex:
	curl -s -X POST http://localhost:8000/flex \
		-H "Content-Type: application/json" \
		-d '{"text": "o advogado constitui"}' | python -m json.tool

# ─── Limpeza ──────────────────────────────────────────────────────────────────

## Remove binários e outputs gerados
clean:
	rm -rf bin/
	find output -name "*.json" -delete
