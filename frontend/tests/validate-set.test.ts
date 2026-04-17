import { validateMatchScore } from "@/lib/utils";

// ---------------------------------------------------------------------------
// Valid standard wins (21-X where X < 20)
// ---------------------------------------------------------------------------

describe("validateMatchScore - standard wins", () => {
  const validScores = Array.from({ length: 20 }, (_, i) => i); // 0..19

  test.each(validScores)("21-%d is valid", (loser) => {
    expect(validateMatchScore(21, loser)).toBeNull();
  });

  test.each(validScores)("%d-21 (reversed) is valid", (loser) => {
    expect(validateMatchScore(loser, 21)).toBeNull();
  });
});

// ---------------------------------------------------------------------------
// Valid deuce wins (both >= 20, win by exactly 2)
// ---------------------------------------------------------------------------

describe("validateMatchScore - deuce wins", () => {
  const deuces: [number, number][] = [
    [22, 20], [23, 21], [24, 22], [25, 23], [30, 28], [40, 38],
  ];

  test.each(deuces)("%d-%d is valid deuce", (a, b) => {
    expect(validateMatchScore(a, b)).toBeNull();
  });

  test.each(deuces)("%d-%d (reversed) is valid deuce", (a, b) => {
    expect(validateMatchScore(b, a)).toBeNull();
  });
});

// ---------------------------------------------------------------------------
// Invalid: negative
// ---------------------------------------------------------------------------

describe("validateMatchScore - negative", () => {
  test("negative player1 rejected", () => {
    expect(validateMatchScore(-1, 21)).toBe("Scores must be non-negative");
  });
  test("negative player2 rejected", () => {
    expect(validateMatchScore(21, -5)).toBe("Scores must be non-negative");
  });
});

// ---------------------------------------------------------------------------
// Invalid: ties
// ---------------------------------------------------------------------------

describe("validateMatchScore - ties", () => {
  test.each([0, 5, 10, 20, 21, 25])("tie at %d-%d is invalid", (score) => {
    expect(validateMatchScore(score, score)).toBe("Scores cannot be tied");
  });
});

// ---------------------------------------------------------------------------
// Invalid: winner below 21
// ---------------------------------------------------------------------------

describe("validateMatchScore - winner below 21", () => {
  const cases: [number, number][] = [[20, 5], [18, 3], [5, 0], [20, 18], [19, 17]];
  test.each(cases)("%d-%d invalid (winner < 21)", (a, b) => {
    expect(validateMatchScore(a, b)).toBe("Winner needs at least 21");
  });
});

// ---------------------------------------------------------------------------
// Invalid: win by less than 2
// ---------------------------------------------------------------------------

describe("validateMatchScore - win by less than 2", () => {
  const cases: [number, number][] = [[21, 20], [22, 21], [25, 24], [30, 29]];
  test.each(cases)("%d-%d invalid (margin < 2)", (a, b) => {
    expect(validateMatchScore(a, b)).toBe("Must win by 2");
  });
});

// ---------------------------------------------------------------------------
// Invalid: deuce win by more than 2
// ---------------------------------------------------------------------------

describe("validateMatchScore - deuce win by more than 2", () => {
  const cases: [number, number][] = [[24, 20], [25, 21], [30, 25]];
  test.each(cases)("%d-%d invalid (deuce win by > 2)", (a, b) => {
    expect(validateMatchScore(a, b)).toBe("In deuce, must win by exactly 2");
  });
});

// ---------------------------------------------------------------------------
// Invalid: winner above 21 when loser below 20
// ---------------------------------------------------------------------------

describe("validateMatchScore - winner > 21 with loser < 20", () => {
  const cases: [number, number][] = [[22, 5], [23, 18], [25, 19], [30, 0]];
  test.each(cases)("%d-%d invalid (winner != 21 with loser < 20)", (a, b) => {
    expect(validateMatchScore(a, b)).toBe("Score must be 21 when opponent < 20");
  });
});
