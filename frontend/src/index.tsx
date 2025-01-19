import { Config } from "./config/config";
import React from "react";
import { render } from "react-dom";
import { App } from "./App";

Config.GetInstance();

render(
	<React.StrictMode>
		<App />
	</React.StrictMode>,
	document.getElementById("root"),
);
