package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadThreeAndFourDigitClaimIDs(t *testing.T) {
	dir := t.TempDir()
	const ledger = "| ID | Claim |\n|---|---|\n| TOWN-186 | old claim |\n| SAV-CASTCONT-1006 | new claim |\n| SPR16A-FONT-018 | numbered topic |\n"
	if err := os.WriteFile(filepath.Join(dir, "sample.md"), []byte(ledger), 0600); err != nil {
		t.Fatal(err)
	}
	rows, err := load(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"TOWN-186", "SAV-CASTCONT-1006", "SPR16A-FONT-018"}
	if len(rows) != len(want) {
		t.Fatalf("got %d claims, want %d", len(rows), len(want))
	}
	for i := range want {
		if rows[i].id != want[i] {
			t.Fatalf("row %d = %q, want %q", i, rows[i].id, want[i])
		}
	}
}
