# Scrutiny Report — `internal/modules/elasticsearch`

Date: 2026-08-14
Reviewer: opencode (scrutinize skill)
Scope: `internal/modules/elasticsearch/{usecase,delivery/http,presenter}` against active wiring in `internal/server/handlers.go` and `pkg/elasticsearch/client.go`.

> **Status update (same day):** all three MAJORs were fixed per user decision ("fix all 3"). See below.

## Intent

Expose HTTP endpoints that embed text via the configured LLM and store/search vectors in a dedicated Elasticsearch index. Claimed surface: health, create-index, embed+index, index-precomputed-vector, semantic search.

Simpler-alternative pass: **this module duplicates the `vectordata` module.** Both implement the same five operations over the same ES cluster, with two different request/response shapes and two different indexes (`cfg.Elasticsearch.Index` vs `cfg.VectorDB.Index`). The abstraction layer (`pkg/vectordb`) exists precisely so callers do not talk to ES directly. This module bypasses it (`usecase.go:8` imports `esclient` directly). Unless a consumer exists that needs the ES-specific shape, the module should be folded into `vectordata` or deleted — keeping both is more surface to maintain, and only `vectordata` is documented (`README_VectorData.md`).

## Findings (severity order)

### MAJOR-1 — `SearchVector` accepts unbounded `k` — ✅ FIXED
- **Finding:** `SearchVector` only defaults `k` when `<= 0`; no upper clamp (`usecase.go:126-129`), and `esclient.SearchKNN` passes `k` straight through as ES `size` and knn `k` with `num_candidates = k*10` (`client.go:195-211`).
- **Why it matters:** an authenticated client sending `k=100000` issues a search with `size=100000`, `k=100000`, `num_candidates=1,000,000` — ES rejects it (or thrashes) and every call 500s. `vectordata` already caps at 100; this module diverged.
- **Evidence:** trace `POST /api/elasticsearch/search` → `SearchVector` → `es.SearchKNN` → query `size`/`num_candidates`.
- **Fix applied:** `k` clamped to `[1,100]` in `SearchVector` (mirrors `vectordata`).

### MAJOR-2 — `CreateIndex` dereferences a nil `llm` — ✅ FIXED
- **Finding:** `CreateIndex` calls `u.llm.Dims()` with no nil guard (`usecase.go:39`), while `EmbedAndIndex` (`:52-54`) and `SearchVector` (`:122-124`) guard `u.llm == nil`.
- **Why it matters:** a nil `llm` panics on the first request to `POST /api/elasticsearch/create-index`. The recovery middleware converts it to a 500, but the panic is avoidable and the guard inconsistency is a landmine for future wiring/tests.
- **Evidence:** wiring at `server.go:122` always constructs `llmClient` today, so practical risk is low — but the code advertises nil-safety in two of three methods.
- **Fix applied:** `CreateIndex` now returns an error when `u.llm == nil` instead of panicking.

### MAJOR-3 — Index mapping dims come from config, not the actual embedding — ✅ FIXED
- **Finding:** `EmbedAndIndex` ensures the index with `u.llm.Dims()` (`usecase.go:56`) **before** embedding, and `dense_vector.dims` is baked immutably into the mapping (`client.go:117-122`). `llm.NewClient` defaults `dims` to **1536** when unset (`llm/client.go:35-38`); `nomic-embed-text` outputs 768.
- **Why it matters:** if `cfg.LLM.Dims` is unset or wrong, the index is created at 1536 dims and **every** subsequent `Index` call fails with a vector-dimension mismatch — permanently, until someone manually deletes the index. `vectordata` avoids this by embedding first and ensuring with `len(embedding)` (`vectordata/usecase.go:68-79`).
- **Evidence:** trace `EmbedAndIndex` → `EnsureIndex(dims=config)` → `Index(doc)` with real embedding; existing-index short-circuit at `client.go:140` never re-validates dims.
- **Fix applied:** `EmbedAndIndex` embeds **first**, then calls `EnsureIndex(len(embedding))`, so the index mapping always matches the actual embedding size.

### NIT-1 — Handler leaks raw backend errors — ✅ FIXED
- **Finding:** every usecase error is rendered via `ErrInternal(err)` which sends `err.Error()` verbatim (`handler.go:90,118,148`), i.e. full ES error bodies.
- **Why it matters:** exposes cluster internals (index names, shard errors) to API consumers.
- **Fix applied:** all five handlers route through `internalError(w, r, op, err)` — logs the full error server-side, returns a fixed `"internal server error"` message. Verified live: a 4-dim vector into the 768-dim index returned a generic 500 with no ES detail.

### NIT-2 — `validate:"required"` tags are decorative — ✅ FIXED
- **Finding:** `presenter.go:6,13,31` tag fields `required`/`min=1` but no validator runs; handlers re-implement checks manually (`handler.go:84,113,142`).
- **Why it matters:** misleading — same problem removed from `vectordata` presenter earlier.
- **Fix applied:** dropped the decorative tags; manual checks in handlers remain the source of validation (matching `vectordata`).

### NIT-3 — Metadata silently dropped — ✅ FIXED
- **Finding:** `EmbeddingRequest`/`IndexVectorRequest`/`SearchHit` have no `metadata`; the underlying `esclient` supports it (`client.go:247`). Search results drop `h.Metadata` (`usecase.go:143-148`).
- **Why it matters:** inconsistent with `vectordata`, which round-trips metadata.
- **Fix applied:** `metadata` added to `EmbeddingRequest`, `IndexVectorRequest`, and `SearchHit`; usecase sets `doc.Metadata` on index and maps `h.Metadata` into hits. Covered by new `usecase_test.go` cases and verified live (`POST /api/elasticsearch/embeddings` + `search` returned `"metadata":"meta-live-1"`).

### NIT-4 — `u.index` is log-only and can drift — ✅ FIXED
- **Finding:** the usecase stores `index` for logging (`usecase.go:26,31,46`), but the real index is baked into the ES client at construction. Passing a mismatched param silently logs the wrong index name.
- **Fix applied:** the `index` field (and constructor param) were removed; logging/responses use the new `esclient.Client.IndexName()` accessor, so the reported index can never drift from the one actually queried. Live response confirmed `"index":"vector_embeddings"`.

## Verification notes

- Wiring confirmed live: `handlers.go:203-211` guards `esClient != nil && llmClient != nil`; `vectordata` (handlers.go:213-222) shares the same `esClient` via `WithIndex`, so the two modules write to disjoint indexes. No cross-contamination observed.
- `Health` (usecase.go:34-36) delegates to `es.Health` (`GET /_cluster/health`) — verified against the fake ES server in the new unit tests.

## Verdict

**fix-then-ship** if the module is kept, but the recommendation is **rework**: the single biggest reason is MAJOR-1/2/3 fixability is trivial, yet the module's existence duplicates `vectordata`. Decide fold-vs-delete before investing in the three MAJOR fixes.

> **Status update:** user chose "fix all 3" — MAJOR-1 (k clamp), MAJOR-2 (nil-llm guard), and MAJOR-3 (embed-then-ensure-index with `len(embedding)`) are all applied. All four NITs (raw error leak, decorative validate tags, dropped metadata, log-only index) are also fixed, with new `internal/modules/elasticsearch/usecase/usecase_test.go` coverage (metadata round-trip, k clamp, nil-llm guards) and live verification against the running stack. The fold-into-`vectordata` recommendation still stands.
