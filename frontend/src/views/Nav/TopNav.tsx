import { Col, Row } from "antd";
import { ShowBreadcrumbs } from "./ShowBreadcrumbs";
import { DateTimeSelector } from "./DateTimeSelector";

export const TopNav = () => {
	return (
		<Row>
			<Col span={16}>
				<ShowBreadcrumbs />
			</Col>
			<Col span={8}>
				<DateTimeSelector />
			</Col>
		</Row>
	);
};
