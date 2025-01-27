import { Dispatch } from "@reduxjs/toolkit";
import { IClient } from "../../clients/query/client";
import {
	Action,
	ActionTypes,
	Endpoint,
	EndpointsActionFailure,
	EndpointsActionSuccess,
	MaxMinTime,
	OperationNamesActionFailure,
	OperationNamesActionSuccess,
	ServiceMetricItem,
	ServiceMetricsActionFailure,
	ServiceMetricsActionSuccess,
	TagsActionFailure,
	TagsActionSuccess,
} from "../domain";
import { Query } from "../../clients/query/v1/query";

export class ServiceUpdater {
	private static initialEndpointsState = [
		{
			name: "",
			p50: 0,
			p95: 0,
			p99: 0,
			numCalls: 0,
		},
	];

	private static errorEndpointsState = [
		{
			name: "endpoints fetch failed",
			p50: NaN,
			p95: NaN,
			p99: NaN,
			numCalls: NaN,
		},
	];

	private static initialServiceMetricsState = [
		{
			timestamp: 0,
			p50: 0,
			p95: 0,
			p99: 0,
			numCalls: 0,
			callRate: 0.0,
			numErrors: 0,
			errorRate: 0.0,
		},
	];

	private static errorServiceMetricsState = [
		{
			timestamp: NaN,
			p50: NaN,
			p95: NaN,
			p99: NaN,
			numCalls: NaN,
			callRate: NaN,
			numErrors: NaN,
			errorRate: NaN,
		},
	];

	static Tags(
		client: IClient,
		serviceName: string,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetTags<string[]>(client, serviceName);

				dispatch<TagsActionSuccess>({
					type: ActionTypes.tagsSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<TagsActionFailure>({
					type: ActionTypes.tagsFailure,
					payload: [],
				});
			}
		};
	}

	static TagsReducer(state: Array<string> = [], action: Action): Array<string> {
		switch (action.type) {
			case ActionTypes.tagsSuccess:
				return action.payload;
			case ActionTypes.tagsFailure:
				return [];
			default:
				return state;
		}
	}

	static OperationNames(
		client: IClient,
		serviceName: string,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetOperationNames<string[]>(client, serviceName);

				dispatch<OperationNamesActionSuccess>({
					type: ActionTypes.operationNamesSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<OperationNamesActionFailure>({
					type: ActionTypes.operationNamesFailure,
					payload: [],
				});
			}
		};
	}

	static OperationNamesReducer(
		state: Array<string> = [],
		action: Action,
	): Array<string> {
		switch (action.type) {
			case ActionTypes.operationNamesSuccess:
				return action.payload;
			case ActionTypes.operationNamesFailure:
				return [];
			default:
				return state;
		}
	}

	static Endpoints(
		client: IClient,
		maxMinTime: MaxMinTime,
		serviceName: string,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetEndpoints<Endpoint[]>(
					client,
					maxMinTime.minTime,
					maxMinTime.maxTime,
					serviceName,
				);

				dispatch<EndpointsActionSuccess>({
					type: ActionTypes.endpointsSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<EndpointsActionFailure>({
					type: ActionTypes.endpointsFailure,
					payload: [],
				});
			}
		};
	}

	static EndpointsReducer(
		state: Array<Endpoint> = ServiceUpdater.initialEndpointsState,
		action: Action,
	): Array<Endpoint> {
		switch (action.type) {
			case ActionTypes.endpointsSuccess:
				return action.payload;
			case ActionTypes.endpointsFailure:
				return ServiceUpdater.errorEndpointsState;
			default:
				return state;
		}
	}

	static ServiceMetrics(
		client: IClient,
		maxMinTime: MaxMinTime,
		serviceName: string,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetServiceMetrics<ServiceMetricItem[]>(
					client,
					maxMinTime.minTime,
					maxMinTime.maxTime,
					serviceName,
				);

				dispatch<ServiceMetricsActionSuccess>({
					type: ActionTypes.serviceMetricsSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<ServiceMetricsActionFailure>({
					type: ActionTypes.serviceMetricsFailure,
					payload: [],
				});
			}
		};
	}

	static ServiceMetricsReducer(
		state: Array<ServiceMetricItem> = ServiceUpdater.initialServiceMetricsState,
		action: Action,
	): Array<ServiceMetricItem> {
		switch (action.type) {
			case ActionTypes.serviceMetricsSuccess:
				return action.payload;
			case ActionTypes.serviceMetricsFailure:
				return ServiceUpdater.errorServiceMetricsState;
			default:
				return state;
		}
	}
}
