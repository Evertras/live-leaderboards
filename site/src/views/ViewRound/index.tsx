import React, { useEffect } from "react";
import { useLoaderData, useRevalidator } from "react-router-dom";

import { getRoundByID } from "../../lib/client";

import { Round } from "../../lib/api";
import Scorecard from "../../components/Scorecard";

export async function loadRound({ params }: any) {
  const { id } = params;

  return getRoundByID(id);
}

const ViewRound = () => {
  const round = useLoaderData() as Round;

  const revalidator = useRevalidator();

  useEffect(() => {
    const minute = 60 * 1000;
    let id = setInterval(() => revalidator.revalidate(), minute);
    return () => clearInterval(id);
  }, [revalidator]);

  return (
    <React.Fragment>
      <h1>{round.title ? round.title : round.course.name}</h1>
      <Scorecard round={round} />
    </React.Fragment>
  );
};

export default ViewRound;
