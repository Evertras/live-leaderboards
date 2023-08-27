import React from "react";

import styles from "./MatchplayLeader.module.css";
import { Round } from "../../lib/client";
import { MatchplayResult } from "../../lib/winner";

export interface MatchplayLeaderProps {
  round: Round;
  currentResult: MatchplayResult;
}

const MatchplayLeader = ({ round, currentResult }: MatchplayLeaderProps) => {
  let text = "All square";

  if (currentResult.currentWinnerPlayerIndex !== null) {
    text = `${
      round.players[currentResult.currentWinnerPlayerIndex].name
    } is up by ${currentResult.upBy}`;
  }

  return (
    <React.Fragment>
      <div className={styles.text}>{text}</div>
    </React.Fragment>
  );
};

export default MatchplayLeader;
