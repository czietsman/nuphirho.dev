def electricity_cost(kwh, tiers):
    """Calculate the electricity cost for a given consumption.

    Args:
        kwh: total kilowatt-hours consumed in the billing period
        tiers: list of (upper_limit, rate_per_kwh) tuples, sorted by
               upper_limit ascending. The last tier may use float('inf')
               as the upper limit.

    Returns:
        total cost as a float rounded to 2 decimal places
    """
    if kwh <= 0:
        return 0.0

    for limit, rate in tiers:
        if kwh <= limit:
            return round(kwh * rate, 2)

    return round(kwh * tiers[-1][1], 2)
