import React, { useState } from "react";
import { useLoaderData } from "react-router-dom";

import { getRoundByID } from "../../lib/client";

import { Round } from "../../lib/api";
import Scorecard from "../../components/Scorecard";
import AddScore from "../../components/AddScore";

export async function loadRound({ params }: any) {
  const { id } = params;

  return getRoundByID(id);
}

const UpdateRound = () => {
  const round = useLoaderData() as Round;

  const [selected, setSelected] = useState({ hole: 1, playerIndex: 0 });

  const selectedPlayerName = round.players[selected.playerIndex].name;

  const onSubmitScore = (playerIndex: number, hole: number, score: number) => {
    const playerScores: any[] | null = round.players[playerIndex].scores;

    if (
      playerScores &&
      playerScores.some((s: any) => s.hole === hole && s.score === score)
    ) {
      console.log("Score didn't change, ignoring");
    }

    console.log(
      "Submitting score for player",
      playerIndex,
      "on hole",
      hole,
      "as",
      score,
    );
  };

  return (
    <React.Fragment>
      <h1>{round.title ? round.title : round.course.name}</h1>

      <Scorecard
        round={round}
        onSelect={(playerIndex, holeNumber) =>
          setSelected({ hole: holeNumber, playerIndex: playerIndex })
        }
      />

      <AddScore
        playerName={selectedPlayerName}
        playerIndex={selected.playerIndex}
        hole={selected.hole}
        par={round.course.holes[selected.hole - 1].par}
        onSubmit={onSubmitScore}
      />
    </React.Fragment>
  );
};

export default UpdateRound;
