import React from "react";
import ReactDOM from "react-dom";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import App from "./App";
import ChamberFormView from "./components/ChamberFormView";
import ChambersView from "./components/ChambersView";
import ChamberView from "./components/ChamberView";
import LoginFormView from "./components/LoginFormView";
import NotFoundView from "./components/NotFoundView";
import PrivateRoute from "./components/PrivateRoute";
// import Dashboard from "./components/DashboardView";
import SettingsFormView from "./components/SettingsFormView";
import "./index.css";
import reportWebVitals from "./reportWebVitals";

// TODO: add 404 not found route

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter basename="/">
      <Routes>
        <Route path="/" element={<App />}>
          <Route path="*" element={<NotFoundView />} />
          {/* <Route index element={<Dashboard />} /> */}
          <Route index element={<Navigate replace to="/chambers" />} />
          {/* <Route index element={<ChambersView />} /> */}
          <Route path="chambers">
            <Route
              index
              element={
                <PrivateRoute>
                  <ChambersView />
                </PrivateRoute>
              }
            />
            <Route path=":chamberId">
              <Route
                index
                element={
                  <PrivateRoute>
                    <ChamberView />
                  </PrivateRoute>
                }
              />
              <Route
                path="edit"
                element={
                  <PrivateRoute>
                    <ChamberFormView />
                  </PrivateRoute>
                }
              />
            </Route>
            <Route
              path="new"
              element={
                <PrivateRoute>
                  <ChamberFormView />
                </PrivateRoute>
              }
            />
          </Route>
          <Route
            path="settings"
            element={
              <PrivateRoute>
                <SettingsFormView />
              </PrivateRoute>
            }
          />
          <Route path="login" element={<LoginFormView />} />
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
