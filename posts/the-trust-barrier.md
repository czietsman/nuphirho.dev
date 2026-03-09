---
title: "The Trust Barrier"
slug: "the-trust-barrier"
draft: true
tags: [ai, process, bdd, trust, software-engineering]
---

The biggest obstacle to AI in software engineering isn't the technology. It's trust.

Teams don't trust AI-generated code. And they shouldn't, not blindly. But the conversation in most engineering organisations has stalled at that point. The distrust is treated as a conclusion rather than a starting point. AI generates code we can't trust, therefore we shouldn't use AI. End of discussion.

That's the wrong framing. The right question isn't whether to trust AI. It's how to build a process that makes trust measurable.

## The carpenter and the power tools

Uncle Bob at [Uncle Bob Consulting LLC](https://www.linkedin.com/company/cleancoder/) has been working through this in public, and his evolving position is worth following. Back in September 2025, he was clear that AI would change programming for the better but not replace programmers.

%[https://twitter.com/unclebobmartin/status/1962636247769530650]

By January 2026, after working with AI more directly, he observed something important about the relationship between programmer and AI: it's not like a pilot and autopilot. His conclusion was striking.

%[https://twitter.com/unclebobmartin/status/2008879916301898134]

"In the end the AI might write 90% of the code, but the programmer will have put in 90% of the thought and effort."

Then in February, after six weeks of daily use, he compared the experience to a skilled carpenter being handed power tools for the first time. The power is undeniable, but so are the risks. He's still not sure his project wouldn't be just as far along without it.

%[https://twitter.com/unclebobmartin/status/2019025982863069621]

That progression resonates because it captures where most experienced engineers are right now. They can see the potential. They can also see the ways it could go wrong. And without a clear framework for managing those risks, the rational response is caution.

But here's where it gets really interesting. Uncle Bob, the person who formalised the three laws of TDD, recently acknowledged that TDD as traditionally practised is inefficient for AI. The principles remain, but the techniques need to adapt.

%[https://twitter.com/unclebobmartin/status/2023158252700066287]

That's a significant shift from someone who has spent decades advocating for specific engineering disciplines. And it's exactly the point. The principles of verification, rigour, and specification don't change. The techniques do.

I've been through the same journey. I spent time being cautious, reviewing everything line by line, second-guessing outputs, wondering whether I was actually saving time or just shifting the effort from writing to reviewing. It felt productive, but it wasn't fundamentally different from writing the code myself.

Here's what I found on the other side of that phase: the question isn't "can I trust this code?" It's "does this code meet the spec?"

That reframe changes everything.

## From gut feeling to process

When trust is a feeling, it's subjective, inconsistent, and unscalable. One developer might trust AI-generated code after a quick scan. Another might rewrite it entirely. Neither approach is reliable, and neither scales to a team.

When trust is a process, it becomes a solvable engineering problem. You define what correct looks like before any code gets written. You build verification layers that catch failures regardless of who or what produced the code. You measure outcomes rather than gut-checking inputs.

This isn't a new idea. It's the scientific method applied to software delivery. Hypothesise what the system should do. Build it. Test whether it does what you hypothesised. If it doesn't, iterate.

We call it computer science for a reason, but somewhere along the way the discipline drifted from its roots. We stopped treating software delivery as an empirical process and started treating it as a craft, something learned through apprenticeship, refined through experience, and validated by the judgement of senior practitioners. That model worked well enough when humans were writing all the code. Experienced judgement was the best verification layer we had.

But craft-based verification doesn't scale to AI-generated output. You can't rely on an experienced developer's intuition when the volume of generated code outpaces any individual's ability to read it. What scales is the scientific method: define the expected behaviour precisely, run the experiment, measure the result. If the result matches the hypothesis, the code is fit for purpose. If it doesn't, you have a specific, measurable failure to investigate.

The difference now is that AI has made this approach economically viable in ways it wasn't before.

Writing behavioural specifications, running mutation tests, enforcing contract-driven APIs, these practices have existed for years. They were always the right thing to do. But the effort required to implement them properly meant that most teams cut corners. When you're manually writing all the code, spending additional time writing comprehensive specifications and mutation test suites feels like a luxury.

AI changes that equation. When the code generation is fast and cheap, the verification layer becomes the primary deliverable. The spec is no longer documentation you write after the fact. It's the thing you write first, and it's the thing that determines whether the output is fit for purpose.

## What this looks like in practice

The process I've built works like this:

Write behavioural specifications before any code gets written, by human or AI. Use BDD to define what the system should do, not how it should do it. These specifications are executable. They're not documentation sitting in a wiki. They run as part of the pipeline and they fail when the code doesn't match the expected behaviour.

Run mutation tests to validate that your tests actually catch failures. A passing test suite means nothing if the tests would still pass with broken code. Mutation testing introduces deliberate faults and checks whether your tests detect them. It's the verification layer for your verification layer.

Enforce contract-driven APIs so integrations are verifiable independently. When services agree on contracts, you can validate that each side honours the agreement without spinning up the entire system. This matters even more when AI is generating the integration code, because the contracts become the source of truth that the generated code is measured against.

Review outputs against specifications, not line by line. This is the fundamental shift. Instead of reading every line of AI-generated code and asking "does this look right?", you run the specification suite and ask "does this behave correctly?" The first approach doesn't scale. The second does.

## The personal proof point

I don't write code by hand anymore. I direct AI to write it. I've never formally learned Vue or Go, but I'm shipping production-ready prototypes in both. Not because I blindly trust AI, but because I've built the process that lets me verify the output.

That statement makes some people uncomfortable. It sounds reckless to say "I ship code in languages I haven't formally learned." But the discomfort reveals the assumption: that understanding every line of code is a prerequisite for shipping reliable software. It's not. Understanding the specification, the verification layer, and the failure modes is the prerequisite. The language is an implementation detail.

This doesn't mean the code doesn't matter. It means the specification matters more. And when you have a robust specification, the question of whether a human or an AI wrote the implementation becomes less important than whether the implementation meets the spec.

## The competitive advantage is process

The teams that figure this out first won't just be faster. They'll be more reliable than teams still reviewing every line by hand. That's counterintuitive, because we associate human review with quality. But human review is inconsistent, fatiguing, and subject to all the cognitive biases that make us miss things we've seen a hundred times before.

A specification-driven process with automated verification doesn't get tired. It doesn't skip steps on a Friday afternoon. It catches the same class of failures at 5pm that it catches at 9am. It scales to any volume of AI-generated code without degrading in quality.

The trust barrier is real. I'm not dismissing the concern. Engineers who are cautious about AI-generated code are being responsible. But caution without a path forward is just avoidance. The path forward is process.

Build the verification layer. Let the spec be the source of truth. And stop asking whether you can trust the code. Start asking whether it meets the spec.
