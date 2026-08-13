---
title: "Adding the quality layer to your AGENTS.md"
slug: "agents-md-adding-quality-layer"
tags: [ai-governance, software-engineering, agents, promptq]
series: "AGENTS.md: the functional and quality layers"
series_part: 3
publish_date: 2026-08-03
stop_slop: 41/50
toulmin: Track A 6/6, Track B 5/6
format: article
cover_image: agents-md-adding-quality-layer.png
cover_image_prompt: |
  Two versions of the same text document shown side by side, the right one visibly longer and more complete. Plain background, viewed straight on. Muted palette, slate blues and off-whites. Clinical, technical. No people, no screens, no colour accents.
cover_post: |
  The previous post scored a real claude init output at 2/7 on PromptQ's structural quality framework. Five principles were either absent or partially addressed.

  This post shows what adding them looks like in practice. Same repository. Same CLAUDE.md. This time the agent discussed each principle with the codebase in hand before touching anything.

  The sharpest moment came at P3. An Apply had run with no visible approval pause, even though the docs required it. P3 does not investigate what happened. It decides what the next session will do when it reaches that step.

  A governance document governs forward behaviour. The past is the past.

  New piece on a real implementation.
---

The previous post scored a real `claude init` output at 2/7 against the PromptQ structural quality framework. Five principles were either absent or only partially addressed. This post shows what adding them looks like, not as a set of instructions, but as a session account.

Same repository: Terraform-managed Cloudflare records for a company's estate, deployed through a build pipeline. Same `CLAUDE.md`. This time the agent was asked to apply the PromptQ framework, explore the codebase, and discuss each of the seven principles before making any changes. The constraint was explicit: discuss all seven first, write nothing until the design was settled.

---

Before the discussion, the agent mapped the current file against each principle:

**P1 Purpose.** What was there: an "Overview" section describing the repo. What was missing: no agent charter, no definition of done.

**P2 Measurement.** What was there: `fmt`/`validate`/CI plan. What was missing: not framed as acceptance criteria.

**P3 Authority.** What was there: a PR workflow with a self-approve note. What was missing: no statement of what the agent may do alone.

**P4 Data Integrity.** What was there: the API token marked sensitive. What was missing: no provenance rule for record values.

**P5 Quality Control.** What was there: "Review the Plan output." What was missing: no active-engagement gate.

**P6 Consistency.** What was there: nothing. What was missing: live contradictions already present.

**P7 Adaptation.** What was there: nothing. What was missing: no staleness or versioning at all.

Then the discussion started.

---

**P3 produced the sharpest moment.** The authority question was not theoretical. During the session, Apply had run to completion with no visible approval pause, even though the docs said approval was required. The root cause was not investigated: it may have been a branch policy gap, an onboarding shortcut, or something else. That investigation was not the point. The point was: what should every future agent session do when it reaches that step? The agreed answer was to treat every privileged action (merge, approve, apply, delete) as requiring explicit escalation and per-action consent, regardless of what any mechanism does or does not enforce. The governing document tells the next session what to do. The past is the past.

A single "yes" does not cover a list. That rule came from the session itself. The Architect said: when there are multiple questions, don't assume consent is given with a yes; make it clear what is being agreed to. From that point on, every multi-part question in the session was confirmed item by item. A rule stated at P3 changed how the rest of the session ran.

---

**P6 asks what the next session will do when documented obligations and observed reality diverge.** The first case was straightforward: the file said records are ordered alphabetically; the actual file had a jumbled block and a stray key. The agreed resolution was to keep the rule and tell future agents to follow it for new records and flag existing drift.

The second case came back to the approval question from P3. The document asserted a manual approval requirement. An apply had run without a visible pause. The question P6 raised was not what caused that: it was what the next agent session should do when it encounters the same situation. Should it treat an absent gate as permission to proceed? The answer was no. The Known Open Gap section now instructs any future agent to flag the discrepancy and not rely on or exploit what it observes. The document governs forward behaviour. Investigating the past is a separate task.

---

**P2 named the right thing to measure.** `terraform fmt -check` passes. `terraform validate` passes. The plan succeeds. None of those are the real criterion. The real criterion is that the diff matches the intent.

The agreed approach: state the expected diff before starting, explicitly including what must not change ("0 destroyed"), and re-verify at plan and again post-apply. A succeeding plan is a proxy. Diff equals intent is the criterion. The original file had the validation commands. It did not have this distinction.

---

With all seven discussed, the agent presented the consolidated design and asked for confirmation before proceeding, not bundled, honouring the rule that had emerged at P3. Then it applied the new specification to itself. The PR description it opened was an evidence packet: stated intent, expected diff, rollback plan, plan-versus-intent result, conventions met, safety notes. It verified that the diff showed only the `CLAUDE.md` modified, nothing else, nothing destroyed. It did not merge or self-approve.

The improved file is roughly double the length of the original. It adds sections for authority, the required change workflow, data integrity, consistency, and adaptation. Two things the session could not resolve are named in a "Known open gaps" section rather than omitted: the approval-gate enforcement question, and the definition of what counts as a sensitive change. The file names what it does not yet govern.

---

That last point is worth holding. A governance document that omits what it cannot yet specify is overclaiming. One that names its open gaps is making a weaker but more honest set of claims. The session produced a file that is stronger as a governing instrument precisely because it is more explicit about its limits.

The quality layer is not abstract. Applying a structural framework to a real file surfaces questions the original file had not asked: not about what has happened, but about how the next session will behave. What does the next agent do when it finishes a task? What does it verify before opening a PR? What does it refuse? What does it do when observed reality does not match the document's assumptions? The answers are in the improved file. The questions are what the framework provided.

The functional layer describes the repository. The quality layer governs what the agent does in it. The PromptQ framework is available at promptq.ai. The seven principles applied in this series are P1 through P7 of the v1.3.1 framework, which is still being refined as it is applied to real repositories.
