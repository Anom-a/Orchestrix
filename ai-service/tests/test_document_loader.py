from pathlib import Path

from app.services.document_loader import load_document_from_path


def test_load_document_from_path_text_file(tmp_path: Path):
    file_path = tmp_path / "sample.txt"
    file_path.write_text("hello orchestrix", encoding="utf-8")

    text = load_document_from_path(str(file_path))

    assert text == "hello orchestrix"


def test_load_document_from_path_missing_file_raises(tmp_path: Path):
    missing = tmp_path / "missing.txt"

    try:
        load_document_from_path(str(missing))
        raised = False
    except FileNotFoundError:
        raised = True

    assert raised
