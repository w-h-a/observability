import { IClient } from "../client";

export class Client implements IClient {
	async get<T = any>(url: string): Promise<{ data: T }> {
		return {
			data: [] as T,
		};
	}
}
