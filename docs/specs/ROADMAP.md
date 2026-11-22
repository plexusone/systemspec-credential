# systemspec-credential: Roadmap

## Phase 1: Registry

- [x] **RMI-OMNICREDENTIAL-001** Credential types: the entry, its parts, well-known and closed values
- [x] **RMI-OMNICREDENTIAL-002** Validation of a registry, reporting every problem with its entry
- [x] **RMI-OMNICREDENTIAL-003** Event log: types, validation, sorted append, records
- [x] **RMI-OMNICREDENTIAL-004** State: assessment of each credential at a time, with a warning window
- [x] **RMI-OMNICREDENTIAL-005** JSON load and save; JSON Schemas generated from the types

## Phase 2: Reports and Command Line

- [x] **RMI-OMNICREDENTIAL-006** Report rows and the JSON report
- [x] **RMI-OMNICREDENTIAL-007** Self-contained HTML report with search, filters, sorting and details
- [x] **RMI-OMNICREDENTIAL-008** Markdown reference, independent of time
- [x] **RMI-OMNICREDENTIAL-009** Mountable commands: validate, status, report, record, export, schema
- [x] **RMI-OMNICREDENTIAL-010** Standalone command over JSON files
- [x] **RMI-OMNICREDENTIAL-020** Handling: whether a person or a system generates, distributes and rotates each credential, counted in every report
- [x] **RMI-OMNICREDENTIAL-021** Reach: what a holder can get (`reaches`, with `dependsOn`), what a leak exposes in all and what exposes each credential
- [x] **RMI-OMNICREDENTIAL-022** Human instructions: install steps and consoles per store, accounts with protections and recovery, and a by-place checklist
- [x] **RMI-OMNICREDENTIAL-023** Management plan: summary, launch checklist in dependency order, rotation calendar, account protections and exposure
- [x] **RMI-OMNICREDENTIAL-024** Calendar export: rotations and expiries as iCalendar events with reminders
- [x] **RMI-OMNICREDENTIAL-025** Work progress: to create, to deploy, deployed, from a deployed event, with each credential's next step in every report
- [x] **RMI-OMNICREDENTIAL-026** Agent actor: an AI agent acting for an operator as a third kind of handler, counted in every report

## Phase 3: Setup Guidance

- [x] **RMI-OMNICREDENTIAL-027** Need and order: when each credential must exist (`need`), what is set up first (`requires`) and a setup order over both
- [x] **RMI-OMNICREDENTIAL-028** Provisioners: the tool and command that make a value, who runs it, and whether it stores the value
- [x] **RMI-OMNICREDENTIAL-029** Prerequisites: what is set up at an issuer before a credential, with steps and a done state
- [x] **RMI-OMNICREDENTIAL-030** Setup details: entitlements granted by, issuer mode, configuration changes, rotation secondaries, store backups and the process environment store
- [x] **RMI-OMNICREDENTIAL-031** Advice: valid entries likely to mislead, printed by validate, with a strict mode

## Phase 4: Observed State

- [ ] **RMI-OMNICREDENTIAL-011** Observer interface: what a store reports about a path without its value
- [ ] **RMI-OMNICREDENTIAL-012** Checks: missing, unlisted, and expiry from the issuer
- [ ] **RMI-OMNICREDENTIAL-013** Observer for a cloud parameter store
- [ ] **RMI-OMNICREDENTIAL-014** Observer for CI secrets and token expiry
- [ ] **RMI-OMNICREDENTIAL-015** Entitlement checks where the issuer reports them

## Phase 5: Lifecycle Automation

- [x] **RMI-OMNICREDENTIAL-032** Value format: `generation.format`, the exact form of a random value or key, so a tool can make it
- [ ] **RMI-OMNICREDENTIAL-016** Generate random values straight into a store
- [ ] **RMI-OMNICREDENTIAL-017** Rotation reminders from the due and expiring states
- [ ] **RMI-OMNICREDENTIAL-018** Findings as OCSF events
- [ ] **RMI-OMNICREDENTIAL-019** MCP server over the registry and the log
