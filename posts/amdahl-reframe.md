---
title: "Amdahl's Law does not apply here"
slug: amdahl-reframe
publish_date: 2026-05-29
linkedin_url: https://www.linkedin.com/feed/update/urn:li:share:7463661099272384512
tags:
  - ai-assisted-development
  - software-engineering
---

Someone applied Amdahl's Law to AI-assisted development in the comments this week. The argument was that if coding is 20% of the job, the maximum speedup from AI assistance is 25%. It is a good framing. It is also the wrong model.

Amdahl's Law describes a fixed task with a fixed scope. You speed up one stage. The rest stays the same. The ceiling is determined by the proportion you cannot change.

But that is not what is happening.

The 80% that Amdahl treats as fixed is thinking clearly about what to build, communicating it precisely, and verifying it was built correctly. Those activities are also compressing. The spec iteration loop gets shorter because you can pressure-test a design before a single line is written. The verification cycle tightens because test scaffolding is generated alongside the code. The architecture conversation moves faster because a model that pushes back is available at any hour.

The bottleneck is not fixed. It is shrinking too.

But the deeper problem with the Amdahl framing is that it models a single pipeline. What I am actually experiencing is capacity expansion across an entire life. I am running more projects at work simultaneously, at higher quality. The communication is sharper. The planning holds better. I am running a research programme and publishing on it in parallel. The gym routine actually sticks.

None of that is a 25% speedup on a fixed task. It is a different operating mode.

The compounding effect is what the speedup model misses. Better architecture means less rework. Less rework means faster verification. Faster verification means more confidence to ship. Better pipelines mean fewer production incidents. Better onboarding means the next person gets productive sooner. Better security review means fewer vulnerabilities in the first place. Each improvement reduces the drag on everything downstream.

Amdahl models throughput on a single process. It has nothing to say about what happens when freed cognitive load gets redirected into entirely new activities, or when quality improvements compound through a system over time.

The question is not how much faster you are doing the same job. The question is what you are doing with everything AI gives back.
