# Purpose

- Use AI to assist with writing code without changing the system's existing behavior.
- Keep the scope of every change narrow, reviewable, and safe for production.

## Communication language

- These instructions are written in English to save input tokens, but ALWAYS reply to
  the user in Thai so the team can read easily. Keep code, identifiers, commands, and
  file paths unchanged.

## Before writing any code, always:

1. Summarize the current behavior of the code as you understand it.
2. Identify edge cases that may be affected.
3. Explain only the necessary, minimal-impact approach.

Do not write code until you receive explicit confirmation.
If anything is unclear, always ask first — never assume.

## Rules

- Preserve the system's existing behavior.
- Make the smallest change necessary (minimal diff).
- Avoid designs or changes more complex than needed.
- Strictly follow the existing code style of the codebase.
- Do not change logic, data flow, or performance characteristics.
- Do not add abstractions, helpers, libraries, or refactor structure without an explicit instruction.
- Do not guess; if an assumption is required, ask first every time.

## Acceptance Criteria

- The system's behavior must remain identical except where explicitly specified.
- Changes must be narrow in scope and verifiable from the diff.
- Do not add unnecessary logic, queries, loops, or dependencies.
- Must not affect overall system performance.
- Every change must have a clearly explainable rationale and impact.
- If any acceptance criterion cannot be met, report it before writing code.
- Do not modify code until the reviewer explicitly confirms.
