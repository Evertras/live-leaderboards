import React from "react";

import styles from "./ScoreNumber.module.css";

export interface ScoreNumberProps {
  total: number;
  par: number;
}

const ScoreNumber = ({ total, par }: ScoreNumberProps) => {
  const diff = total - par;

  const classes = [styles.score];

  if (diff < -1) {
    classes.push(styles.scoreEagle);
  } else if (diff === -1) {
    classes.push(styles.scoreBirdie);
  } else if (diff === 1) {
    classes.push(styles.scoreBogey);
  } else if (diff > 1) {
    classes.push(styles.scoreDoubleBogey);
  }

  return (
    <React.Fragment>
      {total > 0 ? (
        <div className={classes.join(" ")}>{total}</div>
      ) : (
        <div></div>
      )}
    </React.Fragment>
  );
};

export default ScoreNumber;
