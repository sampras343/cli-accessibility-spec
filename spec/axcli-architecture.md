# axcli: Accessible CLI Wrapper -- Definitive Architecture

## Project Name and Scope

**axcli** -- a Python middleware between the user and any CLI binary that accepts natural language or raw commands via keyboard (primary) or voice (optional), translates intent into CLI invocations, enforces safety guardrails, and presents results as structured progressive-disclosure text optimized for screen readers, TTS, and braille displays.

Scope: axcli makes ANY CLI tool accessible regardless of its CLI-ACS conformance level. The CLI-ACS conformance suite (Go, at `/home/sacm/Documents/Study/my_projects/cli-accessibility-spec/`) evaluates tools; axcli compensates for what they lack.

axcli is NOT a shell replacement. It is a prefix wrapper (`axcli gh pr list`) or interactive REPL (`axcli`) that enhances individual command interactions. It preserves the user's shell, aliases, tmux sessions, and screen reader configuration.

---

## Design Principles

1. **Text first, speech second.** All output goes to stdout as labeled plain text before any TTS fires. Braille displays, screen readers, piped output, and log files all work without special modes.
2. **Keyboard primary, voice optional.** The tool must be fully usable with zero audio hardware. Voice is an accelerator, not a requirement.
3. **Screen reader deference.** When a screen reader is detected, axcli disables its own TTS. It never double-speaks.
4. **Allowlist safety, not blocklist.** Known-safe read-only commands execute without confirmation. Everything else confirms. The surface of destructive commands is unbounded; the surface of safe commands is small and auditable.
5. **Local first.** Commands may contain secrets. The LLM runs on localhost by default. Cloud is opt-in.
6. **Offline graceful degradation.** Without an LLM, axcli still runs raw commands with output formatting, ANSI stripping, screen reader labels, and progressive disclosure. The AI layer is an enhancement, not a dependency.
7. **Latency budgets are hard constraints.** Raw command passthrough: <100ms overhead. LLM-powered translation: first audio/text feedback within 1 second.

---

## Relationship to CLI-ACS Conformance Suite

The CLI-ACS suite at `/home/sacm/Documents/Study/my_projects/cli-accessibility-spec/` is a Go binary (`cli-acs`) that probes CLI tools and evaluates them against accessibility criteria. axcli depends on it in two ways:

**Dependency 1: ProbeResult for tool discovery.** axcli needs to know a target binary's subcommands, flags, flag types, and capabilities (--json, --quiet, --color, --screen-reader). The cli-acs codebase already has this in `internal/probe/discovery.go` via the `Probe()` function, with parser cascade (`CobraParser`, `ClapParser`, `ArgparseParser`, `GenericParser`) in `internal/probe/parsers.go`. However, there is no CLI subcommand that exports a raw `ProbeResult` as JSON. **Prerequisite work: add a `cli-acs probe <binary> --format json` command to the Go binary** (approximately 40 lines of Go: call `probe.Probe()`, marshal the result, write to stdout). This is the single most important prerequisite. Without it, axcli must either shell out to `cli-acs check` and parse the conformance report (which contains Results, not ProbeResult), or rewrite the parsers in Python.

**Dependency 2: ConformanceReport for gap compensation.** On first use with a target binary, axcli can optionally run `cli-acs check <binary> --format json` to get a conformance report. This tells axcli which accessibility gaps to compensate for: if CV (color) domain fails, always strip color; if HD (help) domain fails, use LLM to explain help text; if the tool lacks --json, parse text tables with heuristics.

**Known gaps in cli-acs that axcli exposes:**
- `Subcommand.Flags` is always empty (discovery.go line 100-105 populates `Flags: []Flag{}`). axcli needs subcommand-level flags for constrained LLM generation. **Fix required: recurse into subcommand --help and parse flags.**
- `TakesValue` is always false (parsers.go line 387: `TakesValue: false // TODO`). **Fix required: detect value placeholders (`<file>`, `FILE`, `=VALUE`) in flag descriptions.**
- `IsBlockedSubcommand()` in safety.go is a blocklist. axcli inverts this to an allowlist model.

**Feedback loop:** axcli usage data can propose new CLI-ACS criteria: `--screen-reader` flag as a spec criterion, labeled output sections as a universal pattern, streaming output guidelines, voice interface metadata format.

---

## Component Architecture

```
                         USER INPUT
                    (keyboard or voice)
                            |
                    +-------v--------+
                    |  Session Shell |  (REPL or prefix mode)
                    |  (entry point) |
                    +-------+--------+
                            |
               raw command? | NL input?
              +-------------+-------------+
              |                           |
     +--------v--------+        +--------v---------+
     | Passthrough Path |        |   Intent Engine  |
     | (no LLM needed) |        |   (LLM layer)    |
     +--------+---------+       +--------+---------+
              |                           |
              +-------------+-------------+
                            |
                   +--------v--------+
                   |   Safety Gate   |
                   | (allowlist + confirm)
                   +--------+--------+
                            |
                   +--------v--------+
                   |    Executor     |
                   | (subprocess)    |
                   +--------+--------+
                            |
                   +--------v--------+
                   |    Formatter    |
                   | (progressive   |
                   |  disclosure)    |
                   +--------+--------+
                            |
              +-------------+-------------+
              |                           |
     +--------v--------+        +--------v---------+
     |  stdout labels  |        | TTS Queue (opt)  |
     | (screen reader) |        | (no SR detected) |
     +-----------------+        +------------------+
```

### 1. Session Shell (`axcli/__main__.py`, `axcli/shell.py`)

The entry point. Two modes:

- **Prefix mode** (default): `axcli gh pr list` wraps a single command, formats output, returns to parent shell. Works with pipes (`axcli kubectl get pods | grep auth`), tmux, ssh, aliases, shell scripts. When stdout is not a TTY, Formatter outputs raw text (no labels, no progressive disclosure) for pipe compatibility.
- **REPL mode**: `axcli` with no args opens an interactive session. Uses Python `readline` (NOT prompt_toolkit -- it fights screen reader buffer tracking) for input with history, tab completion against ToolProfile subcommands/flags. Maintains session context across turns.

Screen reader detection at startup:
- Linux: check if `orca` process is running (`pgrep -x orca`), or check AT-SPI via `dbus-send --session --dest=org.a11y.Bus --print-reply /org/a11y/bus org.a11y.Bus.GetAddress` (non-zero response means AT-SPI active). Also check for `FENRIR_RUNNING` env, `speakup` kernel module loaded.
- macOS: `defaults read com.apple.universalaccess voiceOverOnOffKey` or check process `VoiceOver`.
- Windows: check processes `nvda.exe`, `jfw.exe` (JAWS), `narrator.exe`.
- Override: `AXCLI_SCREEN_READER=1` env var (always on) or `AXCLI_SCREEN_READER=0` (always off).
- Re-check every 30 seconds in REPL mode (not just startup) to handle screen reader start/stop mid-session.

Honors: `NO_COLOR`, `TERM=dumb`, `AXCLI_SCREEN_READER`, `AXCLI_TTS`, `AXCLI_VOICE`, `AXCLI_LLM_URL`, `AXCLI_LLM_MODEL`.

**Technology:** Python 3.11+ stdlib (`readline`, `subprocess`, `os`, `signal`, `json`, `tomllib`). Zero external dependencies for this component.

### 2. Tool Registry (`axcli/registry.py`)

On first use with a target binary (e.g., `gh`), discovers its capabilities and caches the result.

**Discovery path (preferred):** Shell out to `cli-acs probe gh --format json` (once the `probe` subcommand is added to the Go binary). Parse the JSON ProbeResult. Cache at `~/.cache/axcli/tools/{binary_name}_{version_hash}.json`.

**Discovery path (fallback):** If `cli-acs` is not installed, run `{binary} --help` directly and parse with a minimal Python port of the `GenericParser` logic (~50 lines: detect flag lines starting with `-`, detect indented subcommand lines under section headers containing "command"). This fallback is less accurate but provides basic subcommand/flag lists.

**Cache invalidation:** Compare binary mtime. User can force refresh via `axcli refresh gh`.

**From ToolProfile, generate:**
- LLM tool schema (JSON function definitions constraining flag names and types).
- STT custom vocabulary list (subcommand names, flag names) for Vosk/Whisper.
- Allowlist of read-only subcommands (list, get, show, status, log, diff, describe, inspect, cat, head, view, version, help, info, whoami, config get).

**Technology:** Python stdlib (`subprocess`, `json`, `hashlib`, `pathlib`). Optional: `cli-acs` Go binary.

### 3. Intent Engine (`axcli/intent.py`)

Translates natural language to structured `CommandIntent`. Only invoked for NL input; raw commands (detected by starting with a known binary name or path) bypass this entirely.

**Input:** User text + ToolProfile (subcommands, flags) + SessionContext (last N commands with structured output).

**Output:** `CommandIntent` dataclass:
```python
@dataclass
class CommandIntent:
    binary: str
    args: list[str]        # ["pr", "list", "--state", "open"]
    explanation: str       # "List your open pull requests"
    is_read_only: bool     # True if subcommand is in the read-only allowlist
    confidence: float      # 0.0-1.0
```

**LLM backend:** Provider-agnostic. Accepts any OpenAI-compatible API endpoint.
- Default: Ollama on `localhost:11434` with Qwen2.5-Coder 7B (Q4_K_M quantized, ~4GB RAM, good function-calling support).
- Cloud fallback: Any OpenAI-compatible endpoint (Claude via proxy, GPT-4.1, etc.) configured by `AXCLI_LLM_URL` and `AXCLI_LLM_MODEL`.
- Library: `httpx` (single dependency for HTTP, already widely installed) or stdlib `urllib.request` for zero-dep mode.

**System prompt structure** (under 300 words, safety constraints at start AND end):
```
SAFETY: Never generate commands not in the tool schema. If unsure, say so.
You translate natural language into CLI commands for {binary} on {os}/{shell}.
Available subcommands and flags: {tool_schema_json}
Session context: {last_3_commands_structured}
Output JSON: {"args": [...], "explanation": "...", "is_read_only": bool, "confidence": float}
Few-shot examples: {3_examples_for_this_tool}
Temperature: 0
SAFETY: Never invent flags not in the schema. Express uncertainty rather than guessing.
```

**Latency mitigation:**
- Pre-warm model on first launch (send a dummy request so Ollama loads the model into VRAM/RAM). Measure TTFT; warn user if >2 seconds.
- For common patterns (`list X`, `show X`, `get X`, `status`), use regex pattern matching against ToolProfile subcommands BEFORE falling back to LLM. This handles 60-70% of commands with zero LLM latency.
- Stream LLM JSON response; parse args as soon as the args array closes (do not wait for explanation field).

**When LLM is unavailable:** Fall back to fuzzy matching user input against ToolProfile subcommand names (stdlib `difflib.get_close_matches`). "show pull requests" -> fuzzy matches "pr" + "list". Lower confidence, but functional.

**Technology:** Python stdlib (`json`, `urllib.request`, `dataclasses`, `difflib`). Optional: `httpx` for cleaner HTTP.

### 4. Safety Gate (`axcli/safety.py`)

**Allowlist model** (inverted from all three proposals' blocklist approach):

A curated allowlist of read-only subcommand verbs that skip confirmation entirely:
```python
READ_ONLY_VERBS = {
    "list", "get", "show", "status", "log", "diff", "describe",
    "inspect", "cat", "head", "tail", "view", "version", "help",
    "info", "whoami", "config-get", "env", "top", "events", "count",
    "search", "find", "check", "verify", "test", "lint", "explain",
    "completion", "api-versions", "api-resources",
}
```

Everything NOT on this list requires plain-language confirmation:
```
status: I will run: gh pr close 287
status: This closes pull request 287 "update dependencies". Proceed? [y/N]
```

**High-severity patterns** (static regex, additional warning):
```python
HIGH_SEVERITY = [
    r"rm\s+(-[rRf]+\s+|.*--force)",
    r"(DROP|DELETE|TRUNCATE)\s+(TABLE|DATABASE|FROM)",
    r"git\s+push\s+.*--force(?!-with-lease)",
    r"kubectl\s+delete\s+.*--all",
    r"chmod\s+(000|777)\s",
    r"mkfs\.", r"dd\s+if=", r"shutdown", r"reboot",
    r"curl\s+.*\|\s*(ba)?sh",
    r"find\s+.*-delete",
]
```

These require typing `yes` (not just `y`) and include an explicit consequences warning.

**Command echo:** ALWAYS echo the full generated command before execution, including inferred flags. `status: Running: gh pr list --state open --author @me`. This lets the user (who may be an expert) catch LLM errors from the command, not just the explanation.

**No LLM-based safety classification in the hot path.** The LLM generated the command; asking the same LLM "is this safe?" is circular. The allowlist + high-severity regex + human confirmation is the safety stack. LLM semantic classification is deferred to Phase 4 as an optional enhancement using a second, smaller model.

**Technology:** Python stdlib (`re`, `dataclasses`). Zero dependencies.

### 5. Executor (`axcli/executor.py`)

Runs approved commands as subprocesses.

```python
proc = subprocess.Popen(
    [binary] + args,          # NEVER shell=True; args as list
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    env={**os.environ, "NO_COLOR": "1", "TERM": "dumb"},
    timeout=None,             # User sends SIGINT via Ctrl+C
)
```

**Critical design decisions:**
- `shell=False` always. Arguments passed as list. This prevents command injection from LLM-generated argument values.
- `subprocess.Popen` (not `run`) with line-by-line stdout/stderr reading on background threads. Enables: (a) streaming output, (b) heartbeat ("still running, 10 seconds elapsed"), (c) cancellation via Ctrl+C forwarded as SIGINT.
- Environment: set `NO_COLOR=1`, `TERM=dumb`. If ToolProfile shows the binary supports `--screen-reader` or `--accessible`, append those flags. If ToolProfile has `HasJSONFlag` and the command produces tabular output, append `--json` for structured capture.
- Default timeout: 120 seconds (not 30 -- real commands like `docker build` need time). Configurable per tool.
- Heartbeat: every 5 seconds during execution, print `status: still running, {N} seconds elapsed` to stdout. This prevents the "did it crash?" anxiety for screen reader users.

**Interactive command detection:** Maintain a list of commands that require stdin/TTY: `vi`, `vim`, `nano`, `less`, `more`, `ssh`, `sudo`, `python`, `irb`, `node`, `docker login`, `kubectl edit`. When detected, warn: `status: This command opens an interactive session. Passing through to terminal.` Then exec with `os.execvp` (replaces the wrapper process) or allocate a PTY on Unix via `os.openpty`. On Windows, document as unsupported for TTY commands until Phase 4.

**Technology:** Python stdlib (`subprocess`, `threading`, `os`, `signal`).

### 6. Formatter (`axcli/formatter.py`)

The core accessibility component. Transforms raw output into labeled, progressive-disclosure structured text.

**ANSI stripping:** Port the regex from `internal/probe/ansi.go`:
```python
ANSI_RE = re.compile(r'\x1b\[[0-9;]*[a-zA-Z]')
def strip_ansi(text: str) -> str:
    return ANSI_RE.sub('', text)
```

**Output labels** (screen reader navigation landmarks):
```
you: show me my pull requests
status: Running: gh pr list --state open --author @me --json number,title,updatedAt
result: 3 open pull requests.
result: 1. Number 123: Fix login bug, updated yesterday.
result: 2. Number 456: Add dark mode, updated 3 days ago.
result: 3. Number 789: Refactor auth, updated last week.
error: Push failed. Exit code 1.
```

**Progressive disclosure tiers:**

- **Tier 1 -- Headline** (always shown): One sentence. Exit code + primary result. `result: 3 open pull requests.` or `error: Build failed with 2 errors.`
- **Tier 2 -- Summary** (auto-shown for <20 items, on-demand for more): Key items as `Label: value` pairs, one per line. No tables. Tables converted to `Header: value` sentence pairs per Emacspeak pattern.
- **Tier 3 -- Detail** (on-demand: user types `more` or `detail 2`): Full output, paginated at 10 items. `next`, `prev`, `search <term>` for navigation.

**Content-aware classification** (NOT line-count-only):
- Exit code != 0 or non-empty stderr: ALWAYS show error lines first, regardless of output length.
- JSON output (from --json flag): Parse into indexed array, store in SessionContext.
- Tabular output (2+ whitespace-separated columns with a header line): Parse into indexed array. Announce column count and row count. Render as `Header: value` pairs.
- Large output (>50 lines): Headline only. Announce line count. `result: 847 lines of output. Type 'more' to see first page, or 'search ERROR' to find errors.`

**Compact mode** for braille displays (`AXCLI_COMPACT=1`): Use abbreviations, skip full sentences. `auth-svc CrashLoop R:47 3d` instead of `Pod auth-service. Status: CrashLoopBackOff. Restarts: 47. Age: 3 days.`

**LLM summarization** (optional, for Tier 2 when output is complex): Send output + command to Intent Engine with a TTS-optimized prompt: `Summarize this CLI output in one paragraph for a screen reader user. Lead with the result. Use cardinal numbers. No special characters. No ANSI codes.` Only used when heuristic formatting produces poor results (e.g., unstructured verbose logs).

**Pipe detection:** When stdout is not a TTY (`not sys.stdout.isatty()`), skip all labels and progressive disclosure. Pass raw output through for pipe compatibility.

**Technology:** Python stdlib (`re`, `json`, `csv`, `textwrap`). Zero dependencies.

### 7. Session Context (`axcli/session.py`)

In-memory state for multi-turn reference resolution.

```python
@dataclass
class Turn:
    command: list[str]               # ["gh", "pr", "list", ...]
    stdout: str                      # raw output (truncated to 4000 chars)
    stderr: str
    exit_code: int
    structured: list[dict] | None    # parsed table/list as indexed array
    timestamp: float

@dataclass
class SessionContext:
    turns: list[Turn]                # rolling window, last 10
    target_binary: str
    tool_profile: dict | None        # cached ToolProfile
    summary: str                     # one-sentence session summary
```

**Structured data extraction:** After each command, attempt to parse stdout:
1. If JSON (starts with `[` or `{`): parse with `json.loads`, store as `structured`.
2. If tabular (heuristic: header line with 2+ space-separated words, followed by rows with same column alignment): parse into list of dicts.
3. Otherwise: split into numbered lines for line-reference (`detail 7`).

**Reference resolution:** When the Intent Engine receives "close the second one", it gets `SessionContext.turns[-1].structured[1]` injected into the LLM prompt. The LLM resolves the reference against the actual data, not its memory.

**Token budget:** When injecting into LLM context, truncate older turns to their headline summary only. Keep the last 3 turns with full structured data. Total session context injection: under 2000 tokens.

**Persistence:** In-memory only by default (matches terminal session lifecycle). On session exit, optionally write to `~/.cache/axcli/session.json` for resume. **Security:** Never persist if output contained patterns matching API keys, tokens, or passwords (regex: `(sk-|ghp_|token|password|secret|key)[\s=:]+\S+`).

**Technology:** Python stdlib (`dataclasses`, `json`, `re`, `time`).

### 8. Audio Subsystem (`axcli/audio.py`) -- OPTIONAL

Only active when: (a) no screen reader detected, (b) user opts in (`AXCLI_TTS=1` or `AXCLI_VOICE=1`), (c) audio dependencies are installed.

**TTS output:**
- Primary: espeak-ng via subprocess (`espeak-ng "text"`). Instant (<50ms), clear at high speed, familiar to blind users. Respects system speech rate.
- Upgrade: Piper (`piper-tts` v1.6.0, pip install) for neural quality.
- macOS fallback: `say "text"`.
- Architecture: Background worker thread consuming from `queue.Queue`. Interrupt: `queue.Queue.clear()` + kill espeak-ng subprocess. New speech always interrupts current speech (Emacspeak pattern).

**STT input (voice mode):**
- Activation: Toggle key (press Ctrl+Space to start recording, press again to stop). NOT hold-to-talk (terminals cannot detect key release; hold-to-talk requires global hotkey via platform APIs). In REPL mode, keybinding registered via `readline`.
- Engine: Vosk 0.3.45+ (offline, ~50MB model, streaming, ~200ms latency). Partial results echoed to terminal as `hearing: show me my...`.
- Upgrade: faster-whisper 1.2.1+ for higher accuracy at GPU cost.
- Custom vocabulary: subcommand/flag names from ToolProfile injected into Vosk grammar.

**Earcons** (short audio cues via `sounddevice`):
- Command received: short click (50ms, 800Hz).
- Success: rising tone (100ms, 600->900Hz).
- Failure: falling tone (100ms, 600->400Hz).
- Working heartbeat: subtle tick every 5 seconds.
- These convey status faster than speech and are sub-100ms.

**STT misrecognition handling:**
- Echo-back before execute for all non-read-only commands: `status: I heard: kubectl get pods. Run this?`
- Low confidence (<0.7): present alternatives using `difflib.get_close_matches` against ToolProfile vocabulary: `status: I heard 'get parts'. Did you mean: get pods, get ports? Say the number or repeat.`
- Spelling fallback: `spell kubectl` -> user spells `k-u-b-e-c-t-l`.
- Custom aliases in config: `[aliases]\nkube-cuddle = "kubectl"`.

**Technology:** `sounddevice` >=0.5.3 (pip, wraps PortAudio), `vosk` >=0.3.45 (pip), `espeak-ng` (system package), optionally `piper-tts` (pip). All optional -- `pip install axcli[audio]`.

### 9. Configuration (`axcli/config.py`)

TOML at `~/.config/axcli/config.toml`. Editable with any text editor (accessible to blind users with screen readers).

```toml
[llm]
url = "http://localhost:11434/v1"
model = "qwen2.5-coder:7b-instruct-q4_K_M"

[tts]
engine = "espeak-ng"   # espeak-ng | piper | say | none
rate = 300             # WPM

[stt]
engine = "vosk"        # vosk | whisper | none
model = "vosk-model-small-en-us-0.15"

[safety]
confirmation = "non-read-only"   # all | non-read-only | high-severity-only
timeout = 120

[display]
compact = false
page_size = 10

[aliases]
kube-cuddle = "kubectl"
```

Environment variable precedence: env > config file > defaults. All settings accessible via `axcli config get <key>` / `axcli config set <key> <value>`.

**Technology:** Python 3.11+ `tomllib` (stdlib) for reading. `tomli_w` (pip) or manual string writing for `config set`.

---

## Data Flow: Complete Single-Turn Interaction

```
User types: "show me my open pull requests"
  |
  v
[Session Shell] Detects: not a raw command (no binary name at start)
  |                      -> route to Intent Engine
  |
  v
[Tool Registry] Loads cached ToolProfile for "gh"
  |              (or probes via cli-acs probe gh --format json)
  |              ToolProfile: {subcommands: [pr, issue, repo, ...],
  |                            flags: {pr: {list: [--state, --author, ...]}},
  |                            has_json: true}
  |
  v
[Intent Engine] Pattern match: "show|list" + "pull requests" -> pr list (no LLM needed)
  |              Or LLM call if pattern match fails:
  |              System prompt + ToolProfile schema + SessionContext
  |              -> {args: ["pr","list","--state","open","--author","@me"],
  |                  explanation: "List your open PRs",
  |                  is_read_only: true, confidence: 0.95}
  |
  v
[Safety Gate] is_read_only=true, "list" in READ_ONLY_VERBS -> ALLOW
  |            Output: status: Running: gh pr list --state open --author @me
  |
  v
[Executor] Popen(["gh","pr","list","--state","open","--author","@me","--json","number,title,updatedAt"],
  |          env={NO_COLOR=1, TERM=dumb}, shell=False)
  |         Captures stdout (JSON array), stderr, exit_code=0
  |
  v
[Formatter] Detects JSON output. Parses into indexed array:
  |          [{number:123, title:"Fix login", updatedAt:"2026-09-10"}, ...]
  |          Stores in SessionContext.turns[-1].structured
  |          Headline: "result: 3 open pull requests."
  |          Summary:  "result: 1. Number 123: Fix login bug, updated yesterday."
  |                    "result: 2. Number 456: Add dark mode, updated 3 days ago."
  |                    "result: 3. Number 789: Refactor auth, updated last week."
  |
  v
[stdout] All labeled text written to terminal (screen reader reads it)
[TTS]    If enabled: speaks headline + summary
```

---

## Example Interactions

### 1. Basic read-only command (no LLM, no confirmation)
```
you: git status
status: Running: git status
result: On branch main. 3 files modified: src/auth.py, src/login.py, tests/test_auth.py. 1 untracked file: notes.txt.
```
Raw command detected (starts with `git`). Bypasses Intent Engine. Allowlist: `status` is read-only. Executor runs directly. Formatter strips ANSI, labels output.

### 2. Natural language with LLM translation
```
you: show me my open pull requests
status: Running: gh pr list --state open --author @me --json number,title,updatedAt
result: 3 open pull requests.
result: 1. Number 123: Fix login bug, updated yesterday.
result: 2. Number 456: Add dark mode, updated 3 days ago.
result: 3. Number 789: Refactor auth, updated last week.
```

### 3. Multi-turn anaphoric reference
```
you: close the second one
status: I will run: gh pr close 456
status: This closes pull request 456 "Add dark mode". Proceed? [y/N]
you: y
result: Pull request 456 closed successfully.
```
Intent Engine resolves "the second one" against `SessionContext.turns[-1].structured[1]` -> PR 456. Safety Gate: `close` is NOT in `READ_ONLY_VERBS` -> confirm.

### 4. Destructive command with high-severity warning
```
you: delete all pods in production
status: I will run: kubectl delete pods --all -n production
warning: HIGH SEVERITY. This will permanently delete all running pods in the production namespace. This action cannot be undone and will cause service downtime. Type 'yes' to proceed or 'n' to cancel.
you: yes
result: Deleted 12 pods in namespace production.
```

### 5. Command failure with error interpretation
```
you: push my changes
status: Running: git push origin main
error: Push failed. Exit code 1. The remote branch has commits you do not have locally.
error: To fix: 1. Run git pull --rebase origin main. 2. Resolve any conflicts. 3. Push again.
error: Type 'pull main' to do step 1 now.
```
Executor captures exit code 1 + stderr. If LLM available, passes to Intent Engine for plain-language explanation. If LLM unavailable, falls back to local error pattern database (top 50 common errors for git, kubectl, gh, docker).

### 6. Large output with progressive disclosure
```
you: git log
status: Running: git log --oneline -50
result: 50 commits. Most recent: "fix HD-5 section detection" by Sachin, today.
result: Type 'more' for next page, 'detail 3' for full commit info, or 'search fix' to filter.
you: detail 3
result: Commit 625c748. Author: Sachin Sampras M. Date: 2026-09-09.
result: Message: fix: detect flush-left subcommands and bare help/version forms
result: Files: 2 changed, 45 insertions, 12 deletions.
```

### 7. Offline mode (no LLM available)
```
you: list my repos
warning: LLM unavailable (Ollama not running). Fuzzy-matching against gh commands.
status: Best match: gh repo list. Run this? [Y/n]
you: y
status: Running: gh repo list
result: 24 repositories. First 10 shown.
result: 1. cli-accessibility-spec, updated 2 hours ago.
result: 2. gitsign, updated yesterday.
```
Intent Engine detects Ollama connection failure. Falls back to `difflib.get_close_matches("list repos", tool_profile.subcommands)`.

### 8. Text-only mode with screen reader active
```
$ AXCLI_SCREEN_READER=1 axcli gh
axcli: Screen reader detected. TTS disabled. All output via terminal.
axcli: Type natural language or raw commands. Type 'help' for options.
you: show my notifications
status: Running: gh api notifications --jq '.[].subject.title'
result: 5 notifications.
result: 1. Review requested: Add dark mode support
result: 2. Mentioned: Fix CI pipeline timeout
result: 3. Assigned: Update deployment docs
result: 4. Review requested: Refactor auth middleware
result: 5. CI failed: Integration test suite
```
NVDA reads all `result:` and `status:` labels from the terminal buffer. User navigates scrollback with NVDA review cursor. No TTS conflicts.

### 9. Tool discovery without LLM
```
you: what can I do with kubectl
result: kubectl is the Kubernetes CLI, version 1.31.0. 42 subcommands in 8 groups.
result: Basic: create, expose, run, set.
result: Deploy: rollout, scale, autoscale.
result: Cluster: certificate, cluster-info, top, cordon, uncordon, drain, taint.
result: Troubleshoot: describe, logs, attach, exec, port-forward, proxy, cp, auth, debug, events.
result: Type 'more about deploy' to explore a group, or 'help rollout' for subcommand detail.
```
Reads directly from cached ToolProfile. No LLM call needed.

### 10. Voice input with misrecognition correction
```
[Ctrl+Space pressed]
hearing: get parts
status: I heard "get parts" but confidence is low (0.4).
status: Did you mean: 1. get pods  2. get ports  3. get parts
status: Say the number or repeat your command.
you: 1
status: Running: kubectl get pods
result: 12 pods. 11 running, 1 CrashLoopBackOff.
```

### 11. Pipe-through mode (raw command, stdout not a TTY)
```
$ axcli kubectl get pods -o name | xargs -I {} axcli kubectl logs {}
```
When stdout is not a TTY, Formatter passes raw output through without labels or formatting. Works transparently with pipes and redirects.

### 12. Conformance gap compensation
```
you: show docker containers
status: Note: docker lacks NO_COLOR support (CV-1 fail). Stripping color automatically.
status: Running: docker ps --format json
result: 4 running containers.
result: 1. nginx, image nginx:latest, up 3 days, ports 80, 443.
result: 2. postgres, image postgres:15, up 3 days, port 5432.
```
ConformanceReport from `cli-acs check docker --format json` indicated CV-1 failure. Formatter auto-strips ANSI and announces the compensation.

---

## Latency Budget

| Step | Target | Method |
|---|---|---|
| Raw command passthrough (total overhead) | <100ms | No LLM, just ANSI strip + label |
| Screen reader detection | <50ms | Cached, re-checked every 30s |
| Tool registry lookup (cached) | <10ms | JSON file read |
| Pattern-match intent (regex, no LLM) | <20ms | Covers `list`, `show`, `get`, `status` |
| LLM intent resolution (local 7B) | <1500ms total | Pre-warmed model, temperature 0 |
| First text feedback to user | <1000ms | Print `status: Running:` before LLM summary completes |
| Safety gate check | <10ms | Regex + set lookup |
| Executor startup | <50ms | Popen, no shell |
| Heartbeat during execution | Every 5s | Background thread |
| Formatter (heuristic, no LLM) | <50ms | Regex + line counting |
| Formatter (LLM summary) | <2000ms | Only for complex non-JSON output |
| TTS first word (espeak-ng) | <100ms | Subprocess, non-blocking |
| Total voice round-trip | <3s | STT (200ms) + intent (1500ms) + exec + format + TTS (100ms) |

---

## Security Model

1. **Never `shell=True`.** Commands always passed as list to `subprocess.Popen`. LLM-generated argument values cannot inject shell metacharacters.
2. **Allowlist, not blocklist.** Known-safe read-only commands skip confirmation. Everything else confirms. The user always sees the exact command before execution.
3. **LLM does not self-classify safety.** The allowlist + regex patterns determine safety. The LLM only generates commands.
4. **Output sanitization before LLM re-injection.** When feeding command output back to LLM for summarization or next-turn context: truncate to 4000 chars, wrap in a structured JSON envelope (not raw text insertion), strip prompt-injection patterns (`ignore previous`, `system:`, `assistant:`).
5. **Secret detection.** Regex patterns for API keys, tokens, passwords in command output. Detected secrets: redacted in TTS output, not persisted to session file, logged as `[REDACTED]`.
6. **No audit log by default** (privacy). Optional `--audit` flag writes append-only JSON-lines to `~/.local/share/axcli/audit.jsonl` with timestamp, command, exit code. Never logs output content.
7. **Credential passthrough.** axcli never stores, logs, or speaks credentials. If a command requires auth (e.g., `docker login`), it is handed through to the terminal directly.

---

## Implementation Phases

### Phase 1: Keyboard wrapper with formatted output (Weeks 1-4)
**Deliverable:** `pip install axcli` -- pure keyboard, no LLM, no audio. Independently useful.

- Session Shell with prefix mode (`axcli kubectl get pods`) and REPL mode (`axcli`).
- Formatter: ANSI stripping, labeled output sections (`result:`, `error:`, `status:`), table-to-sentence conversion, progressive disclosure (headline/summary/detail with `more`/`next`/`prev`/`search`).
- Executor: `subprocess.Popen` with `NO_COLOR=1`, `TERM=dumb`, heartbeat, Ctrl+C forwarding.
- Screen reader detection (Linux/macOS/Windows) with `AXCLI_SCREEN_READER` env override.
- Tool Registry: fallback parser (Python GenericParser port, ~50 lines). Cache at `~/.cache/axcli/tools/`.
- Session Context: in-memory, last 10 turns, structured data extraction for JSON/tabular output.
- Config: `~/.config/axcli/config.toml` with env var overrides.
- Navigation commands: `more`, `next`, `prev`, `detail N`, `search TERM`, `repeat`, `help`.
- Pipe compatibility: raw output when stdout is not a TTY.

**Dependencies:** Python 3.11+ stdlib only. Zero pip dependencies.

### Phase 2: LLM intent engine with safety gate (Weeks 5-8)
**Deliverable:** Natural language commands. `axcli gh` then type "show my PRs".

- Intent Engine with provider-agnostic LLM backend (Ollama default, cloud fallback).
- Pattern-match fast path for common verbs (list/show/get/status) -- no LLM needed.
- Auto-generated JSON tool schemas from ToolProfile.
- Safety Gate: allowlist + high-severity regex + plain-language confirmation.
- Command echo before execution for all non-read-only commands.
- Error interpretation: local error pattern database (top 50 errors for git/kubectl/gh/docker) + optional LLM explanation.
- Few-shot example loading from `~/.config/axcli/examples/{tool}.json`.
- Offline degradation: fuzzy match against ToolProfile when LLM unavailable.
- `cli-acs probe` integration: add `probe` subcommand to Go binary, use it as primary tool discovery path.

**Dependencies:** `httpx` (optional, stdlib `urllib` works). Ollama (system install, optional).

### Phase 3: Voice input/output (Weeks 9-12)
**Deliverable:** `axcli --voice gh` for voice interaction.

- Audio Subsystem: push-to-toggle STT via Vosk, TTS via espeak-ng with background thread + interrupt queue.
- Earcons for status (success/failure/working/received).
- Custom STT vocabulary from ToolProfile subcommand/flag names.
- Echo-back for voice commands before execution.
- Misrecognition correction: alternatives via fuzzy match, spelling fallback.
- Screen reader deference: TTS auto-disabled when screen reader detected.
- Voice navigation: `repeat`, `next`, `previous`, `tell me more`, `stop`, `cancel`, `help`.
- Hybrid input: voice for commands, keyboard for free-text arguments (commit messages, search queries).

**Dependencies:** `sounddevice` >=0.5.3, `vosk` >=0.3.45, `espeak-ng` (system). Install group: `pip install axcli[audio]`.

### Phase 4: Intelligence and context (Weeks 13-16)
**Deliverable:** Multi-turn context, LLM summarization, advanced error help.

- Session Context with anaphoric reference resolution ("the second one", "that error", "its logs").
- LLM-powered output summarization for complex non-JSON output (log files, stack traces).
- Streaming output handling with periodic TTS checkpoints for long-running commands.
- Custom voice aliases (`kube-cuddle` = `kubectl`).
- Per-command verbosity memory (headline-only for frequently-run commands).
- `axcli setup` command: checks all prerequisites (Ollama, models, mic, TTS), reports in labeled text.
- `axcli tutorial` command: interactive guided walkthrough.
- `axcli whatsnew` command: speaks changes after updates.
- faster-whisper as high-accuracy STT option. Piper as neural TTS option.
- Conformance gap reporting: run `cli-acs check`, announce compensations.

**Dependencies:** Optional: `faster-whisper`, `piper-tts`.

### Phase 5: Hardening and cross-platform (Weeks 17-20)
**Deliverable:** Production-ready, cross-platform.

- Windows support: `pywinpty` for TTY commands, SAPI TTS fallback, process detection for NVDA/JAWS.
- macOS support: `say` TTS, CoreAudio via sounddevice, VoiceOver detection.
- Interactive command detection and passthrough (git commit, ssh, sudo).
- Secret detection and redaction in TTS output and session persistence.
- Optional audit log (`--audit`).
- Compact mode for braille displays (`AXCLI_COMPACT=1`).
- Shell integration: bash/zsh completion script, shell function wrapper for transparent `axcli` prefix.
- openWakeWord as optional hands-free mode.
- Error recovery: graceful handling of axcli's own crashes (catch exceptions, print screen-reader-friendly error, clean up audio devices).

**Dependencies:** Optional: `pywinpty` (Windows), `openWakeWord`.

### Phase 6: Ecosystem and community (Weeks 21-24)
**Deliverable:** Published package with onboarding.

- Zero-config first run: auto-detect capabilities, announce what is available.
- `axcli cheatsheet` command: top 10 commands in 30 seconds of speech.
- Contextual help: `what can I say` lists available commands for current context.
- Local error pattern database expanded to 200+ errors across 10 tools.
- Feed spec proposals back to CLI-ACS: `--screen-reader` flag criterion, labeled output sections criterion, streaming output criterion, voice interface metadata format.
- Package: PyPI (`pip install axcli`, `pip install axcli[audio]`, `pip install axcli[full]`).
- Platform packages: Homebrew formula, Fedora COPR, AUR PKGBUILD.
- Documentation: man page, `--help`, plain-text README (screen-reader-optimized: no ASCII art, no wide tables).
- Accessibility testing with real screen reader users (NVDA, VoiceOver, Orca).

---

## File Structure

```
axcli/
  __init__.py
  __main__.py          # entry point: python -m axcli
  shell.py             # Session Shell (REPL + prefix mode)
  registry.py          # Tool Registry (probe + cache)
  intent.py            # Intent Engine (LLM + pattern match)
  safety.py            # Safety Gate (allowlist + regex + confirm)
  executor.py          # Command Executor (Popen, heartbeat)
  formatter.py         # Output Formatter (ANSI strip, labels, tiers)
  session.py           # Session Context (turns, structured data)
  audio.py             # Audio Subsystem (TTS + STT, optional)
  config.py            # Configuration (TOML + env vars)
  errors.py            # Local error pattern database
  screen_reader.py     # Screen reader detection
  cli.py               # CLI argument parsing (click or argparse)
pyproject.toml
README.txt             # plain text, not markdown (screen-reader-friendly)
```

---

## Dependency Summary

**Phase 1 (core):** Python 3.11+ stdlib only. Zero pip dependencies.

**Phase 2 (LLM):** `httpx` (optional, stdlib `urllib` suffices). External: Ollama (optional).

**Phase 3 (audio):** `sounddevice` >=0.5.3, `vosk` >=0.3.45. System: `espeak-ng`, PortAudio.

**Phase 4 (upgrades):** `faster-whisper` >=1.2.1 (optional), `piper-tts` (optional).

**External tools:** `cli-acs` Go binary (optional, enhances tool discovery). Ollama (optional, enhances NL translation).

Install groups in `pyproject.toml`:
```toml
[project.optional-dependencies]
audio = ["sounddevice>=0.5.3", "vosk>=0.3.45"]
llm = ["httpx>=0.27"]
full = ["sounddevice>=0.5.3", "vosk>=0.3.45", "httpx>=0.27", "piper-tts>=1.6.0"]
```

---

## Prerequisites in cli-acs (Go changes needed)

Before Phase 2, the following changes to the Go codebase at `/home/sacm/Documents/Study/my_projects/cli-accessibility-spec/` are required:

1. **Add `cli-acs probe <binary> --format json` command** (~40 lines in `cmd/cli-acs/probe.go`): Call `probe.Probe()`, marshal `ProbeResult` to JSON, write to stdout. No `os.Exit()` -- return normally so callers can parse output.

2. **Populate `Subcommand.Flags`** (in `internal/probe/discovery.go`): After discovering top-level subcommands, run `{binary} {subcommand} --help` on each (up to `maxSubcmds`) and parse flags into `Subcommand.Flags` using the same parser cascade.

3. **Fix `TakesValue` detection** (in `internal/probe/parsers.go`, `parseFlag` function): Detect value placeholders in flag descriptions (`<file>`, `FILE`, `<value>`, `=VALUE`, `string`, `int`, `path`) and set `TakesValue: true`. Heuristic: if the token after the flag name is an UPPER_CASE word or angle-bracketed placeholder, the flag takes a value.