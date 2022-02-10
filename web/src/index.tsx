import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import App from "./App";
import ChamberFormView from "./components/ChamberFormView";
import ChambersView from "./components/ChambersView";
import ChamberView from "./components/ChamberView";
import Dashboard from "./components/DashboardView";
import SettingsView from "./components/SettingsView";
import "./index.css";
import reportWebVitals from "./reportWebVitals";

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter basename="/ui">
      <Routes>
        <Route path="/" element={<App />}>
          <Route index element={<Dashboard />} />
          <Route path="chambers">
            <Route index element={<ChambersView />} />
            <Route path=":chamberId">
              <Route index element={<ChamberView />} />
              <Route path="edit" element={<ChamberFormView />} />
            </Route>
            <Route path="new" element={<ChamberFormView />} />
          </Route>
          <Route path="settings" element={<SettingsView />} />
        </Route>
      </Routes>
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById("root")
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();
