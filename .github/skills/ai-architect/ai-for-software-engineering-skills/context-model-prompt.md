---
title: "Context, Model & Prompt"
summary: "Give AI coding tools the right information, select a capable model, and request a verifiable outcome"
---

# Context, Model & Prompt

> Better AI coding starts with the right context, not a cleverer sentence.

![Context, model and prompt](/img/context-model-prompt.png)

## TL;DR
- **Context** is the task, relevant code, constraints, and evidence available to the model.
- **Model** choice should match task complexity, latency, cost, and risk.
- **Prompt** turns intent into steps, guardrails, and an output format that can be checked.

## Quickstart (Do this now)
1. Write one sentence describing the desired outcome and one sentence naming what must not change.
2. Select only the relevant files, interfaces, tests, and recent decisions.
3. Use a fast model for focused boilerplate; use a reasoning-capable model for ambiguous, security-sensitive, or architectural work.
4. Ask for a plan, assumptions, and validation steps before asking for edits.
5. Require a structured answer: changed files, tests run, remaining risks, and questions.

## The Idea (Slightly deeper)
Most unreliable AI output is an input-design problem. Too little context forces guessing; too much unrelated context dilutes attention and raises cost. Good context is a deliberately curated working set, not an entire repository pasted into a chat.

The Context → Model → Prompt sequence is a practical control loop. First decide what the tool may know; then choose the least expensive model that can perform the work reliably; finally specify the work and proof expected. Revisit the context as the task changes instead of letting conversation history become an unreviewed source of truth.

## Deep Dive

| Layer | Ask yourself | Good engineering outcome |
| --- | --- | --- |
| Context | Which source files, interfaces, tests, and constraints matter? | Relevant, traceable changes |
| Model | Is this routine, complex, risky, or time-sensitive? | Suitable quality at sensible cost |
| Prompt | What steps, boundaries, and proof should the tool produce? | A reviewable plan or diff |

Use the **plan–editor pattern** for larger work: a stronger reasoning model designs the approach, while a faster implementation model executes a reviewed plan. The separation reduces expensive reasoning on routine edits and makes the handoff artifact explicit.

### Prompt pattern

```text
Goal: Add <outcome>.
Context: Read <files/spec/tests>; this service uses <constraints>.
Steps: First explain the plan. Then implement only after the plan is accepted.
Constraints: Do not change <non-goals>; do not add dependencies; preserve <invariant>.
Proof: Run <commands> and report changed files, results, assumptions, and open risks.
```

## Core Skills
- **Context budgeting**: include the minimum authoritative material; summarize the rest.
- **Task-to-model routing**: use capability and risk, not brand preference, to choose a model.
- **Structured prompting**: name goal, steps, constraints, and output format.
- **Uncertainty handling**: instruct the tool to identify missing information rather than invent it.

## When to Use This
- **Use when**: onboarding to a new codebase, debugging multi-file behavior, or designing a change.
- **Use when**: agent responses are plausible but inconsistent.
- **Do not use**: a giant prompt as a replacement for a repository contract or durable specification.

## Common Pitfalls
- **Context dumping**: more files are not automatically more signal.
- **Model overkill**: reserve premium reasoning for decisions and difficult analysis.
- **Vague acceptance criteria**: “make it better” cannot be reliably evaluated.
- **Hidden assumptions**: ask the model to list assumptions and stop at decision boundaries.

## Resources & Links
- **[GitHub Copilot custom instructions](https://docs.github.com/en/copilot/how-tos/configure-custom-instructions-in-your-ide/add-repository-instructions-in-your-ide)** — durable repository guidance for AI tools.
- **[OpenAI prompting guide](https://platform.openai.com/docs/guides/prompt-engineering)** — techniques for clear, structured instructions.
- **[Agents.md](https://agents.md/)** — an open convention for repository-level agent instructions.

## Next Steps
- Turn prompts into durable contracts with [Specification-Driven Development](specification-driven-development.md).
- Package recurring guidance in [Secure Model Selection, Privacy & AI Skills](secure-model-selection-and-skills.md).
