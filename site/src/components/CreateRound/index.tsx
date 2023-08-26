import React from "react";
import asiaToride from "../../data/courses/asia-toride-in-west";
import { Configuration, RoundApi, RoundRequest } from "../../lib/api";

const CreateRound = () => {
  const createRound = async () => {
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

    return response;
  };

  const clickCreate = async () => {
    const response = await createRound();

    console.log(response.id);
  };

  return (
    <React.Fragment>
      <button onClick={clickCreate}>Create</button>
    </React.Fragment>
  );
};

export default CreateRound;
