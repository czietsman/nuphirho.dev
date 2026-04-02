package tags_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/czietsman/nuphirho.dev/internal/tags"
)

type tagsContext struct {
	glossary    tags.Glossary
	canonical   []string
	mapResult   tags.MapResult
	validations []tags.TagValidation
}

func (tc *tagsContext) aTagGlossary(doc *godog.DocString) error {
	var err error
	tc.glossary, err = tags.ParseGlossary([]byte(doc.Content))
	return err
}

func (tc *tagsContext) canonicalTags(tagsJSON string) error {
	return json.Unmarshal([]byte(tagsJSON), &tc.canonical)
}

func (tc *tagsContext) tagsAreMappedFor(platform string) error {
	tc.mapResult = tc.glossary.MapTags(tc.canonical, tags.Platform(platform))
	return nil
}

func (tc *tagsContext) theMappedTagsAre(expected string) error {
	var expectedTags []string
	if err := json.Unmarshal([]byte(expected), &expectedTags); err != nil {
		return fmt.Errorf("invalid expected tags JSON: %v", err)
	}
	if len(tc.mapResult.Tags) != len(expectedTags) {
		return fmt.Errorf("expected %d tags %v, got %d tags %v", len(expectedTags), expectedTags, len(tc.mapResult.Tags), tc.mapResult.Tags)
	}
	for i, tag := range expectedTags {
		if tc.mapResult.Tags[i] != tag {
			return fmt.Errorf("tag %d: expected %q, got %q", i, tag, tc.mapResult.Tags[i])
		}
	}
	return nil
}

func (tc *tagsContext) aWarningIsProduced(expected string) error {
	for _, w := range tc.mapResult.Warnings {
		if strings.Contains(w, expected) {
			return nil
		}
	}
	return fmt.Errorf("expected warning containing %q, got: %v", expected, tc.mapResult.Warnings)
}

func (tc *tagsContext) noWarningsAreProduced() error {
	if len(tc.mapResult.Warnings) > 0 {
		return fmt.Errorf("expected no warnings, got: %v", tc.mapResult.Warnings)
	}
	return nil
}

func (tc *tagsContext) tagsAreValidatedFor(platform string) error {
	tc.validations = tc.glossary.ValidateTags(tc.canonical, tags.Platform(platform))
	return nil
}

func (tc *tagsContext) tagIsValidBecauseItHasAGlossaryMapping(tag string) error {
	for _, v := range tc.validations {
		if v.Tag == tag {
			if !v.Valid {
				return fmt.Errorf("expected %q to be valid, but it was invalid: %s", tag, v.Reason)
			}
			if !v.HasMapping {
				return fmt.Errorf("expected %q to have a glossary mapping, but it doesn't", tag)
			}
			return nil
		}
	}
	return fmt.Errorf("tag %q not found in validations", tag)
}

func (tc *tagsContext) tagIsInvalidBecauseItContainsHyphens(tag string) error {
	for _, v := range tc.validations {
		if v.Tag == tag {
			if v.Valid {
				return fmt.Errorf("expected %q to be invalid, but it was valid", tag)
			}
			if !strings.Contains(v.Reason, "hyphens") {
				return fmt.Errorf("expected reason to mention hyphens, got %q", v.Reason)
			}
			return nil
		}
	}
	return fmt.Errorf("tag %q not found in validations", tag)
}

func (tc *tagsContext) tagIsValid(tag string) error {
	for _, v := range tc.validations {
		if v.Tag == tag {
			if !v.Valid {
				return fmt.Errorf("expected %q to be valid, but it was invalid: %s", tag, v.Reason)
			}
			return nil
		}
	}
	return fmt.Errorf("tag %q not found in validations", tag)
}

func (tc *tagsContext) allTagsAreValid() error {
	for _, v := range tc.validations {
		if !v.Valid {
			return fmt.Errorf("expected all tags valid, but %q is invalid: %s", v.Tag, v.Reason)
		}
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	tc := &tagsContext{}

	ctx.Step(`^a tag glossary:$`, tc.aTagGlossary)
	ctx.Step(`^canonical tags (\[.*\])$`, tc.canonicalTags)
	ctx.Step(`^tags are mapped for "([^"]*)"$`, tc.tagsAreMappedFor)
	ctx.Step(`^the mapped tags are (\[.*\])$`, tc.theMappedTagsAre)
	ctx.Step(`^a warning is produced: "([^"]*)"$`, tc.aWarningIsProduced)
	ctx.Step(`^no warnings are produced$`, tc.noWarningsAreProduced)
	ctx.Step(`^tags are validated for "([^"]*)"$`, tc.tagsAreValidatedFor)
	ctx.Step(`^"([^"]*)" is valid because it has a glossary mapping$`, tc.tagIsValidBecauseItHasAGlossaryMapping)
	ctx.Step(`^"([^"]*)" is invalid because it contains hyphens$`, tc.tagIsInvalidBecauseItContainsHyphens)
	ctx.Step(`^"([^"]*)" is valid$`, tc.tagIsValid)
	ctx.Step(`^all tags are valid$`, tc.allTagsAreValid)
}

// BDD: specs/tag_glossary.feature :: Feature: Tag glossary
func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/tag_glossary.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
