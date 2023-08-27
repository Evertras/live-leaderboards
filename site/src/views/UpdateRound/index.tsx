import React from "react";
import { useLoaderData } from "react-router-dom";

import { getRoundByID } from "../../lib/client";

import { Round } from "../../lib/api";
import Scorecard from "../../components/Scorecard";

export async function loadRound({ params }: any) {
  const { id } = params;

  return getRoundByID(id);
}

const UpdateRound = () => {
  const round = useLoaderData() as Round;

  return (
    <React.Fragment>
      <h1>{round.title ? round.title : round.course.name}</h1>
      <Scorecard round={round} />
    </React.Fragment>
  );
};

export default UpdateRound;
