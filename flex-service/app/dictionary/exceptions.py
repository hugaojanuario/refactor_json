"""Dicionário manual de formas irregulares em português.

Cobre casos que as regras morfológicas regulares não conseguem derivar
(heterônimos, irregularidades históricas, etc.).

Formato:
    { palavra_masculino_singular: { "FS": fem_sing, "MP": masc_plur, "FP": fem_plur } }
"""

# ---------------------------------------------------------------------------
# Substantivos com gênero heterônimo (palavras completamente diferentes)
# ---------------------------------------------------------------------------
EXCEPTIONS: dict[str, dict[str, str]] = {
    "ator": {"FS": "atriz", "MP": "atores", "FP": "atrizes"},
    "cão": {"FS": "cadela", "MP": "cães", "FP": "cadelas"},
    "homem": {"FS": "mulher", "MP": "homens", "FP": "mulheres"},
    "rei": {"FS": "rainha", "MP": "reis", "FP": "rainhas"},
    "imperador": {"FS": "imperatriz", "MP": "imperadores", "FP": "imperatrizes"},
    "herói": {"FS": "heroína", "MP": "heróis", "FP": "heroínas"},
    "pai": {"FS": "mãe", "MP": "pais", "FP": "mães"},
    "boi": {"FS": "vaca", "MP": "bois", "FP": "vacas"},
    "cavalo": {"FS": "égua", "MP": "cavalos", "FP": "éguas"},
    "galo": {"FS": "galinha", "MP": "galos", "FP": "galinhas"},
    "carneiro": {"FS": "ovelha", "MP": "carneiros", "FP": "ovelhas"},
    "bode": {"FS": "cabra", "MP": "bodes", "FP": "cabras"},
    # Termos jurídicos comuns
    "réu": {"FS": "ré", "MP": "réus", "FP": "rés"},
}

# ---------------------------------------------------------------------------
# Verbos irregulares (3ª pessoa singular → 3ª pessoa plural)
# Cobre os verbos mais usados em documentos jurídicos e cartoriais.
# ---------------------------------------------------------------------------
VERB_PLURALS: dict[str, str] = {
    # Irregulares clássicos
    "é": "são",
    "foi": "foram",
    "era": "eram",
    "será": "serão",
    "está": "estão",
    "estava": "estavam",
    "tem": "têm",
    "tinha": "tinham",
    "vem": "vêm",
    "vai": "vão",
    "faz": "fazem",
    "diz": "dizem",
    "traz": "trazem",
    "quer": "querem",
    "sabe": "sabem",
    "pode": "podem",
    "deve": "devem",
    # Verbos terminados em -ui (grupo constitui/contribui)
    "constitui": "constituem",
    "contribui": "contribuem",
    "distribui": "distribuem",
    "possui": "possuem",
    "exclui": "excluem",
    "inclui": "incluem",
    "conclui": "concluem",
    "atribui": "atribuem",
    # Verbos regulares de alta frequência em documentos jurídicos
    "representa": "representam",
    "assina": "assinam",
    "outorga": "outorgam",
    "nomeia": "nomeiam",
    "autoriza": "autorizam",
    "delega": "delegam",
    "revoga": "revogam",
    "declara": "declaram",
    "assume": "assumem",
    "compromete": "comprometem",
    "responde": "respondem",
    "reconhece": "reconhecem",
    "age": "agem",
    "fica": "ficam",
    "exerce": "exercem",
    "presta": "prestam",
    "subscreve": "subscrevem",
    "assiste": "assistem",
}
