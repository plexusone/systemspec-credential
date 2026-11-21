package credspec

import (
	"slices"
	"strings"
	"testing"
)

func project(id string) Prerequisite {
	return Prerequisite{
		ID: id, Name: "Cloud project", System: "example-app", Status: StatusPlanned, Need: NeedLaunch,
		Issuer: Issuer{Type: IssuerGoogleCloud}, Steps: []string{"Create the project."}, Done: "It is listed in the console.",
	}
}

func TestSetupFieldsValidate(t *testing.T) {
	c := secret("a")
	c.Need = NeedLaunch
	c.Requires = []string{"p"}
	c.Issuer.Mode = ModeTest
	c.ConfigChanges = []string{"session_key: ${SESSION_KEY}"}
	c.Entitlements = []Entitlement{{Resource: "sessions", Access: "seal", GrantedBy: "the app"}}
	c.Generation.Provisioner = &Provisioner{Tool: "infra secrets", Command: "infra secrets ensure", RunBy: ActorAgent, Stores: true}
	r := NewRegistry(c)
	r.Prerequisites = []Prerequisite{project("p")}
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
}

func TestSetupFieldsReportProblems(t *testing.T) {
	c := secret("a")
	c.Need = "soon"
	c.Requires = []string{"nope", "a"}
	c.Generation.Provisioner = &Provisioner{RunBy: "robot"}
	c.Rotation.Secondaries = []string{"missing"}
	c.ConfigChanges = []string{" "}
	p := project("a")
	p.Steps = nil
	r := NewRegistry(c)
	r.Prerequisites = []Prerequisite{p, {ID: "q"}}
	got := problemsOf(t, r.Validate())
	for _, want := range []string{
		`a: need "soon" is not one of`,
		`a: requires names unknown entry "nope"`,
		`a: requires names itself`,
		`a: generation.provisioner needs a tool and a command`,
		`a: generation.provisioner.runBy "robot" is not one of`,
		`a: rotation.secondaries names unknown entry "missing"`,
		`a: configChanges[0] is empty`,
		`a: duplicate id`,
		`a: needs steps or a command`,
		`q: name is required`,
		`q: status is required`,
	} {
		wantProblem(t, got, want)
	}
}

func TestSetupOrder(t *testing.T) {
	client := secret("client")
	client.Requires = []string{"consent"}
	dsn := secret("dsn")
	dsn.DependsOn = []string{"password"}
	pw := secret("password")
	consent := project("consent")
	consent.Requires = []string{"project"}
	r := NewRegistry(client, dsn, pw)
	r.Prerequisites = []Prerequisite{consent, project("project")}
	got, err := r.SetupOrder()
	if err != nil {
		t.Fatal(err)
	}
	var ids []string
	for _, s := range got {
		ids = append(ids, s.ID)
	}
	if want := []string{"project", "consent", "client", "password", "dsn"}; !slices.Equal(ids, want) {
		t.Errorf("order %v, want %v", ids, want)
	}
	if !got[0].Prerequisite || got[2].Prerequisite {
		t.Errorf("prerequisite marks wrong: %+v", got)
	}
}

func TestSetupOrderFindsCycle(t *testing.T) {
	a, b := secret("a"), secret("b")
	a.Requires = []string{"b"}
	b.Requires = []string{"a"}
	got := problemsOf(t, NewRegistry(a, b).Validate())
	wantProblem(t, got, "setup order has a cycle: a -> b -> a")
}

func TestLint(t *testing.T) {
	prov := secret("prov")
	prov.Name = "DB_DSN"
	prov.Entitlements = []Entitlement{{Resource: "db", Access: "rows"}}
	prov.Generation.Provisioner = &Provisioner{Tool: "infra db-roles", Command: "infra db-roles apply", RunBy: ActorHuman, Stores: true}
	prov.Stores[0].Install = `aws ssm put-parameter --value "<value>"`
	prov.Rotation.Procedure = "Copy SESSION_KEY to OLD_SESSION_KEY, then set DB_DSN."
	prov.Rotation.Secondaries = []string{"old"}
	old := secret("old")
	old.Name = "SESSION_KEY_PREVIOUS"
	old.Entitlements = prov.Entitlements
	ok := secret("ok")
	ok.Entitlements = prov.Entitlements
	ok.Rotation.Procedure = "Set SESSION_KEY_PREVIOUS from SESSION_KEY."
	r := NewRegistry(prov, old, ok)
	if err := r.Validate(); err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, p := range r.Lint() {
		got = append(got, p.String())
	}
	for _, want := range []string{
		"prov: infra db-roles makes the value, but handling.generate is human",
		"prov: infra db-roles stores the value, but the aws-ssm store's install asks for one by hand",
		"prov: rotation.procedure names OLD_SESSION_KEY",
		`prov: rotation secondary old has need "", want "rotation"`,
	} {
		wantProblem(t, got, want)
	}
	for _, g := range got {
		if strings.HasPrefix(g, "ok:") || strings.HasPrefix(g, "old: no entitlements") {
			t.Errorf("unexpected advice %q", g)
		}
	}
	if n := len(got); n != 4 {
		t.Errorf("%d pieces of advice, want 4: %q", n, got)
	}
}
