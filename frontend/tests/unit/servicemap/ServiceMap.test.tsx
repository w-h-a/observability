import { Provider } from "react-redux";
import { BrowserRouter } from "react-router-dom";
import { render, RenderResult } from "@testing-library/react";
import { TestCase } from "../testcase";
import { ServiceMap } from "../../../src/views/Services/ServiceMap";
import { store } from "../../../src/updaters/store";
import { IClient } from "../../../src/clients/query/client";
import { Client } from "../../../src/clients/query/mock/client";
import { ClientContext } from "../../../src/clients/query/clientCtx";

vi.mock("../../../src/updaters/time/utils", () => {
	return {
		Now: vi.fn(() => 1736393070 * 1000),
	};
});

// the dependency seems to be bundled wrong. it does not work with vitest
const mockForceGraph = vi.fn();
vi.mock("react-force-graph", () => {
	return {
		ForceGraph2D: (props: any) => {
			mockForceGraph(props);
			return <div></div>;
		},
	};
});

describe("ServiceMap", () => {
	const mockQueryClient: IClient = new Client();

	const testCases: TestCase[] = [
		{
			When: "when: the user goes to the service dependency map",
			GetStub: function <T>(url: string): Promise<{ data: T }> {
				if (url.includes("dependencies")) {
					return new Promise((resolve) => {
						resolve({
							data: [
								{
									parent: "testService1",
									child: "testService2",
									callCount: 100,
								},
							] as T,
						});
					});
				}

				return new Promise((resolve) => {
					resolve({
						data: [
							{
								serviceName: "testService1",
								p99: 4000 * 1000000,
								callRate: 0.5,
								errorRate: 0.0,
							},
							{
								serviceName: "testService2",
								p99: 2000 * 1000000,
								callRate: 0.75,
								errorRate: 0.01,
							},
						] as T,
					});
				});
			},
			Then: "then: we render the service dependency map",
			ClientCalledTimes: 2,
			ClientCalledWith: [
				"http://localhost:4000/api/v1/services?start=1736392170&end=1736393070",
				"http://localhost:4000/api/v1/services/dependencies?start=1736392170&end=1736393070",
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
								<ServiceMap />
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

				expect(mockForceGraph).toHaveBeenCalledWith(
					expect.objectContaining({
						graphData: {
							links: [
								{
									source: "testService2",
									target: "testService1",
									value: 0.99,
								},
							],
							nodes: [
								{
									callRate: "0.50",
									errorRate: "0.00",
									group: 1,
									id: "testService1",
									p99: 4000000000,
								},
								{
									callRate: "0.75",
									errorRate: "0.01",
									group: 2,
									id: "testService2",
									p99: 2000000000,
								},
							],
						},
					}),
				);
			});
		});
	});
});
