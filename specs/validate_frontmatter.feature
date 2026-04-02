Feature: Validate frontmatter
  As the publish pipeline
  I want to parse and validate post frontmatter
  So that only well-formed posts are published

  Scenario: Post with valid frontmatter passes validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "My Post Title"
      slug: my-post-title
      subtitle: "A subtitle for the post"
      tags: [go, testing, bdd]
      ---

      Post content here.
      """
    When the frontmatter is parsed
    Then the title is "My Post Title"
    And the slug is "my-post-title"
    And the subtitle is "A subtitle for the post"
    And the tags are ["go", "testing", "bdd"]
    And the draft flag is false
    And the content is "Post content here."
    And validation passes with no errors

  Scenario: Post with draft flag set
    Given a markdown file with frontmatter:
      """
      ---
      title: "Draft Post"
      slug: draft-post
      tags: [drafts]
      draft: true
      ---

      Draft content.
      """
    When the frontmatter is parsed
    Then the draft flag is true

  Scenario: Post missing title fails validation
    Given a markdown file with frontmatter:
      """
      ---
      slug: no-title
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "missing required field: title"

  Scenario: Post missing slug fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "No Slug"
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "missing required field: slug"

  Scenario: Post missing tags fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "No Tags"
      slug: no-tags
      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "missing required field: tags"

  Scenario: Slug with uppercase characters fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "Bad Slug"
      slug: Bad-Slug
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "slug must be lowercase and hyphenated"

  Scenario: Slug with spaces fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "Bad Slug"
      slug: bad slug
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "slug must be lowercase and hyphenated"

  Scenario: Post with potential secrets fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "Secrets Post"
      slug: secrets-post
      tags: [test]
      ---

      Here is my api_key: "ghp_ABC123DEF456GHI789JKL012MNO"
      """
    When the frontmatter is parsed
    Then validation fails with error "potential secret detected in post content"

  Scenario: Em dash fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "Em Dash Post"
      slug: em-dash-post
      tags: [test]
      ---

      This sentence has an em dash — right here.
      """
    When the frontmatter is parsed
    Then validation fails with error "em dash detected"

  Scenario: Double hyphen as em dash substitute fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "Double Hyphen Post"
      slug: double-hyphen-post
      tags: [test]
      ---

      A title -- with a double hyphen.
      """
    When the frontmatter is parsed
    Then validation fails with error "em dash detected"

  Scenario: Double hyphen allowed when allow_emdash is true
    Given a markdown file with frontmatter:
      """
      ---
      title: "Citation Post"
      slug: citation-double-hyphen
      tags: [test]
      allow_emdash: true
      ---

      A cited paper title -- with a double hyphen.
      """
    When the frontmatter is parsed
    Then validation passes with no errors

  Scenario: Em dash allowed when allow_emdash is true
    Given a markdown file with frontmatter:
      """
      ---
      title: "Citation Post"
      slug: citation-post
      tags: [test]
      allow_emdash: true
      ---

      A cited paper title — with an em dash.
      """
    When the frontmatter is parsed
    Then validation passes with no errors

  Scenario: allow_emdash is optional and defaults to false
    Given a markdown file with frontmatter:
      """
      ---
      title: "Normal Post"
      slug: normal-post
      tags: [test]
      ---

      Clean content.
      """
    When the frontmatter is parsed
    Then validation passes with no errors

  Scenario: Quoted title value is stripped correctly
    Given a markdown file with frontmatter:
      """
      ---
      title: "Quoted Title"
      slug: quoted-title
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the title is "Quoted Title"

  Scenario: Unquoted title value is parsed correctly
    Given a markdown file with frontmatter:
      """
      ---
      title: Unquoted Title
      slug: unquoted-title
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the title is "Unquoted Title"

  Scenario: Single-quoted values are parsed correctly
    Given a markdown file with frontmatter:
      """
      ---
      title: 'Single Quoted'
      slug: single-quoted
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the title is "Single Quoted"

  Scenario: Tags in bracket notation are parsed as list
    Given a markdown file with frontmatter:
      """
      ---
      title: "Bracket Tags"
      slug: bracket-tags
      tags: [go, testing, bdd]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the tags are ["go", "testing", "bdd"]

  Scenario: Tags in quoted bracket notation are parsed as list
    Given a markdown file with frontmatter:
      """
      ---
      title: "Quoted Bracket Tags"
      slug: quoted-bracket-tags
      tags: ["go", "testing", "bdd"]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the tags are ["go", "testing", "bdd"]

  Scenario: Tags in YAML list notation are parsed as list
    Given a markdown file with frontmatter:
      """
      ---
      title: "YAML List Tags"
      slug: yaml-list-tags
      tags:
        - go
        - testing
        - bdd
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the tags are ["go", "testing", "bdd"]

  Scenario: Leading H1 matching title is stripped from content
    Given a markdown file with frontmatter:
      """
      ---
      title: "My Post"
      slug: my-post
      tags: [test]
      ---

      # My Post

      The actual content starts here.
      """
    When the frontmatter is parsed
    Then the content starts with "The actual content starts here."

  Scenario: Leading H1 not matching title is preserved
    Given a markdown file with frontmatter:
      """
      ---
      title: "My Post"
      slug: my-post
      tags: [test]
      ---

      # A Different Heading

      Content after heading.
      """
    When the frontmatter is parsed
    Then the content starts with "# A Different Heading"

  Scenario: Leading H1 matching title with no body is stripped to empty content
    Given a markdown file with frontmatter:
      """
      ---
      title: "My Post"
      slug: my-post
      tags: [test]
      ---

      # My Post
      """
    When the frontmatter is parsed
    Then the content is ""

  Scenario: File without frontmatter delimiters fails
    Given a markdown file with frontmatter:
      """
      title: No Delimiters
      slug: no-delimiters
      tags: [test]

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "missing frontmatter"

  Scenario: Empty frontmatter block still validates required fields
    Given a markdown file with frontmatter:
      """
      ---

      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "missing required field: title"

  Scenario: Opening delimiter may have trailing spaces before frontmatter
    Given a markdown file with frontmatter:
      """
      ---
      title: "Padded"
      slug: padded
      tags: [test]
      ---

      Body.
      """
    When the frontmatter is parsed
    Then the title is "Padded"
    And the content is "Body."
    And validation passes with no errors

  Scenario: Content immediately after the closing delimiter is preserved
    Given a markdown file with frontmatter:
      """
      ---
      title: "Immediate Body"
      slug: immediate-body
      tags: [test]
      ---
      Body.
      """
    When the frontmatter is parsed
    Then the content is "Body."
    And validation passes with no errors

  Scenario: Closing delimiter suffix text is ignored before body content
    Given a markdown file with frontmatter:
      """
      ---
      title: "Commented Closing"
      slug: commented-closing
      tags: [test]
      --- trailing text

      Body.
      """
    When the frontmatter is parsed
    Then the content is "Body."
    And validation passes with no errors

  Scenario: Empty tags list fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "Empty Tags"
      slug: empty-tags
      tags: []
      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "missing required field: tags"

  Scenario: Series field is parsed from frontmatter
    Given a markdown file with frontmatter:
      """
      ---
      title: "Series Post"
      slug: series-post
      tags: [go, testing]
      series: "Process Over Technology"
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the series is "Process Over Technology"
    And validation passes with no errors

  Scenario: Series field is optional
    Given a markdown file with frontmatter:
      """
      ---
      title: "No Series"
      slug: no-series
      tags: [go]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the series is ""
    And validation passes with no errors

  Scenario: publish_date field is parsed from frontmatter
    Given a markdown file with frontmatter:
      """
      ---
      title: "Scheduled Post"
      slug: scheduled-post
      tags: [go, testing]
      publish_date: 2026-03-25
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the publish date is "2026-03-25"
    And validation passes with no errors

  Scenario: publish_date field is optional
    Given a markdown file with frontmatter:
      """
      ---
      title: "No Date"
      slug: no-date
      tags: [go]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the publish date is empty
    And validation passes with no errors

  Scenario: Invalid publish_date fails validation
    Given a markdown file with frontmatter:
      """
      ---
      title: "Bad Date"
      slug: bad-date
      tags: [go]
      publish_date: not-a-date
      ---

      Content.
      """
    When the frontmatter is parsed
    Then validation fails with error "invalid publish_date"

  Scenario: Subtitle is optional
    Given a markdown file with frontmatter:
      """
      ---
      title: "No Subtitle"
      slug: no-subtitle
      tags: [test]
      ---

      Content.
      """
    When the frontmatter is parsed
    Then the subtitle is ""
    And validation passes with no errors

  Scenario: Body matching the title is not stripped without an H1 marker
    Given a markdown file with frontmatter:
      """
      ---
      title: "Body only"
      slug: body-only
      tags: [test]
      ---
      Body only
      """
    When the frontmatter is parsed
    Then the content is "Body only"
    And validation passes with no errors
