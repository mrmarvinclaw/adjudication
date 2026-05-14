# Evidence Handling Note

## Observation

The current AAR filing model separates existing packet evidence from attorney-discovered material.

- `offered_files` are first-class evidence references, but only for case files already visible in the packet. A filing must cite existing `file_id` values.
- `technical_reports` are attorney work product. They can summarize outside investigation or analysis, and the council can consider them, but they do not create preserved evidence files with `file_id` values.

In the Clavicular example, if an attorney finds the original livestream clip, transcript, archive page, Polymarket clarification, or UMA oracle material during an open-record run, the attorney can currently describe that discovery only through `technical_reports`, unless the material was already captured as a case file before the run.

## Consequence

That is procedurally usable for open-record discovery, but it is weaker than first-class evidence handling.

A technical report can say what the attorney found, but the system does not currently:

- capture the discovered artifact;
- store its text or bytes in the run record as evidence;
- assign it a new `file_id`;
- let later filings cite it through `offered_files`;
- preserve it cleanly for audit, review, or later closed-record reruns.

This matters because open-record arbitrations often surface decisive material that was missing from the initial packet. If that material remains only a technical-report summary, a later reader cannot inspect the source with the same confidence, and a later closed-record run may fail because the decisive evidence was never added to the packet.

## Desired Direction

Add a controlled attorney-submitted evidence path.

The likely shape is an AAR client method such as `_aar/submit_evidence` or `_aar/attach_evidence` that accepts structured provenance and content, validates it, stores it, and returns a new `file_id`.

A submission should include, at minimum:

- title;
- source URL or source description;
- captured text and/or file bytes;
- MIME type;
- retrieval timestamp;
- attorney role and phase;
- short explanation of relevance;
- hash or equivalent integrity metadata for stored content.

After acceptance, the evidence should be visible in the case record and citable in later `offered_files` entries like any other exhibit.

## Scope of Change

First-class attorney-submitted evidence is a major change. It is not just a runtime convenience around `technical_reports`.

It affects at least these layers:

- Lean engine state: the formal case state needs either a new evidence object or a broader material model beyond the current `offered_files` and `technical_reports` split.
- Lean transition rules: the engine must define which phases may submit evidence, when submitted evidence becomes visible, and whether later filings may cite it.
- Lean proofs: the existing material-limit, reachability, step-preservation, no-stuck, and record-integrity proofs already reason about `offered_files` and `technical_reports`. A new evidence path would need corresponding preservation and bound arguments.
- Runtime ACP contract: AAR needs a client method such as `_aar/submit_evidence` or `_aar/attach_evidence`, with validation, error reporting, and transcript/event recording.
- Filing schema: `offered_files` can remain citation-only, but it must be able to cite newly assigned `file_id` values after evidence ingestion.
- Rendering and outputs: `digest.md`, `run.json`, `events.ndjson`, transcripts, work-product export, and any publication artifacts must preserve the submitted evidence and its provenance.
- Policy: the system needs size limits, MIME rules, per-side or per-filing caps, phase restrictions, and source/provenance requirements.
- Reproducibility: open-record discoveries need capture, hashing, and closed-record replay semantics so later runs can use the same evidence without relying on live public pages.

The small workaround is to keep using `technical_reports` as summaries. That requires little formal change, but it leaves the evidentiary gap described above. The serious version changes the formal case model.

## Interim Rule

Until first-class attorney-submitted evidence exists:

- use `offered_files` only for existing visible packet files;
- use `technical_reports` for newly discovered outside material;
- after an open-record run, backfill decisive discoveries into the case directory before treating the matter as closed-record or reproducible.
