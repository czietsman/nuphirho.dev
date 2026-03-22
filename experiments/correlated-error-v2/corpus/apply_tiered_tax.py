def apply_tiered_tax(income, brackets):
    """Calculate the tax owed on an income using a tiered bracket system.

    Args:
        income: taxable income amount
        brackets: list of (upper_limit, rate) tuples in ascending order,
                  where upper_limit is None for the top bracket

    Returns:
        total tax owed
    """
    tax = 0.0

    for i, (threshold, rate) in enumerate(brackets):
        if income <= threshold:
            break

        if i + 1 < len(brackets):
            next_threshold = brackets[i + 1][0]
        else:
            next_threshold = income

        taxable_in_band = min(income, next_threshold) - threshold
        tax = taxable_in_band * rate

    return round(tax, 2)
