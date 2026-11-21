// gocooldown updates Go module dependencies like `go get -u ./...`, but only
// to versions that were published at least -days ago (supply chain cooldown).
// Afterwards it checks every requirement in go.mod against the cooldown.
//
// Run it from the module directory:
//
//	go run ../scripts/gocooldown.go -days 7          # update + enforce cooldown
//	go run ../scripts/gocooldown.go -days 7 -check   # only check, change nothing
//
// Requirements in go.mod that are younger than the cooldown (e.g. pulled in by
// a plain `go get -u`) are downgraded to the newest version that is old enough.
//
// Only the standard library is used, so it runs without being part of a module.
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type module struct {
	Path     string
	Version  string
	Time     *time.Time
	Main     bool
	Indirect bool
	Versions []string
	Update   *module
	Error    *struct{ Err string }
}

func main() {
	days := flag.Int("days", 7, "minimum age in days of a version before it is used")
	checkOnly := flag.Bool("check", false, "only check go.mod against the cooldown, do not update")
	flag.Parse()

	cutoff := time.Now().Add(-time.Duration(*days) * 24 * time.Hour)

	if !*checkOnly {
		if err := update(cutoff, *days); err != nil {
			fmt.Fprintln(os.Stderr, "gocooldown:", err)
			os.Exit(1)
		}
	}
	for i := 0; i < 3; i++ {
		young := check(cutoff, *days)
		if young == nil {
			os.Exit(1)
		}
		if len(young) == 0 {
			return
		}
		if *checkOnly {
			os.Exit(1)
		}
		if err := downgrade(young, cutoff, *days); err != nil {
			fmt.Fprintln(os.Stderr, "gocooldown:", err)
			os.Exit(1)
		}
	}
	os.Exit(1)
}

// downgrade moves each too young module to the newest older stable version
// that respects the cooldown.
func downgrade(young []*module, cutoff time.Time, days int) error {
	var args []string
	for _, y := range young {
		m, err := listModule("-versions", y.Path)
		if err != nil {
			return err
		}
		target := ""
		for i := len(m.Versions) - 1; i >= 0; i-- {
			v := m.Versions[i]
			if v == y.Version || strings.Contains(v, "-") {
				continue
			}
			info, err := listModule(y.Path + "@" + v)
			if err != nil || info.Time == nil || !info.Time.Before(*y.Time) {
				continue
			}
			if !info.Time.After(cutoff) {
				target = v
				break
			}
		}
		if target == "" {
			return fmt.Errorf("no version of %s older than %d days found", y.Path, days)
		}
		fmt.Printf("down  %s %s -> %s\n", y.Path, y.Version, target)
		args = append(args, y.Path+"@"+target)
	}
	if err := run(append([]string{"get"}, args...)...); err != nil {
		return err
	}
	return run("mod", "tidy")
}

func update(cutoff time.Time, days int) error {
	// Modules providing packages imported by ./... (what `go get -u ./...` upgrades).
	out, err := goCmd("list", "-deps", "-f", "{{with .Module}}{{if not .Main}}{{.Path}}{{end}}{{end}}", "./...")
	if err != nil {
		return err
	}
	seen := map[string]bool{}
	var paths []string
	for _, p := range strings.Fields(string(out)) {
		if !seen[p] {
			seen[p] = true
			paths = append(paths, p)
		}
	}
	sort.Strings(paths)

	var args []string
	for _, p := range paths {
		m, err := listModule("-u", "-versions", p)
		if err != nil {
			return err
		}
		if m.Update == nil {
			continue
		}
		v, ok := pick(m, cutoff)
		if !ok {
			fmt.Printf("hold  %s %s (newer %s is younger than %d days)\n", p, m.Version, m.Update.Version, days)
			continue
		}
		fmt.Printf("bump  %s %s -> %s\n", p, m.Version, v)
		args = append(args, p+"@"+v)
	}

	if len(args) > 0 {
		if err := run(append([]string{"get"}, args...)...); err != nil {
			return err
		}
	} else {
		fmt.Println("no Go module updates older than the cooldown")
	}
	return run("mod", "tidy")
}

// pick returns the newest stable version newer than the current one that is
// at least as old as the cutoff, never going beyond what `go get -u` would pick.
func pick(m *module, cutoff time.Time) (string, bool) {
	if m.Update.Time != nil && !m.Update.Time.After(cutoff) {
		return m.Update.Version, true
	}
	if len(m.Versions) == 0 {
		return "", false // pseudo-version only module, latest commit too young
	}
	upd, cur := -1, -1
	for i, v := range m.Versions {
		if v == m.Update.Version {
			upd = i
		}
		if v == m.Version {
			cur = i
		}
	}
	if upd < 0 {
		return "", false
	}
	for i := upd - 1; i > cur; i-- {
		v := m.Versions[i]
		if strings.Contains(v, "-") { // skip pre-releases
			continue
		}
		info, err := listModule(m.Path + "@" + v)
		if err != nil || info.Time == nil {
			continue
		}
		if cur < 0 && m.Time != nil && !info.Time.After(*m.Time) {
			break // current is a pseudo-version newer than this tag
		}
		if !info.Time.After(cutoff) {
			return v, true
		}
	}
	return "", false
}

// check returns the go.mod requirements younger than the cutoff, or nil on error.
func check(cutoff time.Time, days int) []*module {
	out, err := goCmd("mod", "edit", "-json")
	if err != nil {
		fmt.Fprintln(os.Stderr, "gocooldown:", err)
		return nil
	}
	var gomod struct {
		Require []struct{ Path, Version string }
	}
	if err := json.Unmarshal(out, &gomod); err != nil {
		fmt.Fprintln(os.Stderr, "gocooldown:", err)
		return nil
	}
	required := map[string]string{}
	for _, r := range gomod.Require {
		required[r.Path] = r.Version
	}

	out, err = goCmd("list", "-m", "-json", "all")
	if err != nil {
		fmt.Fprintln(os.Stderr, "gocooldown:", err)
		return nil
	}
	young := []*module{}
	dec := json.NewDecoder(bytes.NewReader(out))
	for {
		var m module
		if err := dec.Decode(&m); err == io.EOF {
			break
		} else if err != nil {
			fmt.Fprintln(os.Stderr, "gocooldown:", err)
			return nil
		}
		if v, in := required[m.Path]; !in || v != m.Version || m.Time == nil {
			continue
		}
		if m.Time.After(cutoff) {
			mm := m
			young = append(young, &mm)
			fmt.Printf("YOUNG %s %s published %s (< %d days)\n", m.Path, m.Version, m.Time.Format("2006-01-02"), days)
		}
	}
	if len(young) == 0 {
		fmt.Printf("all go.mod requirements are at least %d days old\n", days)
	}
	return young
}

func listModule(args ...string) (*module, error) {
	out, err := goCmd(append([]string{"list", "-m", "-json"}, args...)...)
	if err != nil {
		return nil, err
	}
	var m module
	if err := json.Unmarshal(out, &m); err != nil {
		return nil, err
	}
	return &m, nil
}

func goCmd(args ...string) ([]byte, error) {
	cmd := exec.Command("go", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go %s: %w", strings.Join(args, " "), err)
	}
	return out, nil
}

func run(args ...string) error {
	fmt.Println("go", strings.Join(args, " "))
	cmd := exec.Command("go", args...)
	cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
	return cmd.Run()
}
