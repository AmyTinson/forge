package triage

type TriageData struct {
	FileName  string
	RawReport string
}

const triageTemplate = `
# Bug Triage

## Input File
{{.FileName}}

## Raw Bug Report
{{.RawReport}}

## Summary
TODO

## Unknowns
- TODO

## Repro Steps
- TODO
`
