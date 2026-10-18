package waf_test

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// This checks the configured predicates, not Cloudflare's remote expression parser.
func TestProbeRuleScope(t *testing.T) {
	data, err := os.ReadFile("waf.tf")
	if err != nil {
		t.Fatal(err)
	}
	config := string(data)
	if !strings.Contains(config, `phase       = "http_request_firewall_custom"`) || !strings.Contains(config, `action      = "block"`) {
		t.Fatal("expected a blocking custom WAF rule")
	}
	hosts := regexp.MustCompile(`http.host in \{([^}]+)\}`).FindStringSubmatch(config)
	paths := regexp.MustCompile(`lower\(http.request.uri.path\) in \{([^}]+)\}`).FindStringSubmatch(config)
	prefixes := regexp.MustCompile(`starts_with\(lower\(http.request.uri.path\), "([^"]+)"\)`).FindAllStringSubmatch(config, -1)
	if len(hosts) != 2 || len(paths) != 2 || len(prefixes) != 3 {
		t.Fatal("unexpected predicate structure")
	}
	expression := regexp.MustCompile(`(?s)<<-EOT\s*(.*?)\s*EOT`).FindStringSubmatch(config)
	if len(expression) != 2 {
		t.Fatal("missing expression")
	}
	shape := regexp.MustCompile(`http.host in \{[^}]+\}`).ReplaceAllString(expression[1], "H")
	shape = regexp.MustCompile(`lower\(http.request.uri.path\) in \{[^}]+\}`).ReplaceAllString(shape, "P")
	shape = regexp.MustCompile(`starts_with\(lower\(http.request.uri.path\), "[^"]+"\)`).ReplaceAllString(shape, "P")
	shape = strings.Join(strings.Fields(shape), "")
	if shape != "(Hand(PorPorPorP))" {
		t.Fatalf("unexpected expression structure: %s", shape)
	}
	contains := func(set, value string) bool { return strings.Contains(set, `"`+value+`"`) }
	blocked := func(host, path string) bool {
		if !contains(hosts[1], host) {
			return false
		}
		path = strings.ToLower(path)
		if contains(paths[1], path) {
			return true
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(path, prefix[1]) {
				return true
			}
		}
		return false
	}
	for _, host := range []string{"nuphirho.dev", "www.nuphirho.dev"} {
		for _, path := range []string{"/.env", "/.env.local", "/.env.production", "/.env.backup", "/.git", "/.git/config", "/wp-login.php", "/wp-admin", "/wp-admin/", "/wp-admin/install.php", "/xmlrpc.php", "/.ENV", "/WP-LOGIN.PHP"} {
			if !blocked(host, path) {
				t.Errorf("probe allowed: %s%s", host, path)
			}
		}
		for _, path := range []string{"/", "/about", "/privacy", "/cookies", "/papers", "/roadmap", "/business-card/contact.vcf", "/words-of-meaning", "/novel-findings", "/login", "/logout", "/auth/callback", "/auth/logout", "/home", "/evaluate", "/api/waitlist", "/api/waitlist/verify", "/api/waitlist/unsubscribe", "/turnstile-config", "/.well-known/acme-challenge/token", "/.well-known/security.txt", "/_app/immutable/app.js", "/robots.txt", "/sitemap.xml", "/unknown", "/.hidden", "/.environment", "/.github/config", "/wp-admin-guide", "/articles/wp-login.php", "/api/.env"} {
			if blocked(host, path) {
				t.Errorf("legitimate or unrelated path blocked: %s%s", host, path)
			}
		}
	}
	// Verify every concrete route and static asset in this checkout remains outside the rule.
	for _, dir := range []string{"../src/routes", "../static", "../main-site/src/routes", "../main-site/static"} {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			continue
		}
		err := filepath.WalkDir(dir, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if entry.IsDir() {
				return nil
			}
			relative, err := filepath.Rel(dir, path)
			if err != nil {
				return err
			}
			if strings.Contains(dir, "src/routes") {
				if !strings.HasPrefix(filepath.Base(path), "+") {
					return nil
				}
				relative = filepath.Dir(relative)
			}
			segments := []string{}
			for _, segment := range strings.Split(filepath.ToSlash(relative), "/") {
				if segment == "." || strings.HasPrefix(segment, "(") {
					continue
				}
				if strings.Contains(segment, "[") {
					segment = "example"
				}
				segments = append(segments, segment)
			}
			route := "/" + strings.Join(segments, "/")
			if blocked("nuphirho.dev", route) {
				t.Errorf("repository route or asset blocked: %s (%s)", route, path)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	for _, host := range []string{"blog.nuphirho.dev", "auth.nuphirho.dev", "api.nuphirho.dev", "preview.pages.dev", "unrelated.example"} {
		if blocked(host, "/.env") {
			t.Errorf("unrelated host blocked: %s", host)
		}
	}
}
