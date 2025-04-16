import { combineReducers, configureStore } from "@reduxjs/toolkit";
import { TimeUpdater } from "./time/time";
import { ServicesUpdater } from "./services/services";
import { ServiceUpdater } from "./service/service";
import { TracesUpdater } from "./traces/traces";
import { SpansUpdater } from "./spans/spans";
import { MetricsUpdater } from "./metrics/metrics";

const reducer = combineReducers({
	maxMinTime: TimeUpdater.MaxMinTimeReducer,
	services: ServicesUpdater.ServicesReducer,
	serviceNames: ServicesUpdater.ServiceNamesReducer,
	serviceDependencies: ServicesUpdater.ServiceDependencyReducer,
	tags: ServiceUpdater.TagsReducer,
	operationNames: ServiceUpdater.OperationNamesReducer,
	endpoints: ServiceUpdater.EndpointsReducer,
	serviceMetrics: ServiceUpdater.ServiceMetricsReducer,
	traces: TracesUpdater.TracesReducer,
	spanMatrix: SpansUpdater.SpanMatrixReducer,
	customMetrics: MetricsUpdater.CustomMetricsReducer,
});

export const store = configureStore({
	reducer: reducer,
});

export type RootState = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
