import React from "react";
import { Round } from "../../lib/client";
import {
  calcCurrentMatchplayWinner,
  getHoleWinnerPlayerIndex,
  MatchplayNoWinnerResult,
  MatchplayResult,
} from "../../lib/winner";
import ScoreNumber from "../ScoreNumber";

import styles from "./Scorecard.module.css";

export interface ScorecardProps {
  round: Round;
  onSelect?: (playerIndex: number, holeNumber: number) => void;
}

// TODO: Figure out how to automate this better, it's based off API
// spec but the generator doesn't properly make it...
interface PlayerData {
  name: string;
  scores: { hole: number; score: number }[];
}

interface PlayerScoresForAllHoles {
  name: string;
  scores: (number | null)[];
}

const Scorecard = ({ round, onSelect }: ScorecardProps) => {
  const playerData: PlayerScoresForAllHoles[] = round.players.map(
    (p: PlayerData) => {
      const holeScores = round.course.holes.map((): number | null => null);

      if (p.scores) {
        for (const score of p.scores) {
          holeScores[score.hole - 1] = score.score;
        }
      }

      return {
        name: p.name,
        scores: holeScores,
      };
    },
  );

  const currentMatchplayResult: MatchplayResult =
    calcCurrentMatchplayWinner(round);
  const holeWinners: (number | MatchplayNoWinnerResult)[] =
    round.course.holes.map((h) => getHoleWinnerPlayerIndex(round, h.hole));

  const onClick = (playerIndex: number, hole: number) => {
    if (onSelect) {
      onSelect(playerIndex, hole);
    }
  };

  const playerScoreRows = playerData.map((p: any, i: number) => (
    <tr key={i}>
      <td
        className={
          currentMatchplayResult.currentWinnerPlayerIndex === i
            ? styles.winner
            : ""
        }
      >
        {p.name}
      </td>
      {p.scores.map((s: number | null, h: number) => (
        <td
          key={h}
          className={holeWinners[h] === i ? styles.winner : ""}
          onClick={() => onClick(i, h + 1)}
        >
          <ScoreNumber total={s ?? 0} par={round.course.holes[h].par} />
        </td>
      ))}
    </tr>
  ));

  return (
    <React.Fragment>
      <div className={styles.Scorecard}>
        <table>
          <thead>
            <tr>
              <th>Hole</th>
              {round.course.holes.map((_: any, i: number) => (
                <th key={i}>{i + 1}</th>
              ))}
            </tr>
            <tr>
              <th>{round.course.tees}</th>
              {round.course.holes.map((h: any, i: number) => (
                <th key={i}>{h.distanceYards}</th>
              ))}
            </tr>
            <tr>
              <th>Par</th>
              {round.course.holes.map((h: any, i: number) => (
                <th key={i}>{h.par}</th>
              ))}
            </tr>
          </thead>
          <tbody>{playerScoreRows}</tbody>
        </table>
      </div>
    </React.Fragment>
  );
};

export default Scorecard;
