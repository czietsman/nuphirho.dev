import ast
import sys
import os
from datetime import date, datetime

# Add project root to path so we can import from corpus/
PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, PROJECT_DIR)

from behave import given, when, then
from corpus.calculate_final_reserve_fuel import calculate_final_reserve_fuel
from corpus.get_gas_day import get_gas_day
from corpus.validate_diagnosis_sequence import validate_diagnosis_sequence
from corpus.validate_imo_number import validate_imo_number
from corpus.electricity_cost import electricity_cost


# ---------- calculate_final_reserve_fuel ----------

@given(u'a "{engine_type}" engine with a fuel burn rate of {rate} kg per hour')
def step_given_engine(context, engine_type, rate):
    context.engine_type = engine_type
    context.fuel_burn_rate = float(rate)


@when("I calculate the final reserve fuel")
def step_when_reserve_fuel(context):
    context.reserve_fuel = calculate_final_reserve_fuel(
        context.engine_type, context.fuel_burn_rate
    )


@then("the reserve fuel should be {expected:f} kg")
def step_then_reserve_fuel(context, expected):
    assert abs(context.reserve_fuel - expected) < 0.01, \
        f"Expected {expected} kg, got {context.reserve_fuel} kg"


# ---------- get_gas_day ----------

@given('a timestamp of "{timestamp}"')
def step_given_timestamp(context, timestamp):
    context.timestamp = timestamp


@when("I determine the gas day")
def step_when_gas_day(context):
    context.gas_day = get_gas_day(context.timestamp)


@then('the gas day should be "{expected}"')
def step_then_gas_day(context, expected):
    expected_date = date.fromisoformat(expected)
    assert context.gas_day == expected_date, \
        f"Expected gas day {expected_date}, got {context.gas_day}"


# ---------- validate_diagnosis_sequence ----------

@given("a diagnosis sequence of {codes_str}")
def step_given_diagnosis_codes(context, codes_str):
    context.diagnosis_codes = ast.literal_eval(codes_str)


@when("I validate the diagnosis sequence")
def step_when_validate_diagnosis(context):
    context.diagnosis_result = validate_diagnosis_sequence(context.diagnosis_codes)


@then("the sequence should be valid")
def step_then_sequence_valid(context):
    assert context.diagnosis_result["valid"] is True, \
        f"Expected valid, got {context.diagnosis_result}"


@then("the sequence should be invalid")
def step_then_sequence_invalid(context):
    assert context.diagnosis_result["valid"] is False, \
        f"Expected invalid, got {context.diagnosis_result}"


@then('the errors should include "{expected_error}"')
def step_then_errors_include(context, expected_error):
    errors = context.diagnosis_result["errors"]
    expected_lower = expected_error.lower()
    found = any(expected_lower in e.lower() for e in errors)
    assert found, \
        f"Expected error containing '{expected_error}', got {errors}"


# ---------- validate_imo_number ----------

@given('an IMO number "{imo_number}"')
def step_given_imo_number(context, imo_number):
    context.imo_number = imo_number


@when("I validate the IMO number")
def step_when_validate_imo(context):
    context.imo_valid = validate_imo_number(context.imo_number)


@then("the IMO number should be valid")
def step_then_imo_valid(context):
    assert context.imo_valid is True, \
        f"Expected IMO number to be valid"


@then("the IMO number should be invalid")
def step_then_imo_invalid(context):
    assert context.imo_valid is False, \
        f"Expected IMO number to be invalid"


# ---------- electricity_cost ----------

@given("electricity tiers {tiers_str}")
def step_given_electricity_tiers(context, tiers_str):
    import re
    tiers = []
    for match in re.finditer(r'\(([^,]+),\s*([^)]+)\)', tiers_str):
        limit_str = match.group(1).strip()
        rate_str = match.group(2).strip()
        limit = float('inf') if limit_str == 'inf' else float(limit_str)
        rate = float(rate_str)
        tiers.append((limit, rate))
    context.electricity_tiers = tiers


@given("consumption of {kwh:d} kWh")
def step_given_consumption(context, kwh):
    context.consumption_kwh = kwh


@when("I calculate the electricity cost")
def step_when_electricity_cost(context):
    context.electricity_cost = electricity_cost(
        context.consumption_kwh, context.electricity_tiers
    )


@then("the electricity cost should be {expected:f}")
def step_then_electricity_cost(context, expected):
    assert abs(context.electricity_cost - expected) < 0.01, \
        f"Expected {expected}, got {context.electricity_cost}"
