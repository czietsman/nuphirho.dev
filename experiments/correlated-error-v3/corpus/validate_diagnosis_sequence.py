def validate_diagnosis_sequence(codes):
    """Validate a sequence of diagnosis codes for a medical claim.

    Checks that the diagnosis codes are well-formed and that the
    sequence is valid for submission.

    Args:
        codes: list of diagnosis code strings

    Returns:
        dict with "valid" bool and "errors" list
    """
    errors = []

    if not codes:
        return {"valid": False, "errors": ["No diagnosis codes provided"]}

    for i, code in enumerate(codes):
        if not isinstance(code, str) or len(code) < 3:
            errors.append(f"Code at position {i} is malformed: {code}")
            continue

        if not code[0].isalpha():
            errors.append(f"Code at position {i} must start with a letter: {code}")
            continue

        if not code[1:3].isdigit():
            errors.append(f"Code at position {i} has invalid format: {code}")

    if errors:
        return {"valid": False, "errors": errors}

    return {"valid": True, "errors": []}
