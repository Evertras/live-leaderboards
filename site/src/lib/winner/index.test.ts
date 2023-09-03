import { calcCurrentMatchplayWinner } from ".";
import { Round } from "../client";

describe("matchplay result calculations", () => {
  test("no one is winning an empty round", () => {
    const round = genRound();
    const result = calcCurrentMatchplayWinner(round);

    expect(result.upBy).toEqual(0);
    expect(result.currentWinnerPlayerIndex).toBeNull();
    expect(result.totalHolesWonByPlayerIndex).toEqual(
      new Array(round.players.length).fill(0),
    );
  });
});

function genRound(): Round {
  return {
    id: "test-round",
    title: "Test Round",
    course: {
      name: "Test Course",
      holes: [
        {
          hole: 1,
          par: 4,
          distanceYards: 300,
        },
        {
          hole: 2,
          par: 5,
          distanceYards: 500,
        },
      ],
    },
    players: [
      {
        name: "Player A",
        scores: [],
      },
      {
        name: "Player B",
        scores: [],
      },
    ],
  };
}
