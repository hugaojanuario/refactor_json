"""Serviço FastAPI — Motor de Flexão Linguística para Português.

Expõe um único endpoint POST /flex que recebe uma frase no masculino
singular e retorna as quatro formas flexionadas (MS, FS, MP, FP).
"""

from fastapi import FastAPI, HTTPException

from app.models.schema import FlexRequest, FlexResponse
from app.services.flex import flex_phrase

app = FastAPI(
    title="Flex Engine",
    description="Motor de flexão linguística para português brasileiro",
    version="1.0.0",
)


@app.post("/flex", response_model=FlexResponse, summary="Flexionar frase")
async def flex_endpoint(request: FlexRequest) -> FlexResponse:
    """Recebe uma frase no masculino singular e retorna as quatro formas flexionadas.

    - **MS**: masculino singular (retorna a entrada original)
    - **FS**: feminino singular
    - **MP**: masculino plural
    - **FP**: feminino plural
    """
    try:
        return flex_phrase(request.text)
    except Exception as exc:
        raise HTTPException(status_code=500, detail=str(exc)) from exc


@app.get("/health", summary="Health check")
async def health() -> dict:
    """Verifica se o serviço está no ar."""
    return {"status": "ok"}
