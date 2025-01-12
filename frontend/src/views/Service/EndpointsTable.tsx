import { useHistory, useParams } from "react-router-dom";
import { useSelector } from "react-redux";
import { Button, Table } from "antd";
import { Endpoint } from "../../updaters/domain";
import { RootState } from "../../updaters/store";

interface EndpointsTableProps {
	data: Endpoint[];
}

export const EndpointsTable = (props: EndpointsTableProps) => {
	const history = useHistory();

	const { service } = useParams<{ service: string }>();

	const { maxTime, minTime } = useSelector(
		(state: RootState) => state.maxMinTime,
	);

	const handleOnClick = (endpoint: string) => {
		const query = new URLSearchParams();

		query.set("start", String(Number(minTime)));
		query.set("end", String(Number(maxTime)));
		query.set("service", service);
		query.set("span", endpoint);

		history.push(`/spans?${query.toString()}`);
	};

	const columns: any = [
		{
			title: "Name",
			dataIndex: "name",
			key: "name",

			render: (text: string) => (
				<Button type="link" onClick={() => handleOnClick(text)}>
					{text}
				</Button>
			),
		},
		{
			title: "P50  (in ms)",
			dataIndex: "p50",
			key: "p50",
			sorter: (a: Endpoint, b: Endpoint) => a.p50 - b.p50,
			// sortDirections: ['descend', 'ascend'],
			render: (value: number) => (value / 1000000).toFixed(2),
		},
		{
			title: "P95  (in ms)",
			dataIndex: "p95",
			key: "p95",
			sorter: (a: Endpoint, b: Endpoint) => a.p95 - b.p95,
			// sortDirections: ['descend', 'ascend'],
			render: (value: number) => (value / 1000000).toFixed(2),
		},
		{
			title: "P99  (in ms)",
			dataIndex: "p99",
			key: "p99",
			sorter: (a: Endpoint, b: Endpoint) => a.p99 - b.p99,
			// sortDirections: ['descend', 'ascend'],
			render: (value: number) => (value / 1000000).toFixed(2),
		},
		{
			title: "Number of Calls",
			dataIndex: "numCalls",
			key: "numCalls",
			sorter: (a: Endpoint, b: Endpoint) => a.numCalls - b.numCalls,
		},
	];

	return (
		<div>
			<h6>Endpoints</h6>
			<Table
				dataSource={props.data}
				columns={columns}
				pagination={false}
				rowKey="name"
			/>
		</div>
	);
};
