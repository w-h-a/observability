import React from "react";
import ReactDOM from "react-dom/client";
import App from "./App";
import { Provider } from "react-redux";
import { store } from "./store";
import { BrowserRouter } from "react-router-dom";

const root = ReactDOM.createRoot(
	document.getElementById("root") as HTMLElement,
);

root.render(
	<Provider store={store}>
		<React.StrictMode>
			<BrowserRouter basename="/">
				<App />
			</BrowserRouter>
		</React.StrictMode>
	</Provider>,
);
