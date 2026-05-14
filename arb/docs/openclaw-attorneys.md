# OpenClaw Attorneys

This note describes how to run `arb` with OpenClaw agents as the plaintiff and defendant attorneys. It is intended to make the local OpenClaw-attorney workflow reproducible from a fresh checkout.

## Architecture

`aar case` talks to attorneys through ACP. The normal local path starts a Pi ACP wrapper. The OpenClaw path replaces that attorney process with an OpenClaw-backed ACP server.

There are two pieces:

1. `.bin/aar-openclaw-attorney` is a stdio ACP adapter. It receives AAR `session/prompt` requests, asks AAR for the visible case through `_aar/get_case`, loads text-readable case files through `_aar/list_case_files` and `_aar/read_case_text_file`, asks an OpenClaw agent for one filing JSON, and submits that filing through `_aar/submit_decision`.
2. `tools/openclaw-acp-tcp-bridge.js` exposes that stdio adapter as a TCP ACP endpoint. `aar case` connects to this endpoint with `--plaintiff-acp-endpoint` and `--defendant-acp-endpoint`.

The bridge starts a new `.bin/aar-openclaw-attorney` process for each ACP connection. One bridge process can serve both sides of a run.

## Capability boundary

AAR does not select the OpenClaw model. A role using `--*-acp-endpoint` cannot also use `--*-attorney-model`. That flag belongs to the local Pi/xproxy attorney path.

For OpenClaw attorneys:

- model selection belongs to the OpenClaw agent configuration or the OpenClaw runtime invoked by `openclaw agent`;
- tool availability belongs to that OpenClaw environment;
- `run.json` records the ACP endpoint, not an AAR-selected attorney model;
- the bridge deliberately removes `AAR_OPENCLAW_AGENT_MODEL` from the adapter environment.

This matters for reproducibility. A closed-record run is reproduced by using an OpenClaw attorney agent that does not search or by instructing it to stay within the provided record. An open-record run is reproduced by using an OpenClaw attorney agent with search/browser/fetch tools available and by giving explicit open-record instructions through `AAR_OPENCLAW_AGENT_EXTRA_PROMPT`.

## Prerequisites

From `arb/`:

```bash
make build
node --check tools/openclaw-acp-tcp-bridge.js
```

The run also needs:

- `.bin/aar` and `.bin/aarengine`, built by `make build`;
- `.bin/aar-openclaw-attorney`, built by `make build`;
- an OpenClaw CLI on `PATH`, or `AAR_OPENCLAW_CLI` set to the CLI path;
- a dedicated OpenClaw lawyer agent, normally `aar-lawyer`;
- provider credentials required by the OpenClaw lawyer agent;
- `OPENROUTER_API_KEY` for the council models used by AAR.

If the shell environment does not already contain the council key, source it before the run:

```bash
source "$HOME/keys.txt"
```

Do not use a personal default OpenClaw agent. Set `AAR_OPENCLAW_AGENT_ID` explicitly so arbitration work runs in the dedicated lawyer context.

## Closed-record run

This run lets OpenClaw attorneys argue from the packet that AAR provides. It is the reproducible baseline when the case directory is self-contained.

Terminal 1, start the bridge:

```bash
cd arb
AAR_OPENCLAW_AGENT_ID=aar-lawyer \
AAR_OPENCLAW_ATTORNEY_TIMEOUT_SECONDS=900 \
tools/openclaw-acp-tcp-bridge.js --host 127.0.0.1 --port 19701
```

Terminal 2, prepare and run the case:

```bash
cd arb
source "$HOME/keys.txt"

case_dir=examples/clavicular-pregnancy-credible-announcement-condition-simple
out_dir=out/clavicular-openclaw-both-closed-$(date +%Y%m%d-%H%M%S)

.bin/aar complain \
  --situation "$case_dir/situation.md" \
  --out "$case_dir/complaint.md"

.bin/aar case \
  --complaint "$case_dir/complaint.md" \
  --out-dir "$out_dir" \
  --plaintiff-acp-endpoint tcp://127.0.0.1:19701 \
  --defendant-acp-endpoint tcp://127.0.0.1:19701 \
  --council-size 3 \
  --acp-timeout-seconds 900 \
  --invalid-attempt-limit 5
```

Expected output is a single JSON object on stdout, for example:

```json
{"status":"ok","result":"not_demonstrated","votes_for":0,"votes_against":3,"run_id":"run-...","out_dir":"out/clavicular-openclaw-both-closed-..."}
```

## Open-record run

An open-record run uses the same ACP bridge, but gives the OpenClaw attorneys an extra instruction to investigate public sources. The OpenClaw agent must have the relevant tools available. AAR itself does not enable those tools.

Terminal 1, start the bridge with an open-record instruction:

```bash
cd arb
export AAR_OPENCLAW_AGENT_ID=aar-lawyer
export AAR_OPENCLAW_ATTORNEY_TIMEOUT_SECONDS=1200
export AAR_OPENCLAW_AGENT_EXTRA_PROMPT='This is an open-record arbitration. Use available public search, web fetch, browser, transcript, or equivalent tools when they can materially improve the filing. Prefer primary sources over commentary. Preserve URLs, direct excerpts, and uncertainty. If external material matters, include it in technical_reports with enough provenance for later packet backfill. Do not cite external material in offered_files unless AAR exposed it as a case file.'

tools/openclaw-acp-tcp-bridge.js --host 127.0.0.1 --port 19702
```

Terminal 2, run the same case against that endpoint:

```bash
cd arb
source "$HOME/keys.txt"

case_dir=examples/clavicular-pregnancy-credible-announcement-condition-simple
ts=$(date +%Y%m%d-%H%M%S)
out_dir=out/clavicular-openclaw-both-open-$ts
batch_dir=out/_batch-clavicular-openclaw-both-open-$ts
mkdir -p "$batch_dir/logs"

.bin/aar complain \
  --situation "$case_dir/situation.md" \
  --out "$case_dir/complaint.md"

.bin/aar case \
  --complaint "$case_dir/complaint.md" \
  --out-dir "$out_dir" \
  --plaintiff-acp-endpoint tcp://127.0.0.1:19702 \
  --defendant-acp-endpoint tcp://127.0.0.1:19702 \
  --council-size 3 \
  --acp-timeout-seconds 1200 \
  --invalid-attempt-limit 5 \
  2>&1 | tee "$batch_dir/logs/run.log"

cp "$batch_dir/logs/run.log" "$out_dir/run.log"
```

Use `bash` for scripts that inspect pipeline status. Do not use zsh-only or bash-only status variables interchangeably after piping through `tee`.

## What to inspect after the run

A completed run should contain:

```text
run.json
state.json
digest.md
council.json
events.ndjson
transcript.md
run.log
```

Check these before reporting the result:

```bash
jq '.status? // empty' "$out_dir/run.json" 2>/dev/null || true
jq '.result? // empty' "$out_dir/run.json" 2>/dev/null || true
grep -n "Resolution:" "$out_dir/digest.md"
grep -n "technical report\|technical_reports\|http" "$out_dir/digest.md" | head -40
```

For an open-record run, the important question is not only the vote split. Inspect whether the attorneys obtained the decisive public evidence, preserved enough provenance, and distinguished primary evidence from secondary reporting. If the attorneys rely on external evidence that is not in the case packet, backfill the case directory before using later closed-record runs as reproducibility evidence.

## Clavicular runs used during development

The Clavicular case used this workflow:

```text
examples/clavicular-pregnancy-credible-announcement-condition-simple
```

Closed-record OpenClaw attorneys:

```text
out/clavicular-openclaw-both-closed-20260514-160055
```

Open-record OpenClaw attorneys:

```text
out/clavicular-openclaw-both-open-20260514-183235
out/_batch-clavicular-openclaw-both-open-20260514-183235/summary.md
```

The open-record run completed successfully and produced `not_demonstrated`, 0 demonstrated / 3 not_demonstrated. The attorneys searched and filed technical reports. They found third-party public material, including KickChamp/X, a YouTube Short, Polymarket/API material, Forbes, Domahhhh, and fact-check sources. They did not retrieve the original Clavicular VOD, a clean Clavicular transcript, a representative statement, or a separate Polymarket clarification. That result is a useful evidence-discovery audit: the OpenClaw attorneys found the public rumor trail but did not close the primary-source gap.

## Current limitations

- The adapter asks OpenClaw for one strict JSON filing per AAR opportunity. It does not stream intermediate reasoning back to AAR.
- The adapter can submit filings through `_aar/submit_decision` and can read text case files. It does not create durable new case files from external evidence.
- External findings currently enter through `technical_reports`, which are attorney work product. They do not become preserved evidence with new `file_id` values.
- A first-class evidence-submission path would need AAR methods such as `_aar/submit_evidence` or `_aar/submit_case_file`, validation, storage, hashing, and a returned `file_id` that later filings can cite in `offered_files`.
- AAR endpoint metadata deliberately does not claim OpenClaw search capability. The operator must configure and document the OpenClaw agent environment used for open-record runs.
