import { maxBy, uniq, uniqBy } from "lodash";
import { Dispatch } from "@reduxjs/toolkit";
import {
	Action,
	ActionTypes,
	Graph,
	GraphLink,
	GraphNode,
	MaxMinTime,
	Service,
	ServiceDependenciesActionFailure,
	ServiceDependenciesActionSuccess,
	ServiceDependency,
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

	private static initialServiceDependencyState = [
		{ parent: "", child: "", callCount: 0 },
	];

	private static errorServiceDependencyState = [
		{ parent: "failed", child: "failed", callCount: 0 },
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

	static ServiceDependency(
		client: IClient,
		maxMinTime: MaxMinTime,
	): (dispatch: Dispatch) => Promise<void> {
		return async (dispatch: Dispatch) => {
			try {
				const rsp = await Query.GetServiceDependencies<ServiceDependency[]>(
					client,
					maxMinTime.minTime,
					maxMinTime.maxTime,
				);

				dispatch<ServiceDependenciesActionSuccess>({
					type: ActionTypes.serviceDependenciesSuccess,
					payload: rsp.data,
				});
			} catch (_: unknown) {
				dispatch<ServiceDependenciesActionFailure>({
					type: ActionTypes.serviceDependenciesFailure,
					payload: [],
				});
			}
		};
	}

	static ServiceDependencyReducer(
		state: ServiceDependency[] = ServicesUpdater.initialServiceDependencyState,
		action: Action,
	): ServiceDependency[] {
		switch (action.type) {
			case ActionTypes.serviceDependenciesSuccess:
				if (action.payload.length === 0) {
					return ServicesUpdater.initialServiceDependencyState;
				}
				return action.payload;
			case ActionTypes.serviceDependenciesFailure:
				return ServicesUpdater.errorServiceDependencyState;
			default:
				return state;
		}
	}

	static Graph(services: Service[], dependencies: ServiceDependency[]): Graph {
		const highestCallCount = maxBy(
			dependencies,
			(dep) => dep.callCount,
		)!.callCount;

		const div = Number(
			String(1).padEnd(highestCallCount?.toString().length, "0"),
		);

		const links: GraphLink[] = structuredClone(dependencies).map((node) => {
			const { parent, child, callCount } = node;
			// HACK: the graph is being rendered in the wrong direction
			return {
				source: child,
				target: parent,
				value: (100 - callCount / div) * 0.01,
			};
		});

		const uniqParent = uniqBy(structuredClone(dependencies), "parent").map(
			(dep) => dep.parent,
		);

		const uniqChild = uniqBy(structuredClone(dependencies), "child").map(
			(dep) => dep.child,
		);

		const uniqNodes = uniq([...uniqParent, ...uniqChild]);

		const nodes: GraphNode[] = uniqNodes.map((node, idx) => {
			const service = services.find((service) => service.serviceName === node);

			if (!service) {
				return {
					id: node,
					group: idx + 1,
					p99: 0,
					callRate: "0",
					errorRate: "0",
				};
			}

			return {
				id: node,
				group: idx + 1,
				p99: service.p99,
				callRate: service.callRate.toFixed(2),
				errorRate: service.errorRate.toFixed(2),
			};
		});

		return {
			nodes: nodes,
			links: links,
		};
	}
}
