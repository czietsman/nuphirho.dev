package probe_test

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/czietsman/nuphirho.dev/internal/probe"
)

// --- Fakes ---

type fakeHashnodeProber struct {
	credentialsErr    error
	publicationErr    error
	createDraftErr    error
	deleteDraftErr    error
	createdDraftID    string
	deletedDraftIDs   []string
	createDraftCalled bool
}

func (f *fakeHashnodeProber) CheckCredentials() error {
	return f.credentialsErr
}

func (f *fakeHashnodeProber) CheckPublication() error {
	return f.publicationErr
}

func (f *fakeHashnodeProber) CreateProbeDraft() (string, error) {
	f.createDraftCalled = true
	if f.createDraftErr != nil {
		return "", f.createDraftErr
	}
	f.createdDraftID = "probe-draft-001"
	return f.createdDraftID, nil
}

func (f *fakeHashnodeProber) DeleteDraft(id string) error {
	f.deletedDraftIDs = append(f.deletedDraftIDs, id)
	return f.deleteDraftErr
}

type fakeDevtoProber struct {
	credentialsErr     error
	createArticleErr   error
	deleteArticleErr   error
	createdArticleID   int
	deletedArticleIDs  []int
	createArticleCalled bool
}

func (f *fakeDevtoProber) CheckCredentials() error {
	return f.credentialsErr
}

func (f *fakeDevtoProber) CreateProbeArticle() (int, error) {
	f.createArticleCalled = true
	if f.createArticleErr != nil {
		return 0, f.createArticleErr
	}
	f.createdArticleID = 99999
	return f.createdArticleID, nil
}

func (f *fakeDevtoProber) DeleteArticle(id int) error {
	f.deletedArticleIDs = append(f.deletedArticleIDs, id)
	return f.deleteArticleErr
}

// --- godog context ---

type probeContext struct {
	hn     *fakeHashnodeProber
	dt     *fakeDevtoProber
	result *probe.Result
	output string
}

func (pc *probeContext) reset() {
	pc.hn = &fakeHashnodeProber{}
	pc.dt = &fakeDevtoProber{}
	pc.result = nil
	pc.output = ""
}

// --- Given steps ---

func (pc *probeContext) validCredentialsAreConfiguredForAllPlatforms() error {
	// Default fakes return no errors -- all probes pass
	return nil
}

func (pc *probeContext) validHashnodeCredentialsAreConfigured() error {
	// Hashnode fakes default to no errors
	return nil
}

func (pc *probeContext) validDevToCredentialsAreConfigured() error {
	// Dev.to fakes default to no errors
	return nil
}

func (pc *probeContext) aValidHashnodeTokenIsConfigured() error {
	// Credentials pass, but publication may fail
	return nil
}

func (pc *probeContext) anInvalidHashnodeTokenIsConfigured() error {
	pc.hn.credentialsErr = fmt.Errorf("authentication failed")
	return nil
}

func (pc *probeContext) anUnknownPublicationIDIsConfigured() error {
	pc.hn.publicationErr = fmt.Errorf("publication not found")
	return nil
}

func (pc *probeContext) theHashnodeAPIIsUnreachable() error {
	pc.hn.credentialsErr = fmt.Errorf("connection refused")
	return nil
}

func (pc *probeContext) anInvalidDevToAPIKeyIsConfigured() error {
	pc.dt.credentialsErr = fmt.Errorf("authentication failed")
	return nil
}

func (pc *probeContext) theDevToAPIIsUnreachable() error {
	pc.dt.credentialsErr = fmt.Errorf("connection refused")
	return nil
}

func (pc *probeContext) hashnodeDraftDeletionWillFail() error {
	pc.hn.deleteDraftErr = fmt.Errorf("delete draft failed: 500 internal server error")
	return nil
}

func (pc *probeContext) devToArticleDeletionWillFail() error {
	pc.dt.deleteArticleErr = fmt.Errorf("delete article failed: 500 internal server error")
	return nil
}

// --- When steps ---

func (pc *probeContext) theProbeCommandRuns() error {
	var buf bytes.Buffer
	pc.result = probe.ProbeAll(&buf, pc.hn, pc.dt)
	pc.output = buf.String()
	return nil
}

// --- Then steps ---

func (pc *probeContext) allProbesPass() error {
	for _, c := range pc.result.Checks {
		if c.Status == probe.StatusFail {
			return fmt.Errorf("expected all probes to pass, but %s/%s failed: %s", c.Platform, c.Name, c.Detail)
		}
	}
	return nil
}

func (pc *probeContext) theExitCodeIs(code int) error {
	if pc.result.ExitCode != code {
		return fmt.Errorf("expected exit code %d, got %d", code, pc.result.ExitCode)
	}
	return nil
}

func (pc *probeContext) theOutputContains(expected string) error {
	if !strings.Contains(pc.output, expected) {
		return fmt.Errorf("expected output to contain %q, got:\n%s", expected, pc.output)
	}
	return nil
}

func (pc *probeContext) theHashnodeCredentialsProbeFails() error {
	for _, c := range pc.result.Checks {
		if c.Platform == "hashnode" && c.Name == "credentials" && c.Status == probe.StatusFail {
			return nil
		}
	}
	return fmt.Errorf("expected Hashnode credentials probe to fail")
}

func (pc *probeContext) theHashnodePublicationProbeFails() error {
	for _, c := range pc.result.Checks {
		if c.Platform == "hashnode" && c.Name == "publication" && c.Status == probe.StatusFail {
			return nil
		}
	}
	return fmt.Errorf("expected Hashnode publication probe to fail")
}

func (pc *probeContext) theDevToCredentialsProbeFails() error {
	for _, c := range pc.result.Checks {
		if c.Platform == "devto" && c.Name == "credentials" && c.Status == probe.StatusFail {
			return nil
		}
	}
	return fmt.Errorf("expected Dev.to credentials probe to fail")
}

func (pc *probeContext) theHashnodeDraftCreatedByTheProbeIsDeleted() error {
	if !pc.hn.createDraftCalled {
		return fmt.Errorf("createProbeDraft was never called")
	}
	if len(pc.hn.deletedDraftIDs) == 0 {
		return fmt.Errorf("no drafts were deleted")
	}
	for _, id := range pc.hn.deletedDraftIDs {
		if id == pc.hn.createdDraftID {
			return nil
		}
	}
	return fmt.Errorf("expected draft %q to be deleted, deleted: %v", pc.hn.createdDraftID, pc.hn.deletedDraftIDs)
}

func (pc *probeContext) theDevToArticleCreatedByTheProbeIsDeleted() error {
	if !pc.dt.createArticleCalled {
		return fmt.Errorf("createProbeArticle was never called")
	}
	if len(pc.dt.deletedArticleIDs) == 0 {
		return fmt.Errorf("no articles were deleted")
	}
	for _, id := range pc.dt.deletedArticleIDs {
		if id == pc.dt.createdArticleID {
			return nil
		}
	}
	return fmt.Errorf("expected article %d to be deleted, deleted: %v", pc.dt.createdArticleID, pc.dt.deletedArticleIDs)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	pc := &probeContext{}

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		pc.reset()
		return ctx, nil
	})

	// Given
	ctx.Step(`^valid credentials are configured for all platforms$`, pc.validCredentialsAreConfiguredForAllPlatforms)
	ctx.Step(`^valid Hashnode credentials are configured$`, pc.validHashnodeCredentialsAreConfigured)
	ctx.Step(`^valid Dev\.to credentials are configured$`, pc.validDevToCredentialsAreConfigured)
	ctx.Step(`^a valid Hashnode token is configured$`, pc.aValidHashnodeTokenIsConfigured)
	ctx.Step(`^an invalid Hashnode token is configured$`, pc.anInvalidHashnodeTokenIsConfigured)
	ctx.Step(`^an unknown publication ID is configured$`, pc.anUnknownPublicationIDIsConfigured)
	ctx.Step(`^the Hashnode API is unreachable$`, pc.theHashnodeAPIIsUnreachable)
	ctx.Step(`^an invalid Dev\.to API key is configured$`, pc.anInvalidDevToAPIKeyIsConfigured)
	ctx.Step(`^the Dev\.to API is unreachable$`, pc.theDevToAPIIsUnreachable)
	ctx.Step(`^Hashnode draft deletion will fail$`, pc.hashnodeDraftDeletionWillFail)
	ctx.Step(`^Dev\.to article deletion will fail$`, pc.devToArticleDeletionWillFail)

	// When
	ctx.Step(`^the probe command runs$`, pc.theProbeCommandRuns)

	// Then
	ctx.Step(`^all probes pass$`, pc.allProbesPass)
	ctx.Step(`^the exit code is (\d+)$`, pc.theExitCodeIs)
	ctx.Step(`^the output contains "([^"]*)"$`, pc.theOutputContains)
	ctx.Step(`^the Hashnode credentials probe fails$`, pc.theHashnodeCredentialsProbeFails)
	ctx.Step(`^the Hashnode publication probe fails$`, pc.theHashnodePublicationProbeFails)
	ctx.Step(`^the Dev\.to credentials probe fails$`, pc.theDevToCredentialsProbeFails)
	ctx.Step(`^the Hashnode draft created by the probe is deleted$`, pc.theHashnodeDraftCreatedByTheProbeIsDeleted)
	ctx.Step(`^the Dev\.to article created by the probe is deleted$`, pc.theDevToArticleCreatedByTheProbeIsDeleted)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../specs/contract_probes.feature"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
