# Prompt: Convert LinkedIn Post to Blog Post

## Context

Read the following file before starting:

- `docs/STYLE_GUIDE.md`

## Task

Convert the LinkedIn post provided below into a blog post for nuphirho.dev.

## Style and voice

- British English. No em dashes. No emoji. Oxford commas.
- Direct, friendly, enthusiastic tone. No fluff, no exaggeration.
- First person singular for personal experience. First person plural for shared industry challenges.
- Honest about what works and what does not.
- AI-assisted writing is acknowledged openly as a philosophy. It does not need restating in every post.

## Conversion guidelines

- **Expand to blog depth.** LinkedIn rewards punchy, short-form writing. The blog allows more reasoning, context, and nuance. Add the "why" behind claims and positions. Show the thinking, not just the conclusion.
- **Preserve the original voice.** The LinkedIn post has the author's natural tone, which is direct and slightly informal. Do not sanitise, over-polish, or make it sound generic.
- **Target length: 1,200 to 1,800 words.** Say what needs saying and stop. Do not pad.
- **Paragraphs over bullet points.** Restructure any LinkedIn-style bullet lists into flowing prose unless enumeration is genuinely required.
- **Link to sources.** Add references, attributions, or links where the post makes claims or references other people's work.
- **Link to GitHub for code.** If the post references code, implementations, or configurations, link to the relevant file in the public repo rather than embedding large code blocks.

## Output format

Output the post as a markdown file with the following frontmatter:

```yaml
---
title: ""
slug: ""
draft: true
tags: []
---
```

- **title**: a clear, descriptive title. Not clickbait. Not clever for its own sake.
- **slug**: lowercase, hyphenated, derived from the title.
- **draft**: always true on first output. The author publishes deliberately.
- **tags**: 3 to 5 relevant tags.

## Canonical URL

The blog at blog.nuphirho.dev is the permanent home for all content. LinkedIn is the amplification channel. The blog post is the canonical version.

## Quality check before output

- [ ] British English spelling throughout
- [ ] No em dashes
- [ ] No emoji
- [ ] No exaggeration or fluff
- [ ] Original voice and energy preserved from the LinkedIn post
- [ ] Reasoning and context added, not just the original claims restated at greater length
- [ ] All claims attributed where appropriate
- [ ] Within 1,200 to 1,800 words
- [ ] Frontmatter present and complete

---

## LinkedIn post to convert

[Paste the LinkedIn post here]
