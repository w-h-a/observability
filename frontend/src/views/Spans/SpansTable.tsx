import { useContext, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { NavLink } from "react-router-dom";
import { Table } from "antd";
import { AppDispatch, RootState } from "../../updaters/store";
import { SpansUpdater } from "../../updaters/spans/spans";
import { StartTime } from "../../updaters/time/utils";
import { ClientContext } from "../../clients/query/clientCtx";

interface SpanView {
	key: string;
	startTime: number;
	spanid: string;
	parentspanid: string;
	traceid: string;
	serviceName: string;
	operationName: string;
	kind: string;
	code: string;
	duration: number;
}

const columns = [
	{
		title: "Start Time",
		dataIndex: "startTime",
		key: "startTime",
		sorter: (a: SpanView, b: SpanView) => a.startTime - b.startTime,
		render: StartTime,
	},
	{
		title: "Duration (in ms)",
		dataIndex: "duration",
		key: "duration",
		sorter: (a: SpanView, b: SpanView) => a.duration - b.duration,
		render: (value: number) => (value / 1000000).toFixed(2),
	},
	{
		title: "Service Name",
		dataIndex: "serviceName",
		key: "serviceName",
	},
	{
		title: "Operation",
		dataIndex: "operationName",
		key: "operationName",
	},
	{
		title: "Status Code",
		dataIndex: "code",
		key: "code",
	},
	{
		title: "TraceID",
		dataIndex: "traceid",
		key: "traceid",
		render: (id: string) => (
			<NavLink to={`/traces/${id}`}>{id.slice(-16)}</NavLink>
		),
	},
];

export const SpansTable = () => {
	const { queryClient } = useContext(ClientContext);

	const maxMinTime = useSelector((state: RootState) => state.maxMinTime);
	const spanMatrix = useSelector((state: RootState) => state.spanMatrix);

	const spans: SpanView[] = spanMatrix[0].events.map((span, idx) => {
		return {
			key: String(idx),
			startTime: span[0],
			spanid: span[1],
			parentspanid: span[2],
			traceid: span[3],
			serviceName: span[4],
			operationName: span[5],
			kind: span[6],
			code: span[7],
			duration: Number(span[8]),
		};
	});

	const dispatch: AppDispatch = useDispatch();

	useEffect(() => {
		dispatch(SpansUpdater.Spans(queryClient, maxMinTime));
	}, [dispatch, queryClient, maxMinTime]);

	return (
		<div>
			<div>List of spans</div>
			<div>
				<Table dataSource={spans} columns={columns} size="middle" />
			</div>
		</div>
	);
};
