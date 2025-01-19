import { Config } from "../../../src/config/config";

Config.GetInstance();

import { Provider } from "react-redux";
import { BrowserRouter } from "react-router-dom";
import { render, RenderResult } from "@testing-library/react";
import { TestCase } from "../testcase";
import { TraceGraph } from "../../../src/views/Trace/TraceGraph";
import { store } from "../../../src/updaters/store";
import { IClient } from "../../../src/clients/query/client";
import { Client } from "../../../src/clients/query/mock/client";
import { ClientContext } from "../../../src/clients/query/clientCtx";

const traceId = "8fb945b2df69da477a98886d79c9a26d";

vi.mock("react-router-dom", async () => {
	const mod = await vi.importActual("react-router-dom");
	return {
		...mod,
		useParams: () => ({
			id: traceId,
		}),
	};
});

describe("TraceGraph", () => {
	const mockQueryClient: IClient = new Client();

	const testCases: TestCase[] = [
		{
			When:
				"when: the user goes to the spans page and calls to the backend are a success",
			GetStub: function <T>(url: string): Promise<{ data: T }> {
				return new Promise((resolve) => {
					resolve({
						data: [
							{
								columns: [
									"Time",
									"SpanId",
									"ParentSpanId",
									"TraceId",
									"ServiceName",
									"Name",
									"Kind",
									"StatusCode",
									"Duration",
									"Tags",
								],
								events: [
									[
										1737306333967,
										"100ab2aea075c0d5",
										"bfcf37f4c5b5568c",
										"8fb945b2df69da477a98886d79c9a26d",
										"route",
										"/route",
										"Server",
										"Ok",
										"44705750",
										[
											["http.target", "/route"],
											["http.status_code", "200"],
										],
									],
									[
										1737306333966,
										"bfcf37f4c5b5568c",
										"",
										"8fb945b2df69da477a98886d79c9a26d",
										"frontend",
										"HTTP GET",
										"Client",
										"Ok",
										"45147416",
										[
											["http.status_code", "200"],
											["http.method", "GET"],
											[
												"http.url",
												"http://0.0.0.0:8083/route?dropoff=728%2C326&pickup=947%2C38",
											],
										],
									],
								],
							},
						] as T,
					});
				});
			},
			Then: "then: we render the spans table",
			ClientCalledTimes: 1,
			ClientCalledWith: [
				`http://localhost:4000/api/v1/spans/trace?traceId=${traceId}`,
			],
		},
	];

	testCases.forEach((testCase) => {
		describe(testCase.When, () => {
			let view: RenderResult;

			beforeEach(async () => {
				mockQueryClient.get = vi.fn().mockImplementation(testCase.GetStub);

				view = render(
					<Provider store={store}>
						<BrowserRouter>
							<ClientContext.Provider value={{ queryClient: mockQueryClient }}>
								<TraceGraph />
							</ClientContext.Provider>
						</BrowserRouter>
					</Provider>,
				);
			});

			afterEach(() => {
				vi.clearAllMocks();
			});

			test(testCase.Then, async () => {
				expect(mockQueryClient.get).toHaveBeenCalledTimes(
					testCase.ClientCalledTimes,
				);

				for (let count = 1; count <= testCase.ClientCalledTimes; count++) {
					expect(mockQueryClient.get).toHaveBeenCalledWith(
						testCase.ClientCalledWith[count - 1],
					);
				}

				expect(view.container).toMatchSnapshot();
			});
		});
	});
});
