import React from "react";
import { Form, redirect } from "react-router-dom";
import asiaToride from "../../data/courses/asia-toride-in-west";
import { RoundRequest } from "../../lib/api";
import { createRound } from "../../lib/client";

export async function createAction(): Promise<Response> {
  const req: RoundRequest = {
    course: asiaToride,
    players: [
      {
        name: "Brandon",
      },
      {
        name: "Ryo",
      },
    ],
  };

  const id = await createRound(req);

  return redirect(`/round/${id}`);
}

const CreateRound = () => {
  return (
    <React.Fragment>
      <Form method="post">
        <button type="submit">Create Asia Toride round</button>
      </Form>
    </React.Fragment>
  );
};

export default CreateRound;
