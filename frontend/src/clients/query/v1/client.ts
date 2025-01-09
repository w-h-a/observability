import axios, { AxiosInstance } from "axios";
import { IClient } from "../client";

export class Client implements IClient {
	private client: AxiosInstance = axios.create();

	async get<T = any>(url: string): Promise<{ data: T }> {
		const result = await this.client.get<T>(url);
		return result;
	}
}
