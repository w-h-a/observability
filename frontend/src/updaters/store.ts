import { combineReducers, configureStore } from "@reduxjs/toolkit";
import { ServicesUpdater } from "./services/servicesTable";
import { TimeUpdater } from "./time/time";
import { ServiceUpdater } from "./service/service";

const reducer = combineReducers({
	maxMinTime: TimeUpdater.MaxMinTimeReducer,
	services: ServicesUpdater.ServicesReducer,
	endpoints: ServiceUpdater.EndpointsReducer,
});

export const store = configureStore({
	reducer: reducer,
});

export type RootState = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
