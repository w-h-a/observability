// Store

export interface StoreState {
	maxMinTime: MaxMinTime;
	services: Array<Service>;
	serviceNames: Array<string>;
	serviceDependencies: Array<ServiceDependency>;
	tags: Array<string>;
	operationNames: Array<string>;
	endpoints: Array<Endpoint>;
	serviceMetrics: Array<ServiceMetricItem>;
	traces: SpanMatrix;
	spanMatrix: SpanMatrix;
	customMetrics: CustomMetric[];
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

export interface ServiceDependency {
	parent: string;
	child: string;
	callCount: number;
}

export interface GraphNode {
	id: string;
	group: number;
	p99: number;
	callRate: string;
	errorRate: string;
}

export interface GraphLink {
	source: string;
	target: string;
	value: number;
}

export interface Graph {
	nodes: GraphNode[];
	links: GraphLink[];
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
	string, // code
	string, // duration
	string[][], // tags
	Tree[], // children
];

export interface Tree {
	id: string;
	startTime: number;
	name: string;
	code: string;
	value: number;
	time: number;
	tags: Tag[];
	children: Tree[];
}

export interface Tag {
	key: string;
	value: string;
}

export interface CustomMetric {
	timestamp: number;
	value: number;
}

// Actions

export enum ActionTypes {
	maxMinTime = "MAX_MIN_TIME",
	servicesSuccess = "SERVICES_SUCCESS",
	servicesFailure = "SERVICES_FAILURE",
	serviceNamesSuccess = "SERVICE_NAMES_SUCCESS",
	serviceNamesFailure = "SERVICE_NAMES_FAILURE",
	serviceDependenciesSuccess = "SERVICE_DEPENDENCIES_SUCCESS",
	serviceDependenciesFailure = "SERVICE_DEPENDENCIES_FAILURE",
	tagsSuccess = "TAGS_SUCCESS",
	tagsFailure = "TAGS_FAILURE",
	operationNamesSuccess = "OPERATION_NAMES_SUCCESS",
	operationNamesFailure = "OPERATION_NAMES_FAILURE",
	endpointsSuccess = "ENDPOINTS_SUCCESS",
	endpointsFailure = "ENDPOINTS_FAILURE",
	serviceMetricsSuccess = "SERVICE_METRICS_SUCCESS",
	serviceMetricsFailure = "SERVICE_METRICS_FAILURE",
	tracesSuccess = "TRACES_SUCCESS",
	tracesFailure = "TRACES_FAILURE",
	spanMatrixSuccess = "SPAN_MATRIX_SUCCESS",
	spanMatrixFailure = "SPAN_MATRIX_FAILURE",
	customMetricsSuccess = "CUSTOM_METRICS_SUCCESS",
	customMetricsFailure = "CUSTOM_METRICS_FAILURE",
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

export type ServiceNamesActionSuccess = {
	type: ActionTypes.serviceNamesSuccess;
	payload: string[];
};

export type ServiceNamesActionFailure = {
	type: ActionTypes.serviceNamesFailure;
	payload: string[];
};

export type ServiceDependenciesActionSuccess = {
	type: ActionTypes.serviceDependenciesSuccess;
	payload: ServiceDependency[];
};

export type ServiceDependenciesActionFailure = {
	type: ActionTypes.serviceDependenciesFailure;
	payload: ServiceDependency[];
};

export type TagsActionSuccess = {
	type: ActionTypes.tagsSuccess;
	payload: string[];
};

export type TagsActionFailure = {
	type: ActionTypes.tagsFailure;
	payload: string[];
};

export type OperationNamesActionSuccess = {
	type: ActionTypes.operationNamesSuccess;
	payload: string[];
};

export type OperationNamesActionFailure = {
	type: ActionTypes.operationNamesFailure;
	payload: string[];
};

export type EndpointsActionSuccess = {
	type: ActionTypes.endpointsSuccess;
	payload: Endpoint[];
};

export type EndpointsActionFailure = {
	type: ActionTypes.endpointsFailure;
	payload: Endpoint[];
};

export type ServiceMetricsActionSuccess = {
	type: ActionTypes.serviceMetricsSuccess;
	payload: ServiceMetricItem[];
};

export type ServiceMetricsActionFailure = {
	type: ActionTypes.serviceMetricsFailure;
	payload: ServiceMetricItem[];
};

export type TracesActionSuccess = {
	type: ActionTypes.tracesSuccess;
	payload: SpanMatrix;
};

export type TracesActionFailure = {
	type: ActionTypes.tracesFailure;
	payload: SpanMatrix;
};

export type SpanMatrixActionSuccess = {
	type: ActionTypes.spanMatrixSuccess;
	payload: SpanMatrix;
};

export type SpanMatrixActionFailure = {
	type: ActionTypes.spanMatrixFailure;
	payload: SpanMatrix;
};

export type CustomMetricsActionSuccess = {
	type: ActionTypes.customMetricsSuccess;
	payload: CustomMetric[];
};

export type CustomMetricsActionFailure = {
	type: ActionTypes.customMetricsFailure;
	payload: CustomMetric[];
};

export type Action =
	| MaxMinTimeAction
	| ServicesActionSuccess
	| ServicesActionFailure
	| ServiceNamesActionSuccess
	| ServiceNamesActionFailure
	| ServiceDependenciesActionSuccess
	| ServiceDependenciesActionFailure
	| TagsActionSuccess
	| TagsActionFailure
	| OperationNamesActionSuccess
	| OperationNamesActionFailure
	| EndpointsActionSuccess
	| EndpointsActionFailure
	| ServiceMetricsActionSuccess
	| ServiceMetricsActionFailure
	| TracesActionSuccess
	| TracesActionFailure
	| SpanMatrixActionSuccess
	| SpanMatrixActionFailure
	| CustomMetricsActionSuccess
	| CustomMetricsActionFailure;
