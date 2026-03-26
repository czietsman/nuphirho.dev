# Papers

## The Specification as Quality Gate

Source: ../posts/04-practitioner-paper-v2.md (markdown, single source of truth)
Style: arxiv.sty (modified from upstream, tracked)
Output: specification-as-quality-gate.pdf (not tracked, build locally)

### System dependencies

On Ubuntu 22.04 or 24.04, install the required TeX Live packages:

    sudo apt-get install -y \
      texlive-latex-base \
      texlive-latex-extra \
      texlive-fonts-recommended \
      texlive-fonts-extra \
      cm-super

These provide:
- `texlive-latex-base` -- pdflatex and core LaTeX packages
- `texlive-latex-extra` -- booktabs, listings, microtype, hyperref, fancyhdr
- `texlive-fonts-recommended` -- standard font families
- `texlive-fonts-extra` -- extended font support
- `cm-super` -- scalable Type 1 Computer Modern fonts, required by microtype for font expansion

### Build

    # Convert markdown to LaTeX then compile PDF
    cd papers
    python3 convert_to_latex.py
    ./build.sh

    # Clean auxiliary files
    ./build.sh --clean

The build script runs two pdflatex passes (cross-references require a second pass).
`arxiv.sty` must be present in the same directory as the `.tex` file (it is committed to the repo).

Alternatively, build via Docker if no local TeX installation is available:

    docker run --rm -v "$(pwd)":/work -w /work texlive/texlive:latest \
      sh -c 'pdflatex -interaction=nonstopmode specification-as-quality-gate.tex && pdflatex -interaction=nonstopmode specification-as-quality-gate.tex'

### Notes

The markdown file is published as a blog post on Hashnode with frontmatter
and series navigation links. The conversion script strips these before
rendering to LaTeX.

arxiv.sty is a modified copy of the upstream arxiv-style template with
small-caps removed from the title and abstract heading. It is tracked
because the modifications are local. The generated .tex file is not
tracked (gitignored).

When submitting to arXiv, build the PDF from this directory. After arXiv
assigns a DOI, update the arXiv link in the blog post frontmatter and commit.
