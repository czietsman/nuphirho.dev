# nuphirho.dev Style Guide

## Purpose

This document defines the voice, tone, formatting, and editorial standards for nuphirho.dev. It serves as the definition of done for every post and as the brief for AI-assisted drafting.

---

## Voice and Identity

### Who is writing

Christo is a technologist with 20+ years of experience spanning mine seismology, antenna design expert systems, enterprise backup, and cybersecurity. He holds a Masters in Continuum Mechanics (Cum Laude, Stellenbosch University). He currently leads technology innovation at CyberSentriq.

He thinks in systems and architecture. He does not always communicate those ideas clearly on the first pass. AI helps him bridge that gap. The thinking is his. The clarity is a collaboration.

### The domain name

nuphirho takes its name from the Greek letters nu, phi, and rho, rooted in physics and mathematics. It reflects intellectual curiosity and engineering rigour, not a brand exercise.

### Core philosophy

Process matters more than technology. AI has changed the economics of rigorous engineering practices, making things like BDD specifications, mutation testing, and executable verification layers economically viable in ways they were not before. Security is a non-negotiable default, not an afterthought.

---

## Tone

### Do

- Be direct and plain. Say what you mean.
- Be friendly, enthusiastic, and upbeat.
- Maintain a positive outlook for the future of software engineering.
- Be honest. If something did not work, say so. If something is hard, say so.
- Be considered. Posts should feel thought through, not reactive or dashed off.
- Show the reasoning, not just the conclusion.
- Be comfortable with uncertainty. "I do not know yet" is a valid position.

### Do not

- Exaggerate. No superlatives unless they are earned and evidenced.
- Use fluff or filler to reach a word count.
- Be performative. No "I brilliantly discovered..." or false modesty.
- Write to impress. Write to communicate.
- Present conclusions without showing the path to them.

---

## Formatting

### Typography and punctuation

- No em dashes. Use commas, full stops, or restructure the sentence.
- No emoji. Ever.
- Use British English spelling (colour, organised, behaviour).
- Use Oxford commas for clarity.

### Structure

- Clean and minimal. No decorative formatting.
- Use headings to create clear structure, but do not over-nest.
- Paragraphs over bullet points where possible. Use lists only when the content genuinely requires enumeration.
- Keep paragraphs focused. One idea per paragraph.

### External links

- Keep external links minimal. Link only where the reference adds genuine value for the reader.
- Avoid linking to social media profiles (e.g. LinkedIn) inline. Names are sufficient; readers can search if they want to connect.
- Company or product links are acceptable when they provide useful context.
- Excessive outbound links can trigger spam filters on publishing platforms and dilute the reading experience.

### Code

- Do not embed large code blocks in posts. Link to the public GitHub repository.
- Short inline code references are fine for illustrating a point.
- Code exists to support the argument, not to fill space.

### Diagrams and images

- Must work on both light and dark backgrounds.
- Prefer SVG with transparent backgrounds or theme-aware colours.
- Mermaid diagrams are preferred where possible (version-controllable, theme-aware, supported by Hashnode).
- Screenshots should be used sparingly and considered for both theme contexts.
- Provide alt text for all images.

---

## Post Standards

### Length

- Target range: 1,200 to 1,800 words.
- This is a guideline, not a rule. Say what needs saying and stop.
- Complex topics may run longer. Simple topics may run shorter.
- Do not pad to reach a word count. Do not cut substance to stay under one.

### Point of view

- First person singular for personal experience and opinion ("I found...", "I implemented...").
- First person plural when speaking about shared industry challenges or the team context ("We are seeing...", "We adopted...").
- Never third person or passive voice where first person is more direct.

### Structure of a typical post

1. A clear opening that states what the post is about and why it matters.
2. Context and reasoning. Show the thinking, not just the result.
3. The substance. What happened, what was learned, what changed.
4. An honest reflection. What worked, what did not, what remains uncertain.
5. A forward-looking close. Where this leads next.

This is a guide, not a rigid template. Let the content shape the structure.

---

## Integrity

### Attribution

- Always attribute other people's work. Cite sources, credit tools, acknowledge inspiration.
- Link to original sources, not summaries or aggregators.
- If an idea came from a conversation, a talk, or a colleague, say so.

### AI assistance

- AI assists in research, drafting, and refinement for this blog.
- The thinking, decisions, direction, and accountability are the author's.
- This is stated openly on the about page and does not need repeating in every post.
- The framing: "I think in systems and architecture. I do not always communicate those ideas clearly on the first pass. AI helps me bridge that gap. The thinking is mine. The clarity is a collaboration."

### CyberSentriq content

- General approaches, lessons learned, and patterns may be shared.
- Proprietary details, internal data, and confidential information are never shared.
- Approval from the CPTO is required before publishing content that references CyberSentriq work directly.
- When in doubt, generalise. The principle matters more than the specific implementation.

### Corrections

- If a post contains an error, correct it transparently. Note what changed and when.
- Do not silently edit published content.

---

## Security

- HTTPS enforced at all layers: TLD (.dev), CDN (Cloudflare), platform (Hashnode).
- No credentials, tokens, or secrets in post content or screenshots.
- All infrastructure managed as code with secrets in GitHub Secrets, never in the repository.

---

## Design

### Theme

- Support both light and dark modes. Respect the reader's system preference.
- Default to light with a dark mode toggle available.
- All visual content must be tested against both themes.

### Reading experience

- Distraction-free. The content is the product.
- No ads. No popups. No newsletter gates blocking content.
- Fast loading. Minimal external dependencies.

---

## Distribution

### Source of truth

- blog.nuphirho.dev (Hashnode, custom subdomain)

### Cross-post targets

- Dev.to
- Medium

### Canonical URL

- Always set to blog.nuphirho.dev. Every cross-posted article must reference the canonical URL on the blog subdomain to protect SEO.

### Amplification

- LinkedIn (primary professional network, home of the 11-post series)

### Future consideration

- X (not yet, revisit once the content pipeline is established)

---

## Guardrails (keeping the author honest)

These are drawn from self-knowledge and psychometric assessment. They are reminders, not criticisms.

### Watch for tangents

Interest score: 86/100. The temptation to go deep on fascinating details is real. The word count target is a guardrail. Say the interesting thing, link out for depth, move on.

### Show the reasoning

Solo decision-making is a natural tendency. In blog posts, this shows up as presenting conclusions without the path that led to them. Always show the "why", not just the "what".

### Push the boat out

Publishing is a form of risk-taking. The blog exists partly to build the confidence muscle in a controlled environment. If a post feels slightly uncomfortable to publish, that is probably a good sign.

### Take time to process

Strong ideas need time to form. Draft, step away, return. The CI/CD pipeline supports this with draft states and review stages. Use them.

### Seek feedback before publishing

The tendency to work independently means posts could benefit from a second perspective. The pipeline should include a review step, even if the reviewer is AI-assisted.
