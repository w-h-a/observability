import { Dispatch } from "@reduxjs/toolkit";
import { Action, ActionTypes, MaxMinTime, MaxMinTimeAction } from "../domain";
import { Now } from "./utils";

export class TimeUpdater {
	private static initialState: MaxMinTime = {
		maxTime: Math.floor(Now() / 1000),
		minTime: Math.floor((Now() - 15 * 60 * 1000) / 1000),
	};

	static Time(
		interval: string,
		datetimeRange?: [number, number],
	): (dispatch: Dispatch) => void {
		return (dispatch: Dispatch) => {
			let maxTime = 0;
			let minTime = 0;

			switch (interval) {
				case "15min":
					maxTime = Math.floor(Now() / 1000);
					minTime = Math.floor((Now() - 15 * 60 * 1000) / 1000);
					break;

				case "30min":
					maxTime = Math.floor(Now() / 1000);
					minTime = Math.floor((Now() - 30 * 60 * 1000) / 1000);
					break;

				case "1hr":
					maxTime = Math.floor(Now() / 1000);
					minTime = Math.floor((Now() - 1 * 60 * 60 * 1000) / 1000);
					break;

				case "6hr":
					maxTime = Math.floor(Now() / 1000);
					minTime = Math.floor((Now() - 6 * 60 * 60 * 1000) / 1000);
					break;

				case "1day":
					maxTime = Math.floor(Now() / 1000);
					minTime = Math.floor((Now() - 24 * 60 * 60 * 1000) / 1000);
					break;

				case "1week":
					maxTime = Math.floor(Now() / 1000);
					minTime = Math.floor((Now() - 7 * 24 * 60 * 60 * 1000) / 1000);
					break;

				case "custom":
					if (datetimeRange !== undefined) {
						maxTime = Math.floor(datetimeRange[1] / 1000);
						minTime = Math.floor(datetimeRange[0] / 1000);
					}
					break;

				default:
					console.log("not found matching case");
			}

			dispatch<MaxMinTimeAction>({
				type: ActionTypes.maxMinTime,
				payload: { maxTime: maxTime, minTime: minTime },
			});
		};
	}

	static MaxMinTimeReducer(
		state: MaxMinTime = TimeUpdater.initialState,
		action: Action,
	): MaxMinTime {
		switch (action.type) {
			case ActionTypes.maxMinTime:
				return action.payload;
			default:
				return state;
		}
	}
}
