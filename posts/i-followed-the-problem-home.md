---
title: "I Followed the Problem Home"
slug: "i-followed-the-problem-home"
tags: [ai-governance, epistemology, software-verification, collaboration]
canonical: "https://blog.nuphirho.dev/i-followed-the-problem-home"
publish_date: 2026-04-04
edited_at: 2026-04-11T00:00:00Z
---

[James Bach](https://www.linkedin.com/in/james-bach-6188a811/) wrote that failing to detect a problem is not a measurement of non-problemness. He was responding to me. That exchange sent me somewhere I did not expect to go.

I am not a philosopher. I am not a mathematician. I am an engineer who has spent twenty years trying to verify that software does what it is supposed to do, and I kept running into the same wall from different directions.

On a Sunday evening in my office in Technopark, Stellenbosch, I was preparing a presentation on the effect of backfill within deep mine stopes. The numerical model was producing beautiful results. I tried to speed things up, changed something in the solver, and introduced a bug. We used CVS at the time. When you are in the zone and tweaking, you skip the commit. The stress erased the memory of what I had changed. I presented at Tau Tona, one of the deepest gold mines in the world, and the results looked exactly as I expected them to look. It was only at the follow-up workshop, where we tried to make the theory practical, that the new build of the tool revealed itself to be producing garbage. The presentation had looked fine. The practice exposed it. The version control was there. The process existed. I ran straight past it because it felt like friction in the moment. At another company I watched a feature nobody would touch because the selection logic had grown beyond anyone's ability to reason about it. At every step, the problem was the same: how do you know that what you built is actually correct?

I rebuilt the selection model from scratch over a Christmas break using test-driven development. Not because I was told to. Because I needed something external to the code to tell me whether the code was right. The test suite was not the product. But it was the only thing I trusted.

That instinct was twenty years old before I had a name for it.

---

When large language models started writing code, I thought the problem would get worse. A system that produces plausible-sounding output, confident in tone and occasionally incorrect in fact, is a system that makes verification harder, not easier. The code looks right. The tests pass. And yet.

I published a paper on this in March 2026. The core claim: AI reviewing AI-generated code without an external specification shares blind spots. The correlated error hypothesis. Three models reviewed code containing domain-specific bugs none of their training had prepared them for, and they agreed it was correct. The confidence was not evidence of correctness. It was evidence of shared blindness.

The specification was the only thing that did not share the blind spot. The executable BDD scenario, written to describe what the system must do independently of how it does it, could catch what the reviewers could not. Because it was external. Because it was deterministic. Because it did not care how plausible the code looked.

The specification is the quality gate.

---

I did not expect what happened next.

James Bach posted a question about whether past correctness justifies trust in a tool. I replied with the mutation testing argument. He pushed back: failing to detect a problem is not a measurement of non-problemness, unless you have a guaranteed method of detecting every possible kind of problem. He was right. That is the oracle problem, the formal boundary of what executable specifications can catch. I had already published a taxonomy of it. His pushback named the limit more precisely than I had.

It is not. That is Hume's problem of induction. I had stumbled into philosophy.

Simon Wardley, who thinks about strategic positioning and the evolution of capabilities, engaged with the specification-at-scale argument. He has a framework for understanding how practices become commodities. He pushed on whether executable specifications could survive at that scale, and whether the ecosystem would commoditise the verification layer before the engineering discipline caught up.

A paper by Catalini, Hui and Wu at MIT Sloan, published independently, arrived at the same structural diagnosis from a different discipline entirely. The binding constraint on AI-driven growth is not execution speed. It is human verification bandwidth. The specification is the mechanism that makes verification machine-executable rather than biologically bottlenecked. Two disciplines, same conclusion.

I had been following a problem. The problem had been following a path that connected software engineering to epistemology to economics without asking my permission.

---

This is where I need to be honest about what I do not know.

The philosophical grounding of the verification problem goes back much further than I had realised. Sextus Empiricus, writing around 200 CE, described what he called the criterion problem: how do you establish the criterion by which you judge something to be true, without already assuming the criterion is correct? The Pyrrhonists suspended judgement precisely because they could not resolve it.

Software verification has the same structure. To verify that the test is testing the right thing, you need a meta-test. To verify the meta-test, you need a meta-meta-test. At some point, you need an external reference. The specification. The oracle. Something that stands outside the system and says: this is what correct looks like.

I recognised this structure from my own work before I had read a word of ancient philosophy. But recognising a structure and understanding its formal treatment are different things. I am working with a philosopher, Jaco Louw, who specialises in Pyrrhonian scepticism, to think through whether the formal epistemological argument holds. He is the expert. I am the engineer who followed the problem far enough to need one.

This is not false modesty. It is how good thinking works. You follow the problem until it takes you somewhere you cannot go alone. Then you find the people who live there.

---

I have been asking questions in public for several months now. The questions have attracted people I did not expect: practitioners, researchers, strategists, and investors. Not because I have all the answers. Because the questions are the right ones, and the right questions draw out the people who have been thinking about the same things from different angles.

What I am looking for now is more of that.

If you work in formal verification, computability theory, or proof theory and you have thought about what undecidability means for software correctness in practice, I want to hear from you.

If you work in epistemology, philosophy of science, or sceptical traditions and you see the oracle problem as something your field has already mapped, I want to hear from you.

If you work in economics, organisational theory, or complexity science and you see the verification bandwidth argument from a different vantage point, I want to hear from you.

If you are a practitioner who has been living inside this problem and has not yet found the vocabulary for what you keep running into, I especially want to hear from you. The vocabulary matters less than the pattern.

I am not building this alone. I do not want to. The work is better when it is challenged, extended, and connected to things I have not read yet.

Follow the problem with me.

---

The research programme is at nuphirho.dev. The conversation is in the comments.
