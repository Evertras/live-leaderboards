import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import reportWebVitals from "./reportWebVitals";
import { createBrowserRouter, RouterProvider } from "react-router-dom";

import App from "./App";
import CreateRound, { createAction } from "./views/CreateRound";
import ViewRound, { loadRound } from "./views/ViewRound";
import UpdateRound from "./views/UpdateRound";

const router = createBrowserRouter([
  {
    path: "/",
    element: <App />,
    children: [
      {
        path: "/create",
        element: <CreateRound />,
        action: createAction,
      },
      {
        path: "/round/:id",
        element: <ViewRound />,
        loader: loadRound,
      },
      {
        path: "/round/:id/update",
        element: <UpdateRound />,
        loader: loadRound,
      },
    ],
  },
]);

const root = ReactDOM.createRoot(
  document.getElementById("root") as HTMLElement,
);
root.render(
  <React.StrictMode>
    <RouterProvider router={router} />
  </React.StrictMode>,
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
