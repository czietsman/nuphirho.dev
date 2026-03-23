from datetime import datetime, timedelta


def get_gas_day(timestamp):
    """Determine the gas day for a given timestamp.

    A gas day is the standard trading day used in natural gas markets.
    Returns the date label for the gas day that the timestamp falls within.

    Args:
        timestamp: a datetime object representing the point in time

    Returns:
        a date object representing the gas day
    """
    if isinstance(timestamp, str):
        timestamp = datetime.fromisoformat(timestamp)

    return timestamp.date()
