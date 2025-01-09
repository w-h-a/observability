export interface IClient {
	get<T = any>(url: string): Promise<{ data: T }>;
}
