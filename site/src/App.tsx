import React from "react";
import CreateRound from "./components/CreateRound";
import "./App.css";

function App() {
  return (
    <React.Fragment>
      <div className="App">
        <header className="App-header">
          <CreateRound />
        </header>
      </div>
    </React.Fragment>
  );
}

export default App;
