import { Provider } from "react-redux";
import { BrowserRouter } from "react-router-dom";
import { render, RenderResult } from "@testing-library/react";
import { TestCase } from "../testcase";
import { Spans } from "../../../src/views/Spans/Spans";
import { store } from "../../../src/updaters/store";
import { IClient } from "../../../src/clients/query/client";
import { Client } from "../../../src/clients/query/mock/client";
import { ClientContext } from "../../../src/clients/query/clientCtx";

vi.mock("../../../src/updaters/time/utils", () => {
	return {
		Now: vi.fn(() => 1736393070 * 1000),
		StartTime: vi.fn(() => "2025-01-19 12:05:33"),
	};
});

describe("Spans", () => {
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
										"44705750",
										[
											["http.target", "/route"],
											["http.status_code", "200"],
										],
									],
									[
										1737306333966,
										"bfcf37f4c5b5568c",
										"8e45426a4194bb99",
										"8fb945b2df69da477a98886d79c9a26d",
										"frontend",
										"HTTP GET",
										"Client",
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
				"http://localhost:4000/api/v1/spans?start=1736392170&end=1736393070",
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
								<Spans />
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
