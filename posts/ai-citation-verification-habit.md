---
title: "AI gave me a citation. It does not exist. Here is what I do now."
slug: "ai-citation-verification-habit"
tags: [research, ai-assisted-development]
publish_date: 2026-08-21
linkedin_url:
---

AI citation fabrication is not an edge case. A study published on arXiv in February 2026 examined one hundred fabricated citations found in fifty-three papers accepted to NeurIPS 2025, roughly one per cent of the conference's accepted submissions. Each of those papers went through three to five expert peer reviewers. The fabricated citations were not caught in review.

That figure is worth holding onto before I get to the habit it produced in me, because it does something useful: it removes the story in which this is a problem only for careless users who do not know what they are dealing with. Domain experts, reviewing work in their own field, with all the motivation to check that peer review provides, missed fabricated citations at scale.

The taxonomy Mohammad Samar Ansari developed for his analysis breaks the failure modes into five types. Total fabrication, where the cited paper simply does not exist, accounts for sixty-six per cent of cases. Partial attribute corruption, where the paper exists but key metadata (author, year, title, journal) is wrong, accounts for another twenty-seven per cent. Identifier hijacking, the most structurally interesting failure mode, accounts for four per cent: the citation carries a valid identifier that resolves to a real paper, but the metadata fabricated around it does not match that paper. The identifier looks verifiable. It resolves. The paper on the other end is not the one being cited. The remaining three per cent covers placeholder hallucinations and semantic hallucinations, where a real paper is cited but does not actually support the attributed claim.

The pattern that makes this hard to catch without primary-source verification is compound failure. Most fabricated citations combine a plausible surface appearance with at least one mismatch below the surface. You can read the title, recognise the author name, and see the year fits, and still be looking at something that does not say what the citation claims.

---

The first time I encountered a fabricated citation in my own work I noticed it by accident, not by method. I was drafting something that referenced a piece of research, I happened to follow the identifier link out of curiosity, and the paper that came back was unrelated to the claim I was making. The AI had produced a plausible-looking reference pointing to a real paper about a completely different topic. The identifier resolved. The paper was real. The attribution was false.

That is identifier hijacking, and I would not have caught it through any review of the citation text alone.

Since then the habit I have developed is to treat AI citation output as a research direction rather than a reference. An AI that says "this claim is supported by Smith et al. 2023 in the Journal of Applied X" is giving me a starting point, not a source. The verification (the step that turns a starting point into something I can actually use) is fetching the primary source and reading what it says.

The steps I now follow, in order:

First, I ask for the claim separately from the citation. If I need both a claim and a source, I ask for the claim first and then separately for a source that might support it. This keeps the two questions distinct and prevents the citation from anchoring my reading of the claim.

Second, when a citation is provided, I fetch the primary source before I use it. This means: find the paper, find its abstract, confirm the paper exists and the identifier resolves to the right paper. This step alone catches total fabrication and a large fraction of partial attribute corruption.

Third, I read what the primary source actually says about the point being cited. Not the abstract only, but the section or passage that is supposed to support the attribution. This step catches semantic hallucination, where a real paper is used to support a claim it does not make.

For identifier hijacking specifically, the check is to confirm that the metadata the AI attributed (author, year, title, journal or venue) matches the metadata at the resolved identifier. If the DOI resolves but the paper on the other end is not by Smith, is not from 2023, and is not in the Journal of Applied X, the citation is fabricated regardless of whether the DOI itself is real.

This takes time. Not a great deal of time, but enough that skipping it is tempting, especially when the AI's output is fluent and the citation is plausible. The Ansari data is what keeps me from skipping it: three to five expert reviewers, in their own field, looking at papers in their area, missed fabricated citations. Fluency and plausibility are exactly what the failure looks like from the outside.

---

If you are using AI to research something and you intend to cite what you find, the citation is not verified until you have fetched and read the primary source. The AI's confidence in the citation is not evidence of the citation's accuracy. Plausible metadata is not evidence of the citation's accuracy. Even a resolving identifier is not, by itself, evidence of the citation's accuracy.

The standard that follows from this is not onerous: any claim that you intend to attribute to a source should be verifiable against that source before it is used. That was always true. What the AI context changes is that the gap between "the citation looks right" and "the citation has been verified" is no longer obvious from the output. The fabricated citation and the accurate citation look identical until you follow the link.

The harder question is what to do about claims where primary-source verification is difficult: papers behind paywalls, data from unpublished sources, assertions that you cannot independently check. My position on that is conservative: if I cannot verify it, I do not cite it as though I have. If I can get to a paper's abstract but not its full text, I can cite what the abstract says and note the limitation. What I do not do is allow fluency in the AI's output to stand in for verification I did not do.

This is a habit, not a system. The Ansari taxonomy and the NeurIPS data give it more weight than it would otherwise have, but the underlying logic precedes the data: a reference is only useful if it actually says what you say it says, and the only way to know that is to check.
