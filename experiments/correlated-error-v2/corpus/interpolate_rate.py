import math


def interpolate_rate(curve, tenor):
    """Interpolate an interest rate for a given tenor from a rate curve.

    Args:
        curve: list of (tenor, rate) tuples representing known curve points
        tenor: the tenor for which to interpolate a rate

    Returns:
        interpolated rate for the given tenor
    """
    if not curve:
        raise ValueError("Empty curve")

    # Exact match
    for t, r in curve:
        if t == tenor:
            return r

    # Find surrounding points
    lower = None
    upper = None
    for t, r in curve:
        if t < tenor:
            lower = (t, r)
        elif t > tenor and upper is None:
            upper = (t, r)

    if lower is None:
        return curve[0][1]
    if upper is None:
        return curve[-1][1]

    t0, r0 = lower
    t1, r1 = upper

    weight = (tenor - t0) / (t1 - t0)
    rate = r0 + weight * (r1 - r0)

    return round(rate, 6)
