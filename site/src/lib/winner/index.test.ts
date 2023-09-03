import { calcCurrentMatchplayWinner, MatchplayResult } from ".";
import { Round } from "../client";

describe("matchplay 1v1 result calculations", () => {
  const testCases: {
    name: string;
    scoresA: number[];
    scoresB: number[];
    expectedResult: MatchplayResult;
  }[] = [
    {
      name: "empty round",
      scoresA: [],
      scoresB: [],
      expectedResult: {
        totalHolesWonByPlayerIndex: [0, 0],
        currentWinnerPlayerIndex: null,
        upBy: 0,
      },
    },
    {
      name: "tied first",
      scoresA: [3],
      scoresB: [3],
      expectedResult: {
        totalHolesWonByPlayerIndex: [0, 0],
        currentWinnerPlayerIndex: null,
        upBy: 0,
      },
    },
    {
      name: "player A won first",
      scoresA: [3],
      scoresB: [4],
      expectedResult: {
        totalHolesWonByPlayerIndex: [1, 0],
        currentWinnerPlayerIndex: 0,
        upBy: 1,
      },
    },
    {
      name: "player B won first",
      scoresA: [4],
      scoresB: [3],
      expectedResult: {
        totalHolesWonByPlayerIndex: [0, 1],
        currentWinnerPlayerIndex: 1,
        upBy: 1,
      },
    },
    {
      name: "each player wins one",
      scoresA: [4, 4],
      scoresB: [3, 5],
      expectedResult: {
        totalHolesWonByPlayerIndex: [1, 1],
        currentWinnerPlayerIndex: null,
        upBy: 0,
      },
    },
    {
      name: "all holes tied",
      scoresA: [8, 2],
      scoresB: [8, 2],
      expectedResult: {
        totalHolesWonByPlayerIndex: [0, 0],
        currentWinnerPlayerIndex: null,
        upBy: 0,
      },
    },
    {
      name: "all holes won by player B",
      scoresA: [9, 3],
      scoresB: [8, 2],
      expectedResult: {
        totalHolesWonByPlayerIndex: [0, 2],
        currentWinnerPlayerIndex: 1,
        upBy: 2,
      },
    },
  ];

  testCases.forEach((testCase) => {
    test(testCase.name, () => {
      const round = gen1v1Round(testCase.scoresA, testCase.scoresB);
      const result = calcCurrentMatchplayWinner(round);

      expect(result).toEqual(testCase.expectedResult);
    });
  });
});

function gen1v1Round(scoresA: number[], scoresB: number[]): Round {
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
        scores: scoresA.map((n, i) => ({ hole: i + 1, score: n })),
      },
      {
        name: "Player B",
        scores: scoresB.map((n, i) => ({ hole: i + 1, score: n })),
      },
    ],
  };
}
