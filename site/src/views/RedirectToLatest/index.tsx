import React from "react";
import { Form, redirect } from "react-router-dom";
import { getLatestRoundID } from "../../lib/client";

export async function redirectLatestAction() {
  const id = await getLatestRoundID();
  return redirect(`/round/${id}`);
}

const RedirectToLatest = () => {
  return (
    <React.Fragment>
      <Form method="post">
        <button type="submit">View latest</button>
      </Form>
    </React.Fragment>
  );
};

export default RedirectToLatest;
