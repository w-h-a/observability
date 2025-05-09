import { useContext } from "react";
import { Provider } from "react-redux";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";
import { Layout } from "antd";
import { SideNav } from "./views/Nav/SideNav";
import { TopNav } from "./views/Nav/TopNav";
import { ServicesTable } from "./views/Services/ServicesTable";
import { Service } from "./views/Service/Service";
import { ServiceMap } from "./views/Services/ServiceMap";
import { Spans } from "./views/Spans/Spans";
import { TraceGraph } from "./views/Traces/TraceGraph";
import { Explore } from "./views/Explore/Explore";
import { store } from "./updaters/store";
import { ClientContext } from "./clients/query/clientCtx";

export const App = () => {
	return (
		<Provider store={store}>
			<BrowserRouter basename="/">
				<Layout style={{ minHeight: "100vh" }}>
					<SideNav />
					<Layout className="site-layout">
						<Layout.Content style={{ margin: "0 16px" }}>
							<TopNav />
							<ClientContext.Provider value={useContext(ClientContext)}>
								<Switch>
									<Route path="/explore" component={Explore} />
									<Route path="/traces/:id" component={TraceGraph} />
									<Route path="/traces" component={Spans} />
									<Route path="/service-map" component={ServiceMap} />
									<Route path="/application/:service" component={Service} />
									<Route path="/application" component={ServicesTable} />
									<Route
										path="/"
										exact
										render={() => {
											return <Redirect to="/application" />;
										}}
									/>
								</Switch>
							</ClientContext.Provider>
						</Layout.Content>
						<Layout.Footer>Observability {new Date().getFullYear()}</Layout.Footer>
					</Layout>
				</Layout>
			</BrowserRouter>
		</Provider>
	);
};
