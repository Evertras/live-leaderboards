import React from "react";
import { Round } from "../../lib/api/models";

import styles from "./Scorecard.module.css";

export interface ScorecardProps {
  round: Round;
}

// TODO: Figure out how to automate this better, it's based off API
// spec but the generator doesn't properly make it...
interface PlayerData {
  name: string;
  scores: { hole: number; score: number }[];
}

const Scorecard = ({ round }: ScorecardProps) => {
  const playerData = round.players.map((p: PlayerData) => {
    const holeScores = round.course.holes.map(() => null);

    if (p.scores) {
      for (const score of p.scores) {
        holeScores[score.hole - 1] = score.score;
      }
    } else {
      // Testing
      holeScores[0] = p.name.length;
    }

    return {
      name: p.name,
      scores: holeScores,
    };
  });
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
                {p.scores.map((s: number | null, i: number) => {
                  return <td key={i}>{s}</td>;
                })}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </React.Fragment>
  );
};

export default Scorecard;
