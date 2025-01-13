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

			return await client.get<T>(
				`${Config.GetInstance().get("baseUrl")}${path}${query}`,
			);
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

			return await client.get<T>(
				`${Config.GetInstance().get("baseUrl")}${path}${query}`,
			);
		} catch (err: unknown) {
			console.log(`query client failed to retrieve endpoints: ${err}`);
			throw err;
		}
	}
}
