import markdown
import re
from weasyprint import HTML, CSS

with open("04-practitioner-paper-v2.md") as f:
    raw = f.read()

raw = re.sub(r'^---.*?---\s*', '', raw, flags=re.DOTALL)

md = markdown.Markdown(extensions=['tables', 'fenced_code'])
body = md.convert(raw)

html = f"""<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<style>
  @page {{ margin: 2.8cm 3cm; }}
  body {{
    font-family: Georgia, serif;
    font-size: 10.5pt;
    line-height: 1.6;
    color: #111;
    max-width: 100%;
  }}
  h1 {{ font-size: 14pt; margin-top: 1.8em; margin-bottom: 0.4em; border-bottom: 1px solid #ccc; padding-bottom: 0.2em; }}
  h2 {{ font-size: 11.5pt; margin-top: 1.4em; margin-bottom: 0.3em; }}
  h3 {{ font-size: 10.5pt; margin-top: 1.2em; margin-bottom: 0.2em; font-style: italic; }}
  p {{ margin: 0.5em 0 0.7em 0; text-align: justify; }}
  code, pre {{
    font-family: "Courier New", monospace;
    font-size: 8.5pt;
    background: #f5f5f5;
    border-radius: 3px;
  }}
  pre {{
    padding: 0.7em 0.9em;
    margin: 0.8em 0;
    white-space: pre-wrap;
    line-height: 1.4;
  }}
  table {{
    border-collapse: collapse;
    width: 100%;
    margin: 1em 0;
    font-size: 8.5pt;
    table-layout: fixed;
  }}
  th, td {{
    border: 1px solid #ccc;
    padding: 0.3em 0.5em;
    text-align: left;
    word-wrap: break-word;
    overflow-wrap: break-word;
  }}
  th {{ background: #f0f0f0; font-weight: bold; }}
  hr {{ border: none; border-top: 1px solid #ccc; margin: 1.5em 0; }}
  em {{ font-style: italic; }}
  strong {{ font-weight: bold; }}
  .title-block {{
    text-align: center;
    margin-bottom: 2em;
    padding-bottom: 1em;
    border-bottom: 2px solid #333;
  }}
  .title-block h1 {{
    font-size: 17pt;
    border: none;
    margin-bottom: 0.2em;
  }}
</style>
</head>
<body>
<div class="title-block">
  <h1>The Specification as Quality Gate:<br>Three Hypotheses on AI-Assisted Code Review</h1>
  <p>Christo Zietsman &mdash; Independent Researcher &mdash; nuphirho.dev</p>
  <p>March 2026</p>
</div>
{body}
</body>
</html>"""

HTML(string=html).write_pdf(
    "specification-as-quality-gate.pdf",
    stylesheets=[CSS(string="@page { size: A4; }")]
)
print("Done")
