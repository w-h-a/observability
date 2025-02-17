import { Card, Space, Tabs, Typography } from "antd";
import { Tree } from "../../updaters/domain";

interface SelecedSpanDetailsProps {
	data: Tree | undefined;
}

export const SelectedSpanDetails = (props: SelecedSpanDetailsProps) => {
	const node = props.data;

	if (!node) {
		return <div></div>;
	}

	const service = node.name.split(":")[0];
	const operation = node.name.split(":")[1];
	const code = node.code;
	const spanTags = node.tags;

	return (
		<Card
			style={{ border: "none", background: "transparent", padding: 0 }}
			bodyStyle={{ padding: 0 }}
		>
			<Space direction="vertical">
				<strong>Details of selected span</strong>
				<Space direction="vertical" size={2}>
					<Typography.Text style={{ marginTop: "18px" }}>Service</Typography.Text>
					<Typography.Text style={{ color: "#2D9CDB", fontSize: "12px" }}>
						{service}
					</Typography.Text>
				</Space>
				<Space direction="vertical" size={2}>
					<Typography.Text>Operation</Typography.Text>
					<Typography.Text style={{ color: "#2D9CDB", fontSize: "12px" }}>
						{operation}
					</Typography.Text>
				</Space>
			</Space>
			<Tabs defaultActiveKey="1">
				<Tabs.TabPane tab="Status Code" key="1">
					<div>
						<Typography.Text
							style={{ color: "#BDBDBD", fontSize: "12px", marginBottom: "8px" }}
						>
							Status Code
						</Typography.Text>
						<div
							style={{
								background: "#4F4F4F",
								color: "#2D9CDB",
								fontSize: "12px",
								padding: "6px 8px",
								wordBreak: "break-all",
								marginBottom: "16px",
							}}
						>
							{code}
						</div>
					</div>
				</Tabs.TabPane>
				<Tabs.TabPane tab="Tags" key="2">
					{spanTags &&
						spanTags.map((tag, idx) => {
							return (
								<div key={idx}>
									<Typography.Text
										style={{ color: "#BDBDBD", fontSize: "12px", marginBottom: "8px" }}
									>
										{tag.key}
									</Typography.Text>
									<div
										style={{
											background: "#4F4F4F",
											color: "#2D9CDB",
											fontSize: "12px",
											padding: "6px 8px",
											wordBreak: "break-all",
											marginBottom: "16px",
										}}
									>
										{tag.value}
									</div>
								</div>
							);
						})}
				</Tabs.TabPane>
			</Tabs>
		</Card>
	);
};
