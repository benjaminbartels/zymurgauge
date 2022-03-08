import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import App from "./App";
import ChamberFormView from "./components/ChamberFormView";
import ChambersView from "./components/ChambersView";
import ChamberView from "./components/ChamberView";
// import Dashboard from "./components/DashboardView";
import SettingsFormView from "./components/SettingsFormView";
import "./index.css";
import reportWebVitals from "./reportWebVitals";

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter basename="/ui">
      <Routes>
        <Route path="/" element={<App />}>
          {/* <Route index element={<Dashboard />} /> */}
          <Route index element={<Navigate replace to="/chambers" />} />
          <Route index element={<ChambersView />} />
          <Route path="chambers">
            <Route index element={<ChambersView />} />
            <Route path=":chamberId">
              <Route index element={<ChamberView />} />
              <Route path="edit" element={<ChamberFormView />} />
            </Route>
            <Route path="new" element={<ChamberFormView />} />
          </Route>
          <Route path="settings" element={<SettingsFormView />} />
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
