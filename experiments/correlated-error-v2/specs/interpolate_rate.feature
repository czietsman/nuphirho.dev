Feature: Rate curve interpolation using log-linear convention

  Scenario: Exact tenor match returns the rate directly
    Given a rate curve [(1, 0.02), (2, 0.025), (5, 0.035), (10, 0.04)]
    When I interpolate at tenor 2.0
    Then the interpolated rate should be 0.025000

  Scenario: Interpolation between close tenors
    Given a rate curve [(1, 0.02), (2, 0.025), (5, 0.035), (10, 0.04)]
    When I interpolate at tenor 1.5
    Then the interpolated rate should be 0.022361

  Scenario: Interpolation across wide tenor gap
    Given a rate curve [(1, 0.02), (5, 0.04), (10, 0.06)]
    When I interpolate at tenor 3.0
    Then the interpolated rate should be 0.028284
