---
title: "Secure Model Selection, Privacy & AI Skills"
summary: "Route coding work to an appropriate model, protect sensitive context, and package repeatable engineering workflows"
---

# Secure Model Selection, Privacy & AI Skills

> Use the least powerful model that can safely do the job—and give it only the access it needs.

![Secure model selection, privacy and skills](/img/secure-model-selection-and-skills.png)

## TL;DR
- Choose models by task complexity, reliability, latency, cost, and data policy—not leaderboard rank alone.
- Treat AI-generated code and every tool call as untrusted until validated.
- Keep secrets and sensitive source out of prompts unless an approved environment, retention policy, and access boundary permit it.
- Use small, discoverable AI skills and instruction files to make proven workflows repeatable.

## Quickstart (Do this now)
1. Classify the task: routine edit, complex reasoning, security-sensitive change, or production action.
2. Define an allowlist of approved models and environments for each data classification.
3. Build a private evaluation set from real tasks, expected outcomes, and security cases.
4. Create one reusable skill or instruction for a frequent workflow, such as adding tests or reviewing a migration.
5. Test the skill on representative tasks and version it with its resources.

## The Idea (Slightly deeper)
Public benchmarks such as SWE-bench provide useful comparative signals, but they do not predict how a model will perform on your architecture, code conventions, dependency graph, or data policy. Maintain private evaluations that reproduce your own workflows and score quality, security, latency, cost, and recovery behavior.

Privacy begins before the prompt leaves a workstation. Minimize the information shared, redact secrets and personal data, control tool egress, and make provider retention/training settings part of vendor due diligence. For generated code, use the same secure-development checks you apply to human code—plus defenses for prompt injection and overbroad tool access.

## Deep Dive

| Task class | Default approach | Verification intensity |
| --- | --- | --- |
| Boilerplate or local refactor | Fast, cost-efficient model | Tests and diff review |
| Cross-module analysis | Reasoning-capable model + curated context | Plan review + tests |
| Security, auth, or data handling | Strong model in approved environment | Threat review, tests, human approval |
| External side effect | Bounded agent + policy gateway | Approval, audit, idempotency, rollback |

An **AI skill** is a focused, reusable workflow with a clear trigger, instructions, optional tools, and reference material. Use progressive disclosure: short descriptions help discovery; detailed references load only when the skill is activated. Keep skills narrow and test them like code.

## Core Skills
- **Model evaluation**: evaluate models on representative work and failure modes.
- **Data classification**: know what may enter a prompt, tool, log, or vendor environment.
- **Least privilege**: limit tool permissions, time, budget, scope, and network egress.
- **Skill authoring**: name the trigger clearly, keep guidance focused, and provide a verification checklist.

## When to Use This
- **Use when**: buying AI coding tools, setting organization policy, or introducing agent skills.
- **Use when**: code contains credentials, customer data, regulated information, or security controls.
- **Do not use**: a public benchmark score as a standalone production-readiness decision.

## Common Pitfalls
- **Benchmark theater**: test representative internal tasks before rolling out a model.
- **Secret leakage**: prevent tokens, `.env` files, customer data, and production dumps from entering context.
- **Skill bloat**: do not load a giant runbook for every task; split rare detail into references.
- **Trusting generated security code**: review authentication, authorization, injection defenses, and insecure defaults.

## Resources & Links
- **[SWE-bench](https://github.com/SWE-bench/SWE-bench)** — a benchmark of real GitHub issue-resolution tasks; use it as a signal, not a substitute for private evals.
- **[OWASP GenAI Security Project](https://genai.owasp.org/)** — guidance on LLM and agentic security risks.
- **[NIST Secure Software Development Framework](https://csrc.nist.gov/Projects/ssdf)** — a practical secure-development baseline.


## Next Steps
- Supply disciplined inputs with [Context, Model & Prompt](context-model-prompt.md).
- Apply least privilege to [Model Context Protocol](model-context-protocol.md) integrations.
