import { Dispatch } from "@reduxjs/toolkit";
import { IClient } from "../../clients/query/client";
import {
	Action,
	ActionTypes,
	CustomMetric,
	CustomMetricsActionFailure,
	CustomMetricsActionSuccess,
	MaxMinTime,
	SpanMatrix,
	SpanMatrixActionFailure,
	SpanMatrixActionSuccess,
} from "../domain";
import { Query } from "../../clients/query/v1/query";
import { FilteredQuery } from "../../clients/query/filteredQuery";

export class SpansUpdater {
	private static initialState = { "0": { columns: [], events: [] } };

	private static errorState = { "0": { columns: [], events: [] } };

	private static initialCustomMetricsState = [{ timestamp: 0, value: 0 }];

	private static errorCustomMetricsState = [{ timestamp: 0, value: 0 }];

	static Spans(
		client: IClient,
		maxMinTime: MaxMinTime,
		filters?: FilteredQuery,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetSpans<SpanMatrix>(
					client,
					maxMinTime.minTime,
					maxMinTime.maxTime,
					filters,
				);

				dispatch<SpanMatrixActionSuccess>({
					type: ActionTypes.spanMatrixSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<SpanMatrixActionFailure>({
					type: ActionTypes.spanMatrixFailure,
					payload: {},
				});
			}
		};
	}

	static SpanMatrixReducer(
		state: SpanMatrix = SpansUpdater.initialState,
		action: Action,
	): SpanMatrix {
		switch (action.type) {
			case ActionTypes.spanMatrixSuccess:
				return action.payload;
			case ActionTypes.spanMatrixFailure:
				return SpansUpdater.errorState;
			default:
				return state;
		}
	}

	static CustomMetrics(
		client: IClient,
		maxMinTime: MaxMinTime,
		dimension: string,
		aggregation: string,
		filters?: FilteredQuery,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetCustomMetrics<CustomMetric[]>(
					client,
					maxMinTime.minTime,
					maxMinTime.maxTime,
					dimension,
					aggregation,
					filters,
				);

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
		state: CustomMetric[] = SpansUpdater.initialCustomMetricsState,
		action: Action,
	): CustomMetric[] {
		switch (action.type) {
			case ActionTypes.customMetricsSuccess:
				return action.payload;
			case ActionTypes.customMetricsFailure:
				return SpansUpdater.errorCustomMetricsState;
			default:
				return state;
		}
	}
}
