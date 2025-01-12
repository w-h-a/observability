// Store

export interface StoreState {
	maxMinTime: MaxMinTime;
	services: Array<Service>;
	endpoints: Array<Endpoint>;
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

// Actions

export enum ActionTypes {
	maxMinTime = "MAX_MIN_TIME",
	servicesSuccess = "SERVICES_SUCCESS",
	servicesFailure = "SERVICES_FAILURE",
	endpointsSuccess = "ENDPOINTS_SUCCESS",
	endpointsFailure = "ENDPOINTS_FAILURE",
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

export type Action =
	| MaxMinTimeAction
	| ServicesActionSuccess
	| ServicesActionFailure
	| EndpointsSuccess
	| EndpointsFailure;
