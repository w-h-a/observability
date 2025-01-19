// Store

export interface StoreState {
	maxMinTime: MaxMinTime;
	services: Array<Service>;
	endpoints: Array<Endpoint>;
	serviceMetrics: Array<ServiceMetricItem>;
	spanMatrix: SpanMatrix;
}

export interface MaxMinTime {
	maxTime: number;
	minTime: number;
}

export interface Service {
	serviceName: string;
	p99: number;
	avgDuration: number;
	numCalls: number;
	callRate: number;
	numErrors: number;
	errorRate: number;
}

export interface Endpoint {
	name: string;
	p50: number;
	p95: number;
	p99: number;
	numCalls: number;
}

export interface ServiceMetricItem {
	timestamp: number;
	p50: number;
	p95: number;
	p99: number;
	numCalls: number;
	callRate: number;
	numErrors: number;
	errorRate: number;
}

export interface SpanMatrix {
	[id: string]: Spans;
}

export interface Spans {
	columns: string[];
	events: Span[];
}

export type Span = [
	number, // time
	string, // spanid
	string, // parentid
	string, // traceid
	string, // service
	string, // name
	string, // kind
	string, // duration
	string[][], // tags
];

// Actions

export enum ActionTypes {
	maxMinTime = "MAX_MIN_TIME",
	servicesSuccess = "SERVICES_SUCCESS",
	servicesFailure = "SERVICES_FAILURE",
	endpointsSuccess = "ENDPOINTS_SUCCESS",
	endpointsFailure = "ENDPOINTS_FAILURE",
	serviceMetricsSuccess = "SERVICE_METRICS_SUCCESS",
	serviceMetricsFailure = "SERVICE_METRICS_FAILURE",
	spanMatrixSuccess = "SPAN_MATRIX_SUCCESS",
	spanMatrixFailure = "SPAN_MATRIX_FAILURE",
}

export type MaxMinTimeAction = {
	type: ActionTypes.maxMinTime;
	payload: MaxMinTime;
};

export type ServicesActionSuccess = {
	type: ActionTypes.servicesSuccess;
	payload: Service[];
};

export type ServicesActionFailure = {
	type: ActionTypes.servicesFailure;
	payload: Service[];
};

export type EndpointsSuccess = {
	type: ActionTypes.endpointsSuccess;
	payload: Endpoint[];
};

export type EndpointsFailure = {
	type: ActionTypes.endpointsFailure;
	payload: Endpoint[];
};

export type ServiceMetricsSuccess = {
	type: ActionTypes.serviceMetricsSuccess;
	payload: ServiceMetricItem[];
};

export type ServiceMetricsFailure = {
	type: ActionTypes.serviceMetricsFailure;
	payload: ServiceMetricItem[];
};

export type SpanMatrixSuccess = {
	type: ActionTypes.spanMatrixSuccess;
	payload: SpanMatrix;
};

export type SpanMatrixFailure = {
	type: ActionTypes.spanMatrixFailure;
	payload: SpanMatrix;
};

export type Action =
	| MaxMinTimeAction
	| ServicesActionSuccess
	| ServicesActionFailure
	| EndpointsSuccess
	| EndpointsFailure
	| ServiceMetricsSuccess
	| ServiceMetricsFailure
	| SpanMatrixSuccess
	| SpanMatrixFailure;
