package credspec

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
)

// Need says when a credential or prerequisite must exist.
type Need string

// Needs.
const (
	// NeedLaunch is needed before the system first runs.
	NeedLaunch Need = "launch"
	// NeedLater is needed, but not for launch.
	NeedLater Need = "later"
	// NeedOptional is needed only when an optional part is turned on.
	NeedOptional Need = "optional"
	// NeedRotation exists only while another credential rotates.
	NeedRotation Need = "rotation"
)

// Needs lists every Need, soonest first.
func Needs() []Need { return []Need{NeedLaunch, NeedLater, NeedOptional, NeedRotation} }

// SetUpBy is who acts to make the value: whoever runs the provisioner's
// command when a tool makes it, else whoever generates it.
func (c Credential) SetUpBy() Actor {
	if p := c.Generation.Provisioner; p != nil {
		return p.RunBy
	}
	return c.Handling.Generate
}

// Prerequisite is something at an issuer that is set up before a
// credential can be made, and is not itself a secret: a cloud project, an
// OAuth consent screen, a webhook endpoint.
type Prerequisite struct {
	// ID is stable and unique among credentials and prerequisites.
	ID string `json:"id"`
	// Name is the display name.
	Name string `json:"name"`
	// Description says what it is for.
	Description string `json:"description,omitempty"`
	// System is the system it belongs to.
	System string `json:"system"`
	// Environment is "prod", "staging" and so on; empty when one serves
	// every environment.
	Environment string `json:"environment,omitempty"`
	// Status is planned until it is set up, then active.
	Status Status `json:"status"`
	// Need says when it must exist.
	Need Need `json:"need,omitempty"`
	// Issuer is the service it is set up at.
	Issuer Issuer `json:"issuer"`
	// Steps are how to set it up, in order.
	Steps []string `json:"steps,omitempty"`
	// Command sets it up, when one command does.
	Command string `json:"command,omitempty"`
	// Done says how to tell it is set up.
	Done string `json:"done,omitempty"`
	// Requires lists the IDs set up before this one.
	Requires []string `json:"requires,omitempty"`
	// Links point to documentation and consoles.
	Links []Link `json:"links,omitempty"`
}

// FindPrerequisite returns the prerequisite with the ID.
func (r *Registry) FindPrerequisite(id string) (Prerequisite, bool) {
	for _, p := range r.Prerequisites {
		if p.ID == id {
			return p, true
		}
	}
	return Prerequisite{}, false
}

func validatePrerequisite(ps *problems, id string, p Prerequisite, ids map[string]bool) {
	if strings.TrimSpace(p.Name) == "" {
		ps.add(id, "name is required")
	}
	if strings.TrimSpace(p.System) == "" {
		ps.add(id, "system is required")
	}
	closed(ps, id, "status", p.Status, Statuses(), true)
	closed(ps, id, "need", p.Need, Needs(), false)
	slug(ps, id, "issuer.type", p.Issuer.Type, true)
	httpsURL(ps, id, "issuer.console", p.Issuer.Console)
	slug(ps, id, "issuer.mode", p.Issuer.Mode, false)
	if len(p.Steps) == 0 && strings.TrimSpace(p.Command) == "" {
		ps.add(id, "needs steps or a command")
	}
	references(ps, id, "requires", p.ID, p.Requires, ids)
	for i, l := range p.Links {
		if strings.TrimSpace(l.Title) == "" {
			ps.add(id, "links[%d].title is required", i)
		}
		httpsURL(ps, id, fmt.Sprintf("links[%d].url", i), l.URL)
	}
}

// references checks that each ID in list is known and is not self.
func references(ps *problems, id, field, self string, list []string, ids map[string]bool) {
	for _, d := range list {
		switch {
		case d == self:
			ps.add(id, "%s names itself", field)
		case !ids[d]:
			ps.add(id, "%s names unknown entry %q", field, d)
		}
	}
}

// SetupStep is one entry of the setup order: a credential or a
// prerequisite.
type SetupStep struct {
	// ID is the entry's ID.
	ID string
	// Prerequisite marks a prerequisite rather than a credential.
	Prerequisite bool
}

// SetupOrder lists every credential and prerequisite so that each comes
// after what it requires and what it depends on, keeping registry order
// otherwise, prerequisites first. It returns an error naming a cycle.
func (r *Registry) SetupOrder() ([]SetupStep, error) {
	type node struct {
		step  SetupStep
		after []string
	}
	var nodes []node
	for _, p := range r.Prerequisites {
		nodes = append(nodes, node{SetupStep{ID: p.ID, Prerequisite: true}, p.Requires})
	}
	for _, c := range r.Credentials {
		nodes = append(nodes, node{SetupStep{ID: c.ID}, append(slices.Clone(c.Requires), c.DependsOn...)})
	}
	index := make(map[string]int, len(nodes))
	for i, n := range nodes {
		index[n.step.ID] = i
	}
	const (
		unseen = iota
		visiting
		done
	)
	mark := make([]int, len(nodes))
	out := make([]SetupStep, 0, len(nodes))
	var visit func(i int, path []string) error
	visit = func(i int, path []string) error {
		switch mark[i] {
		case done:
			return nil
		case visiting:
			return fmt.Errorf("credspec: setup order has a cycle: %s", strings.Join(append(path, nodes[i].step.ID), " -> "))
		}
		mark[i] = visiting
		path = append(path, nodes[i].step.ID)
		for _, a := range nodes[i].after {
			j, ok := index[a]
			if !ok {
				continue // Validate reports unknown IDs.
			}
			if err := visit(j, path); err != nil {
				return err
			}
		}
		mark[i] = done
		out = append(out, nodes[i].step)
		return nil
	}
	for i := range nodes {
		if err := visit(i, nil); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// envVarPattern matches a name like SESSION_KEY or SYSTEMAUTH_DB_DSN:
// capitals and digits with at least one underscore.
var envVarPattern = regexp.MustCompile(`\b[A-Z][A-Z0-9]*(?:_[A-Z0-9]+)+\b`)

// Lint returns advice about the registry: entries that are valid but are
// likely to mislead the person setting them up. Validate must pass first.
func (r *Registry) Lint() []Problem {
	var ps problems
	names := map[string]bool{}
	for _, c := range r.Credentials {
		names[c.Name] = true
		names[c.Name[strings.LastIndex(c.Name, "/")+1:]] = true
		for _, u := range c.Consumers {
			for _, v := range envVarPattern.FindAllString(u.Via, -1) {
				names[v] = true
			}
		}
		for _, ch := range c.ConfigChanges {
			for _, v := range envVarPattern.FindAllString(ch, -1) {
				names[v] = true
			}
		}
	}
	for _, c := range r.Credentials {
		if c.Status == StatusRetired {
			continue
		}
		if p := c.Generation.Provisioner; p != nil && c.Handling.Generate != ActorSystem {
			ps.add(c.ID, "%s makes the value, but handling.generate is %s; set it to system, and provisioner.runBy to who runs the command", p.Tool, c.Handling.Generate)
		}
		if p := c.Generation.Provisioner; p != nil && p.Stores {
			for _, s := range c.Primary() {
				if strings.Contains(s.Install, "<value>") {
					ps.add(c.ID, "%s stores the value, but the %s store's install asks for one by hand", p.Tool, s.Type)
				}
			}
		}
		if len(c.Entitlements) == 0 && c.Generation.Method != MethodFederated {
			ps.add(c.ID, "no entitlements: say what a holder of the value may do")
		}
		for _, v := range envVarPattern.FindAllString(c.Rotation.Procedure, -1) {
			if !names[v] {
				ps.add(c.ID, "rotation.procedure names %s, which no credential, consumer or configuration change does", v)
			}
		}
		for _, s := range c.Rotation.Secondaries {
			if sc, ok := r.Find(s); ok && sc.Need != NeedRotation {
				ps.add(c.ID, "rotation secondary %s has need %q, want %q", s, sc.Need, NeedRotation)
			}
		}
	}
	if _, err := r.SetupOrder(); err != nil {
		ps.add("", "%v", err)
	}
	return ps
}
