Feature: Hashnode client
  As the publish pipeline
  I want to interact with the Hashnode GraphQL API
  So that posts are published, updated, drafted, and restored correctly

  Background:
    Given a Hashnode client configured with publication ID "pub123"

  # --- Publish ---

  Scenario: Publish a new post
    Given no post exists with slug "my-new-post"
    And no deleted post exists with slug "my-new-post"
    When the pipeline publishes a post:
      | title    | My New Post                |
      | slug     | my-new-post                |
      | subtitle | A subtitle                 |
      | content  | Hello world.               |
      | tags     | go,testing                 |
    Then a publishPost mutation is sent
    And the response contains post ID "post-001" and URL "https://blog.example.com/my-new-post"

  Scenario: Update an existing post by slug
    Given a post exists with slug "existing-post" and ID "post-002"
    When the pipeline publishes a post:
      | title    | Updated Post               |
      | slug     | existing-post              |
      | subtitle | Updated subtitle           |
      | content  | Updated content.           |
      | tags     | go                         |
    Then an updatePost mutation is sent with ID "post-002"
    And the response contains post URL "https://blog.example.com/existing-post"

  # --- Drafts ---

  Scenario: Create a new draft
    Given no draft exists with slug "draft-post"
    When the pipeline creates a draft:
      | title    | Draft Post                 |
      | slug     | draft-post                 |
      | subtitle | Draft subtitle             |
      | content  | Draft content.             |
      | tags     | drafts                     |
    Then a createDraft mutation is sent
    And the response contains draft ID "draft-001"

  Scenario: Skip creating a draft that already exists
    Given a draft exists with slug "existing-draft" and ID "draft-002"
    When the pipeline creates a draft:
      | title    | Existing Draft             |
      | slug     | existing-draft             |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | drafts                     |
    Then no mutation is sent
    And the result reports the draft already exists with ID "draft-002"

  # --- Deleted post restoration ---

  Scenario: Restore a deleted post and update it
    Given no post exists with slug "deleted-post"
    And a deleted post exists with slug "deleted-post" and ID "del-001"
    When the pipeline publishes a post:
      | title    | Restored Post              |
      | slug     | deleted-post               |
      | subtitle | Restored subtitle          |
      | content  | Restored content.          |
      | tags     | recovery                   |
    Then a restorePost mutation is sent with ID "del-001"
    And an updatePost mutation is sent with ID "del-001"
    And the response contains post URL "https://blog.example.com/deleted-post"

  Scenario: No deleted post found falls through to create
    Given no post exists with slug "brand-new"
    And no deleted post exists with slug "brand-new"
    When the pipeline publishes a post:
      | title    | Brand New Post             |
      | slug     | brand-new                  |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | new                        |
    Then a publishPost mutation is sent
    And the response contains post ID "post-001" and URL "https://blog.example.com/brand-new"

  # --- Existence checks ---

  Scenario: Check post exists by slug returns the post
    Given a post exists with slug "found-post" and ID "post-010"
    When the client checks for a post with slug "found-post"
    Then the post is found with ID "post-010"

  Scenario: Check post exists by slug returns not found
    Given no post exists with slug "missing-post"
    When the client checks for a post with slug "missing-post"
    Then the post is not found

  Scenario: Check deleted posts finds a match
    Given a deleted post exists with slug "archived-post" and ID "del-010"
    When the client checks for deleted posts with slug "archived-post"
    Then the deleted post is found with ID "del-010"

  Scenario: Check deleted posts finds no match
    Given no deleted post exists with slug "clean-slate"
    When the client checks for deleted posts with slug "clean-slate"
    Then no deleted post is found

  # --- Error handling ---

  Scenario: Invalid token returns authentication error
    Given the Hashnode API returns an authentication error
    When the pipeline publishes a post:
      | title    | Fail Post                  |
      | slug     | fail-post                  |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | test                       |
    Then the error is "authentication failed"

  Scenario: Publication not found returns error
    Given the Hashnode API returns a publication not found error
    When the client checks for a post with slug "any-slug"
    Then the error is "publication not found"

  Scenario: Network failure returns error
    Given the Hashnode API is unreachable
    When the pipeline publishes a post:
      | title    | Fail Post                  |
      | slug     | fail-post                  |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | test                       |
    Then the error contains "connection"

  # --- Dry-run ---

  Scenario: Dry-run for new post builds mutation without sending
    Given no post exists with slug "dry-run-post"
    And no deleted post exists with slug "dry-run-post"
    And dry-run mode is enabled
    When the pipeline publishes a post:
      | title    | Dry Run Post               |
      | slug     | dry-run-post               |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | test                       |
    Then no HTTP request is made
    And the dry-run result action is "publish"
    And the dry-run result slug is "dry-run-post"

  Scenario: Dry-run for existing post builds update mutation
    Given a post exists with slug "dry-run-update" and ID "post-099"
    And dry-run mode is enabled
    When the pipeline publishes a post:
      | title    | Dry Run Update             |
      | slug     | dry-run-update             |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | test                       |
    Then the dry-run result action is "update"
    And the dry-run result existing ID is "post-099"

  Scenario: Dry-run for deleted post builds restore and update
    Given no post exists with slug "dry-run-restore"
    And a deleted post exists with slug "dry-run-restore" and ID "del-099"
    And dry-run mode is enabled
    When the pipeline publishes a post:
      | title    | Dry Run Restore            |
      | slug     | dry-run-restore            |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | test                       |
    Then the dry-run result action is "restore_and_update"
    And the dry-run result deleted ID is "del-099"

  Scenario: Dry-run for draft builds create draft mutation
    Given no draft exists with slug "dry-run-draft"
    And dry-run mode is enabled
    When the pipeline creates a draft:
      | title    | Dry Run Draft              |
      | slug     | dry-run-draft              |
      | subtitle | Subtitle                   |
      | content  | Content.                   |
      | tags     | test                       |
    Then no HTTP request is made
    And the dry-run result action is "create_draft"
