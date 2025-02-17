import { Progress, Table, Tabs } from "antd";
import { Tree } from "../../updaters/domain";
import { TracesUpdater } from "../../updaters/traces/traces";

interface TraceWaterfallProps {
	tree: Tree[];
	setSpanTagsInfo: Function;
}

export const TraceWaterfall = (props: TraceWaterfallProps) => {
	const { tree, setSpanTagsInfo } = props;

	const minGlobal: number = tree[0].startTime;
	const maxGlobal: number = TracesUpdater.MaxEndTime(tree);
	const medianGlobal = (minGlobal + maxGlobal) / 2;

	const tabMinVal = 0;
	const tabMedianVal = (medianGlobal - minGlobal).toFixed(0);
	const tabMaxVal = (maxGlobal - minGlobal).toFixed(0);

	const tabs = document.querySelectorAll(
		"#collapsable .ant-tabs-tab",
	) as NodeListOf<HTMLElement>;
	const tabsContainer = document.querySelector(
		"#collapsable .ant-tabs-nav-list",
	) as HTMLElement;

	const columns = [
		{
			title: "",
			dataIndex: "name",
			key: "name",
		},
		{
			title: (
				<Tabs>
					<Tabs.TabPane tab={tabMinVal + "ms"} key="1" />
					<Tabs.TabPane tab={tabMedianVal + "ms"} key="2" />
					<Tabs.TabPane tab={tabMaxVal + "ms"} key="3" />
				</Tabs>
			),
			dataIndex: "trace",
			name: "trace",
			render: (_: any, record: Tree) => {
				let widths = [];
				let length: string = "";

				if (widths.length < tabs.length) {
					Array.from(tabs).forEach((tab) => {
						widths.push(tab.offsetWidth);
					});
				}

				let paddingLeft = 0;
				let startTime = record.startTime;
				let duration = parseFloat((record.time / 1000000).toFixed(2));

				paddingLeft = ((timeDiff, totalTime, totalWidth) => {
					return (timeDiff / totalTime) * totalWidth;
				})(
					startTime - minGlobal,
					maxGlobal - minGlobal,
					tabsContainer?.offsetWidth,
				);

				let textPadding = paddingLeft;

				if (paddingLeft === tabsContainer?.offsetWidth - 20) {
					textPadding = tabsContainer?.offsetWidth - 40;
				}

				length = ((duration / (maxGlobal - startTime)) * 100).toFixed(2);

				return (
					<>
						<div style={{ paddingLeft: textPadding + "px" }}>{duration}ms</div>
						<Progress
							percent={parseInt(length)}
							showInfo={false}
							style={{ paddingLeft: paddingLeft + "px" }}
						/>
					</>
				);
			},
		},
	];

	const onClickRow = (record: Tree) => {
		const node: Tree = TracesUpdater.Span(
			tree,
			(item: Tree) => item.id === record.id,
		);

		setSpanTagsInfo({ data: node });
	};

	return (
		<>
			<Table
				dataSource={tree}
				columns={columns}
				rowKey="id"
				onRow={(record) => {
					return {
						onClick: () => onClickRow(record),
					};
				}}
				rowClassName="row-styles"
				pagination={false}
				scroll={{ y: 540 }}
				sticky={true}
			/>
		</>
	);
};
