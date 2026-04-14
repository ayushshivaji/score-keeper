package tests

import (
	"reflect"
	"testing"

	"github.com/ayush-sr/score-keeper/backend/internal/service"
)

// ---------------------------------------------------------------------------
// ComputeStreakStats
// ---------------------------------------------------------------------------

func TestComputeStreakStats_Empty(t *testing.T) {
	cur, win, loss := service.ComputeStreakStats(nil)
	if cur != 0 || win != 0 || loss != 0 {
		t.Errorf("expected all zeros, got cur=%d win=%d loss=%d", cur, win, loss)
	}
}

func TestComputeStreakStats_AllWins(t *testing.T) {
	cur, win, loss := service.ComputeStreakStats([]bool{true, true, true, true})
	if cur != 4 {
		t.Errorf("current streak expected 4, got %d", cur)
	}
	if win != 4 {
		t.Errorf("longest win expected 4, got %d", win)
	}
	if loss != 0 {
		t.Errorf("longest loss expected 0, got %d", loss)
	}
}

func TestComputeStreakStats_AllLosses(t *testing.T) {
	cur, win, loss := service.ComputeStreakStats([]bool{false, false, false})
	if cur != -3 {
		t.Errorf("current streak expected -3, got %d", cur)
	}
	if win != 0 {
		t.Errorf("longest win expected 0, got %d", win)
	}
	if loss != 3 {
		t.Errorf("longest loss expected 3, got %d", loss)
	}
}

func TestComputeStreakStats_Mixed(t *testing.T) {
	// W W L W W W L L W — current streak = 1 (most recent is a win),
	// longest win = 3 (the W W W run), longest loss = 2 (L L run).
	results := []bool{true, true, false, true, true, true, false, false, true}
	cur, win, loss := service.ComputeStreakStats(results)
	if cur != 1 {
		t.Errorf("current streak expected 1, got %d", cur)
	}
	if win != 3 {
		t.Errorf("longest win expected 3, got %d", win)
	}
	if loss != 2 {
		t.Errorf("longest loss expected 2, got %d", loss)
	}
}

func TestComputeStreakStats_CurrentAfterLoss(t *testing.T) {
	// W W W L — most recent loss; current = -1.
	cur, win, _ := service.ComputeStreakStats([]bool{true, true, true, false})
	if cur != -1 {
		t.Errorf("current streak expected -1, got %d", cur)
	}
	if win != 3 {
		t.Errorf("longest win expected 3, got %d", win)
	}
}

// ---------------------------------------------------------------------------
// RecentForm
// ---------------------------------------------------------------------------

func TestRecentForm_NewestFirst(t *testing.T) {
	// Oldest first: W L W L L → newest first limited to 5: L L W L W
	form := service.RecentForm([]bool{true, false, true, false, false}, 5)
	expected := []string{"L", "L", "W", "L", "W"}
	if !reflect.DeepEqual(form, expected) {
		t.Errorf("expected %v, got %v", expected, form)
	}
}

func TestRecentForm_CapsToN(t *testing.T) {
	// Oldest first: T T F T F T T (indices 0..6)
	// Last 5 oldest first (indices 2..6): F T F T T
	// Newest first: T T F T F → W W L W L
	form := service.RecentForm([]bool{true, true, false, true, false, true, true}, 5)
	expected := []string{"W", "W", "L", "W", "L"}
	if !reflect.DeepEqual(form, expected) {
		t.Errorf("expected %v, got %v", expected, form)
	}
}

func TestRecentForm_Empty(t *testing.T) {
	form := service.RecentForm(nil, 5)
	if len(form) != 0 {
		t.Errorf("expected empty slice, got %v", form)
	}
}
