import React from "react";
import { Form, redirect } from "react-router-dom";
import asiaToride from "../../data/courses/asia-toride-in-west";
import { Configuration, RoundApi, RoundRequest } from "../../lib/api";

export async function createAction() {
  console.log("hi");
  const configuration = new Configuration({
    // TODO: Config this on build or use
    // basePath: window.location.origin
    basePath: "http://localhost:8037",
  });

  const api = new RoundApi(configuration);
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

  const response = await api.roundPost({ roundRequest: req });

  console.log("Created round:", response.id);

  return redirect(`/round/${response.id}`);
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
