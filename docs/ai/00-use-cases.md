# AI for CLI Accessibility — Use Cases

How AI is applied to make command-line tools accessible, and to help disabled users work with tools that aren't. Two primary use cases, one shared technical foundation (**on-device / local LLM inference**).

**Related design docs:**
- [Use Case 1 — Conformance Evaluation](01-conformance-ai.md) (build / CI time)
- [Use Case 2 — Accessible Runtime Wrapper](02-accessible-wrapper-ai.md) (end-user aid)
- [Scenarios](03-scenarios.md) (persona-driven, concrete commands)

---

## Shared foundation

On-device LLM inference is the enabling primitive across both use cases:

- **Private by default** — source, pre-release binaries, and terminal output (tokens, keys, PII) never leave the machine.
- **Reproducible** — a pinned model + fixed prompts yield deterministic verdicts, essential for conformance.
- **Offline-capable** — works air-gapped in CI and on locked-down or disconnected user machines.

---

## Use Case 1 — Conformance Evaluation (build / CI time)

**Actor:** CLI developer / QA / accessibility auditor.
**Problem:** CLI-ACS has 95 criteria; ~50 are deterministically testable, but ~45 `[SEMI]` criteria depend on *language quality* (is an error actionable, is help clear, is output meaningful read linearly) that regex/exit-code checks cannot judge.
**AI role:** a local LLM acts as the `[SEMI]` judgment engine — scoring captured stderr/help/structured output against versioned per-criterion rubrics, emitting verdicts + self-validating remediation diffs, and drafting the CLI-ACR report and standards crosswalk.
**Boundary:** never touches `[AUTO]` (pure code) or the `[MANUAL]` verdict (real assistive-technology testing).
→ Details: [01-conformance-ai.md](01-conformance-ai.md)

## Use Case 2 — Accessible Runtime Wrapper (end-user aid)

**Actor:** disabled CLI user (blind, low-vision, motor-impaired, cognitive).
**Problem:** the terminal exposes only a character-cell matrix — no accessibility tree — so screen readers linearize columns into noise, color-coded state is unspoken, and cryptic errors bury the one actionable line.
**AI role:** a local capture layer tees streams (raw reaches AT untouched); on failure or request, a structure-inferer + small on-device LLM linearizes tables, explains errors in plain language, and answers "what flag do I need"; an optional voice channel maps speech → command with read-back and destructive-verb gating.
**Boundary:** additive runtime aid, not conformance — cannot add data the CLI never emitted, only best-effort inference, confidence-gated to raw passthrough.
→ Details: [02-accessible-wrapper-ai.md](02-accessible-wrapper-ai.md)

---

## How AI can be applied in these kinds of cases — brainstorm

Beyond the two core use cases, patterns for applying AI to CLI accessibility. Grouped by function; each notes the realistic value and the guardrail.

### A. Perception — reconstruct missing structure
- **Semantic layout inference** — detect tables, headings, prompts, and progress from the raw character matrix and expose a synthesized "accessibility tree." *Guardrail:* deterministic pre-parse first, LLM only for ambiguous cases, confidence-gated.
- **Output summarization / triage** — collapse long output to "what matters" (which pod failed, what changed). *Guardrail:* cite exact identifiers so the user can verify.
- **Diff/delta narration** — describe what changed between two runs or two files, not just print both.

### B. Comprehension — make meaning accessible
- **Error explanation + next step** — translate cryptic errors + exit codes into plain language with a ranked fix. *Guardrail:* suggest, don't auto-run; explain tradeoffs.
- **Jargon / plain-language mode** — rewrite dense output for cognitive/dyslexic users on demand.
- **Man-page / help Q&A** — answer "how do I do X" from the tool's own `--help`/man, no navigation.

### C. Interaction — reduce input cost
- **Intent → command** — natural language (typed or spoken) to a candidate command with read-back. *Guardrail:* propose-not-execute; second gate on destructive verbs.
- **Voice-first control** — local STT + TTS loop for motor-impaired users; keyboard remains fully sufficient (voice never required).
- **Adaptive verbosity** — learn a user's preferred detail level and apply it to summaries (local, per-user).

### D. Authoring — help developers build accessible CLIs
- **`[SEMI]` quality judging** — score error/help text quality against rubrics in CI.
- **Remediation generation** — draft the flag wiring, `--plain`/`--json` mode, or an error rewrite as an advisory diff that must re-pass the deterministic check.
- **Conformance report drafting** — generate CLI-ACR prose and WCAG/EN 301 549/508 crosswalk narrative from raw signals.
- **Test-matrix inference** — read a tool's `--help` and generate the invocation/edge-case matrix the suite should run.
- **Pre-commit accessibility linting** — flag likely violations (color-only output, missing `--json`, vague errors) at author time.

### E. Cross-cutting enablers
- **Local RAG over tool docs** — ground answers in the specific tool's man pages / help to cut hallucination.
- **Structured-output coercion** — when a tool has no `--json`, best-effort parse to structured data for AT consumption (clearly marked inferred, not authoritative).
- **Redaction model** — detect and strip secrets/PII before any text reaches AI (needed for both cloud opt-in and even local logging).

---

## What AI must NOT do (applies across all cases)

- **Not** replace `[MANUAL]` assistive-technology verification — an LLM "pass" is not a screen-reader test.
- **Not** sit in the execution hot path or buffer the raw stream — raw output must reach the user/AT immediately.
- **Not** auto-execute commands — propose, read back, gate destructive actions.
- **Not** emit confident guesses — low confidence falls back to raw passthrough.
- **Not** exfiltrate terminal or source data — local by default, cloud strictly opt-in after redaction.
