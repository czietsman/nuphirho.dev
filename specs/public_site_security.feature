Feature: Public site security baseline
  As a visitor
  I want both public sites to declare browser protections and a security contact
  So that browsers constrain risky capabilities and reports have a maintained route

  Scenario: A visitor opens either public site
    When the site returns a page
    Then the response prevents framing
    And the response restricts browser capabilities
    And the content policy permits the configured Google Fonts

  Scenario: A security researcher checks the standard disclosure path
    When the researcher requests the site's security.txt file
    Then the file names a contact, an expiry, and its canonical URL

  Scenario: A crawler requests the main sitemap
    When the sitemap is generated
    Then it includes the papers, privacy, and cookie pages
