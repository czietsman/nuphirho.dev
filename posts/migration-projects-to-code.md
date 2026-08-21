---
title: "The programme applied its own thesis to itself"
slug: "migration-projects-to-code"
tags: [ai-governance, agentic-development, personal-practice, research]
subtitle: "What moving from Claude Projects to Claude Code revealed about my own operating model"
publish_date: 2026-08-14
format: article
stop_slop: 41/50
toulmin: Track A 6/6, Track B 6/6
cover_image: migration-projects-to-code.png
cover_image_prompt: |
  Two workspaces side by side: the left scattered with loose notes and no clear structure, the right with a single open notebook, organised and indexed. Muted palette, warm grey and off-white. No people, no screens, no digital elements, no colour accents.
cover_post: |
  I have spent months building a research programme that studies what AI governance requires: provenance, versioning, and auditable change history.

  For most of that time, my own AI operating environment had none of those things.

  That is not a confession of hypocrisy. It is a description of how structural requirements become visible: you outgrow the workarounds you have built to compensate for their absence.

  What I found when I moved from Claude Projects to Claude Code: the workarounds were not accidental. They were specifications pointing at what was structurally missing.

  New piece on the transition, the knowledge loss, and why it had to happen.
linkedin_url: https://www.linkedin.com/pulse/programme-applied-its-own-thesis-itself-christo-zietsman-fpoyf
---

The programme I have spent months building studies what AI governance requires: provenance, versioning, and auditable change history. For most of that time, my own AI operating environment had none of those things.

That is not a confession of hypocrisy. It is a description of how structural requirements become visible: you outgrow the workarounds you have built to compensate for their absence.

The research and writing practice ran on Claude Projects, with role definitions in the project knowledge, governing documents available at every session, skill files and voice principles loaded at start. Four distinct roles operated in this environment: the Science Officer for research and evidence quality, the Blogger for all published content, the Dev Lead for code and infrastructure, the PA for scheduling. Each role had quality gates and required reading lists. Sessions produced real work: published posts, submitted arXiv papers, a governance framework with theoretical grounding in sixty years of organisational theory.

What they also produced, alongside the work, was a growing workaround debt.

The project knowledge feature is what I would call trust by assertion. The documents are present. They say what they say. There is no record of what they said last week, who changed them, or when they became stale. A session produces an artifact; the artifact becomes a project document; the provenance of that transition is invisible. You accept the current state because the platform gives you nothing else to compare it against.

For a personal writing workflow, this is fine. For a research programme that studies provenance and structural completeness in AI governance, it is the wrong foundation.

The second friction was simpler to name. Conversation AI and coding AI lived in different environments and did not share a store. The Dev Lead worked in the terminal and the repository. The Blogger worked in Projects, with drafts and the editorial calendar. When a post moved from drafting to publication, someone had to carry the file across the boundary manually. When the Science Officer completed a handover, the Dev Lead received it as a pasted document rather than a shared file. Every transfer across that divide cost coordination effort that should not have been necessary. In Thompson's terms, the two environments had a transfer dependency with no shared medium to resolve it.

Moving to Claude Code meant building a repository that could function as the full operating environment. The CLAUDE.md session bootstrap, previously a single-line stub pointing at the engineering rules, became a five-step procedure: establish the active role, read the role file, read the required reading, read the repository index, read the role handover. Role definitions moved from project knowledge into a roles/ directory with explicit required-reading lists and quality gates written into the file itself. A three-tier knowledge structure classified 119 research artifacts by how settled the knowledge in them was. An index replaced semantic search as the way to locate historical context quickly without opening folders blind.

The workload was not small. The session archive contained 1,341 files including binaries, five roles worth of handover state, and a naming collision problem that required a UUID-based scheme to resolve. It took a dedicated Dev Lead session and a Chrome extension that intercepted the Claude.ai fetch layer to get the data out at all.

But what emerged from the build was better than a fix. It was a pattern.

Every structure I created to make the repository work had a counterpart in the Projects workaround list. The handover protocol that compensated for context loss became a first-class version-controlled artifact. The indexing rules that compensated for unnavigable session history became a structured map. The novel findings registry that compensated for the absence of shared memory became a committed file with authorship and date. The workarounds were not accidental. They were specifications pointing at what was structurally missing.

Coordination theory offers the clearest vocabulary for what changed. Claude Projects coordinated through standardisation of norms: the project knowledge was a shared ideology, and sessions coordinated by holding the same context. Claude Code coordinates through standardisation of outputs: sessions read and write shared files, and coordinate by inspecting each other's artifacts. One session's committed handover is the next session's input. The coordination mechanism shifted from pooled to reciprocal interdependence.

The trust model is the sharpest difference. Project knowledge was trust by assertion. A git repository is trust by provenance: you can see what changed, who changed it, and when. That is not an ergonomic improvement. It is the difference between a governance document that knows when it became stale and one that structurally cannot.

The programme's own evaluation framework has a principle, P7, that requires governance documents to declare their staleness conditions. Project knowledge violated P7 by design, and I did not notice. Not for lack of attention: what began as a blog outgrew that scope faster than I had planned for, and the workarounds absorbed the cost of the gap before I had registered it as a gap. I was also still learning the environment itself, most of that from the mobile app, and had not yet grasped how wide the divide between Claude Projects and Claude Code actually was. The convenience was real, and the detour was not wasted either: it built empathy for how most of my actual audience works, and gave me real knowledge of what a lower-friction entry point actually needs.

One honest caveat on scale. For a single writing practice with occasional AI assistance, Claude Projects is still the right tool. Faster to start, lower friction on mobile, no repository setup required. The case for a repository operating model builds at the point where the work compounds across agents, where decisions need to be auditable, where the output of one session is the formal input to the next. Most business people will never be comfortable working in version control, and that is not where most work needs to be either. It is where this programme had to go.

The honest close is Polanyi: handovers transfer what can be told. What is known beyond telling does not come across.

The calibration about which sources to trust, the reasoning embedded in specific editorial decisions, the feel for preferences that built up across months of sessions: none of that is in any committed file. It is in the session history of a closed platform. When I moved to Claude Code, those conversations stayed where they were. The new sessions build on the committed artifacts. They do not have access to the conversations that produced them.

The transition is a structural upgrade with a permanent knowledge loss embedded in it. I made it because the alternative, running a governance research programme on infrastructure that structurally could not support its own thesis, was not a viable long-term position.

The specification is the product. The operating environment has to live up to the specification.
