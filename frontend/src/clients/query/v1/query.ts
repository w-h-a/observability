import { Config } from "../../../config/config";
import { IClient } from "../client";

export class Query {
	static async GetServices<T = any>(
		client: IClient,
		start: number,
		end: number,
	): Promise<{ data: T }> {
		try {
			const path = `/services`;

			const query = `?start=${start}&end=${end}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get("baseUrl")}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve services: ${err}`);
			throw err;
		}
	}

	static async GetEndpoints<T = any>(
		client: IClient,
		start: number,
		end: number,
		serviceName: string,
	): Promise<{ data: T }> {
		try {
			const path = `/service/endpoints`;

			const query = `?start=${start}&end=${end}&service=${serviceName}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get("baseUrl")}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve endpoints: ${err}`);
			throw err;
		}
	}

	static async GetServiceMetrics<T = any>(
		client: IClient,
		start: number,
		end: number,
		serviceName: string,
	): Promise<{ data: T }> {
		try {
			const path = `/service/overview`;

			const query = `?start=${start}&end=${end}&step=60&service=${serviceName}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get("baseUrl")}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve service metrics: ${err}`);
			throw err;
		}
	}

	static async GetSpans<T = any>(
		client: IClient,
		start: number,
		end: number,
		filters?: string,
	): Promise<{ data: T }> {
		try {
			const path = `/spans`;

			const query = `?start=${start}&end=${end}`;

			const rsp = await client.get<T>(
				`${Config.GetInstance().get("baseUrl")}${path}${query}`,
			);

			return rsp;
		} catch (err: unknown) {
			console.log(`query client failed to retrieve spans: ${err}`);
			throw err;
		}
	}
}
