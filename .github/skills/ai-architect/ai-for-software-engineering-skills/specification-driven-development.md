---
title: "Specification-Driven Development"
summary: "Use written requirements, contracts, plans, and tests as the source of truth for AI-assisted implementation"
---

# Specification-Driven Development

> Give an AI agent a contract to satisfy, not a vague feature request to interpret.

![Specification-driven development](/img/specification-driven-development.png)

## TL;DR
- A specification states the intent, constraints, interfaces, acceptance criteria, and proof before implementation.
- Specifications reduce ambiguity, make AI changes easier to review, and survive context-window limits.
- API-first contracts let teams and agents work in parallel while tests detect drift.

## Quickstart (Do this now)
1. Create a `specs/<feature>.md` with one outcome, non-goals, acceptance criteria, and edge cases.
2. For an API, write or update the OpenAPI contract before implementation.
3. Ask the agent to extract assumptions and propose a plan; review both before edits.
4. Break the plan into independently verifiable tasks.
5. Validate the diff against the spec and contract tests—not only against the agent's summary.

## The Idea (Slightly deeper)
Specification-driven development (SDD) decouples *what must be true* from *how it is coded*. That makes a feature understandable by humans and agents, creates a durable handoff between planning and implementation, and reduces “helpful” scope expansion.

A useful spec has five parts: high-level objective, mid-level objectives, implementation notes, beginning context, and ending context. The last two matter in agentic work: they say what exists before the change and what must be true after it. Keep the document versioned with the code.

## Deep Dive

```mermaid
flowchart LR
  A[Specify outcome and constraints] --> B[Plan and inspect assumptions]
  B --> C[Create tasks and contracts]
  C --> D[Implement small changes]
  D --> E[Validate with tests and contract]
  E --> F[Review against specification]
```

### Minimal feature specification

```markdown
# Customer export
## Objective
Authenticated users can export their own customer records as CSV.
## Acceptance criteria
- Export excludes deleted records and other tenants.
- Request above 10,000 rows becomes an asynchronous job.
- Audit event contains actor, tenant, and export ID.
## Non-goals
- No spreadsheet formatting or cross-tenant reporting.
## Proof
- Unit, authorization, and API contract tests pass.
```

Treat OpenAPI, schemas, and ADRs as executable or reviewable parts of the specification. If implementation and contract disagree, fix the disagreement explicitly—do not quietly make the code the source of truth.

## Core Skills
- **Requirements writing**: make outcomes, boundaries, error cases, and non-goals testable.
- **Contract-first design**: define API and data contracts ahead of service internals.
- **Task decomposition**: create small steps with independent evidence.
- **Traceability**: connect requirement → plan → code → test → review.

## When to Use This
- **Use when**: a change crosses components, has business rules, or will be implemented by an agent.
- **Use when**: front-end and back-end teams need a shared contract.
- **Use a lighter form when**: fixing a contained defect with a clear failing test.

## Common Pitfalls
- **Vague objectives**: an agent cannot verify “modernize” or “improve performance” without a measure.
- **Code-first drift**: update the contract and spec whenever a deliberate behavior change occurs.
- **Missing failure paths**: specify authorization, validation, retries, and idempotency where relevant.
- **Stale instructions**: review `AGENTS.md`, `CLAUDE.md`, and tool-specific instructions as part of the change.

## Frameworks at a Glance
- **Spec Kit**: a lightweight way to make specification-driven work repeatable. It helps teams move from a feature intent to a versioned specification, implementation plan, and small tasks. Use it when you want consistent artifacts that an engineer or agent can review and implement; the specification and tests remain the source of truth.
- **BMAD Method**: a role-based delivery framework for work that needs more product and architecture shaping. It right-sizes the process: Quick Flow for clear small changes; the full method for products and complex features; and an enterprise track for compliance-heavy work. Its path is analysis → planning → solutioning → implementation, with a readiness gate before code is written.
- **Superpowers**: a collection of composable agent skills that enforces engineering discipline inside the implementation loop. Its workflow turns an approved design into a detailed plan, uses isolated workspaces and focused subagents where useful, then requires test-first development, review, systematic debugging, and fresh verification evidence before completion.
- **AI Development Lifecycle (AI DLC)**: the team-level loop that connects the methods: discover the problem → specify outcomes and constraints → plan and design → implement → validate quality, safety, and behavior → release, observe, and improve. Keep people at approval points and feed production lessons back into the next specification.

## Next Steps
- Coordinate specialized implementation steps with [Multi-Agent Coding Workflows](multi-agent-coding-workflows.md).
- Enforce the resulting contract in [AI PR Review Automation](ai-pr-review-automation.md).
