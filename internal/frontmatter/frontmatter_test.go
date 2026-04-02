package frontmatter_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/czietsman/nuphirho.dev/internal/frontmatter"
)

type frontmatterContext struct {
	post   *frontmatter.Post
	result *frontmatter.ValidationResult
}

func (fc *frontmatterContext) aMarkdownFileWithFrontmatter(doc *godog.DocString) error {
	fc.post, fc.result = frontmatter.Parse(doc.Content)
	return nil
}

func (fc *frontmatterContext) theFrontmatterIsParsed() error {
	// parsing already happened in the Given step
	return nil
}

func (fc *frontmatterContext) theTitleIs(expected string) error {
	if fc.post.Title != expected {
		return fmt.Errorf("expected title %q, got %q", expected, fc.post.Title)
	}
	return nil
}

func (fc *frontmatterContext) theSlugIs(expected string) error {
	if fc.post.Slug != expected {
		return fmt.Errorf("expected slug %q, got %q", expected, fc.post.Slug)
	}
	return nil
}

func (fc *frontmatterContext) theSeriesIs(expected string) error {
	if fc.post.Series != expected {
		return fmt.Errorf("expected series %q, got %q", expected, fc.post.Series)
	}
	return nil
}

func (fc *frontmatterContext) theSubtitleIs(expected string) error {
	if fc.post.Subtitle != expected {
		return fmt.Errorf("expected subtitle %q, got %q", expected, fc.post.Subtitle)
	}
	return nil
}

func (fc *frontmatterContext) theTagsAre(expected string) error {
	var expectedTags []string
	if err := json.Unmarshal([]byte(expected), &expectedTags); err != nil {
		return fmt.Errorf("invalid expected tags JSON: %v", err)
	}
	if len(fc.post.Tags) != len(expectedTags) {
		return fmt.Errorf("expected %d tags %v, got %d tags %v", len(expectedTags), expectedTags, len(fc.post.Tags), fc.post.Tags)
	}
	for i, tag := range expectedTags {
		if fc.post.Tags[i] != tag {
			return fmt.Errorf("tag %d: expected %q, got %q", i, tag, fc.post.Tags[i])
		}
	}
	return nil
}

func (fc *frontmatterContext) theDraftFlagIs(expected string) error {
	want := expected == "true"
	if fc.post.Draft != want {
		return fmt.Errorf("expected draft=%v, got %v", want, fc.post.Draft)
	}
	return nil
}

func (fc *frontmatterContext) theContentIs(expected string) error {
	trimmed := strings.TrimSpace(fc.post.Content)
	if trimmed != expected {
		return fmt.Errorf("expected content %q, got %q", expected, trimmed)
	}
	return nil
}

func (fc *frontmatterContext) theContentStartsWith(expected string) error {
	trimmed := strings.TrimSpace(fc.post.Content)
	if !strings.HasPrefix(trimmed, expected) {
		return fmt.Errorf("expected content to start with %q, got %q", expected, trimmed[:min(len(trimmed), len(expected)+20)])
	}
	return nil
}

func (fc *frontmatterContext) validationPassesWithNoErrors() error {
	if !fc.result.Passed() {
		return fmt.Errorf("expected validation to pass, got errors: %v", fc.result.Errors)
	}
	if len(fc.result.Warnings) > 0 {
		return fmt.Errorf("expected no warnings, got: %v", fc.result.Warnings)
	}
	return nil
}

func (fc *frontmatterContext) validationPassesWithWarning(warning string) error {
	if !fc.result.Passed() {
		return fmt.Errorf("expected validation to pass, got errors: %v", fc.result.Errors)
	}
	for _, w := range fc.result.Warnings {
		if strings.Contains(w, warning) {
			return nil
		}
	}
	return fmt.Errorf("expected warning containing %q, got warnings: %v", warning, fc.result.Warnings)
}

func (fc *frontmatterContext) thePublishDateIs(expected string) error {
	if fc.post.PublishDate == nil {
		return fmt.Errorf("expected publish_date %q, got nil", expected)
	}
	actual := fc.post.PublishDate.Format("2006-01-02")
	if actual != expected {
		return fmt.Errorf("expected publish_date %q, got %q", expected, actual)
	}
	return nil
}

func (fc *frontmatterContext) thePublishDateIsEmpty() error {
	if fc.post.PublishDate != nil {
		return fmt.Errorf("expected no publish_date, got %s", fc.post.PublishDate.Format("2006-01-02"))
	}
	return nil
}

func (fc *frontmatterContext) validationFailsWithError(errMsg string) error {
	if fc.result.Passed() {
		return fmt.Errorf("expected validation to fail with %q, but it passed", errMsg)
	}
	for _, e := range fc.result.Errors {
		if strings.Contains(e, errMsg) {
			return nil
		}
	}
	return fmt.Errorf("expected error containing %q, got errors: %v", errMsg, fc.result.Errors)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	fc := &frontmatterContext{}

	ctx.Step(`^a markdown file with frontmatter:$`, fc.aMarkdownFileWithFrontmatter)
	ctx.Step(`^the frontmatter is parsed$`, fc.theFrontmatterIsParsed)
	ctx.Step(`^the title is "([^"]*)"$`, fc.theTitleIs)
	ctx.Step(`^the slug is "([^"]*)"$`, fc.theSlugIs)
	ctx.Step(`^the subtitle is "([^"]*)"$`, fc.theSubtitleIs)
	ctx.Step(`^the series is "([^"]*)"$`, fc.theSeriesIs)
	ctx.Step(`^the tags are (\[.*\])$`, fc.theTagsAre)
	ctx.Step(`^the draft flag is (true|false)$`, fc.theDraftFlagIs)
	ctx.Step(`^the content is "([^"]*)"$`, fc.theContentIs)
	ctx.Step(`^the content starts with "([^"]*)"$`, fc.theContentStartsWith)
	ctx.Step(`^validation passes with no errors$`, fc.validationPassesWithNoErrors)
	ctx.Step(`^validation passes with warning "([^"]*)"$`, fc.validationPassesWithWarning)
	ctx.Step(`^the publish date is "([^"]*)"$`, fc.thePublishDateIs)
	ctx.Step(`^the publish date is empty$`, fc.thePublishDateIsEmpty)
	ctx.Step(`^validation fails with error "([^"]*)"$`, fc.validationFailsWithError)
}

// BDD: specs/validate_frontmatter.feature :: Feature: Validate frontmatter
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/validate_frontmatter.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
