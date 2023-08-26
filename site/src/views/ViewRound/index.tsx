import React from "react";
import { useLoaderData } from "react-router-dom";
import { getRoundByID } from "../../lib/client";

export async function loadRound({ params }: any) {
  const { id } = params;

  return getRoundByID(id);
}

const Round = () => {
  const round = useLoaderData();

  console.log("In render:", round);

  return (
    <React.Fragment>
      <div>Round here</div>
    </React.Fragment>
  );
};

export default Round;
