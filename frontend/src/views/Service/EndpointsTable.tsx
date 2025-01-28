import { Table } from "antd";
import { Endpoint } from "../../updaters/domain";

interface EndpointsTableProps {
	data: Endpoint[];
}

const columns: any = [
	{
		title: "Name",
		dataIndex: "name",
		key: "name",
		render: (text: string) => text,
	},
	{
		title: "P50  (in ms)",
		dataIndex: "p50",
		key: "p50",
		sorter: (a: Endpoint, b: Endpoint) => a.p50 - b.p50,
		render: (value: number) => (value / 1000000).toFixed(2),
	},
	{
		title: "P95  (in ms)",
		dataIndex: "p95",
		key: "p95",
		sorter: (a: Endpoint, b: Endpoint) => a.p95 - b.p95,
		render: (value: number) => (value / 1000000).toFixed(2),
	},
	{
		title: "P99  (in ms)",
		dataIndex: "p99",
		key: "p99",
		sorter: (a: Endpoint, b: Endpoint) => a.p99 - b.p99,
		render: (value: number) => (value / 1000000).toFixed(2),
	},
	{
		title: "Number of Calls",
		dataIndex: "numCalls",
		key: "numCalls",
		sorter: (a: Endpoint, b: Endpoint) => a.numCalls - b.numCalls,
	},
];

export const EndpointsTable = (props: EndpointsTableProps) => {
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
