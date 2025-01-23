import { Dispatch } from "@reduxjs/toolkit";
import {
	Action,
	ActionTypes,
	MaxMinTime,
	Service,
	ServiceNamesActionFailure,
	ServiceNamesActionSuccess,
	ServicesActionFailure,
	ServicesActionSuccess,
} from "../domain";
import { IClient } from "../../clients/query/client";
import { Query } from "../../clients/query/v1/query";

export class ServicesUpdater {
	private static initialState = [
		{
			serviceName: "",
			p99: 0,
			avgDuration: 0,
			numCalls: 0,
			callRate: 0,
			numErrors: 0,
			errorRate: 0,
		},
	];

	private static errorState = [
		{
			serviceName: "service fetch failed",
			p99: NaN,
			avgDuration: NaN,
			numCalls: NaN,
			callRate: NaN,
			numErrors: NaN,
			errorRate: NaN,
		},
	];

	static Services(
		client: IClient,
		maxMinTime: MaxMinTime,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetServices<Service[]>(
					client,
					maxMinTime.minTime,
					maxMinTime.maxTime,
				);

				dispatch<ServicesActionSuccess>({
					type: ActionTypes.servicesSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<ServicesActionFailure>({
					type: ActionTypes.servicesFailure,
					payload: [],
				});
			}
		};
	}

	static ServicesReducer(
		state: Array<Service> = ServicesUpdater.initialState,
		action: Action,
	): Array<Service> {
		switch (action.type) {
			case ActionTypes.servicesSuccess:
				return action.payload;
			case ActionTypes.servicesFailure:
				return ServicesUpdater.errorState;
			default:
				return state;
		}
	}

	static ServiceNames(client: IClient): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetServiceNames<string[]>(client);

				dispatch<ServiceNamesActionSuccess>({
					type: ActionTypes.serviceNamesSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<ServiceNamesActionFailure>({
					type: ActionTypes.serviceNamesFailure,
					payload: [],
				});
			}
		};
	}

	static ServiceNamesReducer(
		state: Array<string> = [],
		action: Action,
	): Array<string> {
		switch (action.type) {
			case ActionTypes.serviceNamesSuccess:
				return action.payload;
			case ActionTypes.serviceNamesFailure:
				return [];
			default:
				return state;
		}
	}
}
