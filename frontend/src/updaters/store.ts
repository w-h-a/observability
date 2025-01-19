import { combineReducers, configureStore } from "@reduxjs/toolkit";
import { ServicesUpdater } from "./services/services";
import { TimeUpdater } from "./time/time";
import { ServiceUpdater } from "./service/service";
import { SpansUpdater } from "./spans/spans";

const reducer = combineReducers({
	maxMinTime: TimeUpdater.MaxMinTimeReducer,
	services: ServicesUpdater.ServicesReducer,
	endpoints: ServiceUpdater.EndpointsReducer,
	serviceMetrics: ServiceUpdater.ServiceMetricsReducer,
	spanMatrix: SpansUpdater.SpanMatrixReducer,
});

export const store = configureStore({
	reducer: reducer,
});

export type RootState = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
