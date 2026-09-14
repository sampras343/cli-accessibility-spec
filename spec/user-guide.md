# CLI Accessibility User Guide

A practical guide for blind, low-vision, and motor-impaired developers who use command-line tools daily. This is not a specification — it tells you how to set up your environment, navigate output, and work effectively with tools like oc, kubectl, gh, git, and docker using assistive technology.

---

## Table of Contents

1. [First Steps: Environment Variables](#1-first-steps-environment-variables)
2. [Windows: NVDA + Terminal](#2-windows-nvda--terminal)
3. [macOS: VoiceOver + Terminal + tdsr](#3-macos-voiceover--terminal--tdsr)
4. [Linux: Orca + GNOME Terminal](#4-linux-orca--gnome-terminal)
5. [Working With CLI Output](#5-working-with-cli-output)
6. [Reading Tables](#6-reading-tables)
7. [Handling Errors](#7-handling-errors)
8. [Using axcli (Accessible Wrapper)](#8-using-axcli-accessible-wrapper)
9. [Remote Work Over SSH](#9-remote-work-over-ssh)
10. [Shell Configuration](#10-shell-configuration)
11. [Common Workflows](#11-common-workflows)
12. [Tools That Help](#12-tools-that-help)
13. [When a Tool Doesn't Work](#13-when-a-tool-doesnt-work)

---

## 1. First Steps: Environment Variables

Before anything else, add these to your shell profile (`~/.bashrc`, `~/.zshrc`, or PowerShell `$PROFILE`):

```bash
export NO_COLOR=1
export PAGER=cat
export EDITOR=nano
export VISUAL=nano
```

What these do:

- **`NO_COLOR=1`** tells CLI tools to suppress color escape codes. Without this, your screen reader may announce raw escape characters like "escape left bracket three one m" instead of the actual text. Over 180 tools respect this variable.

- **`PAGER=cat`** prevents tools from opening `less` or `more` for long output. Pagers use an alternate screen buffer that screen readers cannot scroll back through — when you quit the pager, the output vanishes. Setting `PAGER=cat` keeps all output in your terminal's main scrollback where you can review it.

- **`EDITOR=nano`** and **`VISUAL=nano`** set your default editor to nano, which is simple and usable with a screen reader. If you use Emacs with Emacspeak, set these to `emacs` instead.

After adding these, reload your shell:

```bash
source ~/.bashrc
```

---

## 2. Windows: NVDA + Terminal

### Recommended Setup

- **Screen reader:** NVDA (free, open source) — download from nvaccess.org
- **Terminal:** Windows Terminal (from Microsoft Store) or CMD
- **NVDA add-on:** Install the Console Toolkit add-on from addons.nvda-project.org — it adds real-time speech output and an accessible command editing window

### NVDA Settings for Terminal

Open NVDA Settings (NVDA+N → Preferences → Settings):

- **Speech:** Turn off "Automatic language switching" — it causes the voice to change mid-output when it encounters text it thinks is another language
- Press **NVDA+6** to disable "Caret moves review cursor" — this lets you scroll back through output without losing your typing position

### Navigating Terminal Output

Use the review cursor to read output line by line:

| What you want to do | Numpad (Desktop) | Laptop layout |
|---|---|---|
| Read current line | Numpad 8 | NVDA+Shift+. |
| Previous line | Numpad 7 | NVDA+Up |
| Next line | Numpad 9 | NVDA+Down |
| Previous word | Numpad 4 | NVDA+Ctrl+Left |
| Next word | Numpad 6 | NVDA+Ctrl+Right |
| Previous character | Numpad 1 | NVDA+Left |
| Next character | Numpad 3 | NVDA+Right |
| Read everything from cursor down | NVDA+Down Arrow | NVDA+A |
| Scroll terminal back | Ctrl+Up | Ctrl+Up |

**Console Toolkit shortcuts:**

- **NVDA+E** opens an accessible edit window for the current command — useful when the command is long and you need to review or edit it
- **Shift+Numpad 7** reads the first visible line (press twice for first line in the scrollback buffer)

### Copy Text from Terminal

1. Press **NVDA+F9** to mark the start position
2. Move to the end of the text you want
3. Press **NVDA+F10** twice to copy to clipboard

---

## 3. macOS: VoiceOver + Terminal + tdsr

### Recommended Setup

- **Screen reader:** VoiceOver (built in, activate with Cmd+F5)
- **Terminal:** Terminal.app (built in) — more reliable with VoiceOver than iTerm2
- **For heavy CLI work:** Install tdsr (Tyler's Damned Screen Reader), a Python-based console screen reader that speaks output directly

VoiceOver's terminal support has known limitations. The Lighthouse for the Blind described it as "inadequate for all but the most minimal use cases." For serious CLI work, tdsr is recommended.

### Install tdsr

```bash
brew install uv
uv tool install tdsr
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zprofile
```

To auto-launch tdsr when you open Terminal: Terminal → Settings → Profiles → Shell → set "Run command" to `~/.local/bin/tdsr` and check "Run inside shell."

### Prevent VoiceOver and tdsr from Talking Over Each Other

Open VoiceOver Utility (VO+F8) → Activities → add a "Terminal" activity → set Voices → enable "Mute Speech." This silences VoiceOver in Terminal while tdsr handles the speech.

### Terminal.app Setting

Settings → Profiles → Keyboard → enable "Use Option as meta key" — needed for tdsr's navigation shortcuts.

### tdsr Navigation

| What you want to do | Shortcut |
|---|---|
| Previous line | Option+U |
| Current line | Option+I |
| Next line | Option+O |
| Previous word | Option+J |
| Current word | Option+K |
| Next word | Option+L |
| Previous character | Option+M |
| Current character | Option+, |
| Next character | Option+. |

### VoiceOver Terminal Shortcuts (without tdsr)

- **VO+Shift+Down** to interact with the Terminal area
- **VO+Left/Right** to navigate by character or word
- **VO+Up/Down** to navigate by line
- **VO+A** to read all from current position

---

## 4. Linux: Orca + GNOME Terminal

### Recommended Setup

- **Screen reader:** Orca (included in GNOME, launch with Super+Alt+S)
- **Terminal:** GNOME Terminal, MATE Terminal, or XFCE4 Terminal — these are GTK-based and expose accessibility APIs to Orca

**Do not use:** xterm, urxvt, st, Konsole, or Alacritty — they are completely inaccessible to Orca.

**For split panes:** Use Terminator (not tmux splits) — Terminator uses separate GTK widgets per pane, so Orca reads each pane independently. tmux and screen splits cause Orca to read text from both panes on the same line.

### Orca Flat Review for Terminal Output

| What you want to do | Numpad (Desktop) | Laptop (CapsLock as modifier) |
|---|---|---|
| Previous line | Numpad 7 | CapsLock+U |
| Current line | Numpad 8 | CapsLock+I |
| Next line | Numpad 9 | CapsLock+O |
| Previous word | Numpad 4 | CapsLock+J |
| Current word | Numpad 5 | CapsLock+K |
| Next word | Numpad 6 | CapsLock+L |
| Previous character | Numpad 1 | CapsLock+M |
| Current character | Numpad 2 | CapsLock+, |
| Next character | Numpad 3 | CapsLock+. |
| Read all | Orca+; | Orca+; |
| Learn mode (discover shortcuts) | Orca+H | Orca+H |

### Console Access Without a Desktop (TTY)

If you work on a Linux console without a graphical desktop (servers, SSH, recovery):

- **Fenrir** is a user-space Python screen reader for the Linux console. Install via pip or your distribution's package manager. It uses speech-dispatcher or espeak directly, supports braille via BRLTTY, and is scriptable.
- **Speakup** is a kernel module included in most distributions. Combined with ESpeakup (which bridges Speakup to espeak), it provides speech from boot — before any user-space tools load.

### Emacspeak

If you use Emacs, Emacspeak transforms it into a complete audio desktop. It wraps shell commands with audio formatting — different voice pitches for different content types, earcons for events. Run shell commands via `M-x shell` or `M-x eshell`. Available since 1995.

---

## 5. Working With CLI Output

### The Core Problem

CLI output is unstructured text. A web page has headings, links, and landmarks that your screen reader can navigate. Terminal output has none of that. Your screen reader reads it line by line, character by character.

### Strategies

**Redirect to a file, read in your editor:**

```bash
oc get pods -n myns > /tmp/pods.txt
nano /tmp/pods.txt
```

This gives you full editor navigation — search, jump to line, copy sections. It's the most reliable way to read long or complex output.

**Use --json and jq for structured data:**

```bash
# Instead of reading a table:
kubectl get pods -o json | jq '.items[].metadata.name'

# Get specific fields:
oc get pods -o json | jq '.items[] | {name: .metadata.name, status: .status.phase}'
```

JSON + jq gives you exactly the fields you need, one per line, with no table alignment to decode.

**Use grep to find what you need:**

```bash
oc get pods -n myns | grep -i error
oc get pods -n myns | grep -i pending
git status | grep modified
```

Instead of reading 200 lines of output, grep for the keyword you care about.

**Use wc to count before reading:**

```bash
oc get pods -n myns | wc -l
```

Knowing "there are 47 lines of output" before reading helps you decide whether to read it all, grep it, or redirect to a file.

---

## 6. Reading Tables

Tables are the single hardest format for screen readers. When oc or kubectl outputs:

```
NAME          READY   STATUS    RESTARTS   AGE
nginx-abc     1/1     Running   0          3d
postgres-xy   0/1     Pending   2          1h
```

Your screen reader reads: "NAME READY STATUS RESTARTS AGE nginx-abc one slash one Running zero three d" — a stream of words with no way to tell which value belongs to which column.

### Option 1: axcli (recommended)

```bash
axcli --domain screen-reader oc get pods -n myns
```

Converts tables to labeled sentences:

```
result: 1. NAME: nginx-abc. READY: 1/1. STATUS: Running. RESTARTS: 0. AGE: 3d.
result: 2. NAME: postgres-xy. READY: 0/1. STATUS: Pending. RESTARTS: 2. AGE: 1h.
```

Every value is paired with its column header.

### Option 2: JSON output

```bash
oc get pods -o json | jq '.items[] | "\(.metadata.name): \(.status.phase)"'
```

Output:

```
"nginx-abc: Running"
"postgres-xy: Pending"
```

### Option 3: Custom columns

```bash
oc get pods -o custom-columns=NAME:.metadata.name,STATUS:.status.phase
```

Produces a simpler two-column table that's easier to follow.

### Option 4: Name-only output

```bash
oc get pods -o name
```

One pod per line, no columns:

```
pod/nginx-abc
pod/postgres-xy
```

---

## 7. Handling Errors

### Errors Go to stderr

Most CLI tools send errors to stderr, not stdout. If you're piping output, you might miss errors entirely:

```bash
# Error goes to stderr — you might not hear it if piping stdout:
oc get pods -n nonexistent | grep Running
# The grep runs on empty stdout. The error went to stderr and may have scrolled past.
```

**Fix:** Redirect stderr to stdout so everything goes to one stream:

```bash
oc get pods -n nonexistent 2>&1 | less
```

Or use axcli, which captures both and labels them:

```bash
axcli oc get pods -n nonexistent
# error: Error from server (Forbidden): pods is forbidden...
```

### Exit Codes

After any command, check the exit code:

```bash
echo $?
```

0 means success. Anything else means failure. This works even when you can't read the error message — if `$?` is not 0, something went wrong.

### Common Error Patterns

When you hear an error, look for these keywords:

- **"forbidden"** or **"cannot"** → you don't have permission. Try `oc whoami` to check your user, `oc auth can-i` to check what you can do.
- **"not found"** → the resource or namespace doesn't exist. Check your spelling. Try `oc get namespaces` to see available namespaces.
- **"connection refused"** → the server is down or unreachable. Check your cluster URL with `oc config view`.
- **"timeout"** → the operation took too long. Try again, or check network connectivity.

---

## 8. Using axcli (Accessible Wrapper)

axcli wraps any CLI tool and makes its output accessible. Install it:

```bash
pip install axcli
```

Use it by prefixing any command:

```bash
axcli oc get pods -n myns
axcli gh pr list
axcli kubectl get deployments
```

axcli detects whether you're using a screen reader and adapts:

- **Screen reader detected** (NO_COLOR set, TERM=dumb, or Orca/NVDA running): tables become labeled sentences, Unicode symbols become text, errors get `error:` prefixes
- **No screen reader**: status words get color (Running=green, Error=red, Pending=yellow)

Force screen reader mode explicitly:

```bash
axcli --domain screen-reader oc get pods -n myns
```

See the [axcli documentation](https://github.com/cli-accessibility/axcli) for full details.

---

## 9. Remote Work Over SSH

Screen readers work over SSH because your screen reader reads the local terminal, which displays the remote output. The remote machine doesn't need a screen reader installed.

### Setup

Add the same environment variables to the remote machine's `~/.bashrc`:

```bash
export NO_COLOR=1
export PAGER=cat
export EDITOR=nano
```

### tdsr Over SSH

tdsr can run on the remote machine itself — install it there and use it as your shell. This is more reliable than local screen reader + remote terminal because tdsr sees the character stream directly instead of diffing the screen buffer:

```bash
# On the remote machine:
pip install tdsr
# Set as default shell or run manually:
tdsr
```

### tmux for Persistence

Use tmux on the remote machine to keep your sessions alive when you disconnect:

```bash
# Start a named session:
tmux new -s work

# Detach: Ctrl+B, then D
# Reattach later:
tmux attach -t work
```

**Important:** Do not use tmux split panes — screen readers read text from both panes on the same line, making both unreadable. Use separate tmux windows instead:

```bash
# Create a new window:
Ctrl+B, then C

# Switch windows:
Ctrl+B, then N (next) or P (previous)
```

---

## 10. Shell Configuration

### Bash vs Zsh

Both work fine with screen readers. Recommendations:

- **Bash** for simplicity — it's the default on most systems, well-tested with assistive technology
- **Zsh** if you want better tab completion — `zsh-autosuggestions` can help, but test it with your screen reader first (some users report issues with async prompt updates)
- **Fish** has nice defaults but its async features can confuse some screen readers

### Useful Shell Settings

Add to your `~/.bashrc`:

```bash
# Announce the exit code after every command
PROMPT_COMMAND='echo "exit: $?"'

# Simpler prompt — avoid Unicode symbols, git status indicators, etc.
PS1='\u@\h:\w\$ '

# Bell on command completion (some screen readers announce this)
set bell-style visible
```

### Tab Completion

Tab completion is your best friend. Instead of remembering flag names:

```bash
oc get <Tab><Tab>     # lists all resource types
oc get pods --<Tab>   # lists all flags
gh pr <Tab><Tab>      # lists subcommands
```

Install shell completions for your tools:

```bash
# oc
oc completion bash > /etc/bash_completion.d/oc

# kubectl
kubectl completion bash > /etc/bash_completion.d/kubectl

# gh
gh completion -s bash > /etc/bash_completion.d/gh
```

---

## 11. Common Workflows

### "I need to check the status of my pods"

```bash
# Quick count:
oc get pods -n myns | wc -l

# Find problems:
oc get pods -n myns | grep -v Running | grep -v Completed

# With axcli (labeled output):
axcli --domain screen-reader oc get pods -n myns

# Just the names and status:
oc get pods -n myns -o custom-columns=NAME:.metadata.name,STATUS:.status.phase
```

### "I need to read a pod's logs"

```bash
# Last 50 lines:
oc logs my-pod -n myns --tail=50

# Search for errors:
oc logs my-pod -n myns | grep -i error

# Save to file for reading in editor:
oc logs my-pod -n myns > /tmp/logs.txt
nano /tmp/logs.txt
```

### "I need to check my pull requests"

```bash
# List PRs (with axcli):
axcli --domain screen-reader gh pr list

# Just titles:
gh pr list --json number,title --jq '.[] | "#\(.number): \(.title)"'
```

### "I need to deploy something"

```bash
# Preview first:
oc apply -f manifest.yaml --dry-run=client

# Then apply:
oc apply -f manifest.yaml

# Check the exit code:
echo $?
```

### "Something is wrong and I don't know what"

```bash
# Step 1: What namespace am I in?
oc project

# Step 2: What's running?
oc get pods -n myns | grep -v Completed | grep -v Running

# Step 3: What happened?
oc describe pod <problem-pod> -n myns > /tmp/describe.txt
nano /tmp/describe.txt
# Search for "Events:" section — it tells you what went wrong

# Step 4: Check logs:
oc logs <problem-pod> -n myns --tail=100 | grep -i "error\|fatal\|panic"
```

---

## 12. Tools That Help

| Tool | What it does | Install |
|---|---|---|
| **axcli** | Wraps any CLI, converts tables to labeled text, adds color semantics | `pip install axcli` |
| **jq** | Extracts specific fields from JSON output | `sudo dnf install jq` / `brew install jq` |
| **tdsr** | Console screen reader for macOS/Linux, better than VoiceOver for terminals | `uv tool install tdsr` |
| **Console Toolkit** | NVDA add-on for better terminal support | NVDA Add-on Store |
| **Emacspeak** | Complete audio desktop via Emacs | `sudo dnf install emacspeak` |
| **Fenrir** | Console screen reader for Linux TTY | `pip install fenrir-screenreader` |
| **nano** | Simple accessible text editor | Usually pre-installed |

---

## 13. When a Tool Doesn't Work

Not all CLI tools are accessible. When you encounter one that doesn't work:

### Symptoms

- Screen reader reads escape characters ("escape left bracket three one m")
- Output is garbled or has no structure
- The tool hangs and you can't tell if it's working or stuck
- Interactive prompts don't announce options
- Spinners or progress bars produce repeated garbage speech

### Workarounds

1. **Set `NO_COLOR=1`** — this fixes most color-related garbling
2. **Try `TERM=dumb`** — suppresses all escape sequences, not just color
3. **Use axcli** — wraps the tool and cleans its output
4. **Redirect to file** — `tool > output.txt 2>&1` then read in your editor
5. **Use `--json` if available** — structured output is always more accessible
6. **Pipe through `cat -v`** — makes escape sequences visible as text so you can report the issue
7. **Press Ctrl+C** — if the tool seems stuck, this cancels it

### Reporting Accessibility Issues

When you find a CLI tool that doesn't work with your screen reader, file a bug. Include:

- Your operating system and version
- Your screen reader and version
- Your terminal emulator
- The exact command you ran
- What you expected to hear vs what you actually heard
- The output of `tool --version`
- The output of `NO_COLOR=1 TERM=dumb tool --help | head -20`

Many tool maintainers are not aware of accessibility issues. A clear, reproducible bug report is the most effective way to improve the ecosystem.

### Run the CLI-ACS Conformance Suite

If you want to evaluate a tool's accessibility formally:

```bash
go install github.com/sampras343/cli-accessibility-spec/cmd/cli-acs@latest
cli-acs check /path/to/tool
```

This produces a conformance report showing which accessibility criteria the tool meets and which it fails. Share this report with the tool's maintainers.

---

## Sources

This guide draws on:

- [Sampath, Merrick, Macvean: "Accessibility of Command Line Interfaces" (CHI 2021)](https://dl.acm.org/doi/fullHtml/10.1145/3411764.3445544) — academic study of 12 developers using CLIs with screen readers
- [Institute of Accessible Technology: CLI and IDE with NVDA](https://www.instituteofaccesstech.com/using-command-lines-and-ides-with-nvda/)
- [NVDA Console Toolkit add-on](https://addons.nvda-project.org/addons/consoleToolkit.en.html)
- [Marco Salsiccia: Beginner Guide to macOS Terminal for VoiceOver](https://marconius.com/terminal/)
- [tdsr — console screen reader for macOS/Linux](https://github.com/tspivey/tdsr)
- [Orca Screen Reader documentation](https://wiki.gnome.org/Projects/Orca)
- [Fenrir — console screen reader for Linux](https://github.com/chrys87/fenrir)
- [Emacspeak — the complete audio desktop](https://emacspeak.sourceforge.net/)
- [Seirdy: Best practices for inclusive CLIs](https://seirdy.one/posts/2022/06/10/cli-best-practices/)
- [GitHub Blog: Building a more accessible GitHub CLI](https://github.blog/engineering/user-experience/building-a-more-accessible-github-cli/)
- [NO_COLOR standard](https://no-color.org/)
