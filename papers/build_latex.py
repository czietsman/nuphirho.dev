import subprocess
import os

os.chdir(os.path.dirname(os.path.abspath(__file__)))

for _ in range(2):
    subprocess.run(
        ["pdflatex", "-interaction=nonstopmode", "specification-as-quality-gate.tex"],
        check=True
    )

for ext in [".aux", ".log", ".out", ".toc"]:
    path = f"specification-as-quality-gate{ext}"
    if os.path.exists(path):
        os.remove(path)

print("Done: specification-as-quality-gate.pdf")
