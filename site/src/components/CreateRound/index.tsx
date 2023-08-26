import React from "react";
import asiaToride from "../../data/courses/asia-toride-in-west.json";

const CreateRound = () => {
  const createRound = () => {
    return fetch("http://localhost:8037/round", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      mode: "cors",
      body: JSON.stringify(asiaToride),
      redirect: "follow",
    });
  };

  const clickCreate = () => {
    createRound().then((res) => {
      console.log(res);
    });
  };

  return (
    <React.Fragment>
      <button onClick={clickCreate}>Create</button>
    </React.Fragment>
  );
};

export default CreateRound;
