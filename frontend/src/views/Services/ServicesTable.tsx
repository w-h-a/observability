import { useContext, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { NavLink } from "react-router-dom";
import { Table } from "antd";
import { ClientContext } from "../../App";
import { ServicesUpdater } from "../../updaters/services/servicesTable";
import { AppDispatch, RootState } from "../../updaters/store";

const columns = [
	{
		title: "Service",
		dataIndex: "serviceName",
		key: "serviceName",
		render: (name: string) => (
			<NavLink style={{ textTransform: "capitalize" }} to={`/application/${name}`}>
				<strong>{name}</strong>
			</NavLink>
		),
	},
	{
		title: "P99 latency (in ms)",
		dataIndex: "p99",
		key: "p99",
		sorter: (a: any, b: any) => a.p99 - b.p99,
		render: (value: number) => (value / 1000000).toFixed(2),
	},
	{
		title: "Error Rate (in %)",
		dataIndex: "errorRate",
		key: "errorRate",
		sorter: (a: any, b: any) => a.errorRate - b.errorRate,
		render: (value: number) => value.toFixed(2),
	},
	{
		title: "Requests Per Second",
		dataIndex: "callRate",
		key: "callRate",
		sorter: (a: any, b: any) => a.callRate - b.callRate,
		render: (value: number) => value.toFixed(2),
	},
];

export const ServicesTable = () => {
	const { queryClient } = useContext(ClientContext);

	const useAppSelector = useSelector.withTypes<RootState>();
	const services = useAppSelector((state) => state.services);
	const maxMinTime = useAppSelector((state) => state.maxMinTime);

	const useAppDispatch = useDispatch.withTypes<AppDispatch>();
	const dispatch = useAppDispatch();

	useEffect(() => {
		dispatch(ServicesUpdater.Services(queryClient, maxMinTime));
	}, [dispatch, queryClient, maxMinTime]);

	return (
		<Table
			dataSource={services}
			columns={columns}
			pagination={false}
			rowKey="serviceName"
		/>
	);
};
