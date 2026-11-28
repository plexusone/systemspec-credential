# The Registry

The registry lists every credential a set of systems needs, without values.
It validates against rules that are the same whether you write JSON or
build it in Go, and its JSON Schema is generated from the Go types.

## One Entry per Secret, Place and Environment

The staging and production copies of a key are two entries, because they
are created, stored and rotated separately. IDs are stable and unique,
conventionally `system/environment/NAME`, using letters, digits and
`._/:@-`.

## Fields

### Identity

| Field | Required | Meaning |
|---|---|---|
| `id` | yes | Stable, unique ID |
| `name` | yes | Display name, often the variable or parameter name |
| `description` | | What it is for |
| `system` | yes | The system that owns it |
| `environment` | | `prod`, `staging` and so on; empty for one shared by every environment |
| `kind` | yes | What it is; see [Kinds](#kinds) |
| `status` | yes | `planned`, `active` or `retired` |
| `need` | | When it must exist: `launch`, `later`, `optional` (only when an optional part is turned on) or `rotation` (only while another credential rotates) |
| `requires` | | IDs of credentials and [prerequisites](#prerequisites) set up before this one |
| `sensitivity` | | `low`, `medium`, `high` or `critical` |
| `owner` | yes | The person or role who creates and rotates it |
| `tags`, `links` | | Free labels, and `title` plus `url` links |

### Where It Comes From

`issuer` is who issues the value: `type` (such as `self`, `aws`, `github`),
`account`, `console`, the URL where a person manages it, and `mode`, such
as `test` or `live` for a payment provider whose staging keys are test
keys.

`generation` is how the value is made:

| Field | Meaning |
|---|---|
| `method` | `random`, `keypair`, `issued` (by the issuer), `derived` (built from others, listed in `dependsOn`) or `federated` (no stored value) |
| `spec` | The shape, such as "32 random bytes, hex" |
| `where` | Where it is made, such as "the operator's machine, straight into the parameter store" |
| `command` | The command that makes it |
| `steps` | Console steps, in order, when there is no command |
| `provisioner` | The tool that makes the value when its command runs: `tool`, `command`, `runBy` (who runs it) and `stores` (it writes the primary store itself) |

With a `provisioner`, the tool generates the value, so `handling.generate`
is `system`, and `runBy` says whether a person, an agent or a system runs
the command. Credentials that share a command are made by one run of it:
the plan and the review page show the command once.

### Where It Lives and Who Uses It

`stores` lists every place the value is kept. Each has a `type` (such as
`aws-ssm`, `github-actions`, `password-manager`, `host-file`), a `path`,
an `account`, a `console` and a `role`:

| Role | Meaning |
|---|---|
| `primary` | The source of truth. Every credential needs one, unless federated |
| `backup` | An offline copy, such as in a password manager |
| `rendered` | Written from the primary at deploy time, where consumers read it |

`install` says, in a sentence or a command, how the value gets into that
store. The reports put it in each credential's next step, so write it for
the person who will do it. `backup` names what backs the store up, such as
a backup job, when it holds the only copy. The `process-env` store type is
a process's environment, set by whatever starts it, and `none` is a value
never stored: a federated credential, or a token exported in a shell for
one run.

`consumers` lists who reads it: `system`, `component`, and `via`, such as
`env SESSION_KEY`. `configChanges` lists the edits a consumer's
configuration needs before it reads the value, such as a key that
references the variable; turning on an optional credential means making
them.

### What It Can Do

| Field | Meaning |
|---|---|
| `entitlements` | What the issuer grants it: `resource`, `access`, `note`, such as a token's scopes, and `grantedBy`, the tool or policy that grants it, so the entitlement can be checked against it |
| `exposure` | What a holder of the value can do, in a sentence |
| `revocation` | How to revoke it quickly |
| `dependsOn` | IDs it is built from, such as a connection string from a password |
| `reaches` | IDs a holder of this one can read, use or replace, such as every parameter an AWS identity can read |

`dependsOn` and `reaches` together let each report show what a leak
exposes in all, following links through other credentials, and which
credentials expose this one.

### Rotation and Handling

`rotation` says when the value is replaced:

| Policy | Meaning |
|---|---|
| `scheduled` | Every `intervalDays` |
| `on-expiry` | Before the issuer's expiry, recorded with an `expires` event |
| `on-demand` | Only on compromise or need |

`procedure` and `impact` describe how a rotation is done and what it
interrupts. `secondaries` lists the IDs of credentials that exist only
during a rotation, such as the previous value kept for verification; give
each `need: rotation`.

`handling` says who does each stage:

| Actor | Meaning |
|---|---|
| `human` | A person, by hand: a console, or a command they run |
| `agent` | An AI agent acting for an operator, such as a coding agent running a command; the operator approves what it cannot undo |
| `system` | A deterministic program, such as a deploy tool or a CI workflow, with no person or agent involved |

| Stage | Meaning |
|---|---|
| `generate` | Makes the value |
| `distribute` | Gets it from its primary store to its consumers |
| `rotate` | Replaces it |

`automation` names the systems doing the `system` stages. The reports
count how many credentials need a person for at least one stage, how many
need an agent but no person, how many are fully automatic, and how many
rotations a person does in a year. `Handling.Responsible` gives the most
hands-on actor of the three stages.

### Accounts

A person's login that controls other credentials, such as the owner of a
code-hosting organization or a cloud root user, is an entry of kind
`account`. Give it `protections` (such as "hardware security key",
"MFA") and `recovery` (how to regain access, or that it cannot be), and
list in `reaches` what it controls. The plan lists every account's
protections, and a `verified` event ticks them off.

## Prerequisites

`prerequisites`, beside `credentials` at the top of the registry, lists
what is set up at an issuer before some credentials can be made, and
holds no secret: a cloud project, an OAuth consent screen, a webhook
endpoint. Each has an `id` (unique among credentials and prerequisites),
`name`, `system`, `environment`, `status` (`planned` until it is set up,
then `active`), `need`, `issuer`, `steps` or a `command`, `done` (how to
tell it is set up), `requires` and `links`.

A credential lists the prerequisites it needs in `requires`.
`Registry.SetupOrder` puts every entry after what it requires and what it
depends on, and the plan lists the prerequisites still planned before its
launch checklist.

## Kinds

`kind`, `issuer.type` and `stores[].type` are open sets: any lowercase
slug is valid, and these are the well-known values the reports and the
schema suggest.

- **kind:** `random-secret`, `password`, `private-key`, `certificate`,
  `connection-string`, `oauth-client`, `api-token`,
  `personal-access-token`, `access-key`, `encryption-key`, `ssh-key`,
  `federated`, `account`.
- **issuer type:** `self`, `aws`, `cloudflare`, `github`, `google-cloud`,
  `postgres`, `stripe`, `registrar`, `resend`.
- **store type:** `aws-ssm`, `aws-secrets-manager`, `aws-kms`,
  `github-actions`, `password-manager`, `pulumi-config`, `host-file`,
  `local-file`, `process-env`, `none`.
- **issuer mode:** `test`, `live`.

## Validation

`credspec validate` reports every problem at once, each with its entry's
ID:

- required fields are present and IDs are unique;
- closed sets hold only their values, and open sets are slugs;
- every credential has a `primary` store, unless federated;
- `scheduled` has `intervalDays` above zero, and other policies have none;
- `dependsOn` and `reaches` name existing IDs, never the entry itself;
- `requires` and `rotation.secondaries` name existing credentials or
  prerequisites, and the setup order has no cycle;
- a `provisioner` has a tool, a command and who runs it;
- link and console URLs are absolute `https` URLs.

A registry file refuses unknown fields, so a misspelled field is an error
rather than silently ignored.

It then prints advice: entries that are valid but likely to mislead
whoever sets them up. `--strict` fails on advice too.

- a provisioned credential whose `handling.generate` is not `system`, or
  whose provisioner stores the value while an `install` asks for it by
  hand;
- a credential with no entitlements, unless federated;
- a rotation procedure naming a variable that no credential, consumer or
  configuration change defines;
- a rotation secondary whose `need` is not `rotation`.
