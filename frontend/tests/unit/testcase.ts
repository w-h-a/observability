import { ReactElement } from "react";

export interface TestCase {
	When: string;
	GetStub<T = any>(url: string): Promise<{ data: T }>;
	Element: ReactElement;
	Then: string;
	ClientCalledTimes: number;
	ClientCalledWith: string[];
}
