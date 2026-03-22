from datetime import date, timedelta


def schedule_maintenance(last_service_date, flight_hours, cycles,
                         hour_limit=500, cycle_limit=300,
                         calendar_limit_days=365, as_of_date=None):
    """Determine whether maintenance is due based on flight hours and cycles
    accumulated since the last service.

    Args:
        last_service: date of the last maintenance service
        flight_hours: total flight hours since last service
        cycles: total flight cycles since last service
        hours_limit: maximum permitted hours between services
        cycles_limit: maximum permitted cycles between services

    Returns:
        True if maintenance is due, False otherwise
    """
    if isinstance(last_service_date, str):
        last_service_date = date.fromisoformat(last_service_date)

    if as_of_date is None:
        as_of_date = date.today()
    if isinstance(as_of_date, str):
        as_of_date = date.fromisoformat(as_of_date)

    days_since_service = (as_of_date - last_service_date).days

    reasons = []

    if flight_hours >= hour_limit and cycles >= cycle_limit:
        reasons.append("flight_hours")
        reasons.append("cycles")

    if days_since_service >= calendar_limit_days:
        reasons.append("calendar")

    return {
        "due": len(reasons) > 0,
        "reasons": reasons
    }
