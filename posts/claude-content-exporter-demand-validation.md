---
title: "I built a tool to get my work out of Claude Projects. Should I open source it?"
slug: "claude-content-exporter-demand-validation"
tags: [ai-assisted-development, personal-practice, agentic-development]
subtitle: "A Chrome extension that exports conversation data to your repository, built by going around limits the platform never sanctioned, and an honest question about whether it should be public"
publish_date: 2026-08-21
format: article
cover_image: claude-content-exporter-demand-validation.png
cover_image_prompt: |
  Not generated from this field's prompt; the actual image in use is a deliberate exception. It is a full explainer infographic: hand-drawn line art on aged cream paper, a labelled diagram of the extension's request/response flow (Claude Projects, Claude servers, Chrome extension, local filesystem), boxed callouts ("Locked in", "Reality check", "Honest truth", "Should I open source this?"), multiple coloured icons, bullet lists, and text rendered throughout. Produced by accident during a separate generation attempt and adopted deliberately per Director instruction, 16 August 2026. The specific prompt or process that produced it is not recorded; see meta/artifacts/cover-image-style-guide.md's "Alternate style: full explainer infographic" section.
cover_post: |
  Claude Projects has no native export. Months of session transcripts, research artifacts, and governing documents lived there, in a format you cannot version, index, or carry into a different workflow.

  I built a Chrome extension to solve the problem. It intercepts the fetch layer, captures API responses, and exports conversation transcripts, artifact files, and project documents to your local filesystem.

  Chrome-only, manual, and medium-fragile against platform updates. That is an honest description of a personal tool built to solve one person's problem.

  It works by going around limits Claude Projects was never built to expose, which rules out turning it into a product. It does not rule out sharing it.

  Should I open source this?
linkedin_url: https://www.linkedin.com/pulse/i-built-tool-get-my-work-out-claude-projects-should-open-zietsman-ierff
---

When I decided to move my AI research and writing practice from Claude Projects to Claude Code, the first obstacle was straightforward: how do I get the work out?

Months of session transcripts, research artifacts, governing documents, role definitions, and editorial decisions lived inside Claude Projects. The platform has no native export. Conversations are not yours in any format you can version, index, or carry into a different workflow. They live where they live.

I built a Chrome extension to solve the problem.

The extension is called Claude Content Exporter. It works by injecting a script into the same JavaScript execution context as the claude.ai page and wrapping the browser's fetch function, so that every API call the page makes is intercepted as it flows through. A separate content script accumulates the captured data. The popup triggers the export to your local filesystem.

What it captures: conversation transcripts formatted as markdown, the raw JSON response for machine-readable access, artifact files (code, HTML, SVG, markdown) written out by type with a deduplication pass, and project-level exports including project knowledge documents, project memory, and project metadata.

What it does not capture: tool use and tool result blocks (these are stripped), conversations you have not opened in the current browser session (there is no background polling), binary image content in most cases, anything that flows through a streaming connection rather than standard fetch responses, and Claude Code sessions run through claude.ai, a gap worth closing.

The export is always a manual act. You open the extension popup and click export. It captures what has flowed through your browser session since you opened that tab. A background sync would require credentials the extension should not hold.

Chrome-only, and medium-fragile against platform updates. The fetch interception approach works under the current Manifest V3 extension format. It does not work in Firefox, Safari, or Edge. The URL patterns the extension matches against claude.ai's internal API could break on any platform update. The JSON shapes it parses have already changed at least once during development, which is why the parser carries multiple fallbacks. If Anthropic ships a significant interface change, something in the export is likely to break silently until I update the matching logic. There are no automated tests against live API shapes.

The problem is not specific to Claude. It is the structural gap between where AI conversation platforms keep your work and where your work needs to live. Conversations have no version history. You cannot see what changed between session three and session seven. You cannot run a search that returns results across all your sessions without loading each one individually. You cannot reference a specific decision and see what preceded it, because the chronology exists only inside the platform.

For a simple use case, this is acceptable. For a programme where the output of one session is the formal input to the next, where decisions need to be auditable, and where knowledge accumulates across months and roles, it is a structural limitation. The data portability problem for AI platforms is not solved, and the current generation of tools mostly does not try to solve it.

There is also a harder constraint than any of the technical ones. The extension works by intercepting API calls Claude Projects was never built to expose, which is not something Anthropic's terms of service sanction for programmatic use. That closes off turning this into a product, regardless of how common the underlying problem turns out to be. Selling access to something that only works by going around a platform's own terms is not a business. It is a liability.

What is still open is whether to open source it. The problem it solves, work you cannot get out of a platform in a form you can version or audit, is real and not specific to me. Making the extension public could help anyone else stuck behind the same wall. It would also make it easier for more people to do, at scale, the exact thing that was never sanctioned in the first place, which is a different risk than one person quietly solving their own problem.

If you are running a serious AI-assisted workflow and have run into the walled-garden problem with a conversation platform, I want to know. Specifically: have you built something to get structured data out, have you switched platforms because the data portability was insufficient, or have you accepted the constraint as the cost of using the tool?

The extension exists, privately, for now. Should I open source it?
