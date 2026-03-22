from datetime import date

def days_between(date1, date2):
    """Return the absolute number of days between two dates."""
    if isinstance(date1, str):
        date1 = date.fromisoformat(date1)
    if isinstance(date2, str):
        date2 = date.fromisoformat(date2)
    return (date2 - date1).days
