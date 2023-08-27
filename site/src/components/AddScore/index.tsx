import React, { useState } from "react";
import styles from "./AddScore.module.css";

export interface AddScoreProps {
  playerName: string;
  playerIndex: number;
  hole: number;
  par: number;
  initialScore?: number;
  onSubmit: (playerIndex: number, hole: number, score: number) => void;
}

const AddScore = ({
  playerName,
  playerIndex,
  hole,
  par,
  initialScore,
  onSubmit,
}: AddScoreProps) => {
  const [score, setScore] = useState(initialScore ?? par);

  const increment = () => setScore(Math.min(20, score + 1));
  const decrement = () => setScore(Math.max(1, score - 1));

  return (
    <React.Fragment>
      <div className={styles.AddScore}>
        <div className={styles.header}>
          {playerName} on {hole} (par {par})
        </div>
        <div className={styles.adjustment}>
          <div
            className={[styles.button, styles.decrement].join(" ")}
            onClick={decrement}
          >
            -
          </div>
          <div>{score}</div>
          <div
            className={[styles.button, styles.increment].join(" ")}
            onClick={increment}
          >
            +
          </div>
        </div>
        <div
          className={[styles.button, styles.submit].join(" ")}
          onClick={() => onSubmit(playerIndex, hole, score)}
        >
          Submit
        </div>
      </div>
    </React.Fragment>
  );
};

export default AddScore;
