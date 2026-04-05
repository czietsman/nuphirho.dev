// Package tags loads a tag glossary and maps canonical tags to
// platform-specific equivalents.
package tags

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

// Platform identifies a publishing target.
type Platform string

const (
	DevTo    Platform = "devto"
	Hashnode Platform = "hashnode"
)

// devtoMaxTags is the maximum number of tags Dev.to allows per article.
const devtoMaxTags = 4

var devtoValidTag = regexp.MustCompile(`^[a-z0-9]+$`)

// Glossary maps canonical tag names to per-platform overrides.
// Tags not present in the glossary pass through unchanged.
type Glossary map[string]map[string]string

// LoadGlossary reads a tag glossary JSON file.
func LoadGlossary(path string) (Glossary, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading tag glossary: %w", err)
	}
	return ParseGlossary(data)
}

// ParseGlossary parses tag glossary JSON from raw bytes.
func ParseGlossary(data []byte) (Glossary, error) {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parsing tag glossary: %w", err)
	}

	g := make(Glossary)
	for key, val := range raw {
		// Skip non-object entries (e.g. "_comment": "string")
		var mapping map[string]string
		if err := json.Unmarshal(val, &mapping); err != nil {
			continue
		}
		g[key] = mapping
	}
	return g, nil
}

// MapResult holds the mapped tags and any warnings produced during mapping.
type MapResult struct {
	Tags     []string
	Warnings []string
}

// MapTags maps canonical tags to their platform-specific equivalents.
// Tags without a glossary entry pass through unchanged.
// For Dev.to, the result is truncated to 4 tags with a warning listing
// which tags were dropped.
func (g Glossary) MapTags(tags []string, platform Platform) MapResult {
	mapped := make([]string, 0, len(tags))
	for _, tag := range tags {
		if mapping, ok := g[tag]; ok {
			if m, ok := mapping[string(platform)]; ok {
				mapped = append(mapped, m)
				continue
			}
		}
		mapped = append(mapped, tag)
	}

	var warnings []string
	if platform == DevTo && len(mapped) > devtoMaxTags {
		dropped := mapped[devtoMaxTags:]
		warnings = append(warnings, fmt.Sprintf("%d tags dropped for Dev.to: %s", len(dropped), joinTags(dropped)))
		mapped = mapped[:devtoMaxTags]
	}

	return MapResult{Tags: mapped, Warnings: warnings}
}

func joinTags(tags []string) string {
	result := ""
	for i, t := range tags {
		if i > 0 {
			result += ", "
		}
		result += t
	}
	return result
}

// PostTagResult holds the outcome of validating a post's tags for Dev.to.
type PostTagResult struct {
	Skipped  bool
	Valid    bool
	Unmapped []string
}

// ValidatePostTags checks whether a post's tags are valid for Dev.to
// cross-posting. Draft posts are skipped. Returns unmapped hyphenated
// tags that need glossary entries.
func (g Glossary) ValidatePostTags(postTags []string, draft bool) PostTagResult {
	if draft {
		return PostTagResult{Skipped: true, Valid: true}
	}

	validations := g.ValidateTags(postTags, DevTo)
	var unmapped []string
	for _, v := range validations {
		if !v.Valid {
			unmapped = append(unmapped, v.Tag)
		}
	}

	return PostTagResult{
		Valid:    len(unmapped) == 0,
		Unmapped: unmapped,
	}
}

// TagValidation holds the result of validating a single tag.
type TagValidation struct {
	Tag       string
	Valid     bool
	HasMapping bool
	Reason   string
}

// ValidateTags checks whether each tag is valid for the given platform.
// For Dev.to, tags must be alphanumeric (no hyphens) unless the glossary
// provides a mapping.
func (g Glossary) ValidateTags(tags []string, platform Platform) []TagValidation {
	results := make([]TagValidation, len(tags))
	for i, tag := range tags {
		results[i] = g.validateTag(tag, platform)
	}
	return results
}

func (g Glossary) validateTag(tag string, platform Platform) TagValidation {
	v := TagValidation{Tag: tag, Valid: true}

	if platform != DevTo {
		return v
	}

	// Check if glossary has a Dev.to mapping
	if mapping, ok := g[tag]; ok {
		if _, ok := mapping[string(DevTo)]; ok {
			v.HasMapping = true
			return v
		}
	}

	// No mapping: validate the raw tag
	if !devtoValidTag.MatchString(tag) {
		v.Valid = false
		v.Reason = "contains hyphens"
		return v
	}

	return v
}
