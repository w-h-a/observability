import { Dispatch } from "@reduxjs/toolkit";
import { IClient } from "../../clients/query/client";
import {
	Action,
	ActionTypes,
	MaxMinTime,
	SpanMatrix,
	SpanMatrixFailure,
	SpanMatrixSuccess,
} from "../domain";
import { Query } from "../../clients/query/v1/query";

export class SpansUpdater {
	private static initialState = { "0": { columns: [], events: [] } };

	private static errorState = { "0": { columns: [], events: [] } };

	static Spans(
		client: IClient,
		maxMinTime: MaxMinTime,
		filters?: string,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetSpans<SpanMatrix>(
					client,
					maxMinTime.minTime,
					maxMinTime.maxTime,
					filters,
				);

				dispatch<SpanMatrixSuccess>({
					type: ActionTypes.spanMatrixSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<SpanMatrixFailure>({
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
