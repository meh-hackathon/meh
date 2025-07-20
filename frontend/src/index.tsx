/* @refresh reload */
import { render } from "solid-js/web";
import "./index.css";
import { Route, Router } from "@solidjs/router";
import App from "./App";
import Home from "./Home";
import { NotFound } from "./NotFound";
import { QrCodes } from "./QrCodes";

const root = document.getElementById("root");

render(() => (
    <Router root={App}>
      <Route path="/" component={Home} />
      <Route path="/qrcodes" component={QrCodes} />
      <Route path="*paramName" component={NotFound} />
    </Router>
), root!);
