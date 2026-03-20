"""Schemas Pydantic para request/response do endpoint /flex."""
from pydantic import BaseModel, Field


class FlexRequest(BaseModel):
    text: str = Field(..., min_length=1, description="Frase no masculino singular")

    model_config = {"json_schema_extra": {"example": {"text": "o advogado constitui"}}}


class FlexResponse(BaseModel):
    MS: str = Field(..., description="Masculino singular (original)")
    FS: str = Field(..., description="Feminino singular")
    MP: str = Field(..., description="Masculino plural")
    FP: str = Field(..., description="Feminino plural")

    model_config = {
        "json_schema_extra": {
            "example": {
                "MS": "o advogado constitui",
                "FS": "a advogada constitui",
                "MP": "os advogados constituem",
                "FP": "as advogadas constituem",
            }
        }
    }
