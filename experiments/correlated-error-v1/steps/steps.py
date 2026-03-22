import ast
import sys
import os

PROJECT_DIR = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, PROJECT_DIR)

from behave import given, when, then
from corpus.paginate import paginate
from corpus.binary_search import binary_search
from corpus.is_leap_year import is_leap_year
from corpus.truncate_string import truncate_string
from corpus.days_between import days_between


# ---------- paginate ----------

@given("a list of {n:d} items")
def step_given_list_of_n_items(context, n):
    context.items = list(range(1, n + 1))


@when("I request page {page:d} with page_size {page_size:d}")
def step_when_request_page(context, page, page_size):
    context.result = paginate(context.items, page, page_size)


@then("I should get items {start:d} through {end:d}")
def step_then_items_range(context, start, end):
    expected = list(range(start, end + 1))
    assert context.result == expected, f"Expected {expected}, got {context.result}"


@then("I should get an empty list")
def step_then_empty_list(context):
    assert context.result == [], f"Expected [], got {context.result}"


# ---------- binary_search ----------

@given("a sorted list {list_str}")
def step_given_sorted_list(context, list_str):
    context.arr = ast.literal_eval(list_str)


@when("I search for {target:d}")
def step_when_search(context, target):
    context.result = binary_search(context.arr, target)


@then("the result should be index {index:d}")
def step_then_index(context, index):
    assert context.result == index, f"Expected {index}, got {context.result}"


# ---------- is_leap_year ----------

@given("the year {year:d}")
def step_given_year(context, year):
    context.year = year


@when("I check if it is a leap year")
def step_when_check_leap(context):
    context.result = is_leap_year(context.year)


@then("the leap year result should be {expected}")
def step_then_result_bool(context, expected):
    expected_val = expected == "True"
    assert context.result == expected_val, f"Expected {expected_val}, got {context.result}"


# ---------- truncate_string ----------

@given('the string "{s}" with max_len {max_len:d}')
def step_given_string_maxlen(context, s, max_len):
    context.input_string = s
    context.max_len = max_len


@when("I truncate the string")
def step_when_truncate(context):
    context.result = truncate_string(context.input_string, context.max_len)


@then('the truncated result should be "{expected}"')
def step_then_result_string(context, expected):
    assert context.result == expected, f"Expected '{expected}', got '{context.result}'"


# ---------- days_between ----------

@given('date1 is "{d1}" and date2 is "{d2}"')
def step_given_dates(context, d1, d2):
    context.date1 = d1
    context.date2 = d2


@when("I calculate the days between them")
def step_when_days_between(context):
    context.result = days_between(context.date1, context.date2)


@then("the days result should be {expected:d}")
def step_then_result_int(context, expected):
    assert context.result == expected, f"Expected {expected}, got {context.result}"
