import { Dispatch } from "@reduxjs/toolkit";
import { IClient } from "../../clients/query/client";
import {
	Action,
	ActionTypes,
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
}
