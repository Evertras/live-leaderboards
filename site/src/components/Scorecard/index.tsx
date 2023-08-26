import React from "react";
import { Round } from "../../lib/api/models";

import styles from "./Scorecard.module.css";

export interface ScorecardProps {
  round: Round;
}

const Scorecard = ({ round }: ScorecardProps) => {
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
            {round.players.map((p: any, i: number) => (
              <tr key={i}>
                <td>{p.name}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </React.Fragment>
  );
};

export default Scorecard;
