# Como rodar o Refactor Doc

## Pré-requisitos

- Go 1.21+
- Python 3.11 (não use 3.12+ — spaCy ainda não é compatível)
- Homebrew (Mac)

---

## 1. Instalar Python 3.11

```bash
brew install python@3.11
```

Confirme:
```bash
python3.11 --version
# Python 3.11.x
```

---

## 2. Configurar o serviço Python

Abra um terminal e entre na pasta do serviço:

```bash
cd /Users/hug40/refactor_doc/flex-service
```

Crie o ambiente virtual com Python 3.11:

```bash
python3.11 -m venv venv
```

Ative o ambiente virtual:

```bash
source venv/bin/activate
```

> O terminal vai mostrar `(venv)` na frente. Isso confirma que está ativo.

Instale as dependências:

```bash
pip install -r requirements.txt
```

Baixe o modelo de português do spaCy:

```bash
python -m spacy download pt_core_news_sm
```

Suba o servidor:

```bash
uvicorn app.main:app --reload --port 8000
```

> Deve aparecer: `Uvicorn running on http://127.0.0.1:8000`
> **Deixe esse terminal aberto e rodando.**

---

## 3. Testar o serviço Python (Terminal novo)

Abra um **novo terminal** e rode:

```bash
curl -s -X POST http://localhost:8000/flex \
  -H "Content-Type: application/json" \
  -d '{"text": "o advogado constitui"}' | python3 -m json.tool
```

Saída esperada:

```json
{
  "MS": "o advogado constitui",
  "FS": "a advogada constitui",
  "MP": "os advogados constituem",
  "FP": "as advogadas constituem"
}
```

Se aparecer isso, o serviço está funcionando corretamente.

---

## 4. Gerar arquivos de teste

No mesmo terminal (raiz do projeto):

```bash
cd /Users/hug40/refactor_doc
python3 scripts/generate_fixtures.py
```

Isso cria arquivos JSON de exemplo em:
- `input/outorgante/`
- `input/outorgado/`

---

## 5. Rodar a CLI Go

```bash
cd /Users/hug40/refactor_doc
go run ./cmd/cli
```

Siga as instruções na tela:

```
Escolha o tipo de documento:
  1 - Outorgante
  2 - Outorgado

Opção: 1

Digite a frase no masculino singular
(ex: "o advogado constitui"): o advogado constitui
```

Os arquivos processados serão salvos em:
- `output/outorgante/`
- `output/outorgado/`

---

## 6. Verificar o resultado

Para inspecionar o XML dentro do arquivo gerado:

```bash
cat output/outorgante/procuracao_simples.json | python3 -c "
import json, sys, base64, zipfile, io

data = json.load(sys.stdin)
outer = base64.b64decode(data['content'])
with zipfile.ZipFile(io.BytesIO(outer)) as z1:
    docx = z1.read(z1.namelist()[0])
    with zipfile.ZipFile(io.BytesIO(docx)) as z2:
        xml = z2.read('word/document.xml').decode()
        print(xml)
"
```

O XML deve mostrar as flexões no lugar dos placeholders `*MS`, `*FS`, `*MP`, `*FP`.

---

## Resumo dos terminais

| Terminal | O que roda | Fica aberto? |
|---|---|---|
| Terminal 1 | Serviço Python (`uvicorn`) | Sim, sempre |
| Terminal 2 | CLI Go + scripts de teste | Não |

---

## Problemas comuns

**`command not found: uvicorn`**
→ O ambiente virtual não está ativo. Rode: `source venv/bin/activate`

**`spaCy não compatível / erro Pydantic`**
→ Você está usando Python 3.12+. Recrie o venv com Python 3.11:
```bash
rm -rf venv
python3.11 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python -m spacy download pt_core_news_sm
```

**`connection refused` no curl**
→ O serviço Python não está rodando. Volte ao Terminal 1 e suba com `uvicorn`.

**`nenhum arquivo .docx ou .odt encontrado`**
→ Rode `python3 scripts/generate_fixtures.py` antes de rodar a CLI.
