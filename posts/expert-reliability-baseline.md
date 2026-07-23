---
title: "Measuring AI against humans who disagree with themselves"
slug: "expert-reliability-baseline"
subtitle: "What fingerprint examiner data says about the human expert standard"
tags: [research, ai-governance]
format: article
cover_image: expert-reliability-baseline.png
cover_image_prompt: |
  A minimal hand-drawn image on plain paper of two side-by-side measurement readings on the same simple scale, each showing a slightly different value. Muted palette, slate and off-white. No people, no screens, no colour accents. The image reads as a conceptual illustration of inconsistency in measurement, not a technical diagram.
cover_post: |
  The implicit benchmark for AI performance is almost always: would a human expert get this right?

  Reasonable. But how reliable is the human expert standard?

  The clearest data is in fingerprint examination, treated as near-infallible by courts for decades. Two large-scale studies by Ulery and colleagues tested qualified examiners on real casework samples.

  Intra-rater inconsistency: roughly ten per cent. One in ten categorical conclusions changed when the same examiner re-examined the same evidence seven months later, without knowing they had seen it before. Concentrated at the threshold: whether the evidence was clear enough to call.

  If the human expert baseline is itself unstable at borderline cases, checking AI output against expert agreement is measuring two imperfect sources against each other.

  New piece on what that means for validation.
publish_date: 2026-07-29
linkedin_url:
---

When people evaluate whether to trust an AI decision, the implicit comparison is almost always: would a human expert get this right? It is a reasonable place to start. But the question it rests on (how reliable is the human expert standard?) receives far less scrutiny than the AI's performance does.

The fingerprint evidence is the clearest data I have found on this, and it is worth looking at carefully because fingerprint examination is a domain that courts, lawyers, and investigators treated for decades as near-infallible. If expert reliability is measurable and imperfect there, the implication for softer domains is significant.

---

Two large-scale studies by Brandon Ulery and colleagues establish the picture. Both used professional fingerprint examiners working on actual casework samples, not simplified laboratory tasks.

The 2011 study, published in PNAS, tested 169 qualified examiners across a pool of roughly 744 fingerprint pairs, with each examiner evaluating about a hundred pairs. The false positive rate (examiners concluding two different fingers matched when they did not) was 0.1 per cent. That sounds low, and by some comparisons it is. The false negative rate, examiners failing to identify matching fingerprints as matching, was 7.5 per cent, and eighty-five per cent of examiners made at least one false negative across their evaluated pairs. About twenty-three per cent of decisions were threshold disagreements: examiners who concluded there was not enough detail to make a determination, which is a category that is itself a judgment call and that generates its own inter-examiner variation.

The 2012 study, published in PLoS ONE, asked a different question: not whether examiners were accurate but whether they were consistent with themselves. Seventy-two qualified examiners were retested on twenty-five fingerprint pairs they had already evaluated, approximately seven months after their first assessment, without being told they had seen the samples before. The overall repeatability was ninety per cent for mated pairs and eighty-six per cent for non-mated pairs. Read differently: roughly one in ten categorical conclusions ("this is a match" or "this is not a match") changed when the same examiner re-examined the same evidence under the same conditions.

The point about where the inconsistency is concentrated is important. Reversals were not uniformly distributed across difficulty levels. They were concentrated at the sufficiency threshold: decisions about whether the evidence was clear enough to call at all. That is not the easy case, but it is also not a rare category in real casework. And if the threshold decision is itself unstable, the downstream conclusion that follows from it is built on less stable ground than the formal process implies.

---

The natural response to this data is to note that fingerprint examiners are still performing well in absolute terms on clear cases, and that is true. The false positive rate of 0.1 per cent, in a domain where it carries very high consequences, is low. The 2011 study was, when published, generally read as a demonstration that fingerprint examination was more reliable than critics had claimed.

I am not making an argument against fingerprint examination, and I am not making an argument that AI would do better. The point is narrower: a domain that had been treated as providing near-certain expert conclusions showed, under controlled measurement, non-trivial intra-rater inconsistency and inter-rater disagreement, concentrated at the threshold decisions that most require judgment.

That pattern, once you are looking for it, appears in other domains where expert judgment is used as a reference standard. Pathologists interpreting breast biopsy specimens, radiologists double-reading the same scan, peer reviewers assessing the same submission: the research in each of these areas consistently finds that expert agreement is an imperfect proxy for a stable ground truth, particularly at the borderline cases where judgment matters most. I am being deliberately vague about specific figures in those adjacent domains because I have not verified them to the same level as the Ulery papers, but the directional pattern is well documented.

---

The implicit comparison most people make (would a competent human get this right?) embeds an assumption about the stability of the human reference standard that the data does not fully support. Particularly for the cases that are hardest and where expert judgment diverges, "would a competent human agree with this AI output?" is not a clean question.

This is not an argument for trusting AI more. It is an argument for being more precise about what validation means. If you are using an AI tool to assist with evaluation tasks, and you are checking the AI's output by asking whether a human expert would concur, you are measuring agreement between two imperfect and potentially inconsistent sources rather than agreement with a stable ground truth. In clear cases that may not matter. At the borderline, where judgment is doing real work on both sides, it matters considerably.

The alternative frame is to separate the question of accuracy from the question of reliability. Accuracy asks: does the output match the right answer? Reliability asks: does the output hold up under repeated application to the same input, and does it remain stable when the surrounding context changes? For AI as for human experts, these are different properties, and a system can have one without the other.

There is also a contextual-bias dimension worth naming. Research on fingerprint examiners given contextually biasing information (prior knowledge that a suspect had confessed, for example) found that at least some examiners reversed conclusions on evidence they had previously assessed, on the basis of context rather than the physical evidence. This is a known human failure mode: expert judgment is not context-free. AI systems have their own contextual sensitivity, of a different kind, but the general principle (that what surrounds the evidence affects how it is read) applies in both cases.

---

The Ulery data has a scope limit that matters here. It concerns a bounded, well-specified task: given two fingerprint images and a fixed protocol, decide whether they match. The examiner and the comparison are working from the same evidence, evaluated against the same class of question, repeated across many trials. That is what makes the inconsistency measurable at all.

Most expert judgment in engineering does not work that way. When I review a large system I have worked on for months, my assessment draws on more than what is in front of me: prior decisions and why they were made, constraints that never made it into any document, the shape of the whole architecture held in memory rather than stated anywhere. An AI reviewing the same code at the point of the prompt has access to what is in the prompt and the repository, not to that accumulated context. That is a different problem from the one the Ulery studies expose. It is not that the AI is less consistent on a shared task. It is that the task is not actually shared: the human is drawing on information that was never given, and no amount of prompting fully closes that gap, because much of it is not the kind of thing that gets written down.

The reliability question and the context question are separate failure modes, and they call for different responses. Where the task is genuinely bounded and the evidence genuinely shared, comparing AI output to expert agreement is the right frame, with the caveat this piece is making: expert agreement is not a stable ground truth either. Where the task depends on context the AI was never given, the gap is not about reliability at all. It is a scope gap, and the fix is not more scrutiny of the AI's answer. It is making sure the context that lives in the expert's head reaches the prompt in the first place.

---

The practical implication for individual practice is not complicated, but it requires resisting a cognitive shortcut that feels like due diligence. Checking AI output against your own expert judgment is useful. It is not a substitute for independent evaluation of the underlying question. Your expert judgment on the same case may be inconsistent with itself across time, and with other expert judgments across people, in ways that are not obvious from the inside.

For decisions that carry significant consequences, the right check is not "does this match what I would expect?" but "does this hold up against independent evidence from the primary source?" That is a slower and more demanding standard, and it is often not proportionate to the stakes involved. For low-stakes, high-volume tasks, agreement-with-expert is probably the right practical proxy. For the cases where the decision matters and where the evidence is borderline, that is precisely where both human and AI judgment are most likely to be inconsistent, and where primary-source verification is most worth the cost.

The Ulery data is the clearest available demonstration that a near-infallible expert standard, under measurement, is neither. If fingerprint examiners in a high-stakes forensic domain show measurable intra-rater inconsistency at the threshold decisions that matter most, the human expert baseline in softer domains is not a ceiling to be beaten. It is an imperfect reference that requires the same scrutiny we apply to the AI output it is used to evaluate.
