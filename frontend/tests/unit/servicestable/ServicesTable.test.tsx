import { Provider } from "react-redux";
import renderer, { create, act } from "react-test-renderer";
import { IClient } from "../../../src/clients/query/client";
import { Client } from "../../../src/clients/query/mock/client";
import { store } from "../../../src/updaters/store";
import { ServicesTable } from "../../../src/views/Services/ServicesTable";
import { TestCase } from "../testcase";
import { ClientContext } from "../../../src/App";
import { BrowserRouter } from "react-router-dom";

vi.mock("../../../src/updaters/time/utils", () => {
	return {
		Now: vi.fn(() => 1736393070 * 1000),
	};
});

describe("ServicesTable", () => {
	const mockQueryClient: IClient = new Client();

	const testCases: TestCase[] = [
		{
			When: "when: the user goes to the services table and there are services",
			GetStub: function <T>(url: string): Promise<{ data: T }> {
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
			Element: <ServicesTable />,
			Then: "then: we render a table with the services",
			ClientCalledTimes: 1,
			ClientCalledWith: [
				"http://localhost:4000/api/v1/services?start=1736392170&end=1736393070",
			],
		},
		{
			When:
				"when: the user goes to the services table and the query client fails for some reason",
			GetStub: function <T>(url: string): Promise<{ data: T }> {
				return new Promise((_, reject) => {
					reject(new Error("whoops!"));
				});
			},
			Element: <ServicesTable />,
			Then: "then: we render the error table",
			ClientCalledTimes: 1,
			ClientCalledWith: [
				"http://localhost:4000/api/v1/services?start=1736392170&end=1736393070",
			],
		},
	];

	testCases.forEach((testCase) => {
		describe(testCase.When, () => {
			let view: renderer.ReactTestRenderer;

			beforeEach(async () => {
				mockQueryClient.get = vi.fn().mockImplementation(testCase.GetStub);

				await act(async () => {
					view = create(
						<Provider store={store}>
							<BrowserRouter>
								<ClientContext.Provider value={{ queryClient: mockQueryClient }}>
									{testCase.Element}
								</ClientContext.Provider>
							</BrowserRouter>
						</Provider>,
					);
				});
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

				expect(view.toJSON()).toMatchSnapshot();
			});
		});
	});
});
