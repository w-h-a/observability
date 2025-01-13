export interface TestCase {
	When: string;
	GetStub<T = any>(url: string): Promise<{ data: T }>;
	Then: string;
	ClientCalledTimes: number;
	ClientCalledWith: string[];
}
