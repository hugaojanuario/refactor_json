"""Motor de flexão linguística para português brasileiro.

Estratégia:
1. Tokeniza a frase com spaCy (pt_core_news_sm) para obter POS tags.
2. Para cada token, aplica transformações baseadas no tipo (DET, NOUN, VERB, ADJ).
3. Usa o dicionário de exceções para formas irregulares.
4. Fallback conservador: mantém o token inalterado se a regra não se aplicar.

Decisão de design: regras explícitas em vez de modelo de ML para flexão.
Razão: documentos jurídicos têm vocabulário controlado e previsível; regras
determinísticas são mais auditáveis e não dependem de GPU/dados de treino.
"""

import spacy
from app.dictionary.exceptions import EXCEPTIONS, VERB_PLURALS
from app.models.schema import FlexResponse

# Carrega o modelo uma vez no nível do módulo (evita reload a cada request).
# Download: python -m spacy download pt_core_news_sm
try:
    nlp = spacy.load("pt_core_news_sm")
except OSError:
    raise RuntimeError(
        "Modelo spaCy 'pt_core_news_sm' não encontrado.\n"
        "Execute: python -m spacy download pt_core_news_sm"
    )

# ---------------------------------------------------------------------------
# Tabela de artigos/determinantes
# Chave: forma minúscula | Valor: { variante: forma_transformada }
# ---------------------------------------------------------------------------
ARTICLES: dict[str, dict[str, str]] = {
    "o":    {"FS": "a",    "MP": "os",   "FP": "as"},
    "a":    {"MS": "o",    "MP": "os",   "FP": "as"},
    "os":   {"MS": "o",    "FS": "a",    "FP": "as"},
    "as":   {"MS": "o",    "FS": "a",    "MP": "os"},
    "um":   {"FS": "uma",  "MP": "uns",  "FP": "umas"},
    "uma":  {"MS": "um",   "MP": "uns",  "FP": "umas"},
    "uns":  {"MS": "um",   "FS": "uma",  "FP": "umas"},
    "umas": {"MS": "um",   "FS": "uma",  "MP": "uns"},
    "este": {"FS": "esta", "MP": "estes","FP": "estas"},
    "esta": {"MS": "este", "MP": "estes","FP": "estas"},
    "esse": {"FS": "essa", "MP": "esses","FP": "essas"},
    "essa": {"MS": "esse", "MP": "esses","FP": "essas"},
    "aquele": {"FS": "aquela", "MP": "aqueles", "FP": "aquelas"},
}


def flex_phrase(text: str) -> FlexResponse:
    """Ponto de entrada: gera as quatro formas flexionadas de uma frase.

    Args:
        text: Frase no masculino singular (ex: "o advogado constitui").

    Returns:
        FlexResponse com MS, FS, MP, FP.
    """
    doc = nlp(text)

    fs_parts: list[str] = []
    mp_parts: list[str] = []
    fp_parts: list[str] = []

    for token in doc:
        ws = token.whitespace_  # preserva espaçamento original
        fs_parts.append(_transform(token, gender="F", number="S") + ws)
        mp_parts.append(_transform(token, gender="M", number="P") + ws)
        fp_parts.append(_transform(token, gender="F", number="P") + ws)

    return FlexResponse(
        MS=text,
        FS="".join(fs_parts).strip(),
        MP="".join(mp_parts).strip(),
        FP="".join(fp_parts).strip(),
    )


def _transform(token, gender: str, number: str) -> str:
    """Transforma um único token para a combinação de gênero/número alvo.

    Args:
        token: Token spaCy com .text, .pos_, .whitespace_
        gender: "M" = masculino, "F" = feminino
        number: "S" = singular, "P" = plural

    Returns:
        Forma transformada do token.
    """
    pos = token.pos_
    text = token.text
    target_key = _variant_key(gender, number)

    if pos == "DET":
        return _flex_determiner(text, target_key)

    if pos in ("NOUN", "PROPN"):
        return _flex_noun(text, gender, number)

    if pos in ("VERB", "AUX"):
        return _flex_verb(text, number)

    if pos == "ADJ":
        # Adjetivos seguem as mesmas regras morfológicas dos substantivos
        return _flex_word(text, gender, number)

    # PUNCT, ADV, CCONJ, SCONJ, NUM, etc. → mantém inalterado
    return text


def _variant_key(gender: str, number: str) -> str:
    """Converte par gênero/número no formato de chave do dicionário (MS/FS/MP/FP)."""
    g = "M" if gender == "M" else "F"
    n = "S" if number == "S" else "P"
    return f"{g}{n}"


# ---------------------------------------------------------------------------
# Determinantes / Artigos
# ---------------------------------------------------------------------------

def _flex_determiner(text: str, target_key: str) -> str:
    """Transforma artigos e determinantes."""
    lower = text.lower()
    mapping = ARTICLES.get(lower)
    if mapping is None:
        return text

    transformed = mapping.get(target_key, text)
    return _preserve_case(text, transformed)


# ---------------------------------------------------------------------------
# Substantivos
# ---------------------------------------------------------------------------

def _flex_noun(text: str, gender: str, number: str) -> str:
    """Transforma um substantivo: verifica exceções primeiro, depois regras."""
    lower = text.lower()
    target_key = _variant_key(gender, number)

    if lower in EXCEPTIONS:
        result = EXCEPTIONS[lower].get(target_key, text)
        return _preserve_case(text, result)

    return _flex_word(text, gender, number)


def _flex_word(text: str, gender: str, number: str) -> str:
    """Aplica regras morfológicas regulares: primeiro gênero, depois número.

    Preserva o padrão de capitalização do texto original no resultado final.
    Ex: "OUTORGADO" → "OUTORGADAS" (não "OUTORGADas")
    """
    word = _apply_gender(text, gender)
    word = _apply_number(word, number)
    # Passa o texto ORIGINAL (antes de qualquer transformação) para preservar case
    return _preserve_case(text, word)


def _apply_gender(word: str, target_gender: str) -> str:
    """Aplica transformação de gênero (M→F).

    A entrada é sempre masculino singular, então só precisamos M→F.
    """
    lower = word.lower()

    if target_gender == "F":
        if lower.endswith("or"):
            # autor → autora, professor → professora, outorgador → outorgadora
            return word + "a"
        if lower.endswith("ão"):
            # outorgante permanece (não termina em -ão masculino tipicamente)
            # padrão: ão → ã (cidadão → cidadã)
            return word[:-2] + "ã"
        if lower.endswith("o"):
            # advogado → advogada, outorgado → outorgada
            return word[:-1] + "a"
        # Palavras em -e, -nte, consoante: invariantes de gênero (ex: gerente, réu→ré via exceção)
        return word

    # target_gender == "M": a entrada já é masculino, não há transformação
    return word


def _apply_number(word: str, target_number: str) -> str:
    """Aplica transformação de número (singular → plural)."""
    lower = word.lower()

    if target_number == "P":
        if lower.endswith("ão"):
            return word[:-2] + "ões"       # outorgação → outorgações
        if lower.endswith("ã"):
            return word + "s"              # irmã → irmãs, cidadã → cidadãs
        if lower.endswith("m"):
            return word[:-1] + "ns"        # homem → homens (via exceção), item → itens
        if lower.endswith("l"):
            # animal → animais, papel → papéis
            # Simplificação: remove -l e adiciona -is
            return word[:-1] + "is"
        if lower.endswith("r") or lower.endswith("z"):
            return word + "es"             # professor → professores, vez → vezes
        if lower.endswith("s"):
            return word                    # já plural ou invariante
        return word + "s"                  # advogada → advogadas (caso geral)

    # target_number == "S": entrada já é singular
    return word


# ---------------------------------------------------------------------------
# Verbos
# ---------------------------------------------------------------------------

def _flex_verb(text: str, number: str) -> str:
    """Aplica flexão de número ao verbo.

    Usa o dicionário de plurais para formas conhecidas.
    Para verbos regulares desconhecidos, aplica regras de terminação.
    Comportamento conservador: mantém inalterado se nenhuma regra se aplicar.
    """
    if number == "S":
        return text  # entrada já é singular

    lower = text.lower()

    # 1. Verifica dicionário de plurais (formas irregulares e frequentes)
    plural = VERB_PLURALS.get(lower)
    if plural:
        return _preserve_case(text, plural)

    # 2. Regras para verbos regulares com preservação de case
    if lower.endswith("a"):
        return _preserve_case(text, lower[:-1] + "am")
    if lower.endswith("e"):
        return _preserve_case(text, lower + "m")

    # Conservador: desconhecido, retorna inalterado
    return text


# ---------------------------------------------------------------------------
# Utilitário
# ---------------------------------------------------------------------------

def _preserve_case(original: str, transformed: str) -> str:
    """Aplica o padrão de capitalização do original ao transformado.

    Exemplos:
        ("Advogado", "advogada") → "Advogada"
        ("RÉUS", "réus")        → "RÉUS" (já transformado, uppercase)
    """
    if not original or not transformed:
        return transformed
    if original.isupper():
        return transformed.upper()
    if original[0].isupper():
        return transformed.capitalize()
    return transformed
