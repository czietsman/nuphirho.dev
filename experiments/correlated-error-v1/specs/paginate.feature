Feature: Paginate a list of items

  Scenario: First page of results
    Given a list of 10 items
    When I request page 1 with page_size 5
    Then I should get items 1 through 5

  Scenario: Middle page of results
    Given a list of 15 items
    When I request page 2 with page_size 5
    Then I should get items 6 through 10

  Scenario: Last page when items exactly divisible by page_size
    Given a list of 10 items
    When I request page 2 with page_size 5
    Then I should get items 6 through 10

  Scenario: Page beyond total pages returns empty
    Given a list of 10 items
    When I request page 4 with page_size 5
    Then I should get an empty list
