# Papers

## The Specification as Quality Gate

Source: 04-practitioner-paper-v2.md
Output: specification-as-quality-gate.pdf

To rebuild the PDF, run from this directory:

    pip install markdown weasyprint --break-system-packages
    python build_pdf.py

The markdown file is the single source of truth. The PDF is derived from it.
The same markdown file is published as a blog post on Hashnode with frontmatter
and series navigation links included. The PDF build script strips the frontmatter
before rendering.

When submitting to arXiv, use the PDF from this directory. After arXiv assigns
a DOI, update the arXiv link in the blog post frontmatter and commit.
