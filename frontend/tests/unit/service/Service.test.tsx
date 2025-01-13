import { Provider } from "react-redux";
import { BrowserRouter } from "react-router-dom";
import { render, RenderResult } from "@testing-library/react";
import { TestCase } from "../testcase";
import { Service } from "../../../src/views/Service/Service";
import { store } from "../../../src/updaters/store";
import { IClient } from "../../../src/clients/query/client";
import { Client } from "../../../src/clients/query/mock/client";
import { ClientContext } from "../../../src/clients/query/clientCtx";

vi.mock("../../../src/updaters/time/utils", () => {
	return {
		Now: vi.fn(() => 1736393070 * 1000),
	};
});

const service = "frontend";

vi.mock("react-router-dom", async () => {
	const mod = await vi.importActual("react-router-dom");
	return {
		...mod,
		useParams: () => ({
			service: service,
		}),
	};
});

describe("Service", () => {
	const mockQueryClient: IClient = new Client();

	const testCases: TestCase[] = [
		{
			When:
				"when: the user goes to the service page and calls to the backend are a success",
			GetStub: function <T>(url: string): Promise<{ data: T }> {
				return new Promise((resolve) => {
					resolve({
						data: [
							{
								name: "/config",
								p50: 0.07 * 1000000,
								p95: 0.09 * 1000000,
								p99: 0.1 * 1000000,
								numCalls: 8,
							},
							{
								name: "/dispatch",
								p50: 765.25 * 1000000,
								p95: 904.79 * 1000000,
								p99: 937.2 * 1000000,
								numCalls: 6,
							},
						] as T,
					});
				});
			},
			Then: "then: we render the service-specific tabs and tables",
			ClientCalledTimes: 1,
			ClientCalledWith: [
				`http://localhost:4000/api/v1/service/endpoints?start=1736392170&end=1736393070&service=${service}`,
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
								<Service />
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
