
## 2026-04-29 adjudication readiness check

Observation: the operator indicated more adjudication runs may follow. I checked `<repo-root>`, not as a maintainer cleanup task, but as an operator preparing to run the arbitration harness.

Findings:
- Repository branch is `main` tracking `origin/main` with two pre-existing modified example artifacts: `arb/examples/ex1/confession.sig.b64` and `arb/examples/ex1/samantha_public.pem`.
- The expected user-facing runner in this checkout is `arb/arbitrate.sh`; I did not find a file named `adjudicate.sh` in the workspace search depth used.
- `make test` in `arb/` passed.
- `make build` initially failed because `lake` was absent from the default shell `PATH`.
- `lake` and `lean` exist under `$HOME/.elan/bin`; rebuilding with `PATH="$HOME/.elan/bin:$PATH"` succeeded.
- `./arbitrate.sh` already exports that PATH and sources `$HOME/keys.txt` if present.
- `$HOME/keys.txt` is present and defines `OPENAI_API_KEY` and `OPENROUTER_API_KEY`; I did not read or record values.
- Podman is running, and the `agentcourt-pi-sandbox` images are present.

Decision: Do not start an arbitration run without explicit instruction, because it removes the selected output directory and invokes external model services. The repo is ready for a run subject to the operator choosing the input/output/model parameters.

## 2026-04-29 ex2-ex6 three-run plan

Makefile observations:
- `ex2` and `ex3` run `build`, remove `out/exN-demo`, regenerate `examples/exN/complaint.md` from `situation.md`, then run `.bin/aar case` with `openai://gpt-5`.
- `ex4`, `ex5`, and `ex6` run `build`, remove `out/exN-demo`, then run `.bin/aar case` against the checked-in complaint with `openai://gpt-5?tools=search`.
- The Makefile targets are unsuitable for repeated runs without modification because each target reuses and removes the same output directory.

Operational plan:
- Build once with `$HOME/.elan/bin` on `PATH`.
- Source `$HOME/keys.txt` before the batch; do not record key values.
- Generate complaints once for `ex2` and `ex3` before the repeated runs.
- Run all 15 cases serially to avoid xproxy port conflicts, Podman contention, and provider-rate problems.
- Use stable output directories of the form `out/exN-rK` and matching `--run-id exN-rK-YYYYMMDD-HHMMSS`.
- Use `openai://gpt-5` for ex2/ex3 and `openai://gpt-5?tools=search` for ex4/ex5/ex6, matching the Makefile.
- After each run, save stdout/stderr under `out/_batch-ex2-ex6-*/logs/`, verify the expected packet files, and append a manifest row with status, result, votes, run ID, and output directory.

Do not start the batch without explicit approval because it invokes external model services and writes/removes run-output directories.

## 2026-04-29 arbitrate.sh generalization

Observation: `arbitrate.sh` assumed the ex1 input shape by requiring `sign.sh` and `situation.md`. The Makefile examples require a more general runner: ex2/ex3 regenerate complaints from situations, while ex4/ex5/ex6 run from checked-in complaints.

Decision: Generalize `arbitrate.sh` without changing its external argument shape. The script now runs `sign.sh` only when it exists and is executable, regenerates `complaint.md` when `situation.md` exists, otherwise requires a pre-existing `complaint.md`, then invokes `.bin/aar case` as before.

Verification: `bash -n arbitrate.sh` passed. I did not run a case, because that would invoke external model services and remove the requested output directory.

## 2026-04-29 final pre-batch build

Verification: Ran `PATH="$HOME/.elan/bin:$PATH" make build` in `arb/`. Lean `aarengine` built successfully and the Go `aar` binary was rebuilt into `.bin/aar`.

## 2026-04-29 ex2-ex6 batch execution

Execution: Ran ex2 through ex6 three times each using `arbitrate.sh`, serially, with unique output directories `out/exN-rK`. The active batch directory was `out/_batch-ex2-ex6-20260429-215039`.

Results:
- ex2-r1: ok, not_demonstrated, 1-4.
- ex2-r2: ok, not_demonstrated, 0-5.
- ex2-r3: ok, not_demonstrated, 1-4.
- ex3-r1: ok, not_demonstrated, 1-4.
- ex3-r2: ok, demonstrated, 5-0.
- ex3-r3: ok, demonstrated, 3-2.
- ex4-r1: ok, demonstrated, 3-2.
- ex4-r2: ok, demonstrated, 5-0.
- ex4-r3: ok, demonstrated, 4-1.
- ex5-r1: ok, demonstrated, 4-1.
- ex5-r2: ok, not_demonstrated, 0-5.
- ex5-r3: ok, demonstrated, 3-2.
- ex6-r1: error, plaintiff exceeded invalid-attempt limit; incomplete packet, events only.
- ex6-r2: error, defendant exceeded invalid-attempt limit; incomplete packet, events only.
- ex6-r3: ok, not_demonstrated, 2-3.

Verification: There are 15 `out/exN-rK` directories for ex2 through ex6. Thirteen contain the expected packet files. `out/ex6-r1` and `out/ex6-r2` are incomplete because the run terminated after invalid-attempt-limit errors before writing the full packet. No relevant `aar case` or batch-runner processes remained after completion.

## 2026-04-30 13:41 CDT — Ukraine mineral-deal clearer-polarity C1-C3 trial batch

Observe: The repo, input directory, and executable arbitrate.sh are present. No existing output directory matched the requested C1-C3 20260430 prefix during the initial check.
Think: Run the three requested arbitrations sequentially with unique output directories and no retries except trivial directory collisions. Preserve all command output in the terminal/session logs and inspect packet files after completion.
Do: Started with a visible Slack progress update in thread C0ARASSVB45 / 1777565144.132449.
Verify: Will check each output directory after the corresponding run and inspect digest/state/council artifacts before the final Slack report.

C1: Output directory out/ukraine-mineral-deal-clear-gpt5-search-c1-20260430-134145. The run failed before resolution with JSON error: {"status":"error","error":"plaintiff exceeded invalid-attempt limit"}. The wrapper shell also attempted to assign zsh's read-only `status` variable afterward, causing a shell exit code 1; that wrapper error happened after arbitrate.sh had already reported the invalid-attempt failure. No retry performed.

C2: Output directory out/ukraine-mineral-deal-clear-gpt5-search-c2-20260430-134316. The run failed before resolution with JSON error: {"status":"error","error":"defendant exceeded invalid-attempt limit"}. No retry performed.

C3: Output directory out/ukraine-mineral-deal-clear-gpt5-search-c3-20260430-134640. The run failed before resolution with JSON error: {"status":"error","error":"defendant exceeded invalid-attempt limit"}. No retry performed.

Inspection: Each C1-C3 output directory contains only events.ndjson. No digest, state, council packet, transcript, or final result files were produced. All three failed in openings before council voting. No result, resolution, or vote split is available. The visible malformed-output pattern is attorney opening statements exceeding the 4000-character tool limit; no vote/rationale inconsistency can be assessed because voting never began.

## 2026-04-30 five sequential Ukraine clear runs

Observe: The operator requested five strictly sequential runs using ./arbitrate.sh, INPUT_DIR=examples/ukraine-mineral-deal-clear, ATTORNEY_MODEL=openai://gpt-5?tools=search, and fresh output directories with logs.
Think: Run one process at a time, capture combined stdout/stderr into per-run logs, inspect packets/events after each completion, and report concise progress in Slack.
Do: Starting from the arb directory. No commits will be made.

Verify: Completed five runs sequentially. Logs were captured in the batch directory and copied into each output directory as run.log. Inspected output directories, run.json where present, events.ndjson, council votes, and over-limit submissions. Four runs completed demonstrated; one failed before resolution after defendant exhausted the five invalid attempts on opening. Seq3 contains an apparent C2 vote/rationale mismatch: vote not_demonstrated, rationale supports absence of a qualifying government agreement.
Document: Detailed summary written to out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-seq5-20260430-173545/summary.md.

## 2026-04-30 revised Ukraine clear runs — council model/persona behavior

Observation: The operator asked for a cursory analysis of all revised `ukraine-mineral-deal-clear` case runs, limited to council behavior: model mix, persona reuse, vote patterns, and vote/rationale mismatches. I inspected the revised `invalid5` output directories: the earlier `c6-invalid5` run plus the 5 + 5 + 10 sequential batches. Twenty-one revised output directories existed. Nineteen reached council votes. Two failed before council voting and cannot support council-behavior conclusions.

Scope inspected:
- `out/ukraine-mineral-deal-clear-gpt5-search-c6-invalid5-20260430-161809`
- `out/ukraine-mineral-deal-clear-gpt5-search-invalid5-seq*-20260430-*`
- `out/ukraine-mineral-deal-clear-gpt5-search-invalid5-more*-20260430-*`
- `out/ukraine-mineral-deal-clear-gpt5-search-invalid5-tenmore*-20260430-*`
- Batch summaries:
  - `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-seq5-20260430-173545/summary.md`
  - `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-more5-20260430-183739/summary.md`
  - `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-tenmore-20260430-194511/summary.md`

Aggregate council outcomes among the nineteen completed revised runs:
- All completed revised runs resolved `demonstrated`.
- Vote-split distribution:
  - `5-0`: 10 runs.
  - `4-1`: 5 runs.
  - `3-2`: 4 runs.

Main behavior finding: the council outcome was stable, but the vote label was fragile. The completed councils consistently reached `demonstrated`, but many dissenting `not_demonstrated` votes were not substantive disagreements. Several rationales said, in substance, that there was no official qualifying U.S.–Ukraine rare-earth-elements agreement by the deadline and that the market should have resolved `No`. That rationale supports the arbitration proposition and should map to `demonstrated`. The failure mode appears to be polarity confusion between the market outcome label `No` and the arbitration result label `demonstrated`.

Model observations:
- Strongest consistency: Anthropic Opus/Sonnet-family models, Gemini 3 Flash, GPT-5.x, Grok, and most Gemini 2.5 Flash votes. These generally tracked the timing/explicitness theory and voted coherently.
- Gemini 2.5 Flash was mostly stable. It had one substantive `not_demonstrated` vote based on incomplete official-source coverage, but otherwise voted `demonstrated` consistently.
- Weakest label discipline: `meta-llama/llama-4-scout`, `amazon/nova-premier-v1`, and legacy `openai/gpt-4`. Llama Scout voted `not_demonstrated` 4 times in 8 completed votes, and several of those rationales appeared label-inverted rather than actually adverse to the plaintiff. Nova showed the same problem in several `not_demonstrated` votes. GPT-4 appeared twice and both votes looked inverted.
- `z-ai/glm-4.7-flash` produced more genuine burden-skepticism than most models. Its dissents tended to argue that the plaintiff had not exhaustively canvassed official Ukrainian and U.S. sources, so the categorical negative had not been proven. That is a real adjudicative posture rather than a mere label error.
- `openai/gpt-4o` had one apparent inversion, but otherwise behaved normally in this small sample.

Completed-vote counts by model:
- `google/gemini-2.5-flash`: 14 `demonstrated`, 1 `not_demonstrated`.
- `anthropic/claude-opus-4.5`: 9 `demonstrated`, 0 `not_demonstrated`.
- `google/gemini-3-flash-preview`: 9 `demonstrated`, 0 `not_demonstrated`.
- `meta-llama/llama-4-scout`: 4 `demonstrated`, 4 `not_demonstrated`.
- `amazon/nova-premier-v1`: 4 `demonstrated`, 3 `not_demonstrated`.
- `openai/gpt-4o`: 6 `demonstrated`, 1 `not_demonstrated`.
- `openai/gpt-5.4`: 6 `demonstrated`, 0 `not_demonstrated`.
- `z-ai/glm-4.7-flash`: 4 `demonstrated`, 2 `not_demonstrated`.
- `anthropic/claude-sonnet-4.5`: 5 `demonstrated`, 0 `not_demonstrated`.
- `x-ai/grok-4-fast`: 5 `demonstrated`, 0 `not_demonstrated`.
- `google/gemini-3.1-flash-lite-preview`: 4 `demonstrated`, 0 `not_demonstrated`.
- `x-ai/grok-3`: 4 `demonstrated`, 0 `not_demonstrated`.
- `anthropic/claude-sonnet-4.6`: 3 `demonstrated`, 0 `not_demonstrated`.
- `openai/gpt-5.2-chat`: 3 `demonstrated`, 0 `not_demonstrated`.
- `anthropic/claude-opus-4.6`: 2 `demonstrated`, 0 `not_demonstrated`.
- `openai/gpt-4`: 0 `demonstrated`, 2 `not_demonstrated`; both appeared inverted on rationale inspection.

Persona observations:
- The model effect appears stronger than the persona effect. The same persona behaved differently depending on the model. For example, `d715074-9` was coherent with Claude/Sonnet-family models but produced inverted votes when paired with Llama Scout. `d715074-0` was stable with Gemini/GPT/Grok pairings and unstable mostly when paired with Llama/GLM.
- `e50e538-1` is worth watching. It had a high `not_demonstrated` rate in this sample, but the votes came through GPT-4/GPT-4o-style label inversions or burden framing rather than a clear persona-driven disagreement. I would not attribute the pattern to the persona without more controlled pairings.
- The bookkeeping/payroll written-lie persona `c4e1a2b-0` did not consistently distort the case. Its problematic votes were mostly Nova label inversions. With other models it behaved normally.
- Several personas were perfectly stable in this sample: `d715074-2`, `d715074-5`, and `d715074-8` all voted `demonstrated` every time they reached a vote. This may reflect model assignment as much as persona design.

Completed-vote counts by persona:
- `d715074-0.txt`: 16 `demonstrated`, 4 `not_demonstrated`.
- `d715074-9.txt`: 11 `demonstrated`, 2 `not_demonstrated`.
- `d715074-4.txt`: 10 `demonstrated`, 1 `not_demonstrated`.
- `c4e1a2b-0.txt`: 7 `demonstrated`, 2 `not_demonstrated`.
- `d715074-5.txt`: 8 `demonstrated`, 0 `not_demonstrated`.
- `d715074-8.txt`: 8 `demonstrated`, 0 `not_demonstrated`.
- `d715074-2.txt`: 7 `demonstrated`, 0 `not_demonstrated`.
- `e50e538-1.txt`: 3 `demonstrated`, 3 `not_demonstrated`.
- `d715074-1.txt`: 4 `demonstrated`, 1 `not_demonstrated`.
- `d715074-7.txt`: 4 `demonstrated`, 0 `not_demonstrated`.
- `e50e538-0.txt`: 3 `demonstrated`, 0 `not_demonstrated`.
- `d715074-6.txt`: 1 `demonstrated`, 0 `not_demonstrated`.

Council-composition issue: Persona duplication was common. Some councils reused the same persona across multiple seats, sometimes three seats in one run. That weakens the apparent diversity of the council. Different models still produced variation, but the persona layer did not always produce five distinct human priors.

Interpretation: The council was not struggling with the merits of this revised case as much as it was struggling with polarity. The complaint proposition says the market should have resolved `No`; the arbitration result label for proving that proposition is `demonstrated`. Several models wrote the correct rationale and then chose the wrong arbitration label.

Recommended next analysis step: separate council votes into three buckets before using raw vote splits:
1. coherent `demonstrated` votes;
2. true `not_demonstrated` dissents based on burden/completeness;
3. label-inverted `not_demonstrated` votes whose rationales support `demonstrated`.

This separation will matter more than the raw vote split for evaluating model/persona behavior.

## 2026-04-30 simplest resolution-right Ukraine variant

Observation: the operator proposed stripping away another layer of indirection and testing the simplest version where the plaintiff argues the market resolution was right.

Decision: Create a new example rather than overwrite `examples/ukraine-mineral-deal-clear`, preserving the earlier negative/no-resolution version and prior run data.

Do: Created `examples/ukraine-mineral-deal-resolution-right-simple`. It copies the source record files from `examples/ukraine-mineral-deal-clear` and replaces the situation/complaint proposition with a direct pro-resolution claim: the Polymarket final `Yes` resolution was right because official U.S. or Ukrainian government information showed a qualifying U.S.–Ukraine rare-earth-elements deal by the deadline.

Verify: The new `situation.md` and `complaint.md` are identical and contain only the simplified proposition. No case run has completed yet for this variant at the time of this note.

## 2026-04-30 simplified pro-Yes sequential run batch

Resumed the five-run sequential batch for `examples/ukraine-mineral-deal-resolution-right-simple` using `openai://gpt-5?tools=search` after seq1 and seq2 had completed and seq3 was active. Monitored seq3 until exit, copied/confirmed `run.log`, then ran seq4 and seq5 strictly sequentially.

Batch directory: `out/_batch-ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739`.

Output directories:
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq1-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq2-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq3-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq4-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739`

All five runs reached `closed` with resolution `not_demonstrated`. Final council votes parsed from `digest.md` council-vote sections were:
- seq1: 1 `demonstrated`, 4 `not_demonstrated`
- seq2: 0 `demonstrated`, 5 `not_demonstrated`
- seq3: 0 `demonstrated`, 5 `not_demonstrated`
- seq4: 1 `demonstrated`, 4 `not_demonstrated`
- seq5: 0 `demonstrated`, 5 `not_demonstrated`

Invalid submission/recovery summary: seq1 had 5 character-limit failures and 2 ENOENT read failures; seq2 had 3 character-limit failures and 2 ENOENT read failures; seq3 had 4 character-limit failures and 3 ENOENT read failures; seq4 had 4 character-limit failures, 3 ENOENT read failures, and 1 validation failure; seq5 had 6 character-limit failures and 3 ENOENT read failures. None exhausted the invalid-attempt limit or failed before resolution.

No obvious vote/rationale polarity mismatches were found. The simplified pro-Yes proposition appears to have reduced polarity confusion: `demonstrated` consistently meant affirming the Polymarket Yes resolution, and `not_demonstrated` consistently meant rejecting it. The two minority `demonstrated` votes, seq1 C3 and seq4 C4, both had pro-Yes rationales.

Wrote batch summary at `out/_batch-ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739/summary.md`. Note: seq4's local wrapper exited 1 after completion because the zsh wrapper attempted to assign to read-only variable `status`; seq5's wrapper recorded `SEQ5_EXIT=1` despite final run-log JSON `status: ok`. In both cases, `run.json`, `digest.md`, and `run.log` were present and complete.

## 2026-04-30 five sequential runs: resolution-right-simple with GPT-5 search attorneys

Task: Run five sequential trials for `examples/ukraine-mineral-deal-resolution-right-simple` using attorney model `openai://gpt-5?tools=search`, fresh output directories, and per-run logs. No commits.

Batch artifacts:
- Batch directory: `out/_batch-ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739`
- Per-run logs: `out/_batch-ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739/logs/seq1.log` through `seq5.log`
- Run outputs:
  - `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq1-20260430-213739`
  - `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq2-20260430-213739`
  - `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq3-20260430-213739`
  - `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq4-20260430-213739`
  - `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739`

Run outcomes:
- `seq1`: `not_demonstrated`, 4-1. The lone `demonstrated` vote relied on official presidential remarks, the February White House-hosted poll, and March 20 materials as enough to show a pre-deadline announcement.
- `seq2`: `not_demonstrated`, 5-0. The council rejected stitching a February Ukrainian rare-earth readout to a March 25 White House security/agreement readout.
- `seq3`: `not_demonstrated`, 5-0. The council treated the White House-hosted Harvard CAPS/Harris poll as a hypothetical survey question, not an announcement of bilateral assent.
- `seq4`: `not_demonstrated`, 4-1. The lone `demonstrated` vote accepted the February Ukrainian rare-earth partnership statement plus March 5 Senate remarks about a minerals deal at the signing stage.
- `seq5`: `not_demonstrated`, 5-0. The council focused on official April 30/May 1 agreement records and found them post-deadline and not rare-earth-explicit.

Aggregate: all five runs resolved `not_demonstrated`; final council votes were 23 `not_demonstrated` and 2 `demonstrated`.

Main finding: the simplified pro-Yes proposition flipped the council polarity in the expected direction. In the earlier negative/no-resolution framing, most councils resolved `demonstrated`, meaning they accepted the proposition that the market should not have resolved Yes. In this simplified variant, the proposition says the final Yes resolution was right, and every completed council rejected that proposition. The dominant rationale remained the same factual theory: no official pre-deadline source showed the same U.S.-Ukraine deal both agreed and explicitly involving Ukrainian rare earth elements; the clearest official agreement date was April 30, after the deadline; generic minerals/natural-resource language did not satisfy the explicit rare-earth requirement.

Vote/rationale mismatch finding: the final council votes in these five runs did not show the label-inversion problem seen in the earlier negative-proposition tests. The 23 `not_demonstrated` votes mostly gave rationales that the Yes resolution was not supported. The 2 `demonstrated` votes gave pro-Yes rationales and therefore were coherent minority disagreements, not vote-label mistakes. I saw no final council vote where a rationale substantively rejected the Yes resolution while voting `demonstrated`, or substantively endorsed Yes while voting `not_demonstrated`.

Residual system observations:
- The attorneys still explored several weak pro-Yes theories: White House-hosted polling language, February Ukrainian readouts, March 5 Senate remarks, March 20 presidential remarks, March 25 unrelated White House agreements, and April 30/May 1 official agreement documents. The council usually rejected these because assent, REE explicitness, and timing did not align on the same deal.
- There were transient invalid submissions in all five runs, mostly character-limit violations and `/home/user/...` file path mistakes. Counts by run: `seq1` 5 character-limit invalids and 2 ENOENTs; `seq2` 3 character-limit invalids and 2 ENOENTs; `seq3` 4 character-limit invalids and 3 ENOENTs; `seq4` 4 character-limit invalids, 3 ENOENTs, and 1 malformed-tool validation failure; `seq5` 6 character-limit invalids and 3 ENOENTs. None exhausted the invalid-submission allowance or caused a run failure.

Interpretation: Reframing the proposition as “the final Yes resolution was right” appears to fix the council label-polarity problem for this case. The councils now map the common “no qualifying official pre-deadline agreement” rationale to `not_demonstrated`, which is the correct arbitration label for a pro-Yes proposition. The remaining variation is substantive: occasional models accept stitched official fragments as enough under preponderance, but those votes are minority and their rationales match their labels.

## 2026-04-30 consolidated AgentCourt / Polymarket arbitration research note

Purpose: record the full day of AgentCourt arbitration work on the Polymarket Ukraine mineral-deal prediction market, including work that may have started in another Slack thread. This note is written as an analyst's lab record rather than a chat summary. It consolidates the case variants, run batches, interventions, artifacts, and behavioral observations about attorneys and councils.

### Market and evidentiary object

Market: Polymarket, "Ukraine agrees to Trump mineral deal before April?"

Market rule captured in `examples/ukraine-mineral-deal*/market-page.txt`:
- Resolve `Yes` if the United States and Ukraine agreed to any deal between February 2 and March 31, 2025, 11:59 PM ET, explicitly involving Ukrainian rare earth elements.
- Otherwise resolve `No`.
- The resolution source was official information from the governments of the United States and Ukraine.
- The displayed final market outcome was `Yes`.

Core evidentiary question isolated during the day: whether official U.S. or Ukrainian government information, by the March 31 cutoff, showed a single U.S.-Ukraine deal that was both agreed and explicitly about Ukrainian rare earth elements. The repeated council rationale across successful negative-framing runs was that the strongest official agreement record dated the agreement to April 30, 2025, after the cutoff, and that generic minerals / natural-resource language did not satisfy the explicit rare-earth-elements requirement.

### Case variants created or used

1. `examples/ukraine-mineral-deal`
   - Original negative challenge formulation.
   - Proposition: the final `Yes` resolution was not justified because, by the deadline, the United States and Ukraine had not agreed to a qualifying official rare-earth-elements deal.
   - This variant was used for the first run series earlier in the day.

2. `examples/ukraine-mineral-deal-clear`
   - Created as a clearer-polarity copy of the original Ukraine mineral-deal example.
   - Proposition: `The Polymarket market should have resolved No`, followed by the direct factual clause that official U.S./Ukrainian information did not show a qualifying deal by the deadline.
   - Purpose: reduce ambiguity in the original complaint language while preserving the plaintiff's anti-Yes position.
   - Result: improved framing, but still left one layer of indirection: `market should resolve No` had to map to arbitration vote `demonstrated`.

3. `examples/ukraine-mineral-deal-resolution-right-simple`
   - Created later as the simplest pro-resolution variant.
   - Proposition: `The Polymarket market's final Yes resolution was right`, followed by the direct factual clause that official U.S. or Ukrainian government information showed a qualifying U.S.-Ukraine rare-earth-elements deal by the deadline.
   - Purpose: test whether stripping away the negative/no-resolution framing eliminated council vote-label polarity errors.
   - Result: yes. Councils rejected the pro-Yes proposition with coherent `not_demonstrated` votes and no obvious final vote/rationale polarity inversions.

### Runner and invalid-submission intervention

Operational issue observed repeatedly: attorney submissions, especially opening statements, exceeded the submission character limit. Early failures ended runs before council voting, producing only `events.ndjson` and no final packet.

Relevant runner changes and upstream changes:
- `arbitrate.sh` had previously been generalized so `sign.sh` is optional, `complaint.md` is regenerated from `situation.md` when present, and checked-in complaints may be used directly.
- Upstream/repo state was synced to commit `f4b9bc6` (`Raise opening cap and update guidance`) and `arb` was rebuilt successfully.
- `arbitrate.sh` was then updated locally to pass a fixed `--invalid-attempt-limit 5` on every `.bin/aar case` invocation.
- `bash -n arb/arbitrate.sh` verified shell syntax after the local change.
- No commits were made.

Effect of intervention:
- Raising caps and giving invalid-submission feedback improved completion rates but did not by itself eliminate attorney over-length behavior.
- Moving to `--invalid-attempt-limit 5` materially improved completion rates. Many runs had one or more invalid submissions and then recovered. A minority still exhausted all five invalid attempts.

### Original variant run series: `examples/ukraine-mineral-deal`

Output directories inspected:
- `out/ukraine-mineral-deal-gpt5-search-r1-20260430-111945`
- `out/ukraine-mineral-deal-gpt5-search-r2-20260430-122933`
- `out/ukraine-mineral-deal-gpt5-search-r3-20260430-123205`
- `out/ukraine-mineral-deal-gpt5-search-r4-20260430-125701`
- `out/ukraine-mineral-deal-gpt5-search-r5-20260430-125915`
- `out/ukraine-mineral-deal-gpt5-search-r6-20260430-130835`

Results:
- r1: completed, `demonstrated`, 4-1, run_id `run-1777566013971234000`.
- r2: failed before resolution; only `events.ndjson` was present. Character-limit failures included plaintiff openings at 4185 and 4007 characters against a 4000 limit, and defendant openings at 4728 and 4539 characters against the same limit.
- r3: completed, `demonstrated`, 3-2, run_id `run-1777570325896487000`.
- r4: failed before resolution. Local note in `out/ukraine-mineral-deal-gpt5-search-r4-20260430-125701/analysis.md` recorded defendant opening failures at 4758 and 4472 characters against a 4000 limit, followed by invalid-attempt exhaustion. No `run.json`, `state.json`, `council.json`, or `digest.md` was produced.
- r5: completed, `demonstrated`, 4-1, run_id `run-1777571955953107000`.
- r6: completed, `demonstrated`, 4-1, run_id `run-1777572515515865000`.

Aggregate for original variant:
- Completed runs: 4 of 6.
- Failed before council resolution: 2 of 6.
- Completed-run results: all `demonstrated`.
- Completed-run final votes: 15 `demonstrated`, 5 `not_demonstrated`.
- Interpretation: councils generally accepted the anti-Yes challenge, but attorney length failures were still a dominant operational failure mode.

### Clear negative/no-resolution variant: early trials before fixed `invalid-attempt-limit 5`

Output directories inspected:
- `out/ukraine-mineral-deal-clear-gpt5-search-c1-20260430-134145`
- `out/ukraine-mineral-deal-clear-gpt5-search-c2-20260430-134316`
- `out/ukraine-mineral-deal-clear-gpt5-search-c3-20260430-134640`
- `out/ukraine-mineral-deal-clear-gpt5-search-c4-20260430-145420`
- `out/ukraine-mineral-deal-clear-gpt5-search-c5-20260430-160901`
- `out/ukraine-mineral-deal-clear-gpt5-search-c6-invalid5-20260430-161809`

Results and operational observations:
- c1: failed before resolution; plaintiff opening exceeded 4000-character cap at 4920 and 4421 characters.
- c2: failed before resolution; defendant opening exceeded 4000-character cap at 4707 and 4638 characters.
- c3: failed before resolution; plaintiff opening exceeded 4000-character cap at 4079 characters; defendant opening exceeded 4000-character cap at 4360 and 4323 characters.
- c4: failed before resolution after new feedback was present but limit remained 3 invalid attempts. Defendant opening submitted 4664, then 4166, then 4361 characters against a 4000 limit. The feedback told the attorney to resubmit at 3500 characters or fewer; the attorney did not comply enough to finish.
- c5: after upstream cap/guidance changes and rebuild, the run still failed before resolution. Plaintiff first submitted 5080 characters against a 5000 cap and recovered; defendant then submitted 5299, 5159, and 5018 characters against the 5000 cap and exhausted three invalid attempts.
- c6-invalid5: after explicitly running `.bin/aar case --invalid-attempt-limit 5`, the run completed successfully: `demonstrated`, 5-0, run_id `run-1777583889484636000`. Invalid-submission recovery improved: defendant opening initially 5379 > 5000 and then accepted at 4646; defendant argument initially 6124 > 6000 and then accepted at 5085. No obvious vote/rationale mismatch was observed.

Interpretation:
- The clearer negative proposition did not by itself solve attorney invalid-submission failures.
- The increased invalid-attempt allowance was the key operational improvement.
- Once the run reached council voting, the council result was stable and strongly anti-Yes / pro-No.

### Revised `invalid5` clear negative/no-resolution batches

Batches inspected:
- `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-seq5-20260430-173545`
- `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-more5-20260430-183739`
- `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-tenmore-20260430-194511`

Run outputs inspected:
- `out/ukraine-mineral-deal-clear-gpt5-search-invalid5-seq1-20260430-173553` through `seq5-20260430-180323`
- `out/ukraine-mineral-deal-clear-gpt5-search-invalid5-more1-20260430-183739` through `more5-20260430-183739`
- `out/ukraine-mineral-deal-clear-gpt5-search-invalid5-tenmore01-20260430-194511` through `tenmore10-20260430-194511`

Batch 1, five sequential runs:
- seq1: failed before resolution after defendant exhausted five invalid opening submissions: 6334, 5285, 5668, 5598, and 5098 characters against a 5000 limit.
- seq2: completed, `demonstrated`, 5-0.
- seq3: completed, `demonstrated`, 4-1. Apparent C2 vote/rationale mismatch: vote `not_demonstrated`, rationale supported absence of a qualifying government agreement, which supports the plaintiff's proposition.
- seq4: completed, `demonstrated`, 5-0. Also recovered from one OpenRouter xAI 502 during council work.
- seq5: completed, `demonstrated`, 5-0.

Batch 2, five more sequential runs:
- more1: failed before resolution after defendant exhausted five invalid opening submissions: 6801, 5969, 5536, 5361, and 5051 characters against a 5000 limit.
- more2: completed, `demonstrated`, 3-2. C2 and C3 appeared to have vote/rationale mismatches: both voted `not_demonstrated`, but their rationales supported the plaintiff/no-resolution proposition.
- more3: completed, `demonstrated`, 5-0.
- more4: completed, `demonstrated`, 5-0.
- more5: completed, `demonstrated`, 5-0.

Batch 3, ten more sequential runs:
- All ten completed successfully with resolution `demonstrated`.
- Vote split distribution: three runs at 5-0, four runs at 4-1, and three runs at 3-2.
- All ten had character-limit invalid submissions and recovered under `--invalid-attempt-limit 5`.
- Apparent vote/rationale mismatches were noted in runs 3, 4, 5, 6, and 7.

Aggregate for the revised `invalid5` clear negative/no-resolution runs, counting c6 plus the 5+5+10 sequential batches:
- Output directories: 21.
- Completed council votes: 19 runs.
- Failed before council voting: 2 runs.
- All 19 completed runs resolved `demonstrated`.
- Completed-run vote split distribution: 10 runs at 5-0, 5 runs at 4-1, and 4 runs at 3-2.
- Completed final votes: 82 `demonstrated`, 13 `not_demonstrated`.

Council-behavior interpretation for clear negative/no-resolution framing:
- The substantive council result was stable: completed councils consistently accepted that the market should have resolved `No`.
- The principal council defect was vote-label polarity confusion. Several `not_demonstrated` votes had rationales stating that no official qualifying agreement existed by the deadline, which supports the proposition and should have mapped to `demonstrated`.
- The suspected cognitive mapping problem was: `market should resolve No` -> proposition true -> arbitration vote `demonstrated`. Several models appeared to get the factual analysis right and the vote label wrong.
- This was most visible in Llama Scout, Nova Premier, and legacy GPT-4 pairings. GLM produced some genuine burden-skeptical dissents, which should be separated from label inversions.

### Council model/persona behavior analysis from revised clear runs

Model-level observations recorded after inspecting council JSON and digest votes:
- Most consistent in this sample: Anthropic Opus/Sonnet-family models, Gemini 3 Flash, GPT-5.x, Grok, and most Gemini 2.5 Flash votes.
- Weakest label discipline: `meta-llama/llama-4-scout`, `amazon/nova-premier-v1`, and legacy `openai/gpt-4`.
- `z-ai/glm-4.7-flash` was notable for real burden-skepticism: it sometimes argued that the plaintiff had not exhaustively canvassed official U.S. and Ukrainian sources, rather than merely inverting the label.
- `openai/gpt-4o` had one apparent inversion but was otherwise normal in the small sample.

Completed-vote counts by model across the revised clear-run council sample:
- `google/gemini-2.5-flash`: 14 `demonstrated`, 1 `not_demonstrated`.
- `anthropic/claude-opus-4.5`: 9 `demonstrated`, 0 `not_demonstrated`.
- `google/gemini-3-flash-preview`: 9 `demonstrated`, 0 `not_demonstrated`.
- `meta-llama/llama-4-scout`: 4 `demonstrated`, 4 `not_demonstrated`.
- `amazon/nova-premier-v1`: 4 `demonstrated`, 3 `not_demonstrated`.
- `openai/gpt-4o`: 6 `demonstrated`, 1 `not_demonstrated`.
- `openai/gpt-5.4`: 6 `demonstrated`, 0 `not_demonstrated`.
- `z-ai/glm-4.7-flash`: 4 `demonstrated`, 2 `not_demonstrated`.
- `anthropic/claude-sonnet-4.5`: 5 `demonstrated`, 0 `not_demonstrated`.
- `x-ai/grok-4-fast`: 5 `demonstrated`, 0 `not_demonstrated`.
- `google/gemini-3.1-flash-lite-preview`: 4 `demonstrated`, 0 `not_demonstrated`.
- `x-ai/grok-3`: 4 `demonstrated`, 0 `not_demonstrated`.
- `anthropic/claude-sonnet-4.6`: 3 `demonstrated`, 0 `not_demonstrated`.
- `openai/gpt-5.2-chat`: 3 `demonstrated`, 0 `not_demonstrated`.
- `anthropic/claude-opus-4.6`: 2 `demonstrated`, 0 `not_demonstrated`.
- `openai/gpt-4`: 0 `demonstrated`, 2 `not_demonstrated`; both appeared inverted on rationale inspection.

Persona-level observations:
- The model effect appeared stronger than the persona effect. The same persona behaved differently depending on model assignment.
- `d715074-9` was coherent with Claude/Sonnet-family models but produced inverted votes when paired with Llama Scout.
- `d715074-0` was stable with Gemini/GPT/Grok pairings and unstable mostly when paired with Llama/GLM.
- `e50e538-1` had a high `not_demonstrated` rate in the sample, but the problematic votes appeared to reflect model/polarity issues rather than a clear persona-driven jurisprudential difference.
- `c4e1a2b-0`, the bookkeeping/payroll written-lie persona, did not consistently distort the case. Its problematic votes were mostly Nova label inversions.
- Persona duplication within councils was common, sometimes reusing the same persona across multiple seats in one run. That weakens apparent council diversity.

Analytic recommendation from this phase:
- Do not rely on raw vote split alone for negative/no-resolution propositions.
- Classify votes into three buckets: coherent `demonstrated`, true burden/completeness `not_demonstrated`, and label-inverted `not_demonstrated` whose rationale supports `demonstrated`.

### Simplified pro-Yes / resolution-right variant

Variant created: `examples/ukraine-mineral-deal-resolution-right-simple`.

Proposition:
- `The Polymarket market's final Yes resolution was right.`
- By the deadline, official U.S. or Ukrainian government information showed a qualifying U.S.-Ukraine rare-earth-elements deal.

Purpose:
- Strip away the `market should resolve No` negative-polarity layer.
- Make arbitration labels direct: `demonstrated` means the final `Yes` resolution was justified; `not_demonstrated` means it was not proven.

Batch inspected:
- `out/_batch-ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739`

Output directories:
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq1-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq2-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq3-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq4-20260430-213739`
- `out/ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739`

Results:
- seq1: `not_demonstrated`, final council 4-1 against the pro-Yes proposition. The lone `demonstrated` vote relied on official presidential remarks, the February White House-hosted poll, and March 20 materials as enough to show a pre-deadline announcement.
- seq2: `not_demonstrated`, 5-0. The council rejected stitching a February Ukrainian rare-earth readout to a March 25 White House security/agreement readout.
- seq3: `not_demonstrated`, 5-0. The council treated the White House-hosted Harvard CAPS/Harris poll as a hypothetical survey question, not an announcement of bilateral assent.
- seq4: `not_demonstrated`, 4-1. The lone `demonstrated` vote accepted the February Ukrainian rare-earth partnership statement plus March 5 Senate remarks about a minerals deal at the signing stage.
- seq5: `not_demonstrated`, 5-0. The council focused on official April 30/May 1 agreement records and found them post-deadline and not rare-earth-explicit.

Aggregate for simplified pro-Yes variant:
- Completed runs: 5 of 5.
- Failed before resolution: 0 of 5.
- Final results: all `not_demonstrated`.
- Final council votes: 23 `not_demonstrated`, 2 `demonstrated`.
- No obvious final vote/rationale polarity mismatches.

Invalid-submission behavior in simplified pro-Yes variant:
- seq1: 5 character-limit invalids and 2 ENOENT read failures.
- seq2: 3 character-limit invalids and 2 ENOENT read failures.
- seq3: 4 character-limit invalids and 3 ENOENT read failures.
- seq4: 4 character-limit invalids, 3 ENOENT read failures, and 1 malformed-tool validation failure.
- seq5: 6 character-limit invalids and 3 ENOENT read failures.
- None exhausted the invalid-attempt allowance or failed before resolution.

Interpretation:
- The simplified pro-Yes proposition flipped council polarity in the expected direction. Councils rejected the final `Yes` resolution as `not_demonstrated`, using the same underlying factual theory as the earlier anti-Yes / pro-No runs.
- The simplified framing appears to fix the final vote-label polarity problem for this case. The 23 `not_demonstrated` votes had anti-Yes rationales, and the 2 `demonstrated` votes had pro-Yes rationales. The minority pro-Yes votes were coherent substantive disagreements rather than label mistakes.
- Residual disagreement now appears to be about evidentiary sufficiency: whether stitched official fragments, presidential remarks, polling text hosted by the White House, or post-deadline official records can satisfy preponderance. The majority answer was no.

### Current scientific conclusions

1. The strongest stable finding is evidentiary, not procedural: completed councils repeatedly converged on the view that the official record did not show a qualifying U.S.-Ukraine rare-earth-elements deal by March 31, 2025, 11:59 PM ET.

2. The original and clear negative variants generally obtained the expected anti-Yes result, but the negative/no-resolution wording created label-polarity risk. Some council members wrote rationales supporting the plaintiff and then voted `not_demonstrated` because they appeared to map `No` in the market to `not_demonstrated` in arbitration.

3. The simplified pro-Yes variant is a better experimental probe of council behavior. It converts the same substantive question into direct label semantics: if official qualifying evidence existed, vote `demonstrated`; if not, vote `not_demonstrated`. In the five-run sample, this eliminated observed final vote/rationale inversions.

4. Attorney invalid submissions remain a separate operational problem. `--invalid-attempt-limit 5` dramatically improved completion rates, but attorney openings and later submissions still frequently exceed limits. The system is recovering rather than preventing the behavior.

5. The remaining council disagreement is substantive and interpretable. Minority pro-Yes votes tend to accept stitched official evidence or ambiguous government statements. Majority not-demonstrated votes require timing, bilateral assent, official-source status, and explicit rare-earth-elements scope to align in the same pre-deadline record.

### Open issues for future analysis

- Separate label inversions from true dissents programmatically. A post-processor should classify each vote/rationale pair as coherent, inverted, or substantively dissenting.
- Reduce persona duplication per council if the purpose is to test persona diversity rather than model diversity.
- Test the simplified pro-Yes formulation across more runs and compare model-specific behavior against the negative formulation.
- Consider making the attorney prompts explicitly state the vote-label mapping for negative market-resolution propositions, or avoid such propositions entirely when evaluating council reliability.
- Investigate why attorneys repeatedly attempt `/home/user/...` file paths and over-length filings. These are operational defects distinct from council reasoning.

### Artifact inventory for this note

Primary analysis notes:
- `analysis.md` in this directory.
- `out/ukraine-mineral-deal-gpt5-search-r4-20260430-125701/analysis.md`.
- `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-seq5-20260430-173545/summary.md`.
- `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-more5-20260430-183739/summary.md`.
- `out/_batch-ukraine-mineral-deal-clear-gpt5-search-invalid5-tenmore-20260430-194511/summary.md`.
- `out/_batch-ukraine-mineral-deal-resolution-right-simple-gpt5-search-invalid5-seq5-20260430-213739/summary.md`.

Case directories:
- `examples/ukraine-mineral-deal`.
- `examples/ukraine-mineral-deal-clear`.
- `examples/ukraine-mineral-deal-resolution-right-simple`.

Runner modified but not committed:
- `arbitrate.sh`, now passing `--invalid-attempt-limit 5` to `.bin/aar case`.

## 2026-05-01 — Zelenskyy suit one-run arbitration

Ran one sequential search-enabled AgentCourt arbitration for `examples/zelenskyy-suit-condition-simple` after preflight confirmed `.bin/aar`, `.bin/aarengine`, `$HOME/keys.txt`, and Podman were present/running.

Command:

```bash
./arbitrate.sh "examples/zelenskyy-suit-condition-simple" "out/zelenskyy-suit-condition-simple-gpt5-search-one-20260501-100251" 'openai://gpt-5?tools=search'
```

Result:

- Run id: `run-1777647771416980000`
- Status: `ok`
- Resolution: `demonstrated`
- Vote tally: 5 demonstrated, 0 not demonstrated
- Started: `2026-05-01T15:02:51Z`
- Finished: `2026-05-01T15:12:35Z`

The council treated identity, authenticity, and timing as established by official NATO and Ukrainian-government records. The decisive issue was attire. All five council members accepted that the NATO E5 imagery, combined with the plaintiff's search-enabled technical report and secondary reporting, made suit attire more likely than not.

Operational notes:

- No provider, xproxy, Podman, or runner failure occurred.
- Non-fatal invalid submissions occurred when a defendant opening and plaintiff closing exceeded character limits; both were resubmitted successfully.
- Some attorney file reads attempted missing relative paths under `/home/user`, but later filings still incorporated relevant case material and the run completed.

Summary file:

- `out/_batch-zelenskyy-suit-gpt5-search-one-20260501-100251/summary.md`

## 2026-05-01 — Israel-Lebanon one-run open-record attempt

the operator approved one open-record/search-enabled arbitration run for `examples/israel-lebanon-invasion-condition-simple`.

Preflight passed:

- `.bin/aar` present.
- `.bin/aarengine` present.
- `$HOME/keys.txt` present.
- Podman running.

Run command:

```bash
./arbitrate.sh "examples/israel-lebanon-invasion-condition-simple" "out/israel-lebanon-invasion-condition-simple-gpt5-search-one-20260501-105747" 'openai://gpt-5?tools=search'
```

Output directory:

- `out/israel-lebanon-invasion-condition-simple-gpt5-search-one-20260501-105747`

Batch directory:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-one-20260501-105747`

Result:

- Failed before council voting during defendant opening.
- Defendant exhausted five invalid submissions for exceeding the 5,000-character opening-statement limit. Attempts were 6,123; 5,650; 5,043; 5,192; and 5,002 characters.
- Plaintiff opening had two invalid over-limit submissions and then succeeded on the third attempt.
- No `run.json`, `state.json`, `digest.md`, or `council.json` was produced. `events.ndjson` and `run.log` are present.
- No retry was attempted because the operator authorized one run only.

Summary:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-one-20260501-105747/summary.md`

## 2026-05-01 — Policy JSON path for 10,000-character filing limit

Observed after the Israel-Lebanon open-record attempt failed during defendant opening because the default opening limit was 5,000 characters.

Findings:

- `aar case` supports `--policy FILE` for a Policy JSON file.
- When `--policy` is omitted, the CLI loads `./etc/policy.json` from the current working directory if present; otherwise it uses built-in defaults.
- `runtime/runner/policy.go` loads the policy by starting from `DefaultPolicy()` and then unmarshalling JSON over it, so a partial JSON file can override only selected fields.
- The current repository policy file `etc/policy.json` sets `max_opening_chars` to 5000, `max_argument_chars` to 6000, `max_rebuttal_chars` to 4000, `max_surrebuttal_chars` to 4000, and `max_closing_chars` to 5000.
- The Lean engine enforces `max_opening_chars` for `record_opening_statement`; the Go ACP bridge also computes the same phase-specific limits for feedback.
- The wrapper `arbitrate.sh` currently does not accept a policy path argument. It invokes `.bin/aar case` without `--policy`, so it uses `arb/etc/policy.json` when run from `arb/`.

Conclusion:

A one-off 10,000-character policy can be supplied by creating a separate JSON file and invoking `.bin/aar case --policy <file>` directly, or by changing the wrapper to accept and pass a policy file. Direct CLI invocation avoids editing the repository default policy.


## 2026-05-01 — Israel-Lebanon one-run open-record attempt with 10,000-character policy

the operator approved one follow-up open-record/search-enabled test using a separate JSON policy file that raised all merits filing text limits to 10,000 characters.

Policy file:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-one-policy10000-20260501-111408/policy-merits-10000.json`

Output directory:

- `out/israel-lebanon-invasion-condition-simple-gpt5-search-one-policy10000-20260501-111408`

Batch directory:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-one-policy10000-20260501-111408`

Command shape:

```bash
.bin/aar case --complaint examples/israel-lebanon-invasion-condition-simple/complaint.md --out-dir out/israel-lebanon-invasion-condition-simple-gpt5-search-one-policy10000-20260501-111408 --attorney-model 'openai://gpt-5?tools=search' --invalid-attempt-limit 5 --policy out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-one-policy10000-20260501-111408/policy-merits-10000.json
```

Result:

- Run id: `run-1777652048300755000`
- Status: `ok`
- Resolution: `demonstrated`
- Vote tally: 3 demonstrated, 2 not demonstrated.
- Effective policy set `max_opening_chars`, `max_argument_chars`, `max_rebuttal_chars`, `max_surrebuttal_chars`, and `max_closing_chars` to 10,000.
- The previous opening-statement failure did not recur. Opening lengths were 9,083 and 8,789 characters.
- No invalid submissions were found in `events.ndjson` during artifact inspection.

Substantive split:

- Majority: September 30 ground raids in Lebanese villages, support fires, and the purpose of pushing Hezbollah back implied temporary local control intent over some Lebanese territory.
- Minority: the record proved limited raids/incursions before the deadline but did not prove a territorial holding/control objective by a preponderance.

Summary:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-one-policy10000-20260501-111408/summary.md`


## 2026-05-01 — Argument analysis observations from 10,000-character policy run

Observed after the successful open-record/search-enabled run using the one-off 10,000-character policy:

- The arguments converged on the correct hinge. Timing, actor, and location mostly dropped out. Both sides effectively accepted that Israeli ground activity inside southern Lebanon began on September 30 before the ET deadline. The live issue became whether those raids or offensive actions implied an intent to establish control over any Lebanese ground.
- Plaintiff's strongest move was the low-threshold reading of `any portion`. Plaintiff framed control as temporary, local tactical control: clearing or denying villages, routes, or approaches near the border long enough to push Hezbollah back. That made the proposition easier to satisfy. The majority accepted this theory: ground troops entering Lebanese villages with fires support and a push-back objective implied some temporary control.
- Defendant's strongest move was the contemporaneous-intent objection. The official in-window language emphasized `limited`, `localized`, `targeted ground raids`, `infrastructure`, and `limited incursions`. None of those sources expressly said hold, occupy, administer, regulate movement, or create a buffer. Defendant also correctly emphasized the absence of in-window objective indicia such as checkpoints, berms, posts, sustained positions, logistics posture, or clear holding orders.
- Plaintiff's weakest inference was that `push Hezbollah back` necessarily means control. That inference is plausible, but not necessary. An actor can push an adversary back through raids, fires, disruption, and repeated incursions without intending to hold terrain in the stronger sense. Plaintiff prevailed because the proposition used `any portion` and did not define control tightly.
- Defendant's weakest inference was demanding too much proof at commencement. The proposition asks about intent to establish control, not completed control infrastructure by hour zero. Requiring checkpoints or fixed positions inside the first evening is probably too strict. Intent can be inferred from operational design, especially if the objective is to clear border villages.
- The council split was substantive rather than a vote-label failure. The three `demonstrated` votes accepted temporary tactical control as enough. The two `not_demonstrated` votes required clearer evidence of a holding or territorial-control mission. The rationales aligned with their votes. C1 and C2 were too conclusory, but C3, C4, and C5 exposed the real dispute cleanly.
- The case is highly sensitive to the definition of `control`. If control means temporary tactical exclusion or denial of ground, plaintiff probably wins. If control means a holding mission, territorial authority, or sustained control measures, defendant probably wins. The market text does not resolve that boundary.
- The run is scientifically useful. It did not produce a random council artifact. The split maps onto a genuine ambiguity in the predicate and shows AgentCourt surfacing the controlling interpretive variable: temporary tactical control versus sustained territorial control.

Operational note: the 10,000-character policy improved process quality for this case. It removed the earlier opening-statement failure without producing invalid submissions, and it let both sides state the real hinge rather than fight the limit.


## 2026-05-01 — Four-run sequential 10,000-character policy batch

The operator requested four additional sequential arbitration runs with the same settings as the successful one-off policy test: Israel-Lebanon case, open-record/search-enabled `openai://gpt-5?tools=search`, invalid-attempt limit 5, and a separate policy JSON setting all merits filing limits to 10,000 characters.

Batch directory:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more4-20260501-113242`

Summary:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more4-20260501-113242/summary.md`

Results:

- Run 01: operational failure before merits result. Error: `acp attorney did not submit a decision`. The run produced `events.ndjson` and `run.log`, but no `run.json`. Last events show the defendant in the argument phase attempted to access `/home/user/official-source-record.txt`, received `ENOENT`, and did not submit the required decision. There was also one recoverable plaintiff opening over-limit attempt at 10,636 characters against the 10,000-character policy limit.
- Run 02: completed, `demonstrated`, 3 demonstrated / 2 not_demonstrated. Run id `run-1777653494900883000`.
- Run 03: completed, `not_demonstrated`, 0 demonstrated / 5 not_demonstrated. Run id `run-1777654015359955000`. There was one recoverable defendant opening over-limit attempt at 10,019 characters against the 10,000-character policy limit.
- Run 04: completed, `demonstrated`, 4 demonstrated / 1 not_demonstrated. Run id `run-1777654590699490000`.

Aggregate for this batch:

- Requested runs: 4.
- Completed merits runs: 3.
- Operational failures: 1.
- Completed-run resolutions: 2 demonstrated / 1 not_demonstrated.
- Completed-run council votes: 7 demonstrated / 8 not_demonstrated.

Aggregate including the earlier successful 10,000-character run:

- Completed 10,000-character merits runs: 4.
- Resolutions: 3 demonstrated / 1 not_demonstrated.
- Council votes: 10 demonstrated / 10 not_demonstrated.

Observation:

The completed runs continue to isolate the same interpretive variable. Councils voting `demonstrated` treat temporary tactical clearing, denial, or local control of Lebanese villages/routes as enough to satisfy `control over any portion of Lebanon`. Councils voting `not_demonstrated` require clearer evidence of an intent to hold or administer terrain, and treat `limited`, `localized`, and `targeted` raid language as insufficient. The batch is therefore useful as a stability test: result labels lean demonstrated by run count, while individual council votes are exactly balanced across the completed 10,000-character runs.


## 2026-05-01 — Five-run sequential 10,000-character policy batch

The operator requested five more sequential arbitration runs with the same settings: Israel-Lebanon case, open-record/search-enabled `openai://gpt-5?tools=search`, invalid-attempt limit 5, and the separate policy JSON setting all merits filing limits to 10,000 characters.

Batch directory:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more5-20260501-123102`

Summary:

- `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more5-20260501-123102/summary.md`

Results:

- Run 01: completed, `not_demonstrated`, 2 demonstrated / 3 not_demonstrated. Run id `run-1777656662484524000`.
- Run 02: completed, `demonstrated`, 3 demonstrated / 2 not_demonstrated. Run id `run-1777657323931777000`.
- Run 03: completed, `demonstrated`, 4 demonstrated / 1 not_demonstrated. Run id `run-1777657993633222000`.
- Run 04: completed, `demonstrated`, 3 demonstrated / 2 not_demonstrated. Run id `run-1777658531644729000`.
- Run 05: completed, `demonstrated`, 4 demonstrated / 1 not_demonstrated. Run id `run-1777659140113700000`.

Aggregate for this batch:

- Requested runs: 5.
- Completed merits runs: 5.
- Operational failures: 0.
- Completed-run resolutions: 4 demonstrated / 1 not_demonstrated.
- Completed-run council votes: 16 demonstrated / 9 not_demonstrated.

Aggregate across completed 10,000-character policy runs so far:

- Completed merits runs: 9.
- Operational failures in attempted 10,000-character runs: 1.
- Resolutions: 7 demonstrated / 2 not_demonstrated.
- Council votes: 26 demonstrated / 19 not_demonstrated.

Operational notes:

Every run in this batch produced `run.json`, `council.json`, `digest.md`, and `events.ndjson`. All five runs had recoverable over-limit opening attempts under the 10,000-character policy. The over-limit attempts ranged from 10,049 to 10,730 characters and were retried successfully.

Observation:

The result distribution now leans toward `demonstrated` by both run count and council-vote count, but the underlying split remains the same. The positive votes read `control over any portion` to include temporary tactical clearing or denial of Lebanese villages/routes. The negative votes require stronger evidence of a holding, occupation, administration, or sustained territorial-control mission and treat official `limited`, `localized`, and `targeted raids` language as insufficient.


## 2026-05-01 — Observations across ten attempted 10,000-character policy runs

Scope:

- Count: 10 attempted 10,000-character policy runs for `israel-lebanon-invasion-condition-simple`.
- Completed merits runs: 9.
- Operational failures: 1.
- Completed-run resolutions: 7 demonstrated / 2 not_demonstrated.
- Completed-run council votes: 26 demonstrated / 19 not_demonstrated.

Observations:

- The completed-result distribution leans `demonstrated`: 7 of 9 completed merits runs. The aggregate council vote is 26 demonstrated / 19 not_demonstrated. This is a real lean, but not a stable consensus.
- The factual record is no longer the main source of variance. Councils generally agree that Israeli ground activity inside Lebanon began on September 30 before the deadline. The dispute has collapsed onto one interpretive variable: what `establish control over any portion of Lebanon` means.
- The `demonstrated` theory is that `any portion` plus `commenced` sets a low threshold. Ground raids into Lebanese villages, combined-arms support, clearing Hezbollah infrastructure, and pushing Hezbollah away from the border imply temporary tactical control or denial over some Lebanese ground.
- The `not_demonstrated` theory is that `control` means more than raid, disruption, or temporary presence. Negative votes treat the official language — `limited`, `localized`, and `targeted ground raids` against infrastructure — as evidence against a holding, occupation, administration, or sustained territorial-control objective.
- The split is coherent rather than a vote-label artifact. The rationales mostly match the votes. AgentCourt is exposing a genuine ambiguity in the predicate rather than producing random output.
- The `demonstrated` wins are often close or moderate: 3-2, 3-2, 4-1, 3-2, 4-1, etc. The strongest negative result was a 0-5 `not_demonstrated` run. That matters because a council adopting the stricter definition of control can reject the proposition cleanly.
- The market wording is under-specified. If `control` includes tactical clearing, local denial, or temporary control during raids, the condition is probably demonstrated. If `control` requires a holding mission or sustained territorial authority, the condition is probably not demonstrated. The market text does not resolve that boundary.
- Operationally, the 10,000-character policy made the runs workable, but it did not remove filing-pressure problems. One attempted run failed before merits because an attorney did not submit a decision after file-access errors. The later five-run batch completed cleanly, but every run still had recoverable over-limit opening attempts.

Bottom line:

This is no longer primarily an evidence-sufficiency problem. It is a definitional-control problem. The experimental signal is that `demonstrated` is the more frequent outcome under open-record AgentCourt, but the predicate remains scientifically unstable unless `control` is defined.


## 2026-05-01 — Ten additional 10,000-character policy runs

Batch: `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more10-20260501-143152`.

Results:

- Run 01: `demonstrated`, 4-1, `run-1777663912784496000`.
- Run 02: `demonstrated`, 3-2, `run-1777664459733122000`.
- Run 03: `not_demonstrated`, 1-4, `run-1777665078745760000`.
- Run 04: `demonstrated`, 4-1, `run-1777665691912179000`.
- Run 05: `demonstrated`, 3-2, `run-1777666333771699000`.
- Run 06: `demonstrated`, 3-2, `run-1777666880804991000`.
- Run 07: `demonstrated`, 3-2, `run-1777667388786878000`.
- Run 08: `demonstrated`, 5-0, `run-1777668086260537000`.
- Run 09: `demonstrated`, 4-1, `run-1777668801431749000`.
- Run 10: `demonstrated`, 4-1, `run-1777669392280714000`.

Aggregate for this batch: 9 `demonstrated` / 1 `not_demonstrated` by resolution, 34 `demonstrated` / 16 `not_demonstrated` by council vote. No operational failures.

Cumulative 10,000-character policy aggregate: 20 attempted runs, 19 completed merits runs, 1 operational failure. Completed-run resolutions: 16 `demonstrated` / 3 `not_demonstrated`. Council votes: 60 `demonstrated` / 35 `not_demonstrated`.

Summary file: `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more10-20260501-143152/summary.md`.
Machine-readable run results: `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more10-20260501-143152/run-results.json`.
Machine-readable model vote stats: `out/_batch-israel-lebanon-invasion-condition-simple-gpt5-search-policy10000-more10-20260501-143152/model-vote-stats.json`.

## 2026-05-01 — Council-model metadata and interpretive-threshold note

Follow-up model check:

The public metadata did not provide a clean release-date or training-regime explanation for the council split. The strongest signal remains behavioral rather than documentary: the models applied different thresholds for the word `control` in the proposition.

Public metadata checked through OpenRouter listed the main council models as follows:

- `openrouter://google/gemini-3-flash-preview`: Dec. 2025 listing, long-context thinking/agentic model. In the cumulative 10,000-character-policy sample it voted 16 `demonstrated` / 0 `not_demonstrated`.
- `openrouter://google/gemini-2.5-flash`: Jun. 2025 listing, reasoning/workhorse model. It voted 8 / 4 and was therefore mixed.
- `openrouter://anthropic/claude-opus-4.6`: Feb. 2026 listing, long-running professional/agentic model. It voted 8 / 3.
- `openrouter://meta-llama/llama-4-scout`: Apr. 2025 listing, MoE model, 17B active / 109B total by OpenRouter metadata. It voted 7 / 1.
- `openrouter://x-ai/grok-3` and `openrouter://x-ai/grok-4-fast`: smaller sample, but leaned `demonstrated`.
- `openrouter://openai/gpt-5.4`, `openrouter://openai/gpt-5.2-chat`, and `openrouter://amazon/nova-premier-v1`: strongly leaned `not_demonstrated`; observed cumulative counts were GPT-5.4 0 / 6, GPT-5.2 Chat 0 / 5, and Nova Premier 0 / 4.

What did not explain the split:

- Recency did not explain it. GPT-5.4 and Claude Opus 4.6 were both listed as 2026-era frontier models, but they split in opposite directions.
- Context length did not explain it. Long-context models appeared on both sides.
- Generic `reasoning` or `agentic` positioning did not explain it. Both sides included models marketed for reasoning or agentic work.
- Public training descriptions were too thin to support a causal training-regime claim. The provider-facing material mostly exposed capability framing, safety/evaluation summaries, and high-level post-training language, not enough to infer why a model used one legal threshold rather than another.

Observed behavioral split:

- The `demonstrated`-leaning models treated documented ground raids, village clearing, Hezbollah-infrastructure destruction, and IDF/U.S. statements as sufficient to infer temporary tactical control over some Lebanese territory.
- The `not_demonstrated`-leaning models required a more explicit record of control: occupation, holding territory, administration, sustained presence, or an intent to control territory rather than conduct raids.

Practical conclusion:

Council composition is a material experimental variable in this scenario. Model family and model version should be recorded as covariates in stability analysis. The models should not be treated as interchangeable jurors when the dispositive issue is an under-specified legal or quasi-legal threshold such as `control`.


## 2026-05-03 Polymarket official token launch one open-record run

Observe: the operator supplied a Polymarket token-launch condition and, after being asked for the required record mode, chose open-record/search-enabled. A new case directory was created at `examples/polymarket-official-token-launch-condition`. The first invocation failed immediately because `situation.md` lacked the required `# Proposition` heading.

Think: The setup failure was an input-format defect, not a merits failure. After correction, one search-enabled run was appropriate because The operator asked for one arbitration and selected open record.

Do: Ran one arbitration with `openai://gpt-5?tools=search`, output `out/polymarket-official-token-launch-condition-gpt5-search-one-20260503-181128`. The arbitration completed `demonstrated`, 3-2. The decisive plaintiff theory was that Polymarket launched pUSD as a standard ERC-20 collateral token on Polygon, documented in official Polymarket pUSD, Contracts, Changelog, and CLOB V2 migration pages. After the run, captured those pages plus Polygonscan corroboration under the case directory's `source-captures/` and updated the case evidence notes.

Verify: Inspected `run.json`, `state.json`, `digest.md`, `council.json`, `events.ndjson`, and `run.log`. Council vote labels matched rationales. The two dissents were coherent burden/semantic dissents, not label inversions. Summary written to `out/_batch-polymarket-official-token-launch-condition-gpt5-search-one-20260503-181128/summary.md`.

Document: The successful arbitration artifacts are complete. Operational caveat: the surrounding shell wrapper exited with code 1 after success because it referenced `PIPESTATUS` under zsh; this did not affect the arbitration result.

## 2026-05-03 Polymarket official token launch nine additional open-record runs

Observe: The operator asked for nine more sequential open-record/search-enabled runs on `examples/polymarket-official-token-launch-condition` after the first run resolved `demonstrated`, 3-2. The case packet had been backfilled with official Polymarket pUSD, Contracts, Changelog, and CLOB V2 migration captures plus Polygonscan corroboration.

Think: Run exactly nine additional attempts sequentially, use fresh output directories, and count failures as failures rather than retrying silently.

Do: A detached batch ran nine sequential attempts in `out/_batch-polymarket-official-token-launch-condition-gpt5-search-nine-more-20260503-183500`. Eight completed and all eight resolved `demonstrated`; one failed before council with `acp attorney did not submit a decision`. Completed vote splits were 5-0, 5-0, 5-0, 4-1, 4-1, 5-0, 3-2, and 5-0. Across the eight completed additional runs: 36 demonstrated votes and 4 not_demonstrated votes.

Verify: Inspected `run.json`, `state.json`, `digest.md`, `council.json`, `events.ndjson`, and logs for completed outputs. The dissenting votes were coherent semantic or burden/record-sufficiency dissents, not vote-label inversions. Summary written to `out/_batch-polymarket-official-token-launch-condition-gpt5-search-nine-more-20260503-183500/summary.md`.

Document: The manifest contains parser errors for completed runs because the wrapper's Python one-liners used unquoted dict keys. Results above are from direct artifact inspection, not the manifest fields. An accidental duplicate continuation was started during diagnosis and stopped; its incomplete events-only directory is `out/polymarket-official-token-launch-condition-gpt5-search-nine-more-r6-20260503-191633` and is not counted.

## 2026-05-11 Russia-Ukraine ceasefire condition single open-record run

Task: The operator asked for one ARB on the Polymarket `Russia x Ukraine ceasefire by May 31, 2026?` resolution, with the proposition set up under the prediction-market arbitration guidelines. After being asked for the required record mode, the operator chose open-record/search-enabled.

Input directory created: `examples/russia-ukraine-ceasefire-condition-may31`.

Proposition:

`By May 31, 2026, 11:59 PM ET, Russia and Ukraine had officially reached a publicly announced and mutually agreed halt in military engagement that constituted a general pause in the conflict, not merely an informal agreement, humanitarian pause, limited sector-specific ceasefire, or future framework.`

Packet scope: the packet preserves the Polymarket market rule and final/disputed-state text. It records that official Russian/Ukrainian source captures and full secondary-source captures remain unresolved gaps, so the case was suitable for open-record first-pass testing rather than closed-record packet-only adjudication.

Run mode: open-record/search-enabled using `openai://gpt-5?tools=search` for both attorneys. One run only. No retries.

Run artifacts:
- Batch directory: `out/_batch-russia-ukraine-ceasefire-condition-may31-gpt5-search-one-20260511-142750`.
- Output directory: `out/russia-ukraine-ceasefire-condition-may31-gpt5-search-one-20260511-142750`.
- Summary: `out/_batch-russia-ukraine-ceasefire-condition-may31-gpt5-search-one-20260511-142750/summary.md`.

Result: completed cleanly with status `ok`, final phase `closed`, resolution `not_demonstrated`, council vote 0 `demonstrated` / 5 `not_demonstrated`.

Merits pattern: the council accepted the defense theory. Russia's on-record acceptance of a May 9-11 truce was not the weak point. The weak point was Ukraine: the council treated Presidential Decree No. 374/2026 as a humanitarian, geofenced Red Square carve-out, treated Zelenskyy's language as prospective, noted no Ukrainian General Staff or MOD stand-down order for May 9-11, and found mutuality and general theater-wide scope unproven by a preponderance.

Operational notes: no invalid-attempt-limit failure occurred. The run log contained no occurrences of `invalid`, `exceeded`, `error`, `ENOENT`, or `validation`. Vote labels matched rationales; I saw no polarity inversion.

## 2026-05-11 Russia-Ukraine ceasefire condition nine additional open-record runs

Task: The operator asked for nine more sequential ARBs after the first single open-record run on `examples/russia-ukraine-ceasefire-condition-may31`. Runs used `openai://gpt-5?tools=search`, fresh output directories, and no retries.

Batch directory: `out/_batch-russia-ukraine-ceasefire-condition-may31-gpt5-search-more9-20260511-150238`.

Results across the nine additional runs: 5 `demonstrated`, 4 `not_demonstrated`. Council votes: 22 `demonstrated`, 23 `not_demonstrated`. Per-run results: seq1 demonstrated 4-1; seq2 demonstrated 4-1; seq3 not_demonstrated 2-3; seq4 demonstrated 3-2; seq5 demonstrated 3-2; seq6 not_demonstrated 1-4; seq7 demonstrated 4-1; seq8 not_demonstrated 1-4; seq9 not_demonstrated 0-5.

Including the prior single run, total observed results are 5 `demonstrated` and 5 `not_demonstrated`; council votes are 22 `demonstrated` and 28 `not_demonstrated`.

Main interpretation: the proposition is unstable under open-record arbitration. Demonstrated votes treat concordant public statements from Zelenskyy and Ushakov, plus credible-media consensus, as enough to show a mutually agreed May 9-11 general ceasefire before the May 31 deadline. Not-demonstrated votes require stronger proof of theater-wide scope and operative state action, especially paired official texts, stand-down orders, or a Ukrainian instrument broader than the Red Square decree. No obvious vote-label polarity inversion was observed.

Operational note: the detached worker announced completion prematurely while the batch was still running. The underlying batch script continued correctly and was monitored to completion from the main session. Summary written to `out/_batch-russia-ukraine-ceasefire-condition-may31-gpt5-search-more9-20260511-150238/summary.md`.

## 2026-05-12 Epstein suicide-note condition open-record run

Task: Run one open-record/search-enabled arbitration for the Polymarket `Epstein suicide note released by May 8?` dispute.

Input: `examples/epstein-suicide-note-condition-simple`. Proposition: By May 8, 2026, 11:59 PM ET, a message or note credibly reported to have been written by Jeffrey Epstein and intended as a suicide note, final message, or equivalent communication had been made publicly available.

Setup: Preserved the Polymarket event page/rule text, Gamma API event and grouped-market JSON, page-visible May 8 dispute/final-review state, and a minimal reporting source map. The packet is intentionally open-record, not closed-record; it lacks the SDNY note/order PDF and complete NYT captures.

Execution: Preflight passed for `.bin/aar`, `.bin/aarengine`, `$HOME/keys.txt`, and Podman. Ran one sequential arbitration using `openai://gpt-5?tools=search`. Output: `out/epstein-suicide-note-condition-simple-gpt5-search-one-20260512-173940`. Batch: `out/_batch-epstein-suicide-note-condition-simple-gpt5-search-one-20260512-173940`.

Result: Completed `ok`, final phase `closed`, resolution `demonstrated`, vote split 5 demonstrated / 0 not_demonstrated. Vote rationales were coherent; no vote-label inversion was found. The council accepted that multiple reputable reports plus court-file release and quoted text satisfied `credibly reported` and public availability before the May 8 deadline, despite authentication and intent caveats.

Operational notes: One nonfatal ENOENT occurred when plaintiff attempted to access `/home/user/analysis.md`; the run recovered. The surrounding zsh wrapper failed after the successful run because it used Bash-style `PIPESTATUS[0]`; `run.log` was copied manually and all expected artifacts were verified.

## 2026-05-12 Epstein suicide-note open-record ten-attempt batch

Execution: After the initial open-record run resolved demonstrated 5-0, The operator requested nine more sequential runs. The additional batch directory is `out/_batch-epstein-suicide-note-condition-simple-gpt5-search-seq9-20260512-180144`. The case directory was `examples/epstein-suicide-note-condition-simple`, with proposition: By May 8, 2026, 11:59 PM ET, a message or note credibly reported to have been written by Jeffrey Epstein and intended as a suicide note, final message, or equivalent communication had been made publicly available.

Results: Ten total attempts including the original first run. Seven completed and all seven resolved `demonstrated`. Aggregate completed-run council votes were 32 demonstrated / 3 not_demonstrated. Vote splits among completed runs: five 5-0 demonstrated, one 4-1 demonstrated, and one 3-2 demonstrated. Three additional attempts failed before council voting because OpenRouter returned 404 for `anthropic/claude-3.7-sonnet`; no silent retries were performed.

Interpretation: Completed open-record runs consistently treated `credibly reported` as satisfied by the court-file release plus multiple reputable reports before the deadline. The coherent dissents focused on `purported` / `alleged` caveats and argued that credible reporting did not prove authorship and suicide-note intent. Summary: `out/_batch-epstein-suicide-note-condition-simple-gpt5-search-seq9-20260512-180144/summary.md`.


## 2026-05-12 Epstein suicide-note seq11-seq13 after council-pool fix

Observe: The operator requested exactly three more sequential open-record/search-enabled arbitrations for `examples/epstein-suicide-note-condition-simple`, using `openai://gpt-5?tools=search`, fresh output directories seq11-seq13, and the current default council pool after removing `openrouter://anthropic/claude-3.7-sonnet`.

Think: Prior seq03, seq08, and seq09 failed before council voting because the old pool selected unavailable `anthropic/claude-3.7-sonnet` through OpenRouter. These three runs should test the same merits condition while checking whether the pool fix eliminates that operational failure.

Do: Ran seq11, seq12, and seq13 sequentially. Batch directory: `out/_batch-epstein-suicide-note-condition-simple-gpt5-search-seq3-after-pool-fix-20260512-195720`. Output directories: `out/epstein-suicide-note-condition-simple-gpt5-search-seq11-20260512-195730`, `out/epstein-suicide-note-condition-simple-gpt5-search-seq12-20260512-200833`, and `out/epstein-suicide-note-condition-simple-gpt5-search-seq13-20260512-201518`.

Verify: All three runs completed. Results were seq11 `demonstrated` 4-1, seq12 `demonstrated` 5-0, and seq13 `demonstrated` 5-0. Aggregate for this batch: 3 completed, 0 failed, 3 demonstrated resolutions, 14 demonstrated / 1 not_demonstrated council votes. No vote-label/rationale inversion was observed. No run selected the removed Claude 3.7 Sonnet model, and no OpenRouter 404 occurred.

Document: Summary written at `out/_batch-epstein-suicide-note-condition-simple-gpt5-search-seq3-after-pool-fix-20260512-195720/summary.md`. Combined with the prior ten attempts, the experiment now has 13 attempts, 10 completed, 3 failed before council voting, 10 completed demonstrated resolutions, and completed council votes of 46 demonstrated / 4 not_demonstrated.

## 2026-05-14 Clavicular pregnancy open-record run

Observation: The operator requested one open-record arbitration for the Polymarket “Clavicular pregnancy in 2026?” dispute. The prepared case directory is `examples/clavicular-pregnancy-credible-announcement-condition-simple`. It preserves the Polymarket rule and final state, Forbes dispute reporting, secondary reporting on the pregnancy rumor/contest context, raw captures, and explicit gaps.

Decision: The proposition was stated as a direct factual claim: “Between market creation on April 20, 2026, and December 31, 2026, 11:59 PM ET, Clavicular made a credible announcement that he and a partner were expecting a baby through pregnancy.” It avoids asking whether Polymarket, UMA, or any oracle was right.

Preflight: `.bin/aar`, `.bin/aarengine`, and `$HOME/keys.txt` were present. `podman info` initially failed because Podman was not connected; `arbitrate.sh` started `podman-machine-default` successfully.

Execution: A first invocation failed before model/API arbitration work because `situation.md` lacked the required `# Proposition` section. I corrected the heading without changing the proposition and then ran the actual open-record arbitration in a fresh output directory with `openai://gpt-5?tools=search`.

Run output: `out/clavicular-pregnancy-credible-announcement-condition-simple-gpt5-search-seq1-20260514-134706`. Batch summary: `out/_batch-clavicular-pregnancy-credible-announcement-condition-simple-gpt5-search-one-20260514-134706/summary.md`. Status ok. Resolution `not_demonstrated`. Council vote 1 demonstrated / 4 not_demonstrated. Vote labels matched rationales.

Interpretation: The majority found that the record lacked the primary Clavicular clip/transcript or representative statement, and that the quoted “I will be a dad” language plus the single Forbes paraphrase did not prove a clear credible announcement of an existing pregnancy with a partner by a preponderance. The run exposed the same packet gap identified before the run: a later closed-record packet should backfill the original livestream/source material and any Polymarket clarification.

## 2026-05-14 OpenClaw-backed attorney development

Observe: The operator asked to start development so an OpenClaw agent can act as counsel in a real `arb` run. Existing `aar case` already supports per-side `--plaintiff-acp-command`, `--defendant-acp-command`, `--plaintiff-acp-endpoint`, and `--defendant-acp-endpoint`. The AAR runner exposes `_aar/get_case`, `_aar/list_case_files`, `_aar/read_case_text_file`, `_aar/write_case_file`, and `_aar/submit_decision` as ACP client methods. The Pi attorney path makes those methods available to the model by injecting a PI extension from `common/submodules/pi-acp/src/acp/ext-method-tools.ts`.

Think: The existing `openclaw acp` bridge is not sufficient by itself for a real AAR attorney. It forwards ACP prompts into an OpenClaw Gateway session, but it does not expose AAR's arbitrary `_aar/*` client methods as tools to that session. A direct substitution would leave the lawyer unable to file through `_aar/submit_decision`.

Decision: Implement the smallest safe bridge first: a stdio ACP attorney server that receives the AAR prompt, calls `_aar/get_case`, delegates the legal decision to a configured OpenClaw/lawyer command, parses a strict JSON decision, and files it through `_aar/submit_decision`. This proves the real AAR filing path without changing OpenClaw configuration or running an external arbitration. A later decision remains whether to use this bridge as the production path or add native dynamic-client-tool support to `openclaw acp`.

Verify: No arbitration run has been started. No commits, pushes, OpenClaw configuration changes, or gateway restart have been made.

Implementation progress: added `arb/runtime/cmd/aar-openclaw-attorney` and `arb/runtime/openclawattorney`. The adapter speaks the AAR stdio ACP subset needed by `common/acp.Client`: `initialize`, `session/new`, and `session/prompt`. On each prompt it calls AAR client methods `_aar/get_case`, `_aar/list_case_files`, and `_aar/read_case_text_file` for text-readable files, then obtains a strict JSON decision either from `AAR_OPENCLAW_ATTORNEY_DECISION_JSON`, a custom `AAR_OPENCLAW_ATTORNEY_COMMAND`, or `openclaw agent` when `AAR_OPENCLAW_AGENT=1`. It files the parsed decision through `_aar/submit_decision`.

Decision: the first implementation preloads the visible record and text-readable case files instead of attempting a full dynamic OpenClaw tool bridge. This makes a real filing path testable now. It does not yet provide dynamic `aar_write_case_file` or binary-file handling to the OpenClaw lawyer. A production-quality version may still need native client-tool bridging if we want OpenClaw to call AAR tools interactively rather than through a preloaded prompt.

Verification: `cd arb/runtime && CGO_ENABLED=0 GOCACHE=$(pwd)/../.cache/go-build go test ./...` passed. `PATH="$HOME/.elan/bin:$PATH" make build` passed and produced `.bin/aar-openclaw-attorney`. No arbitration run was started.

Follow-up progress: the adjudication checkout already has `origin` set to `git@github-mrmarvinclaw:mrmarvinclaw/adjudication.git` and `upstream` set to `https://github.com/agentcourt/adjudication.git`, so a personal fork remote is already configured. I created local branch `openclaw-attorney-adapter` for this work. No commit or push was made.

Added README documentation for `aar-openclaw-attorney`, including the OpenClaw-agent mode and deterministic fixed-decision mode. Ran a subprocess-level dry adapter exercise against `.bin/aar-openclaw-attorney`; the harness exchanged real newline-delimited ACP JSON-RPC over stdio, handled `_aar/get_case`, `_aar/list_case_files`, `_aar/read_case_text_file`, observed `_aar/submit_decision`, and verified the submitted JSON exactly matched the fixed local decision. This still avoided any real arbitration run or external model call.

Verification: `cd arb/runtime && CGO_ENABLED=0 GOCACHE=$(pwd)/../.cache/go-build go test ./...` passed. `PATH="$HOME/.elan/bin:$PATH" make build` passed. The subprocess dry adapter exercise printed `dry-adapter-subprocess-ok`.

Open design issue: `AAR_OPENCLAW_AGENT=1` currently invokes `openclaw agent --message <full prompt> --json`. That is simple and uses the configured Gateway path, but very large case packets can hit OS argument-size limits. The custom command mode reads the full packet from stdin and does not have that limitation. Before making the OpenClaw-agent path production-grade, consider either adding a Gateway RPC client inside the adapter or adding an OpenClaw CLI message-file/stdin option upstream.

Safety adjustment: changed the OpenClaw-agent mode to require `AAR_OPENCLAW_AGENT_ID`. The adapter now fails rather than silently routing arbitration prompts to the default personal OpenClaw agent. README examples use `aar-lawyer` as a dedicated-agent placeholder. Added a focused test for this guardrail. Re-ran `go test ./...`, `make build`, and the subprocess dry adapter exercise; all passed.

Closed Clavicular OpenClaw-lawyer integration run:
- the operator authorized closed-record mode and registering a dedicated OpenClaw lawyer agent.
- Registered OpenClaw agent `aar-lawyer` with workspace `<workspace>/agents/aar-lawyer`. This changed OpenClaw config after explicit approval.
- First OpenClaw-plaintiff/Pi-defendant attempt failed before defendant filing: Pi ACP returned `Internal error: Connection error` after emitting its update notice. The OpenClaw plaintiff had successfully filed an opening through `_aar/submit_decision`.
- Switched to OpenClaw on both sides after the operator approved proceeding.
- All-OpenClaw attempt initially exposed an adapter bug: `_aar/submit_decision` rejected an over-length filing (`6273 characters submitted, 6000 allowed`), and the adapter treated the rejection as fatal instead of re-prompting within the same AAR opportunity.
- Fixed the adapter to retry up to three rejected submissions for the same AAR opportunity, including the exact rejection text in the next lawyer prompt. Added regression coverage for rejected-decision retry. Re-ran `go test ./...` and `make build`; both passed.
- Next all-OpenClaw attempt reached council inference but failed because the subprocess environment lacked `OPENROUTER_API_KEY`; loaded only `OPENAI_API_KEY` and `OPENROUTER_API_KEY` from `~/keys.txt` without printing values.
- Successful run: `out/clavicular-openclaw-both-closed-20260514-160055`, run id `run-1778792455175313000`.
- Configuration: closed-record explicit packet files, plaintiff and defendant both used `<repo-root>/arb/.bin/aar-openclaw-attorney`, attorney model flags were `openai://gpt-5` with search disabled, OpenClaw agent model env was `openai-codex/gpt-5.5`, council size 3.
- Result: `not_demonstrated`, council vote `0 demonstrated / 3 not_demonstrated`.
- Council: C1 `openrouter://google/gemini-3-flash-preview`, C2 `openrouter://openai/gpt-4o`, C3 `openrouter://x-ai/grok-4-fast`.
- Event audit: `events.ndjson` contained 1 `run_initialized`, 8 `attorney_action`, and 3 `council_vote` events; no invalid/error events in the successful run.
- Output files present: `complaint.md`, `council.json`, `digest.md`, `events.ndjson`, `policy.json`, `run.json`, `run.log`, `runtime.json`, `state.json`, and `transcript.md`.
- No commit, push, publication, or gateway restart was performed.

## 2026-05-14 — ACP endpoint attorney model cleanup

Observe: the operator corrected the design direction for OpenClaw lawyers.  AAR should use `--*-acp-endpoint` and assume an OpenClaw ACP attorney is already running there.  OpenClaw should not receive model selection from AAR because it owns model/session selection internally.

Think: The existing implementation made every attorney inherit `--attorney-model`, parsed it as an xproxy URI, derived `search_enabled` from that URI, and recorded that model even for remote ACP endpoints.  That created false metadata and misleading prompt capabilities for OpenClaw.  The correct boundary is generic ACP: AAR controls the transport and `_aar/*` client-method contract; the remote ACP attorney controls its model and native tools.

Do: Changed remote endpoint resolution so endpoint roles omit `model`, omit `search_enabled`, and receive a neutral capability prompt.  Rejected role-specific `--*-attorney-model` when the same role uses `--*-acp-endpoint`.  Reworded CLI help and README to scope `--attorney-model` to local Pi/xproxy attorneys.  Removed `AAR_OPENCLAW_AGENT_MODEL` forwarding from the local OpenClaw compatibility adapter so it no longer calls `openclaw agent --model`.

Verify: Ran `go test ./...` from `arb`; all packages passed.  No commit made.

## 2026-05-14 — OpenClaw ACP endpoint arbitration smoke run

Observation: `aar case` was modified so role-specific `--*-acp-endpoint` attorneys are recorded as TCP ACP endpoints and do not receive AAR local model/search metadata. A temporary TCP stdio bridge on `127.0.0.1:19701` launched `.bin/aar-openclaw-attorney` with `AAR_OPENCLAW_AGENT_ID=aar-lawyer` and no `AAR_OPENCLAW_AGENT_MODEL`.

Action: Ran `examples/clavicular-pregnancy-credible-announcement-condition-simple` with both plaintiff and defendant using `tcp://127.0.0.1:19701`, council size 3. The first run reached council voting but failed because xproxy did not inherit `OPENROUTER_API_KEY`. Reran after sourcing `~/keys.txt` into the arb process environment.

Verification: Second run completed successfully at `out/clavicular-openclaw-acp-endpoint-env-20260514-171249`. `events.ndjson` records both attorneys as `acp_transport=tcp` with endpoint `tcp://127.0.0.1:19701`, with no attorney model or search fields. `digest.md` reports `Resolution: not_demonstrated`. `run.json` status is `ok`; all three council votes were `not_demonstrated`.

Conclusion: The endpoint path works for a full arbitration run. The remaining operational requirement is to ensure xproxy provider keys are present in the arb process environment before council voting.

## 2026-05-14 — Preserve OpenClaw ACP TCP bridge helper

Observation: The operator asked to keep the temporary TCP bridge as a helper. The bridge is useful because `aar case --*-acp-endpoint` expects a TCP ACP endpoint, while `.bin/aar-openclaw-attorney` is a stdio adapter.

Action: Added `tools/openclaw-acp-tcp-bridge.js`. The helper listens on `127.0.0.1:19701` by default, spawns `.bin/aar-openclaw-attorney` per TCP connection, pipes socket data to stdio, defaults `AAR_OPENCLAW_AGENT_ID` to `aar-lawyer`, and strips `AAR_OPENCLAW_AGENT_MODEL` from the child environment. Added `make openclaw-acp-bridge` and README usage notes. The helper is documented as a local smoke-test/simple-integration tool, not a hardened network service.

Verification: `node --check tools/openclaw-acp-tcp-bridge.js` passed. `tools/openclaw-acp-tcp-bridge.js --help` printed the documented options. `make test` passed for all Go runtime packages. A live check on port `19702` confirmed the helper listens, accepts a TCP connection, spawns `.bin/aar-openclaw-attorney`, and exits the child cleanly when the test socket closes.

## 2026-05-14 — Evidence handling note

Observation: the operator clarified that the earlier technical-report discussion was really about evidence handling. Current AAR filings can cite existing packet evidence through `offered_files`, and can summarize newly discovered material through `technical_reports`, but cannot yet ingest newly discovered evidence as first-class case files.

Action: Wrote `../evidence.md` at the adjudication repo root documenting the current state, the consequence for open-record Clavicular-style runs, and the desired `_aar/submit_evidence` / `_aar/attach_evidence` direction.

Follow-up: Added a `Scope of Change` section to `../evidence.md` noting that first-class attorney-submitted evidence affects Lean state, transition rules, proofs, runtime ACP methods, filing schemas, rendering/output artifacts, policy, and reproducibility semantics.

## 2026-05-14 18:32 CDT — Clavicular OpenClaw open-record run start

Observe: The requested case directory exists, AAR/OpenClaw binaries are present, provider keys are available in ~/keys.txt, and Podman is running. The current endpoint design intentionally leaves model selection and native tool availability to OpenClaw; AAR endpoint metadata therefore does not itself mark search_enabled true.
Think: Use the endpoint bridge rather than misleading local xproxy model metadata. Supply an OpenClaw extra prompt making this an open-record evidence-discovery run and asking the lawyers to use any available public-search/browser tools for the original clip, transcript, or clarification. Run exactly one arbitration with both sides using the OpenClaw attorney endpoint and a fresh output directory.
Do: Starting output directory out/clavicular-openclaw-both-open-20260514-183235, batch directory out/_batch-clavicular-openclaw-both-open-20260514-183235, bridge port 19713.

Run command completed for out/clavicular-openclaw-both-open-20260514-183235. The AAR run itself returned ok/not_demonstrated; the wrapper exited nonzero after completion because the shell was zsh and the bash-only PIPESTATUS array reference failed after tee. The run log was copied manually into the output directory.

Verify: Required artifacts are present for out/clavicular-openclaw-both-open-20260514-183235. The OpenClaw bridge process left by the wrapper error was stopped after both attorney children had exited cleanly. Summary written to out/_batch-clavicular-openclaw-both-open-20260514-183235/summary.md. Council result: not_demonstrated, 0 demonstrated / 3 not_demonstrated. Both attorneys filed open-record technical reports; no original Clavicular VOD, clean Clavicular transcript, representative statement, or separate Polymarket clarification was retrieved.

## 2026-05-14 — OpenClaw-attorney reproducibility documentation

Observation: The operator asked that the approach of using OpenClaw attorneys be documented well enough for another person to reproduce it from the repository.

Action: Added `docs/openclaw-attorneys.md` with the architecture, prerequisites, exact closed-record and open-record command sequences, environment variables, capability boundary, post-run artifact checks, Clavicular development-run paths, and current limitations around technical reports versus first-class submitted evidence. Updated `README.md` to link the guide from the layout and OpenClaw attorney sections. Added a `devnotes.md` entry pointing to the guide.

Verification: `node --check tools/openclaw-acp-tcp-bridge.js` passed, `make test` passed, and a local Markdown-link check passed for `README.md`, `docs/openclaw-attorneys.md`, and `devnotes.md`. No commit or push has been made.

## 2026-05-14 evidence support proof repair

Observation: `make prove` now passes `Proofs.StepPreservation`, `Proofs.CaseFrame`, `Proofs.ReachableInvariants`, `Proofs.ReachableMaterialLimits`, and `Proofs.RecordProvenance`. The remaining proof failure is in `Proofs.CouncilIntegrity`, where the public `step` dispatch still falls through on the new `submit_evidence` action.

Thinking: The correct repair is to add explicit preservation for `submit_evidence`. Evidence submission appends to `submitted_evidence` and must leave council votes and council members unchanged, so it should preserve council vote integrity without weakening the invariant.

Do the right thing: Add a specific `step_submit_evidence_preserves_councilVoteIntegrity` lemma and dispatch branch, then rerun `make prove` to find the next invariant file requiring the new action case.

Verification: Pending.
