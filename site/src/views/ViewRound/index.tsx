import React, { useEffect } from "react";
import { useLoaderData, useRevalidator } from "react-router-dom";

import { getRoundByID, Round } from "../../lib/client";

import Scorecard from "../../components/Scorecard";
import MatchplayLeader from "../../components/MatchplayLeader";
import { calcCurrentMatchplayWinner } from "../../lib/winner";

export async function loadRound({ params }: any) {
  const { id } = params;

  return getRoundByID(id);
}

const ViewRound = () => {
  const round = useLoaderData() as Round;

  const revalidator = useRevalidator();

  useEffect(() => {
    // Only poll for updates if the round is still going
    const anyPlayerStillPlaying = round.players.some(
      (p) => p.scores.length !== round.course.holes.length,
    );
    if (!anyPlayerStillPlaying) {
      return;
    }

    const minute = 60 * 1000;
    const id = setInterval(() => revalidator.revalidate(), minute);
    return () => clearInterval(id);
  }, [revalidator, round.players, round.course.holes]);

  return (
    <React.Fragment>
      <h1>{round.title ? round.title : round.course.name}</h1>
      <Scorecard round={round} />
      <MatchplayLeader
        round={round}
        currentResult={calcCurrentMatchplayWinner(round)}
      />
    </React.Fragment>
  );
};

export default ViewRound;
