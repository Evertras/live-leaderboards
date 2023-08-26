import React from "react";
import asiaToride from "../../data/courses/asia-toride-in-west";

const CreateRound = () => {
  const createRound = async () => {
    const response = await fetch("http://localhost:8037/round", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      mode: "cors",
      body: JSON.stringify(asiaToride),
      redirect: "follow",
    });

    return await response.json();
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
