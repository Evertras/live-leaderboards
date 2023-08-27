import React from "react";
import { Round } from "../../lib/client";

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

const Scorecard = ({ round, onSelect }: ScorecardProps) => {
  const playerData = round.players.map((p: PlayerData) => {
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
  });

  const holeWinners = round.course.holes.map((_: any, i: number) => {
    const playerScoresForHole: number[] = playerData.map(
      (p: any) => p.scores[i],
    );

    if (playerScoresForHole.some((s: number | null) => s === null)) {
      return null;
    }

    const lowScore = Math.min(...playerScoresForHole);
    const numWithScore = playerScoresForHole.filter(
      (s: number) => s === lowScore,
    ).length;

    if (numWithScore > 1) {
      return null;
    }

    return playerScoresForHole.findIndex((s: number) => s === lowScore);
  });

  const onClick = (playerIndex: number, hole: number) => {
    if (onSelect) {
      onSelect(playerIndex, hole);
    }
  };

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
          <tbody>
            {playerData.map((p: any, i: number) => (
              <tr key={i}>
                <td>{p.name}</td>
                {p.scores.map((s: number | null, h: number) => (
                  <td
                    key={h}
                    className={holeWinners[h] === i ? styles.winner : ""}
                    onClick={() => onClick(i, h + 1)}
                  >
                    {s}
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </React.Fragment>
  );
};

export default Scorecard;
