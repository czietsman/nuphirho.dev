# Behave discovers step definitions from specs/steps/.
# The canonical step implementations live in the project-level steps/ directory.
# This file loads them so behave can find all decorated step functions.

import importlib.util
import os
import sys

PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))
sys.path.insert(0, PROJECT_DIR)

spec = importlib.util.spec_from_file_location(
    "project_steps", os.path.join(PROJECT_DIR, "steps", "steps.py")
)
mod = importlib.util.module_from_spec(spec)
spec.loader.exec_module(mod)
