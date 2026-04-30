package main

// Page holds the data needed to render a single page template.
type Page struct {
	Title           string
	Description     string
	Canonical       string
	OGTitle         string
	OGDescription   string
	OGURL           string
	OGType          string
	BodyClass       string
	MainID          string
	MainClass       string
	CurrentNav      string   // "blog", "about", "words", or "" for none
	ShowFooter      bool
	PageCSS         []string // additional CSS paths relative to site root
	PageJS          []string // additional JS paths relative to site root
	OutputPath      string   // e.g. "index.html", "about/index.html"
	TemplateName    string   // filename in templates/pages/
	Content         string   // not used by template; reserved for future
}

func pages() []Page {
	return []Page{
		{
			Title:           "nuphirho.dev",
			Description:     "Enterprise-grade engineering at startup speed. Christo Zietsman — Director of Technology Innovation, researcher, and writer.",
			Canonical:       "https://nuphirho.dev/",
			OGTitle:         "Christo Zietsman | nuphirho.dev",
			OGDescription:   "Enterprise-grade engineering at startup speed. Director of Technology Innovation, researcher, and writer.",
			OGURL:           "https://nuphirho.dev/",
			OGType:          "website",
			BodyClass:       "landing-body",
			MainClass:       "landing-main",
			ShowFooter:      true,
			OutputPath:      "index.html",
			TemplateName:    "index.html",
		},
		{
			Title:           "Who am I | nuphirho.dev",
			Description:     "Who I am, what I believe, and what I am building. Christo Zietsman — Director of Technology Innovation, researcher, and writer.",
			Canonical:       "https://nuphirho.dev/about",
			OGTitle:         "About | Christo Zietsman",
			OGDescription:   "Who I am, what I believe, and what I am building. Director of Technology Innovation, researcher, and writer.",
			OGURL:           "https://nuphirho.dev/about",
			OGType:          "website",
			MainID:          "main",
			MainClass:       "site-main",
			CurrentNav:      "about",
			ShowFooter:      true,
			OutputPath:      "about/index.html",
			TemplateName:    "about.html",
		},
		{
			Title:           "Words of Meaning — nuphirho.dev",
			Description:     "A living glossary. Not definitions in the dictionary sense. Explanations of why each word was chosen and what it carries.",
			Canonical:       "https://nuphirho.dev/words-of-meaning",
			OGType:          "website",
			MainID:          "main",
			MainClass:       "site-main wom-content",
			CurrentNav:   "words",
			ShowFooter:   true,
			OutputPath:   "words-of-meaning/index.html",
			TemplateName:    "words-of-meaning.html",
		},
		{
			Title:           "nuphirho.dev — Publishing Roadmap",
			Description:     "Publishing roadmap for nuphirho.dev.",
			Canonical:       "https://nuphirho.dev/roadmap",
			OGType:          "website",
			MainID:          "main",
			MainClass:       "site-main roadmap-content",
			ShowFooter:      true,
			PageCSS:         []string{"css/roadmap.css"},
			PageJS:          []string{"js/roadmap.js"},
			OutputPath:      "roadmap/index.html",
			TemplateName:    "roadmap.html",
		},
                {
			Title:        "Novel Findings | nuphirho.dev",
			Description:  "Confirmed novel contributions from the nuphirho.dev research programme.",
			Canonical:    "https://nuphirho.dev/novel-findings",
			OGTitle:      "Novel Findings | nuphirho.dev",
			OGDescription: "Confirmed novel contributions from the nuphirho.dev research programme.",
			OGURL:        "https://nuphirho.dev/novel-findings",
			OGType:       "website",
			MainID:       "main",
			MainClass:    "site-main",
			CurrentNav:   "novel-findings",
			ShowFooter:   true,
			OutputPath:   "novel-findings/index.html",
			TemplateName: "novel-findings.html",
		},
		{
			Title:           "Cookies | nuphirho.dev",
                        Description:     "Cookie policy for nuphirho.dev. We do not set cookies directly. Cloudflare may set strictly necessary technical cookies.",
                        Canonical:       "https://nuphirho.dev/cookies",
                        OGTitle:         "Cookie policy for nuphirho.dev",
                        OGDescription:   "Cookie policy for nuphirho.dev. We do not set cookies directly. Cloudflare may set strictly necessary technical cookies.",
                        OGURL:           "https://nuphirho.dev/cookies",
                        OGType:          "website",
                        MainID:          "main",
                        MainClass:       "site-main",
                        CurrentNav:      "cookies",
                        ShowFooter:      true,
                        OutputPath:      "cookies/index.html",
                        TemplateName:    "cookies.html",
		},
		{
                        Title:           "Privacy Policy | nuphirho.dev",
                        Description:     "Privacy policy for nuphirho.dev. We do not collect personal data. No account registration required.",
                        Canonical:       "https://nuphirho.dev/privacy",
                        OGTitle:         "Privacy policy for nuphirho.dev",
                        OGDescription:   "Privacy policy for nuphirho.dev. We do not collect personal data. No account registration required.",
                        OGURL:           "https://nuphirho.dev/privacy",
                        OGType:          "website",
                        MainID:          "main",
                        MainClass:       "site-main",
                        CurrentNav:      "privacy",
                        ShowFooter:      true,
                        OutputPath:      "privacy/index.html",
                        TemplateName:    "privacy.html",
                },

	}
}
