import { Round } from "../client";

export enum MatchplayNoWinnerResult {
  Undecided = "u",
  AllSquare = "as",
}

export interface MatchplayResult {
  totalHolesWonByPlayerIndex: number[];

  currentWinnerPlayerIndex: number | null;
  upBy: number;
}

export function getHoleWinnerPlayerIndex(
  round: Round,
  hole: number,
): number | MatchplayNoWinnerResult {
  // TODO: This whole function isn't really efficient, but it works
  if (round.players.some((p) => !p.scores)) {
    return MatchplayNoWinnerResult.Undecided;
  }

  const playerScores = round.players
    .map((p) => p.scores.find((s) => s.hole === hole))
    .filter((p) => !!p)
    .map((p) => p!.score);

  if (playerScores.length !== round.players.length) {
    return MatchplayNoWinnerResult.Undecided;
  }

  const lowestScore = Math.min(...playerScores);

  const tiedPlayers = playerScores.filter((s) => s === lowestScore).length;

  if (tiedPlayers === 1) {
    return playerScores.findIndex((s) => s === lowestScore);
  }

  return MatchplayNoWinnerResult.AllSquare;
}

export function calcCurrentMatchplayWinner(round: Round): MatchplayResult {
  // TODO: This whole function isn't really efficient, but it works
  const results = round.course.holes.map((h) => {
    return getHoleWinnerPlayerIndex(round, h.hole);
  });

  const totalHolesWonByPlayerIndex = results.reduce(
    (prev: number[], current) => {
      if (typeof current === "number") {
        prev[current]++;
      }

      return prev;
    },
    new Array(round.players.length).fill(0),
  );

  const highestTotal = Math.max(...totalHolesWonByPlayerIndex);
  const numTied = totalHolesWonByPlayerIndex.filter(
    (s) => s === highestTotal,
  ).length;

  if (numTied > 1) {
    return {
      totalHolesWonByPlayerIndex,
      currentWinnerPlayerIndex: null,
      upBy: 0,
    };
  }

  const currentWinnerPlayerIndex = totalHolesWonByPlayerIndex.findIndex(
    (s) => s === highestTotal,
  );
  const diffs = totalHolesWonByPlayerIndex
    .filter((s) => s !== highestTotal)
    .map((s) => highestTotal - s);

  const upBy = Math.min(...diffs);

  return {
    totalHolesWonByPlayerIndex,
    currentWinnerPlayerIndex,
    upBy,
  };
}
