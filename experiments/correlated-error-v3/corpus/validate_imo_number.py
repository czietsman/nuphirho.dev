def validate_imo_number(imo_number):
    """Validate an IMO ship identification number.

    An IMO number is a seven-digit identifier assigned to vessels.
    The last digit is a check digit computed from the first six.

    Args:
        imo_number: string in the format "IMO1234567" or "1234567"

    Returns:
        True if the IMO number is valid, False otherwise
    """
    number = imo_number.replace("IMO", "").strip()

    if len(number) != 7 or not number.isdigit():
        return False

    digits = [int(d) for d in number]
    check_digit = digits[6]

    total = 0
    for i in range(6):
        factor = i + 2
        total += digits[i] * factor

    return total % 10 == check_digit
