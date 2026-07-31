package main

import (
	"os"
	"strings"
	"testing"
)

func TestWorkflowActionPins(t *testing.T) {
	t.Parallel()

	checkContains(t, ".github/workflows/validate-pr.yml", []string{
		"actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2",
		"actions/setup-go@4b73464bb391d4059bd26b0524d20df3927bd417 # v6.3.0",
		"cache: true",
		"cache-dependency-path: go.sum",
	})

	checkContains(t, ".github/workflows/terraform.yml", []string{
		"actions/checkout@de0fac2e4500dabe0009e67214ff5f5447ce83dd # v6.0.2",
		"actions/github-script@ed597411d8f924073f98dfc5c65a23a2325f34cd # v8",
	})
}

func TestOnlyFirstPartyBlogIsPublished(t *testing.T) {
	t.Parallel()

	if _, err := os.Stat(".github/workflows/publish.yml"); !os.IsNotExist(err) {
		t.Fatalf("external publishing workflow still exists")
	}

	checkContains(t, ".github/workflows/blog.yml", []string{
		"paths:\n      - 'posts/**'\n      - 'blog/**'",
		"schedule:\n    - cron: '0 5 * * *'",
		"pages deploy .svelte-kit/cloudflare --project-name=nuphirho-blog",
	})
}

func TestBlogDeploymentReportsStatusToTelegram(t *testing.T) {
	t.Parallel()

	checkContains(t, ".github/workflows/blog.yml", []string{
		"name: Build notification tool",
		"go build -o notify ./cmd/notify/",
		"if: always() && github.event_name != 'pull_request'",
		"TELEGRAM_BOT_TOKEN: ${{ secrets.TELEGRAM_BOT_TOKEN }}",
		"TELEGRAM_CHAT_ID: ${{ secrets.TELEGRAM_CHAT_ID }}",
		"DEPLOY_STATUS: ${{ job.status }}",
		"./notify \"Blog deployment ${DEPLOY_STATUS}: ${RUN_URL}\"",
	})
}

func TestTerraformWorkflowCommentsRealPlanOutput(t *testing.T) {
	t.Parallel()

	content := readFile(t, ".github/workflows/terraform.yml")

	if strings.Contains(content, "steps.plan.outputs.stdout") || strings.Contains(content, "steps.plan.outputs.stderr") {
		t.Fatalf("terraform workflow still references non-existent stdout/stderr outputs")
	}

	required := []string{
		"tee plan.txt",
		"echo \"exit_code=$exit_code\" >> \"$GITHUB_OUTPUT\"",
		"const plan = fs.readFileSync('terraform/plan.txt', 'utf8');",
		"steps.plan.outputs.exit_code",
	}
	for _, fragment := range required {
		if !strings.Contains(content, fragment) {
			t.Fatalf("terraform workflow missing %q", fragment)
		}
	}
}

func TestReadmeDescribesFirstPartyPublishing(t *testing.T) {
	t.Parallel()

	content := readFile(t, "README.md")

	for _, externalPlatform := range []string{"Hashnode", "Dev.to", "cross-post"} {
		if strings.Contains(content, externalPlatform) {
			t.Fatalf("README still refers to external publishing platform %q", externalPlatform)
		}
	}

	if !strings.Contains(content, "Posts with `draft: true` in the front matter are excluded from the blog build.") {
		t.Fatalf("README does not describe draft posts as excluded from the blog build")
	}
}

func TestValidateWorkflowRunsMutationTesting(t *testing.T) {
	t.Parallel()

	content := readFile(t, ".github/workflows/validate-pr.yml")

	required := []string{
		"go install github.com/avito-tech/go-mutesting/cmd/go-mutesting@v0.0.0-20251226130216-48d0401f00fb",
		"$(go env GOPATH)/bin/go-mutesting --exec-timeout=20 ./internal/frontmatter",
	}
	for _, fragment := range required {
		if !strings.Contains(content, fragment) {
			t.Fatalf("validate workflow missing %q", fragment)
		}
	}
}

func TestReadmeDocumentsMutationTestingInValidation(t *testing.T) {
	t.Parallel()

	content := readFile(t, "README.md")

	if !strings.Contains(content, "PR validation also runs mutation testing against `internal/frontmatter`.") {
		t.Fatalf("README does not document PR mutation testing")
	}
}

func TestGoVersionSupportsPinnedMutationTool(t *testing.T) {
	t.Parallel()

	content := readFile(t, "go.mod")

	if !strings.Contains(content, "go 1.25.5") {
		t.Fatalf("go.mod does not declare the minimum Go version required for the pinned mutation tool")
	}
}

func checkContains(t *testing.T, path string, fragments []string) {
	t.Helper()

	content := readFile(t, path)
	for _, fragment := range fragments {
		if !strings.Contains(content, fragment) {
			t.Fatalf("%s missing %q", path, fragment)
		}
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(content)
}
