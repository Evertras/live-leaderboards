import React from "react";
import "./App.css";
import { Outlet } from "react-router-dom";

function App() {
  return (
    <React.Fragment>
      <div className="App">
        <header className="App-header">
          <Outlet />
        </header>
      </div>
    </React.Fragment>
  );
}

export default App;
