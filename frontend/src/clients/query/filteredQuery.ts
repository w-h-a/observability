export enum Operator {
	equals = "equals",
	contains = "contains",
}

export enum SpanKind {
	server = "Server",
	client = "Client",
	default = "",
}

// TODO: add status code
export interface FilteredQuery {
	service?: string;
	operation?: string;
	kind: SpanKind;
	duration?: { min: string; max: string };
	tags: { key: string; value: string; operator: Operator }[];
}
