// Store

export interface StoreState {
	services: Array<Service>;
	maxMinTime: MaxMinTime;
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

export interface MaxMinTime {
	maxTime: number;
	minTime: number;
}

// Actions

export enum ActionTypes {
	servicesSuccess = "SERVICES_SUCCESS",
	servicesFailure = "SERVICES_FAILURE",
	maxMinTime = "MAX_MIN_TIME",
}

export type ServicesActionSuccess = {
	type: ActionTypes.servicesSuccess;
	payload: Service[];
};

export type ServicesActionFailure = {
	type: ActionTypes.servicesFailure;
	payload: Service[];
};

export type MaxMinTimeAction = {
	type: ActionTypes.maxMinTime;
	payload: MaxMinTime;
};

export type Action =
	| ServicesActionSuccess
	| ServicesActionFailure
	| MaxMinTimeAction;
