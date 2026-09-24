Since you're working with the **debug-mantra** discipline, here’s exactly how **Git commands** plug into each step of the debugging process.

---

## Git Commands Mapped to Debug-Mantra

### 1. Reproduce Reliably → Pin the Environment
Use Git to capture the **exact state** of the code at failure time.

| Command | Purpose |
|---------|---------|
| `git rev-parse HEAD` | Record the exact commit SHA in your bug report. |
| `git checkout <SHA>` | Jump to that exact commit to test the failure. |
| `git stash` | Stash local uncommitted changes to ensure a **clean, reproducible** state. |
| `git clean -fd` | Remove untracked files that might interfere. |

> **Pro tip**: If the bug exists on `main` but not on `feature`, `git diff main..feature` shows exactly what changed.

---

### 2. Know the Fail Path → Trace History
Find **when** and **who** introduced the failure path.

| Command | Purpose |
|---------|---------|
| `git blame -L <start>,<end> <file>` | See who last changed each line of the failing code. |
| `git log -p -- <file>` | View the full diff history of a file (all changes). |
| `git log -S"<string>"` | Search for a specific string being added/removed (e.g., `git log -S"panic("`). |
| `git log --grep="<pattern>"` | Search commit messages for relevant keywords. |

---

### 3. Question Your Hypothesis → Binary Search (Bisect)
Test your hypothesis by finding the **exact commit** that broke things.

```bash
git bisect start
git bisect bad <commit_with_bug>   # current failing commit
git bisect good <known_good_commit> # last known working commit
```

Then for each step:
```bash
# Build/run your repro script
go test -run TestBroken -v

# Mark the commit
git bisect good   # if test passes
git bisect bad    # if test fails
```

Git will automatically narrow down to the **culprit commit**.  
*(End with `git bisect reset` to return to normal.)*

---

### 4. Every Run is a Breadcrumb → Tag & Branch Each Debug Session
Keep your debugging history organized so you can cross-reference runs.

| Command | Purpose |
|---------|---------|
| `git tag debug/run1-<date>` | Tag a specific reproduction run. |
| `git checkout -b debug/hypothesis-X` | Create a branch for each hypothesis you test. |
| `git diff debug/run1 debug/run2` | Compare two runs to see what changed between them. |
| `git log --oneline --graph --all` | Visualize the branching context of all your debug attempts. |

---

## Quick Cheat Sheet for Debugging

| What you want | Git command |
|---------------|-------------|
| "Is this bug in production code or my local changes?" | `git stash; git reset --hard HEAD` (clean slate) |
| "Did I accidentally revert a fix?" | `git log --oneline -n 10` |
| "What files changed in the last 5 commits?" | `git diff HEAD~5 HEAD --name-only` |
| "Who wrote this suspicious function?" | `git blame -L 100,120 main.go` |
| "When did we stop calling that function?" | `git log -S"oldFunctionName("` |
| "Show me all changes between v1.0 and v1.1" | `git diff v1.0 v1.1` |

---

## Example Debug Session with Git

```bash
# 1. Reproduce
git checkout bug-report-123
git clean -fd
go test -run TestLogin  # fails

# 2. Trace fail path
git blame -L 45,55 auth/login.go
# → line 48 was last changed by "Alice" in commit abc123

# 3. Hypothesis: "Alice's commit broke it"
git bisect start
git bisect bad HEAD
git bisect good abc122  # commit before Alice's
# ... test each bisect step ...
# → Git finds commit abc123 as the culprit

# 4. Breadcrumb
git tag debug/bisect-result-$(date +%s)
git show abc123  # inspect exactly what changed
```

---

**Summary**: Git isn't just for version control—it's your **time-machine debugger**. Use `bisect` to find *when*, `blame` to find *who*, `diff` to find *what*, and `tags/branches` to keep every experiment traceable.