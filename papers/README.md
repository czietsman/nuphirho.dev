# Papers

## The Specification as Quality Gate

Source: ../posts/04-practitioner-paper-v2.md (markdown, single source of truth)
Output: specification-as-quality-gate.pdf (not tracked, build locally)

### Build

    # Download arxiv style (one-time)
    curl -L https://raw.githubusercontent.com/kourgeorge/arxiv-style/master/arxiv.sty -o arxiv.sty

    # Convert markdown to LaTeX
    python convert_to_latex.py

    # Compile PDF (requires pdflatex, or use Docker)
    python build_latex.py

    # Or via Docker if no local TeX installation
    docker run --rm -v "$(pwd)":/work -w /work texlive/texlive:latest \
      sh -c 'pdflatex -interaction=nonstopmode specification-as-quality-gate.tex && pdflatex -interaction=nonstopmode specification-as-quality-gate.tex'
    rm -f specification-as-quality-gate.{aux,log,out,toc}

The markdown file is published as a blog post on Hashnode with frontmatter
and series navigation links. The conversion script strips these before
rendering to LaTeX.

When submitting to arXiv, build the PDF from this directory. After arXiv
assigns a DOI, update the arXiv link in the blog post frontmatter and commit.
