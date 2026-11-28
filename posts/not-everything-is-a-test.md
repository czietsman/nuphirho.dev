---
title: "Not everything is a test"
slug: "not-everything-is-a-test"
tags: [software-verification, ai-governance]
cover_image_prompt: |
  A full explainer infographic in the style of a colourful, cheerful, hand-drawn line-art poster on light cream paper, square format. Title at the top in friendly hand-lettered text: "Not Everything Is A Test". On the left, a panel labelled "A test samples": a large field of dots with only a few circled by a magnifying glass, and a small tag "Says nothing about the rest". On the right, a panel labelled "A checker reads it all": a short chain of linked boxes, each with a tick, and a tag "Complete for the proof, not the world". Below both, a bridge between a block labelled "Maths" and a block labelled "The real world", with the gap under the bridge labelled "The specification: test it here". A closing hand-lettered line at the bottom: "Proof moves the gap to one known place." Bright coral, turquoise, sunny yellow, leaf green, sky blue, and warm orange accents, with dark legible lettering, rounded boxes, lively arrows, and friendly icons. Happy, welcoming, and energetic; not institutional, clinical, or ominous. No people, no photographic elements, no logos, and no watermark. Keep all text well inside the edges.
linkedin_url:
publish_date: 2026-12-02
---

Most of us in software look at checking through one lens: everything is a test. Through that lens a proof checker is one more test that can miss things, and "but you have to trust the checker" sounds like the end of the argument.

I think that misses a change of kind. A test samples. It runs some inputs and says nothing about the rest, so it can show a bug is there but never that none are. A proof is a finite object, and every step in it can be checked locally, so a checker examines all of it, every time. It can be simple and complete at once: complete for the proof, not for the world.

The art of testing is making the test simpler than the thing under test. A proof checker takes that to its limit. HOL Light's kernel is about 400 lines of OCaml, and the s2n-bignum project's soundness notes say no kernel soundness bug has been found since 2003.

Checkers are not infallible. This year an AI model, given expert tips, found seven proofs of False in Rocq's kernel and none in Lean's official kernel. The same search found false proofs in two independent Lean kernels. In July a Lean kernel bug turned up while someone reviewed a short Collatz "proof", and it was fixed the same day. Independent checkers buy diversity, not correctness. But each of these was a bug in a small program that people could read and fix, which is why "who checks the checker" converges instead of running away.

Proof changes the shape of the gap rather than removing it. Testing leaves gaps scattered through the behaviour. Proof leaves one, at the specification: does the formal statement say what we meant? Kallel and El Louadi argue that this check cannot be mechanised, even in principle, and ends in human judgement. OpenAI's proposed proofs from August show why it matters. They were set up to be rechecked by a second, independent kernel, but their statements rest on 218 locally written definitions, and whether those mean what mathematicians mean is not something a kernel decides.

That seam is where maths meets the real world, and checking it is sampling again: examples, translating back, expected consequences. In other words, testing. It is where the tester's craft matters most, and a gap in one known place is a very different risk from gaps everywhere.

So for AI output: where it can be put into a form that is checked exhaustively against a specification, trust in the source stops mattering, and the work moves to writing and testing the specification. Where it cannot, we are back to sampling.
