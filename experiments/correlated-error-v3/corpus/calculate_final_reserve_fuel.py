def calculate_final_reserve_fuel(engine_type, fuel_burn_rate):
    """Calculate the final reserve fuel requirement for an aircraft.

    Args:
        engine_type: either "jet" or "piston"
        fuel_burn_rate: fuel consumption in kg per hour

    Returns:
        required reserve fuel in kg
    """
    if engine_type not in ("jet", "piston"):
        raise ValueError(f"Unknown engine type: {engine_type}")

    if engine_type == "jet":
        reserve_minutes = 45
    else:
        reserve_minutes = 30

    reserve_hours = reserve_minutes / 60
    reserve_fuel = fuel_burn_rate * reserve_hours

    return round(reserve_fuel, 2)
