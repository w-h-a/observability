import { Dispatch } from "@reduxjs/toolkit";
import { IClient } from "../../clients/query/client";
import {
	Action,
	ActionTypes,
	CustomMetric,
	CustomMetricsActionFailure,
	CustomMetricsActionSuccess,
	MaxMinTime,
} from "../domain";
import { Query } from "../../clients/query/v1/query";
import { FilteredQuery } from "../../clients/query/filteredQuery";

export class MetricsUpdater {
	private static initialState = [{ timestamp: 0, value: 0 }];

	private static errorState = [{ timestamp: 0, value: 0 }];

	static CustomMetrics(
		client: IClient,
		maxMinTime: MaxMinTime,
		dimension: string,
		aggregation: string,
		filters?: FilteredQuery,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				let rsp;

				if (dimension.toLowerCase() === "cpu") {
					rsp = await Query.GetMetrics<CustomMetric[]>(
						client,
						maxMinTime.minTime,
						maxMinTime.maxTime,
						dimension,
						aggregation,
						filters,
					);
				} else {
					rsp = await Query.GetAggregation<CustomMetric[]>(
						client,
						maxMinTime.minTime,
						maxMinTime.maxTime,
						dimension,
						aggregation,
						filters,
					);
				}

				dispatch<CustomMetricsActionSuccess>({
					type: ActionTypes.customMetricsSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<CustomMetricsActionFailure>({
					type: ActionTypes.customMetricsFailure,
					payload: [],
				});
			}
		};
	}

	static CustomMetricsReducer(
		state: CustomMetric[] = MetricsUpdater.initialState,
		action: Action,
	): CustomMetric[] {
		switch (action.type) {
			case ActionTypes.customMetricsSuccess:
				return action.payload;
			case ActionTypes.customMetricsFailure:
				return MetricsUpdater.errorState;
			default:
				return state;
		}
	}
}
