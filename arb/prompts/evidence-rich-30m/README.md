# Evidence-rich 30-minute attorney prompts

Alternative attorney prompts for open-record, fact-intensive arbitrations with longer attorney time limits.

Use with explicit longer limits, for example:

```bash
export AAR_OPENCLAW_ATTORNEY_TIMEOUT_SECONDS=1800
.bin/aar case \
  --prompt-dir prompts/evidence-rich-30m \
  --acp-timeout-seconds 1800 \
  ...
```

These prompts preserve the same evidence-discipline rules as the default prompts, but tell attorneys to use the larger budget for staged source retrieval, transcript/OCR work, metadata checks, authenticity checks, and richer evidence-search ledgers. Missing prompt files fall back to `./prompts`.
