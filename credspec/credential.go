// Package credspec is a registry of credentials without their
// values: what each credential is for, where and how it is generated,
// what it is entitled to do, where it is stored, who uses it, who owns
// it and how it rotates. An event log records when each was created,
// rotated, verified, given an expiry or revoked, and Assess turns an
// entry and its record into a lifecycle state.
//
// No type in this package holds a secret value.
package credspec

// FormatVersion is the version of the registry and event log formats.
const FormatVersion = "1"

// Credential is one secret in one place for one environment. The staging
// and production copies of a key are two credentials.
type Credential struct {
	// ID is stable and unique, such as "example-app/prod/SESSION_KEY".
	ID string `json:"id"`
	// Name is the display name, often the variable or parameter name.
	Name string `json:"name"`
	// Description says what the credential is for.
	Description string `json:"description,omitempty"`
	// System is the system that owns the credential.
	System string `json:"system"`
	// Environment is "prod", "staging" and so on; empty when one value
	// serves every environment.
	Environment string `json:"environment,omitempty"`
	// Kind is a well-known kind such as KindRandomSecret, or another slug.
	Kind string `json:"kind"`
	// Status is planned, active or retired.
	Status Status `json:"status"`
	// Need says when the credential must exist: for launch, later, only
	// if an optional part is turned on, or only during a rotation.
	Need Need `json:"need,omitempty"`
	// Requires lists the IDs of credentials and prerequisites that are set
	// up before this one, such as the project an OAuth client lives in.
	// DependsOn is for values this one is built from; Requires is only
	// order.
	Requires []string `json:"requires,omitempty"`
	// Sensitivity is how much harm a leak does.
	Sensitivity Sensitivity `json:"sensitivity,omitempty"`
	// Owner is the person or role who creates and rotates it.
	Owner string `json:"owner"`
	// Issuer is the service that issues or accepts the credential.
	Issuer Issuer `json:"issuer"`
	// Generation says how and where the value is made.
	Generation Generation `json:"generation"`
	// Entitlements are what a holder of the value may do.
	Entitlements []Entitlement `json:"entitlements,omitempty"`
	// Stores are where the value is kept.
	Stores []Store `json:"stores,omitempty"`
	// Consumers are who reads the value, and through what.
	Consumers []Consumer `json:"consumers,omitempty"`
	// Rotation says when and how the value changes, and what that breaks.
	Rotation Rotation `json:"rotation"`
	// Handling says which lifecycle steps a person does and which a system
	// does.
	Handling Handling `json:"handling"`
	// Revocation says how to make the value stop working quickly.
	Revocation string `json:"revocation,omitempty"`
	// Exposure says what someone holding a leaked value can do.
	Exposure string `json:"exposure,omitempty"`
	// DependsOn lists the IDs this credential is derived from, such as
	// the password inside a connection string. A holder of this one also
	// holds those.
	DependsOn []string `json:"dependsOn,omitempty"`
	// Reaches lists the IDs a holder of this credential can read, use or
	// replace, such as the parameters a cloud identity may read. With
	// DependsOn, it gives what a leak exposes in all.
	Reaches []string `json:"reaches,omitempty"`
	// Protections are what guard the credential besides its secrecy, such
	// as "MFA: hardware key" or "SSO only".
	Protections []string `json:"protections,omitempty"`
	// Recovery says how to regain access when the value is lost, or that
	// it cannot be.
	Recovery string `json:"recovery,omitempty"`
	// ConfigChanges are the edits a consumer's configuration needs before
	// it reads the credential, such as a key that references the
	// variable. Turning on an optional credential means making them.
	ConfigChanges []string `json:"configChanges,omitempty"`
	// Tags are free labels.
	Tags []string `json:"tags,omitempty"`
	// Links point to documentation and consoles.
	Links []Link `json:"links,omitempty"`
}

// Well-known credential kinds. Other lowercase slugs are allowed.
const (
	KindRandomSecret        = "random-secret"
	KindPassword            = "password"
	KindPrivateKey          = "private-key"
	KindCertificate         = "certificate"
	KindConnectionString    = "connection-string"
	KindOAuthClient         = "oauth-client"
	KindAPIToken            = "api-token"
	KindPersonalAccessToken = "personal-access-token"
	KindAccessKey           = "access-key"
	KindEncryptionKey       = "encryption-key"
	KindSSHKey              = "ssh-key"
	KindFederated           = "federated"
	// KindAccount is a person's login to a service that controls other
	// credentials, such as a cloud root account or a registrar account.
	KindAccount = "account"
)

// Kinds lists the well-known credential kinds.
func Kinds() []string {
	return []string{
		KindRandomSecret, KindPassword, KindPrivateKey, KindCertificate,
		KindConnectionString, KindOAuthClient, KindAPIToken,
		KindPersonalAccessToken, KindAccessKey, KindEncryptionKey,
		KindSSHKey, KindFederated, KindAccount,
	}
}

// Status is where a credential is in its life.
type Status string

// Statuses.
const (
	// StatusPlanned is a credential that is needed but not created yet.
	StatusPlanned Status = "planned"
	// StatusActive is a credential in use.
	StatusActive Status = "active"
	// StatusRetired is a credential no longer used.
	StatusRetired Status = "retired"
)

// Statuses lists every Status.
func Statuses() []Status { return []Status{StatusPlanned, StatusActive, StatusRetired} }

// Sensitivity is how much harm a leaked value does.
type Sensitivity string

// Sensitivities.
const (
	SensitivityLow      Sensitivity = "low"
	SensitivityMedium   Sensitivity = "medium"
	SensitivityHigh     Sensitivity = "high"
	SensitivityCritical Sensitivity = "critical"
)

// Sensitivities lists every Sensitivity, least first.
func Sensitivities() []Sensitivity {
	return []Sensitivity{SensitivityLow, SensitivityMedium, SensitivityHigh, SensitivityCritical}
}

// Issuer is the service that issues or accepts a credential.
type Issuer struct {
	// Type is a well-known issuer such as IssuerGitHub, or another slug.
	Type string `json:"type"`
	// Account is the organization, project or account at the issuer.
	Account string `json:"account,omitempty"`
	// Console is the https URL where it is managed.
	Console string `json:"console,omitempty"`
	// Mode is the issuer's mode for this environment, such as ModeTest
	// for a payment provider's test keys, or another slug.
	Mode string `json:"mode,omitempty"`
}

// Well-known issuer modes. Other lowercase slugs are allowed.
const (
	// ModeTest is an issuer's test or sandbox mode, where nothing is real.
	ModeTest = "test"
	// ModeLive is an issuer's live mode.
	ModeLive = "live"
)

// Well-known issuer types. Other lowercase slugs are allowed.
const (
	// IssuerSelf is a credential the operator generates, such as a
	// random secret or a key pair.
	IssuerSelf        = "self"
	IssuerAWS         = "aws"
	IssuerCloudflare  = "cloudflare"
	IssuerGitHub      = "github"
	IssuerGoogleCloud = "google-cloud"
	IssuerPostgres    = "postgres"
	IssuerStripe      = "stripe"
	IssuerRegistrar   = "registrar"
)

// Issuers lists the well-known issuer types.
func Issuers() []string {
	return []string{IssuerSelf, IssuerAWS, IssuerCloudflare, IssuerGitHub, IssuerGoogleCloud, IssuerPostgres, IssuerStripe, IssuerRegistrar}
}

// Generation says how and where a value is made.
type Generation struct {
	// Method is how the value comes to exist.
	Method Method `json:"method"`
	// Spec describes the value, such as "32 random bytes, hex" or
	// "RSA 3072, PEM".
	Spec string `json:"spec,omitempty"`
	// Format is the value's exact form, such as FormatHex32, for a random
	// value or a key a tool can make. Spec stays the description people
	// read.
	Format Format `json:"format,omitempty"`
	// Where is the place it is generated, such as "the operator's
	// machine" or a console page.
	Where string `json:"where,omitempty"`
	// Command generates the value, when one command does.
	Command string `json:"command,omitempty"`
	// Steps are the manual steps, in order.
	Steps []string `json:"steps,omitempty"`
	// Provisioner is the tool that makes the value, when one does. A
	// person or an agent runs its command; the tool makes the value, so
	// handling.generate is system.
	Provisioner *Provisioner `json:"provisioner,omitempty"`
}

// Provisioner is a tool that makes a credential's value, and often stores
// it and sets it at the issuer, when its command runs.
type Provisioner struct {
	// Tool names the program, such as "infra db-roles".
	Tool string `json:"tool"`
	// Command is the command line to run. Credentials with the same
	// command are made by one run of it.
	Command string `json:"command"`
	// RunBy is who runs the command: a person, an agent, or a system such
	// as a deploy workflow.
	RunBy Actor `json:"runBy"`
	// Stores says the tool writes the value into the primary store itself,
	// so there is nothing to install by hand.
	Stores bool `json:"stores,omitempty"`
}

// Method is how a value comes to exist.
type Method string

// Methods.
const (
	// MethodRandom is random bytes.
	MethodRandom Method = "random"
	// MethodKeypair is a generated private key.
	MethodKeypair Method = "keypair"
	// MethodIssued is minted by the issuer, in its console or API.
	MethodIssued Method = "issued"
	// MethodDerived is built from other credentials, listed in DependsOn.
	MethodDerived Method = "derived"
	// MethodFederated has no stored value: the holder proves its identity
	// to the issuer each time, as with OIDC role assumption.
	MethodFederated Method = "federated"
)

// Format is the exact form of a value a tool makes.
type Format string

// Formats a tool can make.
const (
	// FormatHex32 is 32 random bytes, hex encoded: 64 characters.
	FormatHex32 Format = "hex32"
	// FormatBase64URL32 is 32 random bytes, base64url without padding.
	FormatBase64URL32 Format = "base64url32"
	// FormatP256PEM is a P-256 private key, PKCS #8, PEM.
	FormatP256PEM Format = "p256-pem"
	// FormatEd25519PEM is an Ed25519 private key, PKCS #8, PEM.
	FormatEd25519PEM Format = "ed25519-pem"
	// FormatRSA3072PEM is an RSA 3072-bit private key, PKCS #8, PEM.
	FormatRSA3072PEM Format = "rsa3072-pem"
)

// Formats lists every Format.
func Formats() []Format {
	return []Format{FormatHex32, FormatBase64URL32, FormatP256PEM, FormatEd25519PEM, FormatRSA3072PEM}
}

// Random reports whether f is random bytes rather than a key.
func (f Format) Random() bool { return f == FormatHex32 || f == FormatBase64URL32 }

// Methods lists every Method.
func Methods() []Method {
	return []Method{MethodRandom, MethodKeypair, MethodIssued, MethodDerived, MethodFederated}
}

// Entitlement is one thing a holder of the value may do.
type Entitlement struct {
	// Resource is what it applies to, such as a repository or a path.
	Resource string `json:"resource"`
	// Access is what it may do there, such as "contents: read".
	Access string `json:"access"`
	// Note adds detail.
	Note string `json:"note,omitempty"`
	// GrantedBy names what grants the access, such as the provisioning
	// tool or the policy that sets it, so the entitlement can be checked
	// against it.
	GrantedBy string `json:"grantedBy,omitempty"`
}

// Store is one place a value is kept.
type Store struct {
	// Type is a well-known store such as StoreAWSSSM, or another slug.
	Type string `json:"type"`
	// Role says whether this is the source of truth, a backup or a copy
	// written at deploy time.
	Role Role `json:"role"`
	// Account is the account, organization or vault holding it.
	Account string `json:"account,omitempty"`
	// Path is the parameter, secret, item or file name.
	Path string `json:"path,omitempty"`
	// Console is the https URL where the store is managed.
	Console string `json:"console,omitempty"`
	// Install says how the value gets here: the command or the steps a
	// person follows, or the system that writes it.
	Install string `json:"install,omitempty"`
	// Note adds detail.
	Note string `json:"note,omitempty"`
	// Backup names what backs this store up, such as a backup job, when
	// it is the only copy.
	Backup string `json:"backup,omitempty"`
}

// Well-known store types. Other lowercase slugs are allowed.
const (
	StoreAWSSSM            = "aws-ssm"
	StoreAWSSecretsManager = "aws-secrets-manager"
	StoreAWSKMS            = "aws-kms"
	StoreGitHubActions     = "github-actions"
	StorePasswordManager   = "password-manager"
	StorePulumiConfig      = "pulumi-config"
	// StoreHostFile is a file on a server.
	StoreHostFile = "host-file"
	// StoreLocalFile is a file on an operator's machine.
	StoreLocalFile = "local-file"
	// StoreProcessEnv is a process's environment, set by whatever starts
	// it, such as a migration step.
	StoreProcessEnv = "process-env"
	// StoreNone is for a value that is never stored: a federated
	// credential, or a token exported in a shell for one run.
	StoreNone = "none"
)

// StoreTypes lists the well-known store types.
func StoreTypes() []string {
	return []string{
		StoreAWSSSM, StoreAWSSecretsManager, StoreAWSKMS, StoreGitHubActions,
		StorePasswordManager, StorePulumiConfig, StoreHostFile, StoreLocalFile, StoreProcessEnv, StoreNone,
	}
}

// Role is what a store is to a credential.
type Role string

// Roles.
const (
	// RolePrimary is the source of truth.
	RolePrimary Role = "primary"
	// RoleBackup is an offline copy kept for recovery.
	RoleBackup Role = "backup"
	// RoleRendered is a copy written from the primary at deploy time.
	RoleRendered Role = "rendered"
)

// Roles lists every Role.
func Roles() []Role { return []Role{RolePrimary, RoleBackup, RoleRendered} }

// Consumer is one reader of a value.
type Consumer struct {
	// System is the system that reads it.
	System string `json:"system"`
	// Component is the part of the system, such as a service or a job.
	Component string `json:"component,omitempty"`
	// Via is how it reads it, such as an environment variable, a file or
	// a configuration key.
	Via string `json:"via,omitempty"`
	// Note adds detail.
	Note string `json:"note,omitempty"`
}

// Rotation says when and how a value changes.
type Rotation struct {
	// Policy is when it rotates.
	Policy Policy `json:"policy"`
	// IntervalDays is the rotation period, for PolicyScheduled only.
	IntervalDays int `json:"intervalDays,omitempty"`
	// Procedure is how to rotate it.
	Procedure string `json:"procedure,omitempty"`
	// Impact is what breaks, or who notices, when it rotates.
	Impact string `json:"impact,omitempty"`
	// Secondaries lists the IDs of credentials that exist only while this
	// one rotates, such as the previous value kept for verification. Their
	// need is NeedRotation.
	Secondaries []string `json:"secondaries,omitempty"`
}

// Policy is when a credential rotates.
type Policy string

// Policies.
const (
	// PolicyScheduled rotates every IntervalDays.
	PolicyScheduled Policy = "scheduled"
	// PolicyOnExpiry rotates before the expiry the issuer sets.
	PolicyOnExpiry Policy = "on-expiry"
	// PolicyOnDemand rotates only on compromise or need.
	PolicyOnDemand Policy = "on-demand"
)

// Policies lists every Policy.
func Policies() []Policy { return []Policy{PolicyScheduled, PolicyOnExpiry, PolicyOnDemand} }

// Handling says who does each step of a credential's lifecycle.
type Handling struct {
	// Generate is who makes the value: a person running a command or
	// using a console, an agent running a command for an operator, or a
	// system such as a provisioning tool.
	Generate Actor `json:"generate"`
	// Distribute is who gets the value from its primary store to its
	// consumers, such as a deploy tool rendering it onto a server.
	Distribute Actor `json:"distribute"`
	// Rotate is who replaces the value when its policy says so.
	Rotate Actor `json:"rotate"`
	// Automation names the systems that do the automated steps.
	Automation string `json:"automation,omitempty"`
}

// Actor is who does a step.
type Actor string

// Actors.
const (
	// ActorHuman is a person, by hand.
	ActorHuman Actor = "human"
	// ActorAgent is an AI agent acting for an operator, such as a coding
	// agent running a command. The operator approves what it cannot undo.
	ActorAgent Actor = "agent"
	// ActorSystem is a deterministic program, such as a deploy tool or a
	// CI workflow, with no person or agent involved.
	ActorSystem Actor = "system"
)

// Actors lists every Actor, from the most hands-on to the least.
func Actors() []Actor { return []Actor{ActorHuman, ActorAgent, ActorSystem} }

// Stage is a step of a credential's lifecycle.
type Stage string

// Stages.
const (
	StageGenerate   Stage = "generate"
	StageDistribute Stage = "distribute"
	StageRotate     Stage = "rotate"
)

// Stages lists every Stage, in lifecycle order.
func Stages() []Stage { return []Stage{StageGenerate, StageDistribute, StageRotate} }

// By returns who does the stage.
func (h Handling) By(s Stage) Actor {
	switch s {
	case StageGenerate:
		return h.Generate
	case StageDistribute:
		return h.Distribute
	case StageRotate:
		return h.Rotate
	}
	return ""
}

// StagesBy lists the stages an actor does, in lifecycle order.
func (h Handling) StagesBy(a Actor) []Stage {
	var out []Stage
	for _, s := range Stages() {
		if h.By(s) == a {
			out = append(out, s)
		}
	}
	return out
}

// Manual lists the stages a person does, in lifecycle order.
func (h Handling) Manual() []Stage { return h.StagesBy(ActorHuman) }

// Agentic lists the stages an agent does, in lifecycle order.
func (h Handling) Agentic() []Stage { return h.StagesBy(ActorAgent) }

// Automatic reports whether a system does every stage.
func (h Handling) Automatic() bool { return len(h.StagesBy(ActorSystem)) == len(Stages()) }

// Responsible is the most hands-on actor of any stage: a person if one
// does any stage, else an agent if one does, else a system.
func (h Handling) Responsible() Actor {
	for _, a := range Actors() {
		if len(h.StagesBy(a)) > 0 {
			return a
		}
	}
	return ""
}

// Link is a titled https URL.
type Link struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

// Primary returns the credential's primary stores.
func (c Credential) Primary() []Store {
	var out []Store
	for _, s := range c.Stores {
		if s.Role == RolePrimary {
			out = append(out, s)
		}
	}
	return out
}
