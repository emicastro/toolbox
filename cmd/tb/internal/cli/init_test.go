package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInit_EmptyRepoCreatesEverything(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	stdout, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "rust-systems"}) })
	if code != 0 {
		t.Fatalf("Init() code = %d, stderr = %q, want 0", code, stderr)
	}
	if !strings.Contains(stdout, "wrote: "+filepath.Join(repo, "toolbox.toml")) {
		t.Errorf("stdout = %q, want wrote: for new toolbox.toml", stdout)
	}
	if !strings.Contains(stdout, "wrote: "+filepath.Join(repo, "AGENTS.md")) {
		t.Errorf("stdout = %q, want wrote: for new AGENTS.md", stdout)
	}
	if strings.Contains(stdout, "unchanged:") {
		t.Errorf("stdout = %q, want no unchanged: on first init", stdout)
	}

	for _, rel := range []string{"toolbox.toml", "AGENTS.md", "docs/requirements.md", "docs/adr/0000-template.md"} {
		if _, err := os.Stat(filepath.Join(repo, rel)); err != nil {
			t.Errorf("expected %s to exist: %v", rel, err)
		}
	}

	pointer, err := os.ReadFile(filepath.Join(repo, "toolbox.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pointer), `profile = "rust-systems"`) {
		t.Errorf("toolbox.toml = %q, want it to contain profile = \"rust-systems\"", pointer)
	}

	agentsMD, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agentsMD), "Profile: rust-systems") {
		t.Errorf("AGENTS.md = %q, want it to contain the rendered profile", agentsMD)
	}
	if !strings.Contains(string(agentsMD), "Skills: spec, adr") {
		t.Errorf("AGENTS.md = %q, want it to contain the rendered skill list", agentsMD)
	}
}

func TestInit_RequiresGitRepo(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	notARepo := t.TempDir()
	t.Chdir(notARepo)

	_, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "rust-systems"}) })
	if code != 1 {
		t.Errorf("Init() code = %d, want 1 outside a git repository", code)
	}
	if !strings.Contains(stderr, "git repository") {
		t.Errorf("stderr = %q, want it to mention git repository", stderr)
	}
}

func TestInit_RequiresProfileFlag(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	_, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, nil) })
	if code != 2 {
		t.Errorf("Init() code = %d, want 2 (usage error) when -p is missing", code)
	}
	if !strings.Contains(stderr, "back-go") {
		t.Errorf("stderr = %q, want usage to list back-go", stderr)
	}
	if !strings.Contains(stderr, "game-bevy") {
		t.Errorf("stderr = %q, want usage to list game-bevy", stderr)
	}
	if !strings.Contains(stderr, "cpp-systems") {
		t.Errorf("stderr = %q, want usage to list cpp-systems", stderr)
	}
}

func TestInit_BackGoProfile(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	_, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "back-go"}) })
	if code != 0 {
		t.Fatalf("Init() code = %d, stderr = %q, want 0", code, stderr)
	}

	pointer, err := os.ReadFile(filepath.Join(repo, "toolbox.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pointer), `profile = "back-go"`) {
		t.Errorf("toolbox.toml = %q, want it to contain profile = \"back-go\"", pointer)
	}
	if !strings.Contains(string(pointer), `toolbox_version = "1.4.0"`) {
		t.Errorf("toolbox.toml = %q, want toolbox_version = \"1.4.0\"", pointer)
	}

	agentsMD, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agentsMD), "Profile: back-go") {
		t.Errorf("AGENTS.md = %q, want it to contain the rendered profile", agentsMD)
	}
}

func TestInit_GameBevyProfile(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	_, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "game-bevy"}) })
	if code != 0 {
		t.Fatalf("Init() code = %d, stderr = %q, want 0", code, stderr)
	}

	pointer, err := os.ReadFile(filepath.Join(repo, "toolbox.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pointer), `profile = "game-bevy"`) {
		t.Errorf("toolbox.toml = %q, want it to contain profile = \"game-bevy\"", pointer)
	}
	if !strings.Contains(string(pointer), `toolbox_version = "1.4.0"`) {
		t.Errorf("toolbox.toml = %q, want toolbox_version = \"1.4.0\"", pointer)
	}

	agentsMD, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agentsMD), "Profile: game-bevy") {
		t.Errorf("AGENTS.md = %q, want it to contain the rendered profile", agentsMD)
	}
}

func TestInit_CppSystemsProfile(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	_, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "cpp-systems"}) })
	if code != 0 {
		t.Fatalf("Init() code = %d, stderr = %q, want 0", code, stderr)
	}

	pointer, err := os.ReadFile(filepath.Join(repo, "toolbox.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pointer), `profile = "cpp-systems"`) {
		t.Errorf("toolbox.toml = %q, want it to contain profile = \"cpp-systems\"", pointer)
	}
	if !strings.Contains(string(pointer), `toolbox_version = "1.4.0"`) {
		t.Errorf("toolbox.toml = %q, want toolbox_version = \"1.4.0\"", pointer)
	}

	agentsMD, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agentsMD), "Profile: cpp-systems") {
		t.Errorf("AGENTS.md = %q, want it to contain the rendered profile", agentsMD)
	}
}

func TestInit_UnknownProfileFails(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	_, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "does-not-exist"}) })
	if code != 1 {
		t.Errorf("Init() code = %d, want 1 for an unknown profile", code)
	}
	if !strings.Contains(stderr, `profile "does-not-exist" not found`) {
		t.Errorf("stderr = %q, want the profile-not-found message", stderr)
	}
}

func TestInit_NeverTouchesExistingADRs(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	existingADR := filepath.Join(repo, "docs", "adr", "0001-something.md")
	writeFile(t, existingADR, "# 0001. Something already decided\n\nDo not touch me.\n")
	before, err := os.ReadFile(existingADR)
	if err != nil {
		t.Fatal(err)
	}

	// First init (creates everything else), then a --force re-run, which
	// is the most aggressive mode init has — the ADR must still survive.
	if _, _, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "rust-systems"}) }); code != 0 {
		t.Fatalf("first Init() code = %d, want 0", code)
	}
	if _, _, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "rust-systems", "--force"}) }); code != 0 {
		t.Fatalf("forced Init() code = %d, want 0", code)
	}

	after, err := os.ReadFile(existingADR)
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Errorf("existing ADR changed:\nbefore = %q\nafter  = %q", before, after)
	}
}

func TestInit_RerunWithoutForceRefusesWhenSomethingWouldChange(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	if _, _, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "rust-systems"}) }); code != 0 {
		t.Fatalf("first Init() code = %d, want 0", code)
	}
	agentsBefore, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	pointerBefore, err := os.ReadFile(filepath.Join(repo, "toolbox.toml"))
	if err != nil {
		t.Fatal(err)
	}

	devNull, err := os.Open(os.DevNull)
	if err != nil {
		t.Fatal(err)
	}
	defer devNull.Close()
	oldStdin := os.Stdin
	os.Stdin = devNull
	t.Cleanup(func() { os.Stdin = oldStdin })

	// A different profile means the rendered AGENTS.md/toolbox.toml would
	// actually change, so this must refuse non-interactively.
	_, stderr, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "infra-go"}) })
	if code != 1 {
		t.Errorf("re-run Init() code = %d, want 1 (non-interactive refuse)", code)
	}
	if !strings.Contains(stderr, "refused") {
		t.Errorf("stderr = %q, want it to say refused", stderr)
	}

	agentsAfter, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	pointerAfter, err := os.ReadFile(filepath.Join(repo, "toolbox.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if string(agentsBefore) != string(agentsAfter) {
		t.Errorf("AGENTS.md changed despite refusal:\nbefore = %q\nafter  = %q", agentsBefore, agentsAfter)
	}
	if string(pointerBefore) != string(pointerAfter) {
		t.Errorf("toolbox.toml changed despite refusal:\nbefore = %q\nafter  = %q", pointerBefore, pointerAfter)
	}
}

func TestInit_ForceIdempotentReportsUnchanged(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	if _, _, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "rust-systems"}) }); code != 0 {
		t.Fatalf("first Init() code = %d, want 0", code)
	}

	stdout, stderr, code := captureOutput(t, func() int {
		return Init(toolboxHome, []string{"-p", "rust-systems", "--force"})
	})
	if code != 0 {
		t.Fatalf("forced Init() code = %d, stderr = %q, want 0", code, stderr)
	}
	if !strings.Contains(stdout, "unchanged: "+filepath.Join(repo, "toolbox.toml")) {
		t.Errorf("stdout = %q, want unchanged: for toolbox.toml when bytes match", stdout)
	}
	if !strings.Contains(stdout, "unchanged: "+filepath.Join(repo, "AGENTS.md")) {
		t.Errorf("stdout = %q, want unchanged: for AGENTS.md when bytes match", stdout)
	}
	if strings.Contains(stdout, "wrote: "+filepath.Join(repo, "toolbox.toml")) {
		t.Errorf("stdout = %q, want no wrote: for unchanged toolbox.toml", stdout)
	}
	if strings.Contains(stdout, "wrote: "+filepath.Join(repo, "AGENTS.md")) {
		t.Errorf("stdout = %q, want no wrote: for unchanged AGENTS.md", stdout)
	}
}

func TestInit_ForceRefreshesOnlyAgentsAndPointer(t *testing.T) {
	toolboxHome := newFixtureToolboxHome(t)
	repo := newGitRepo(t)
	t.Chdir(repo)

	if _, _, code := captureOutput(t, func() int { return Init(toolboxHome, []string{"-p", "rust-systems"}) }); code != 0 {
		t.Fatalf("first Init() code = %d, want 0", code)
	}

	reqPath := filepath.Join(repo, "docs", "requirements.md")
	reqBefore, err := os.ReadFile(reqPath)
	if err != nil {
		t.Fatal(err)
	}

	stdout, stderr, code := captureOutput(t, func() int {
		return Init(toolboxHome, []string{"-p", "infra-go", "--force"})
	})
	if code != 0 {
		t.Fatalf("forced Init() code = %d, stderr = %q, want 0", code, stderr)
	}
	if !strings.Contains(stdout, "wrote: "+filepath.Join(repo, "toolbox.toml")) {
		t.Errorf("stdout = %q, want wrote: for toolbox.toml when the profile changes", stdout)
	}
	if !strings.Contains(stdout, "wrote: "+filepath.Join(repo, "AGENTS.md")) {
		t.Errorf("stdout = %q, want wrote: for AGENTS.md when the profile changes", stdout)
	}

	pointer, err := os.ReadFile(filepath.Join(repo, "toolbox.toml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(pointer), `profile = "infra-go"`) {
		t.Errorf("toolbox.toml = %q, want the forced profile infra-go", pointer)
	}

	agentsMD, err := os.ReadFile(filepath.Join(repo, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(agentsMD), "Profile: infra-go") {
		t.Errorf("AGENTS.md = %q, want the forced profile infra-go", agentsMD)
	}

	reqAfter, err := os.ReadFile(reqPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(reqBefore) != string(reqAfter) {
		t.Errorf("docs/requirements.md changed under --force, want it left alone:\nbefore = %q\nafter  = %q", reqBefore, reqAfter)
	}
}
