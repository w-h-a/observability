import { createContext, useContext } from "react";
import { Provider } from "react-redux";
import { BrowserRouter, Redirect, Route, Switch } from "react-router-dom";
import { Layout } from "antd";
import { Content, Footer } from "antd/es/layout/layout";
import { SideNav } from "./views/Nav/SideNav";
import { TopNav } from "./views/Nav/TopNav";
import { ServicesTable } from "./views/Services/ServicesTable";
import { Client } from "./clients/query/v1/client";
import { IClient } from "./clients/query/client";
import { store } from "./updaters/store";

export const ClientContext = createContext<{ queryClient: IClient }>({
	queryClient: new Client(),
});

export const App = () => {
	return (
		<Provider store={store}>
			<BrowserRouter basename="/">
				<Layout style={{ minHeight: "100vh" }}>
					<SideNav />
					<Layout className="site-layout">
						<Content style={{ margin: "0 16px" }}>
							<TopNav />
							<ClientContext.Provider value={useContext(ClientContext)}>
								<Switch>
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
						</Content>
						<Footer>Trace-Blame 2025</Footer>
					</Layout>
				</Layout>
			</BrowserRouter>
		</Provider>
	);
};
