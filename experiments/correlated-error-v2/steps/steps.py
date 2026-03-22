import ast
import sys
import os

# Add project root to path so we can import from corpus/
PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, PROJECT_DIR)

from behave import given, when, then
from corpus.prorate_premium import prorate_premium
from corpus.apply_tiered_tax import apply_tiered_tax
from corpus.schedule_maintenance import schedule_maintenance
from corpus.calculate_dilution import calculate_dilution
from corpus.interpolate_rate import interpolate_rate


# ---------- prorate_premium ----------

@given("an annual premium of {premium:f}")
def step_given_annual_premium(context, premium):
    context.annual_premium = premium


@when('I prorate from "{start}" to "{end}"')
def step_when_prorate(context, start, end):
    context.prorated = prorate_premium(context.annual_premium, start, end)


@then("the prorated premium should be {expected:f}")
def step_then_prorated_premium(context, expected):
    assert abs(context.prorated - expected) < 0.01, \
        f"Expected {expected}, got {context.prorated}"


# ---------- apply_tiered_tax ----------

@given("tax brackets {brackets_str}")
def step_given_tax_brackets(context, brackets_str):
    context.brackets = ast.literal_eval(brackets_str)


@given("an income of {income:d}")
def step_given_income(context, income):
    context.income = income


@when("I calculate the tax")
def step_when_calculate_tax(context):
    context.tax = apply_tiered_tax(context.income, context.brackets)


@then("the tax should be {expected:f}")
def step_then_tax(context, expected):
    assert abs(context.tax - expected) < 0.01, \
        f"Expected {expected}, got {context.tax}"


# ---------- schedule_maintenance ----------

@given('last service was "{service_date}"')
def step_given_last_service(context, service_date):
    context.last_service_date = service_date


@given("the aircraft has {hours:d} flight hours and {cycles:d} cycles")
def step_given_flight_data(context, hours, cycles):
    context.flight_hours = hours
    context.cycles = cycles


@when('I check maintenance status on "{as_of}"')
def step_when_check_maintenance(context, as_of):
    context.maintenance = schedule_maintenance(
        context.last_service_date,
        context.flight_hours,
        context.cycles,
        as_of_date=as_of,
    )


@then("maintenance should be due")
def step_then_maintenance_due(context):
    assert context.maintenance["due"] is True, \
        f"Expected maintenance due, got {context.maintenance}"


@then("maintenance should not be due")
def step_then_maintenance_not_due(context):
    assert context.maintenance["due"] is False, \
        f"Expected maintenance not due, got {context.maintenance}"


@then('the reasons should include "{reason}"')
def step_then_reasons_include(context, reason):
    assert reason in context.maintenance["reasons"], \
        f"Expected '{reason}' in reasons, got {context.maintenance['reasons']}"


# ---------- calculate_dilution ----------

@given("{shares:d} existing shares")
def step_given_existing_shares(context, shares):
    context.existing_shares = shares


@given("a {pct:d}% option pool")
def step_given_option_pool(context, pct):
    context.option_pool_pct = pct / 100.0


@given("an investment of {investment:d} at {pre_money:d} pre-money valuation")
def step_given_investment(context, investment, pre_money):
    context.investment = investment
    context.pre_money = pre_money


@when("I calculate the dilution")
def step_when_calculate_dilution(context):
    context.dilution = calculate_dilution(
        context.existing_shares,
        context.option_pool_pct,
        context.investment,
        context.pre_money,
    )


@then("founder ownership should be {expected:f}%")
def step_then_founder_pct(context, expected):
    actual = context.dilution["founder_pct"]
    assert abs(actual - expected) < 0.01, \
        f"Expected founder {expected}%, got {actual}%"


@then("investor ownership should be {expected:f}%")
def step_then_investor_pct(context, expected):
    actual = context.dilution["investor_pct"]
    assert abs(actual - expected) < 0.01, \
        f"Expected investor {expected}%, got {actual}%"


@then("pool ownership should be {expected:f}%")
def step_then_pool_pct(context, expected):
    actual = context.dilution["pool_pct"]
    assert abs(actual - expected) < 0.01, \
        f"Expected pool {expected}%, got {actual}%"


# ---------- interpolate_rate ----------

@given("a rate curve {curve_str}")
def step_given_rate_curve(context, curve_str):
    context.curve = ast.literal_eval(curve_str)


@when("I interpolate at tenor {tenor:f}")
def step_when_interpolate(context, tenor):
    context.interpolated_rate = interpolate_rate(context.curve, tenor)


@then("the interpolated rate should be {expected:f}")
def step_then_interpolated_rate(context, expected):
    assert abs(context.interpolated_rate - expected) < 0.000001, \
        f"Expected {expected}, got {context.interpolated_rate}"
