import { useContext, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { NavLink } from "react-router-dom";
import { Table } from "antd";
import { ServicesUpdater } from "../../updaters/services/servicesTable";
import { AppDispatch, RootState } from "../../updaters/store";
import { Service } from "../../updaters/domain";
import { ClientContext } from "../../clients/query/clientCtx";

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
		sorter: (a: Service, b: Service) => a.p99 - b.p99,
		render: (value: number) => (value / 1000000).toFixed(2),
	},
	{
		title: "Error Rate (in %)",
		dataIndex: "errorRate",
		key: "errorRate",
		sorter: (a: Service, b: Service) => a.errorRate - b.errorRate,
		render: (value: number) => value.toFixed(2),
	},
	{
		title: "Requests Per Second",
		dataIndex: "callRate",
		key: "callRate",
		sorter: (a: Service, b: Service) => a.callRate - b.callRate,
		render: (value: number) => value.toFixed(2),
	},
];

export const ServicesTable = () => {
	const { queryClient } = useContext(ClientContext);

	const maxMinTime = useSelector((state: RootState) => state.maxMinTime);
	const services = useSelector((state: RootState) => state.services);

	const dispatch: AppDispatch = useDispatch();

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
