---
title: "What the industry has figured out about AGENTS.md"
slug: "agents-md-industry-pillars"
tags: [ai-governance, software-engineering, agents, promptq]
series: "AGENTS.md: the functional and quality layers"
series_part: 1
publish_date: 2026-07-20
stop_slop: 41/50
toulmin: Track A 6/6, Track B 4/6
cover_image: agents-md-industry-pillars.png
cover_image_prompt: |
  A wall of plain text files stacked in rows, each slightly different, viewed straight on. Muted palette, slate blues and off-whites. One file in the foreground is slightly separated from the rest, not highlighted, just apart. Clean, technical, slightly clinical. No people, no screens, no colour accents.
cover_post: |
  Over the past year, the practitioner community has converged on ten pillars for writing AGENTS.md, CLAUDE.md, and their equivalents: the instruction files that tell an AI coding agent how to operate in your repository.

  The pillars are real. Token budget and organisation. Functional content. Safety and security. All three groups are well-understood and increasingly well-specified.

  The evidence base behind them is thin, though. Two peer-reviewed empirical studies from 2025 provide the only direct evidence about how instruction file content affects agent performance. The remaining nine pillars rest on practitioner norms: vendor documentation, community guides, and blog posts. The advice is reasonable. It is not yet grounded the way a thirty-year certification tradition is grounded.

  And there is a layer the pillars do not address at all.

  First post in a short series on what the industry has figured out about these files, and what it has not.
linkedin_url:
---

Over the past year a quiet consensus has been forming around AGENTS.md, CLAUDE.md, and their equivalents: the instruction files that tell an AI coding agent how to operate in your repository.

---

The practitioner community has converged on ten pillars for writing these files. They cluster into three groups.

The first group is about token budget and organisation. Keep the root file short. Anthropic's own guidance recommends a ceiling of roughly 150-200 lines before splitting into subdirectory-scoped files. Load additional context on demand rather than at initialisation. This is the progressive disclosure pattern codified in the Anthropic Agent Skills standard in December 2025 and since adopted by OpenAI, Google, Cursor, and others: the root file acts as an index and routing entry point; detailed rules load only when the agent needs them for a specific task. Large always-loaded context degrades model performance: important instructions get lost.

The second group is about functional content. State what the agent does, what the architecture is, and what the conventions are. Build and run commands, implementation details, and architecture are the most common content types in production instruction files, appearing in roughly two-thirds of files studied in a 2025 empirical analysis of 2,303 repository context files (arXiv:2511.12884). Specificity matters: instructions written as direct imperatives perform better than generic, verbose paragraphs. A separate 2025 study (arXiv:2602.11988) found that LLM-generated context files, the ones that trade on sounding complete, reduced task success in five of eight tested settings compared to targeted, project-specific instructions. Generic is not safe. Specific is safer.

The third group is about safety and security. Define which tools the agent is permitted to use, distinguish read-only from write-access operations, specify which actions require human approval, and treat all external inputs as untrusted. The OWASP Top 10 for Agentic Applications (December 2025) consolidates the community's security guidance. Maintenance belongs here too: instruction files rot like any other documentation, and a wrong instruction is worse than no instruction.

---

This is a genuine contribution. The functional layer, what the agent does, how the repository is organised, what the agent is permitted to touch, is well-understood and increasingly well-specified.

The evidence base for it is thin, though. Two peer-reviewed empirical studies from 2025 provide the only direct evidence about how instruction file content affects agent performance. The remaining nine pillars rest on practitioner norms: vendor documentation, community guides, and blog posts. The advice is reasonable. It is not yet grounded the way a thirty-year certification tradition is grounded. Practitioner norms produce practitioner knowledge: what can be inferred from how repositories are already built. They cannot produce what requires a deliberate authorship decision to include at all.

The obvious response is that the pillars are sufficient and anything missing is refinement rather than a structural gap. That holds only if the missing element appears somewhere in the existing guidance, underspecified or implied. It does not. And there is a layer the pillars do not address at all.
