import { combineReducers, configureStore } from "@reduxjs/toolkit";
import { ServicesUpdater } from "./services/servicesTable";
import { TimeUpdater } from "./time/time";

const reducer = combineReducers({
	services: ServicesUpdater.ServicesTableReducer,
	maxMinTime: TimeUpdater.MaxMinTimeReducer,
});

export const store = configureStore({
	reducer: reducer,
});

export type RootState = ReturnType<typeof store.getState>;

export type AppDispatch = typeof store.dispatch;
