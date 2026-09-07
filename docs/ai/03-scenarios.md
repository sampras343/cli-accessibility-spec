# Scenarios — AI-Assisted CLI Accessibility

Concrete, technically-grounded scenarios. Each maps a **real disabled-user barrier** to a **specific CLI moment**, the **AI action**, and the **CLI-ACS criteria** it compensates for. `[Wrapper]` = runtime aid, `[Conformance]` = build/CI aid.

Format: **Trigger → Raw behavior today → AI action → Guardrail → Criteria**.

---

## Persona 1 — Maya, blind backend engineer (NVDA + terminal, screen reader)

### S1.1 — `git status` table read as gibberish `[Wrapper]`
- **Trigger:** runs `git status`; output uses columns + color for staged/unstaged.
- **Raw today:** screen reader reads `modified:   src/app.ts` runs of spaces as "modified blank blank blank src slash app dot ts" and color-coded state is silent.
- **AI action:** on hotkey, structure inferer detects the `state: path` pattern; LLM linearizes to: *"3 changes. Staged: src/app.ts (modified). Not staged: README.md (modified), test.log (new, untracked)."*
- **Guardrail:** raw output already spoken; AI is a second, requested view. Low parse confidence → no linearization, stays on raw.
- **Criteria:** OS-4 (linear order), CV-1 (color not sole channel), OS-5 (no layout-only info).

### S1.2 — `kubectl get pods` wide multi-column output `[Wrapper]`
- **Trigger:** 20 pods × 6 columns (`NAME READY STATUS RESTARTS AGE IP`).
- **Raw today:** column alignment is meaningless to a screen reader; user reads 120 cells linearly to find the one failing pod.
- **AI action:** "summarize" → *"19 of 20 Running. 1 not: `api-7f9` — CrashLoopBackOff, 34 restarts."* Follow-up "explain api-7f9" → pulls `kubectl describe`/logs on request.
- **Guardrail:** summary cites exact pod name so the user can verify against raw; never hides the raw table.
- **Criteria:** OS-4, OS-6 (structured output the tool *should* have offered), OS-8 (width).

### S1.3 — cryptic failure `[Wrapper]`
- **Trigger:** `terraform apply` exits 1 with a 200-line stderr stack.
- **Raw today:** user must read the whole wall to find the one actionable line.
- **AI action:** auto (on non-zero exit) announces: *"Failed: resource `aws_s3_bucket.data` already exists. Import it or rename. The error is at line 147."* — pointing INTO the raw, not replacing it.
- **Guardrail:** announced once, politely; user can ignore and read raw. No auto-fix.
- **Criteria:** ES domain (actionable errors), OS-1 (stderr separation).

---

## Persona 2 — Dev, low-vision developer (screen magnifier, 400% zoom)

### S2.1 — wide output overflows the magnifier viewport `[Wrapper]`
- **Trigger:** `docker ps` at 400% zoom — only ~30 chars visible per line; columns scroll off-screen.
- **Raw today:** user pans horizontally per row, losing which column is which.
- **AI action:** "reflow" → narrow vertical `KEY: value` blocks per container that fit the viewport without horizontal panning.
- **Guardrail:** deterministic reflow (no hallucinated values); values copied verbatim from raw cells.
- **Criteria:** OS-8 (width awareness), OS-4.

### S2.2 — color-only diff at low contrast `[Wrapper]`
- **Trigger:** `git diff` red/green with no `+`/`-` emphasis at the user's contrast settings.
- **AI action:** re-emit with explicit `ADDED:` / `REMOVED:` text prefixes alongside color.
- **Criteria:** CV-1 (color not sole channel).

---

## Persona 3 — Sam, RSI / limited hand mobility (voice-first, minimal typing)

### S3.1 — long command by voice `[Wrapper]`
- **Trigger:** says *"show me the last 20 lines of the nginx error log."*
- **AI action:** intent → `tail -n 20 /var/log/nginx/error.log`; reads it back; on "yes" runs it; summarizes result by voice.
- **Guardrail:** read-back before run; read-only command so single confirm.
- **Criteria:** IN domain (input), HD (discoverability without typing).

### S3.2 — destructive command by voice `[Wrapper]`
- **Trigger:** *"delete all the stopped docker containers."*
- **AI action:** candidate `docker container prune`; reads back **what will be lost** ("this permanently removes N stopped containers"); requires a **second** explicit "confirm delete."
- **Guardrail:** destructive-verb gate; STT mishear → repeat, never run.
- **Criteria:** IN domain; wrapper safety invariant #3.

### S3.3 — "what flag do I need?" `[Wrapper]`
- **Trigger:** *"how do I make curl follow redirects and show headers?"*
- **AI action:** answers from `curl --help`/man: *"`-L` follows redirects, `-i` shows headers → `curl -iL <url>`."* No man-page navigation, no typing.
- **Criteria:** HD (help & documentation).

---

## Persona 4 — Rae, cognitive disability (dyslexia + working-memory load)

### S4.1 — jargon error → plain language `[Wrapper]`
- **Trigger:** `npm install` fails with `ERESOLVE unable to resolve dependency tree`.
- **AI action:** *"Two packages need different versions of React. Easiest fix: run `npm install --legacy-peer-deps`. Safer fix: pick one React version."* — ranked, plain, actionable.
- **Guardrail:** suggests, does not run; explains the tradeoff, not just a magic flag.
- **Criteria:** ES domain (human-readable, actionable).

### S4.2 — reduce noise `[Wrapper]`
- **Trigger:** a build tool prints banners, tips, progress spam.
- **AI action:** "quiet summary" → just success/fail + what was produced.
- **Criteria:** OS-7 (quiet mode) — compensating when the tool has none.

---

## Conformance-side scenarios (build/CI) `[Conformance]`

### S5.1 — judge error-message quality (SEMI)
- **Trigger:** suite runs the tool with a bad flag; error goes to stderr, exit non-zero — `[AUTO]` passes structurally.
- **AI action:** scores the *text* against a rubric: does it name the problem? suggest a fix? avoid raw stack-only? → SEMI verdict + reason.
- **Guardrail:** labeled AI-assisted; a human still signs off; never reported as AT verification.
- **Criteria:** ES domain SEMI criteria.

### S5.2 — remediation PR (SEMI/AUTO fail)
- **Trigger:** `OS-6` fails — no `--json`.
- **AI action:** drafts the flag wiring + a linear `--plain` path as an advisory diff; suite re-runs on the patched output to validate the fix actually passes.
- **Guardrail:** advisory only; never auto-merged; must re-pass the deterministic check.
- **Criteria:** OS-6, OS-10.

### S5.3 — draft the CLI-ACR + crosswalk
- **AI action:** turns raw AUTO+SEMI signals into Section 7 report prose and the WCAG/EN 301 549/508 crosswalk narrative.
- **Guardrail:** SEMI marked AI-assisted, MANUAL marked still-required.

---

## Cross-cutting guarantees (apply to every wrapper scenario)

1. **Raw first, always** — AI output is additive; the real bytes already reached the user/AT.
2. **Verifiable** — AI cites exact names/lines/values so the user can check against raw.
3. **Confidence-gated** — low confidence → silent fallback to passthrough, never a confident guess.
4. **No silent execution** — propose → read-back → (2nd gate if destructive) → run.
5. **Local + private** — inference on-device; redaction before AI; cloud strictly opt-in.
