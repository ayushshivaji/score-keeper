import { formatDate, formatDateTime, computeStandingsRow, formatStreak } from "@/lib/utils";

// ---------------------------------------------------------------------------
// formatDate
// ---------------------------------------------------------------------------

describe("formatDate", () => {
  test("formats ISO date string", () => {
    const result = formatDate("2026-04-08T14:00:00Z");
    expect(result).toContain("Apr");
    expect(result).toContain("8");
    expect(result).toContain("2026");
  });

  test("handles different dates", () => {
    const result = formatDate("2025-12-25T00:00:00Z");
    expect(result).toContain("Dec");
    expect(result).toContain("25");
    expect(result).toContain("2025");
  });
});

// ---------------------------------------------------------------------------
// formatDateTime
// ---------------------------------------------------------------------------

describe("formatDateTime", () => {
  test("includes time components", () => {
    const result = formatDateTime("2026-04-08T14:30:00Z");
    expect(result).toContain("2026");
    expect(result).toContain("Apr");
  });
});

// ---------------------------------------------------------------------------
// computeStandingsRow
// ---------------------------------------------------------------------------

describe("computeStandingsRow", () => {
  test("computes losses and win rate for a standard record", () => {
    expect(computeStandingsRow(10, 7)).toEqual({ losses: 3, winRate: 70 });
  });

  test("rounds win rate to nearest integer", () => {
    expect(computeStandingsRow(3, 2)).toEqual({ losses: 1, winRate: 67 });
  });

  test("returns 0% win rate when no matches played", () => {
    expect(computeStandingsRow(0, 0)).toEqual({ losses: 0, winRate: 0 });
  });

  test("handles all wins", () => {
    expect(computeStandingsRow(5, 5)).toEqual({ losses: 0, winRate: 100 });
  });

  test("handles all losses", () => {
    expect(computeStandingsRow(4, 0)).toEqual({ losses: 4, winRate: 0 });
  });

  test("clamps negative losses to zero (defensive)", () => {
    expect(computeStandingsRow(2, 5).losses).toBe(0);
  });
});

// ---------------------------------------------------------------------------
// formatStreak
// ---------------------------------------------------------------------------

describe("formatStreak", () => {
  test("positive shows as Wn", () => {
    expect(formatStreak(3)).toBe("W3");
  });
  test("negative shows as Ln", () => {
    expect(formatStreak(-2)).toBe("L2");
  });
  test("zero shows as dash", () => {
    expect(formatStreak(0)).toBe("—");
  });
  test("one win", () => {
    expect(formatStreak(1)).toBe("W1");
  });
  test("one loss", () => {
    expect(formatStreak(-1)).toBe("L1");
  });
});
