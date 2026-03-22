from datetime import date


def prorate_premium(annual_premium, start_date, end_date):
    """Calculate the prorated premium for a policy that starts or ends mid-year.

    Args:
        annual: annual premium amount
        start: start date of the coverage period
        end: end date of the coverage period

    Returns:
        prorated premium amount for the coverage period
    """
    if isinstance(start_date, str):
        start_date = date.fromisoformat(start_date)
    if isinstance(end_date, str):
        end_date = date.fromisoformat(end_date)

    days_in_period = (end_date - start_date).days

    return round(annual_premium * days_in_period / 365, 2)
