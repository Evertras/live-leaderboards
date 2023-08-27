import React from "react";
import styles from "./AddScore.module.css";

export interface AddScoreProps {
  playerName: string;
  playerIndex: number;
  hole: number;
  par: number;
}

const AddScore = ({ playerName, playerIndex, hole, par }: AddScoreProps) => {
  return (
    <React.Fragment>
      <div className={styles.AddScore}>
        {playerName} ({playerIndex}) on {hole} (par {par})
      </div>
    </React.Fragment>
  );
};

export default AddScore;
