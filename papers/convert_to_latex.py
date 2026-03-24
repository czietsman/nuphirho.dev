#!/usr/bin/env python3
"""Convert 04-practitioner-paper-v2.md to LaTeX using arxiv-style template."""

import re
import os

os.chdir(os.path.dirname(os.path.abspath(__file__)))

with open("../posts/04-practitioner-paper-v2.md") as f:
    raw = f.read()

# Strip YAML frontmatter
raw = re.sub(r'^---\n.*?\n---\n', '', raw, count=1, flags=re.DOTALL)

# Strip series navigation italics at top (Part 4 of... line)
raw = re.sub(r'\*This is Part 4[^\n]*\*\n+', '', raw)

# Strip PDF link italic
raw = re.sub(r'\*A PDF version[^\n]*\*\n+', '', raw)

# Strip everything from final series links onwards
raw = re.sub(r'\n---\n+\*Series:.*', '', raw, flags=re.DOTALL)

# Extract and remove the preamble italic paragraph
preamble_match = re.search(
    r'^\*(This paper develops three.*?here\.\*)\n',
    raw, flags=re.DOTALL | re.MULTILINE
)
preamble_text = ""
if preamble_match:
    preamble_text = preamble_match.group(1).rstrip('*').strip()
    raw = raw[:preamble_match.start()] + raw[preamble_match.end():]

# Strip leading --- separator after preamble
raw = re.sub(r'^\n*---\n+', '\n', raw, count=1)

# Extract references section
refs_match = re.search(r'^## References\n+(.*)', raw, flags=re.DOTALL | re.MULTILINE)
refs_text = ""
if refs_match:
    refs_text = refs_match.group(1).strip()
    raw = raw[:refs_match.start()].strip()

# --- LaTeX escaping helper ---
def escape_latex(text):
    """Escape LaTeX special characters except those we handle separately."""
    text = text.replace('&', r'\&')
    text = text.replace('%', r'\%')
    text = text.replace('#', r'\#')
    text = text.replace('_', r'\_')
    text = text.replace('~', r'\textasciitilde{}')
    return text


def convert_inline(text):
    """Convert inline markdown formatting to LaTeX."""
    # Code spans first (before escaping)
    code_spans = []
    def save_code(m):
        code_spans.append(m.group(1))
        return f"%%CODE{len(code_spans)-1}%%"
    text = re.sub(r'`([^`]+)`', save_code, text)

    # Escape LaTeX special chars
    text = escape_latex(text)

    # Bold before italic
    text = re.sub(r'\*\*(.+?)\*\*', r'\\textbf{\1}', text)
    text = re.sub(r'\*(.+?)\*', r'\\emph{\1}', text)

    # Restore code spans
    for i, code in enumerate(code_spans):
        escaped_code = code.replace('_', r'\_').replace('%', r'\%').replace('#', r'\#').replace('&', r'\&')
        text = text.replace(f"%%CODE{i}%%", f"\\texttt{{{escaped_code}}}")

    # URLs (simple conversion)
    text = re.sub(r'\[([^\]]+)\]\(([^)]+)\)', r'\1', text)

    return text


# --- Process the body into LaTeX blocks ---
lines = raw.split('\n')
latex_blocks = []
i = 0

while i < len(lines):
    line = lines[i]

    # Section headings
    if line.startswith('### '):
        title = line[4:].strip()
        latex_blocks.append(f'\\subsection{{{convert_inline(title)}}}')
        i += 1
        continue
    if line.startswith('## '):
        title = line[3:].strip()
        if title == 'Abstract':
            latex_blocks.append('\\begin{abstract}')
            i += 1
            # Collect abstract paragraphs until next ---
            abstract_lines = []
            while i < len(lines) and lines[i].strip() != '---':
                abstract_lines.append(lines[i])
                i += 1
            abstract_text = '\n'.join(abstract_lines).strip()
            for para in re.split(r'\n\n+', abstract_text):
                para = para.strip()
                if para:
                    latex_blocks.append(convert_inline(para))
                    latex_blocks.append('')
            latex_blocks.append('\\end{abstract}')
            if i < len(lines):
                i += 1  # skip ---
            continue
        else:
            latex_blocks.append(f'\\section{{{convert_inline(title)}}}')
            i += 1
            continue

    # Horizontal rules
    if line.strip() == '---':
        i += 1
        continue

    # Fenced code blocks
    if line.strip().startswith('```'):
        lang = line.strip().lstrip('`').strip()
        i += 1
        code_lines = []
        while i < len(lines) and not lines[i].strip().startswith('```'):
            code_lines.append(lines[i])
            i += 1
        if i < len(lines):
            i += 1  # skip closing ```
        code = '\n'.join(code_lines)
        if lang == 'gherkin':
            latex_blocks.append('\\begin{lstlisting}[language={},basicstyle=\\ttfamily\\small,frame=single,breaklines=true]')
        else:
            latex_blocks.append('\\begin{lstlisting}[basicstyle=\\ttfamily\\small,frame=single,breaklines=true]')
        latex_blocks.append(code)
        latex_blocks.append('\\end{lstlisting}')
        continue

    # Markdown tables
    if '|' in line and i + 1 < len(lines) and re.match(r'\|[-| ]+\|', lines[i + 1]):
        table_lines = []
        while i < len(lines) and '|' in lines[i]:
            table_lines.append(lines[i])
            i += 1

        # Parse header
        header = [c.strip() for c in table_lines[0].split('|')[1:-1]]
        # Skip separator line
        rows = []
        for tl in table_lines[2:]:
            row = [c.strip() for c in tl.split('|')[1:-1]]
            rows.append(row)

        ncols = len(header)

        # Choose column spec based on number of columns
        if ncols == 7:
            colspec = 'p{3.2cm}p{2cm}p{0.8cm}p{0.8cm}p{0.8cm}p{0.8cm}p{0.8cm}'
        elif ncols == 4:
            colspec = 'lllr'
        else:
            colspec = 'l' * ncols

        latex_blocks.append('\\begin{table}[h]')
        latex_blocks.append('\\centering')
        latex_blocks.append('\\small')
        latex_blocks.append(f'\\begin{{tabular}}{{{colspec}}}')
        latex_blocks.append('\\toprule')

        header_tex = ' & '.join(convert_inline(h) for h in header) + ' \\\\'
        latex_blocks.append(header_tex)
        latex_blocks.append('\\midrule')

        for row in rows:
            # Pad row to match header length
            while len(row) < ncols:
                row.append('')
            row_tex = ' & '.join(convert_inline(c) for c in row) + ' \\\\'
            latex_blocks.append(row_tex)

        latex_blocks.append('\\bottomrule')
        latex_blocks.append('\\end{tabular}')
        latex_blocks.append('\\end{table}')
        continue

    # Empty lines
    if line.strip() == '':
        latex_blocks.append('')
        i += 1
        continue

    # Regular paragraph lines
    latex_blocks.append(convert_inline(line))
    i += 1

body = '\n'.join(latex_blocks)

# --- Convert references ---
ref_entries = re.split(r'\n\n+', refs_text.strip())
bib_items = []
for entry in ref_entries:
    if not entry.strip():
        continue
    entry = entry.strip()
    # Generate a citation key from first author surname and year
    # Try to extract author surname
    key_match = re.match(r'([A-Z][a-z]+)', entry)
    year_match = re.search(r'\((\d{4})\)', entry)
    if key_match and year_match:
        key = f"{key_match.group(1).lower()}{year_match.group(1)}"
    elif key_match:
        key = key_match.group(1).lower()
    else:
        key = f"ref{len(bib_items)+1}"

    # Escape the entry text
    entry_tex = escape_latex(entry)
    bib_items.append(f'\\bibitem{{{key}}}\n{entry_tex}')

bib_tex = '\n\n'.join(bib_items)

# --- Assemble full document ---
preamble_tex = convert_inline(preamble_text) if preamble_text else ""

document = rf"""\documentclass{{article}}

\usepackage{{arxiv}}
\usepackage[utf8]{{inputenc}}
\usepackage[T1]{{fontenc}}
\usepackage{{hyperref}}
\usepackage{{url}}
\usepackage{{booktabs}}
\usepackage{{listings}}
\usepackage{{microtype}}

\title{{The Specification as Quality Gate:\\Three Hypotheses on AI-Assisted Code Review}}

\author{{
  Christo Zietsman \\
  Independent Researcher \\
  \texttt{{nuphirho.dev}}
}}

\date{{March 2026}}

\begin{{document}}
\maketitle

\begin{{quote}}
\itshape
{preamble_tex}
\end{{quote}}

{body}

\begin{{thebibliography}}{{99}}

{bib_tex}

\end{{thebibliography}}

\end{{document}}
"""

with open("specification-as-quality-gate.tex", "w", newline='\n') as f:
    f.write(document)

print("Written: specification-as-quality-gate.tex")
