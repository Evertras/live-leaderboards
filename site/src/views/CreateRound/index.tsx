import React from "react";
import { Form, redirect } from "react-router-dom";
// TODO: select between courses
//import asiaToride from "../../data/courses/asia-toride-in-west";
import saitamaKokusai from "../../data/courses/saitama-kokusai-south-west";
import { createRound, RoundRequest } from "../../lib/client";

export async function createAction(): Promise<Response> {
  const req: RoundRequest = {
    title: "Battle at Saitama",
    course: saitamaKokusai,
    players: [
      {
        name: "Brandon",
      },
      {
        name: "Ryo",
      },
    ],
  };

  const createdRound = await createRound(req);

  return redirect(`/round/${createdRound.id}/update`);
}

const CreateRound = () => {
  return (
    <React.Fragment>
      <Form method="post">
        <button type="submit">Create Saitama Kokusai round</button>
      </Form>
    </React.Fragment>
  );
};

export default CreateRound;
