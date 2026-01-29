package operations

import (
	"testing"

	"github.com/gurselcakar/arithmego/internal/game"
)

func TestRegistryContainsAllOperations(t *testing.T) {
	// All 12 operations should be registered via init()
	ops := All()
	if len(ops) != 12 {
		t.Errorf("Expected 12 operations, got %d", len(ops))
	}

	// Check each operation is registered
	expectedOps := []string{
		"Addition", "Subtraction", "Multiplication", "Division",
		"Square", "Cube", "Square Root", "Cube Root",
		"Modulo", "Power", "Percentage", "Factorial",
	}

	for _, name := range expectedOps {
		if _, ok := Get(name); !ok {
			t.Errorf("Operation %q not found in registry", name)
		}
	}
}

func TestBasicOperations(t *testing.T) {
	basic := BasicOperations()
	if len(basic) != 4 {
		t.Errorf("Expected 4 basic operations, got %d", len(basic))
	}

	for _, op := range basic {
		if op.Category() != game.CategoryBasic {
			t.Errorf("Operation %q has category %q, expected %q", op.Name(), op.Category(), game.CategoryBasic)
		}
	}
}

func TestPowerOperations(t *testing.T) {
	power := PowerOperations()
	if len(power) != 4 {
		t.Errorf("Expected 4 power operations, got %d", len(power))
	}

	for _, op := range power {
		if op.Category() != game.CategoryPower {
			t.Errorf("Operation %q has category %q, expected %q", op.Name(), op.Category(), game.CategoryPower)
		}
	}
}

func TestAdvancedOperations(t *testing.T) {
	advanced := AdvancedOperations()
	if len(advanced) != 4 {
		t.Errorf("Expected 4 advanced operations, got %d", len(advanced))
	}

	for _, op := range advanced {
		if op.Category() != game.CategoryAdvanced {
			t.Errorf("Operation %q has category %q, expected %q", op.Name(), op.Category(), game.CategoryAdvanced)
		}
	}
}

func TestGetOperation(t *testing.T) {
	op, ok := Get("Addition")
	if !ok {
		t.Fatal("Failed to get Addition operation")
	}
	if op.Name() != "Addition" {
		t.Errorf("Got operation %q, expected Addition", op.Name())
	}

	_, ok = Get("NonExistent")
	if ok {
		t.Error("Expected false for non-existent operation")
	}
}
