import { applyMiddleware, combineReducers, compose, createStore } from "redux";
import { thunk } from "redux-thunk";

export const reducers = combineReducers({});

export const store = createStore(reducers, compose(applyMiddleware(thunk)));
