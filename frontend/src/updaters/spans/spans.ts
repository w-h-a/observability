import { Dispatch } from "@reduxjs/toolkit";
import { IClient } from "../../clients/query/client";
import {
	Action,
	ActionTypes,
	MaxMinTime,
	Span,
	SpanMatrix,
	SpanMatrixActionFailure,
	SpanMatrixActionSuccess,
	Tag,
	Tree,
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

	static SpansByTraceId(
		client: IClient,
		traceId: string,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetSpansByTraceId<SpanMatrix>(client, traceId);

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

	static SpanMatrixForATraceReducer(
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

	// TODO: rework this and the data structure it assumes
	static SpansToTree = (spans: Span[]): Tree => {
		let tree: Tree = {
			id: "empty",
			startTime: 0,
			name: "default",
			code: "Unset",
			value: 0,
			time: 0,
			tags: [],
			children: [],
		};

		const spanIdToSpan: { [id: string]: Span } = {};

		for (const s of spans) {
			spanIdToSpan[s[1]] = s;
			spanIdToSpan[s[1]][10] = [];
		}

		for (const span of Object.values(spanIdToSpan)) {
			const tags: Tag[] = [];

			for (const tag of span[9]) {
				tags.push({ key: tag[0], value: tag[1] });
			}

			const child: Tree = {
				id: span[1],
				startTime: span[0],
				name: `${span[4]}:${span[5]}`,
				code: span[7],
				value: Number(span[8]),
				time: Number(span[8]),
				tags: tags,
				children: span[10],
			};

			if (span[2]) {
				spanIdToSpan[span[2]][10].push(child);
			} else {
				tree = child;
			}
		}

		return tree;
	};
}
