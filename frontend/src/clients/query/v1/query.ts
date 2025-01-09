import { Config } from "../../../config/config";
import { IClient } from "../client";

export class Query {
	static async GetServices<T = any>(
		client: IClient,
		start: number,
		end: number,
	): Promise<{ data: T }> {
		try {
			const path = `${Config.GetInstance().get("baseUrl")}/services`;

			const query = `?start=${start}&end=${end}`;

			return await client.get<T>(path + query);
		} catch (err: unknown) {
			console.log(`query client failed: ${err}`);
			throw err;
		}
	}
}
