# Forge

A Go-based CLI tool for orchestrating developer workflows.

## Current features

* CLI command parsing
* File input handling
* Initial `bug-triage` workflow scaffold with data artifacts and runs
* YAML workflow definitions with sequential placeholder and LLM steps
* Basic notetaking

## Usage

```bash
forge read <file>
forge run bug-triage --input <file>
forge run workflow-definition --file <workflow.yaml> [--input <text>]
forge note add <note>
```

Workflow definition files use a name, an optional description, and an ordered list of steps:

```yaml
name: release-notes
description: Draft release notes
steps:
  - id: collect
    type: llm
    prompt: "Collect changes from this input: {{.Input}}"
  - id: summarize
    type: llm
    prompt: "Summarize these changes: {{.PreviousOutput}}"
```

Supported step types are `placeholder` and `llm`; omitted types default to `placeholder`. Prompt templates can use `.Input` and `.PreviousOutput`. Workflow execution depends only on the `LLMProvider` interface; the CLI currently injects an OpenAI adapter, configured with `OPENAI_API_KEY`. Other providers can be wired in without changing step execution. Runs are recorded under `data/runs/`, with Markdown artifacts under `data/artifacts/`.

Forge applies best-effort prompt-injection hardening: workflow input and prior step output are HTML-escaped inside explicit untrusted-data blocks, and LLM prompts include guidance to treat those blocks as data. This cannot reliably prevent prompt injection; treat all input and model output as untrusted. Forge also blocks LLM prompts and workflow outputs that match common credential patterns, but this heuristic may miss secrets or reject ordinary text. Never put secrets in prompts. Supply API keys through environment/secret management rather than workflow files; Forge does not persist the API key.

## Next steps
* Run viewer
* Integrate LLM for bug triage
