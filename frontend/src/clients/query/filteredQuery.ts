export enum Operator {
	equals = "equals",
	contains = "contains",
}

// TODO: add status code
export interface FilteredQuery {
	service?: string;
	operation?: string;
	duration?: { min: string; max: string };
	tags: { key: string; value: string; operator: Operator }[];
}
