"""Gera arquivos JSON de exemplo para testar o pipeline.

Cria DOCX mínimos válidos com os placeholders *MS/*FS/*MP/*FP,
empacota em ZIP e codifica em base64 dentro do JSON de entrada.

Uso:
    python scripts/generate_fixtures.py
"""

import base64
import io
import json
import os
import zipfile

# ---------------------------------------------------------------------------
# DOCX mínimo válido
# Um DOCX precisa de pelo menos 4 arquivos internos para ser aberto pelo Word.
# ---------------------------------------------------------------------------

CONTENT_TYPES_XML = """\
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml"
    ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>"""

RELS_XML = """\
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1"
    Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument"
    Target="word/document.xml"/>
</Relationships>"""

WORD_RELS_XML = """\
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships"/>"""


def make_document_xml(title: str) -> str:
    """Gera o word/document.xml com os quatro placeholders."""
    return f"""\
<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p>
      <w:r><w:t>{title}</w:t></w:r>
    </w:p>
    <w:p>
      <w:r><w:rPr><w:b/></w:rPr><w:t>Masculino singular: </w:t></w:r>
      <w:r><w:t>*MS</w:t></w:r>
    </w:p>
    <w:p>
      <w:r><w:rPr><w:b/></w:rPr><w:t>Feminino singular: </w:t></w:r>
      <w:r><w:t>*FS</w:t></w:r>
    </w:p>
    <w:p>
      <w:r><w:rPr><w:b/></w:rPr><w:t>Masculino plural: </w:t></w:r>
      <w:r><w:t>*MP</w:t></w:r>
    </w:p>
    <w:p>
      <w:r><w:rPr><w:b/></w:rPr><w:t>Feminino plural: </w:t></w:r>
      <w:r><w:t>*FP</w:t></w:r>
    </w:p>
    <w:sectPr/>
  </w:body>
</w:document>"""


def build_docx(title: str) -> bytes:
    """Constrói os bytes de um DOCX mínimo válido com os placeholders."""
    buf = io.BytesIO()
    with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
        zf.writestr("[Content_Types].xml", CONTENT_TYPES_XML)
        zf.writestr("_rels/.rels", RELS_XML)
        zf.writestr("word/_rels/document.xml.rels", WORD_RELS_XML)
        zf.writestr("word/document.xml", make_document_xml(title))
    return buf.getvalue()


def build_outer_zip(docx_bytes: bytes, filename: str) -> bytes:
    """Empacota o DOCX dentro de um ZIP externo (estrutura esperada pelo pipeline)."""
    buf = io.BytesIO()
    with zipfile.ZipFile(buf, "w", zipfile.ZIP_DEFLATED) as zf:
        zf.writestr(filename, docx_bytes)
    return buf.getvalue()


def create_fixture(dest_dir: str, name: str, descricao: str, docx_filename: str) -> str:
    """Cria um arquivo JSON de fixture e retorna o caminho."""
    docx_bytes = build_docx(descricao)
    outer_zip = build_outer_zip(docx_bytes, docx_filename)
    content_b64 = base64.b64encode(outer_zip).decode("utf-8")

    payload = {
        "name": name,
        "descricao": descricao,
        "modeloPadrao": False,
        "content": content_b64,
    }

    os.makedirs(dest_dir, exist_ok=True)
    path = os.path.join(dest_dir, f"{name.lower().replace(' ', '_')}.json")
    with open(path, "w", encoding="utf-8") as f:
        json.dump(payload, f, ensure_ascii=False, indent=2)

    return path


def main() -> None:
    fixtures = [
        # (dest_dir, name, descricao, docx_filename)
        (
            "input/outorgante",
            "Procuracao Simples",
            "Instrumento de Procuração — Outorgante",
            "procuracao.docx",
        ),
        (
            "input/outorgante",
            "Substabelecimento",
            "Substabelecimento de Poderes — Outorgante",
            "substabelecimento.docx",
        ),
        (
            "input/outorgado",
            "Contrato de Servico",
            "Contrato de Prestação de Serviços — Outorgado",
            "contrato_servico.docx",
        ),
    ]

    print("Gerando fixtures de exemplo...\n")
    for args in fixtures:
        path = create_fixture(*args)
        print(f"  Criado: {path}")

    print("\nPronto! Execute o pipeline com: make run-go")


if __name__ == "__main__":
    main()
