# Forge

A Go-based CLI tool for orchestrating developer workflows.

## Current features

* CLI command parsing
* File input handling
* Initial `bug-triage` workflow scaffold
* Basic notetaking

## Usage

```bash
forge read <file>
forge run bug-triage --input <file>
forge note add <note>
```

## Next steps

* Persist runs to disk
* Generate markdown artifacts
* Integrate LLM for bug triage
