---
title: "Waiting for perfect data is a permanent excuse"
slug: waiting-for-perfect-data-permanent-excuse
platform: linkedin
tags: [ai-governance, software-engineering, personal-practice]
publish_date: 2026-08-28
cover_image: waiting-for-perfect-data-permanent-excuse.jpg
cover_image_prompt: |
  A full explainer infographic in the style of a hand-drawn, sketchy line-art poster on aged cream paper. Title at the top in hand-lettered text: "Waiting for perfect data is a permanent excuse". A boxed callout titled "Reality check" with three short bullet points: "Always a stale record", "Always an inconsistent field", "Always a drifted schema". Beneath it, a two-path comparison diagram: on the left, a looping arrow labelled "AI re-derives the fix every time" circling back on itself to suggest repeated, wasted effort; on the right, a straight arrow path labelled "Spot it once, then build it in" running through three small steps, "Spot the problem", "Add the missing context", "Build deterministic tooling", ending at a checkmark. A boxed callout titled "Honest caveat" with small text noting that deterministic tooling is not always more controllable, it depends on the case. A closing hand-lettered line at the bottom: "Waiting for perfect data is a permanent excuse. Fixing the specific problem once is the actual work." Warm cream and off-white background, muted slate blue and warm amber accent colours, hand-drawn technical-sketch aesthetic. No people, no photographic elements.
cover_post: |
  An organisation that waits for its data to be fully correct before adopting AI will never adopt. Data quality is not a precondition you clear once. It is an ongoing condition, permanent in exactly the way a fixed launch date is not.

  When someone does spot a specific data problem, the fix is not retraining everyone on the nuance of that schema. It is one person spotting it, adding the missing context once, and building it into deterministic tooling rather than asking AI to reinterpret the same fix from scratch every time.

  Waiting for the data to be perfect is a permanent excuse. Fixing the specific problem once, properly, is the actual work.

  New piece on the habit that actually works.
stop_slop: 45/50
toulmin: Track A 6/6, Track B 6/6
notes: "Source: #592, research-notes/so-synthesis-575-data-imperfection-ai-context-correction-29july2026.md, SO decision #577. Two claims kept separate per the issue: data-quality permanence, and targeted correction crystallised into deterministic tooling as the better response to a specific problem, not a claim that deterministic tooling is unambiguously more controllable (explicitly rejected in that unqualified form per #577). None of the six candidate anti-pattern names used. Smile and CyberSentriq named per the existing Director clearance from #567's thread, reconfirmed applicable here per #577's decision. Reviewed 23 August 2026 against meta/artifacts/voice-principles-blogger.md: cut 'Here is the part I actually think is more interesting than the permanence point', a significance-announcing sentence that deferred the actual point to the next sentence instead of making it. The paragraph opens directly with the content instead. No em dashes, double dashes, or other violations found. Given an explainer-infographic cover_image_prompt per meta/artifacts/cover-image-style-guide.md's opt-in 'Alternate style: full explainer infographic' section (a deliberate departure from the default minimalist convention, not the default for other posts). This is the first attempt at a reproducible prompt for that track; the one prior accepted instance was produced by accident with no recorded process. Briefly converted to format: article on 23 August 2026, then reverted the same day per Director correction: this is a native LinkedIn post with an infographic cover image, not an article. The body already closed on its natural punchy final line, so no body change was needed here."
linkedin_url:
---

An organisation that waits for its data to be fully correct before adopting AI will never adopt. I do not mean that as a criticism. I mean it as a description of how large-organisation data actually works: data quality is not a precondition you clear once. It is an ongoing condition, ordinary, ever-present, and permanent in exactly the way a fixed launch date is not.

If "the data isn't ready yet" is the standard, it will always be true, because there will always be a field that is inconsistent, a record that is stale, a schema that drifted after the system that produced it changed. Treating that as a gate before starting is not caution. It is deferring indefinitely against a bar that was never actually defined and will never be reached, because nobody wrote down what "ready" would even mean.

When someone does spot a specific data problem, a misclassified field, a report that quietly means something different than its label says, the fix is not to retrain everyone in the organisation on the nuance of that schema. The fix is one person spotting it, adding the missing context once, and, where the correction is stable enough to be a rule rather than a judgement call, building it into tested, deterministic tooling rather than asking an AI to reinterpret the same fix from scratch every single time it comes up.

This connects directly to my own background. Smile, the AI governance gateway we are building at CyberSentriq, follows exactly this pattern: use AI once, to help build or update a piece of deterministic companion tooling from a spotted correction, rather than having AI re-derive the same interpretation every time the case recurs. That is not because deterministic tooling is always more controllable than AI judgement in some absolute sense, it depends on the case, and pretending otherwise would overclaim a debate that is genuinely still open. It is because a correction you have already understood once does not need to be re-understood every time, by a human or by AI.

Waiting for the data to be perfect is a permanent excuse. Spotting the specific problem and fixing it once, properly, is the actual work.
