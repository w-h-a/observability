import { useState } from "react";
import { NavLink } from "react-router-dom";
import { Menu } from "antd";
import {
	AlignLeftOutlined,
	BarChartOutlined,
	DeploymentUnitOutlined,
	ExperimentOutlined,
} from "@ant-design/icons";
import Sider from "antd/es/layout/Sider";

// TODO: make routes an enum

export const SideNav = () => {
	const [collapsed, setCollapsed] = useState(false);

	const onCollapse = () => {
		setCollapsed(!collapsed);
	};

	return (
		<Sider collapsible collapsed={collapsed} onCollapse={onCollapse} width={160}>
			<Menu mode="inline">
				<Menu.Item key="/application" icon={<BarChartOutlined />}>
					<NavLink
						to="/application"
						style={{ fontSize: 12, textDecoration: "none" }}
					>
						Services
					</NavLink>
				</Menu.Item>
				<Menu.Item key="/service-map" icon={<DeploymentUnitOutlined />}>
					<NavLink
						to="/service-map"
						style={{ fontSize: 12, textDecoration: "none" }}
					>
						Service Map
					</NavLink>
				</Menu.Item>
				<Menu.Item key="/traces" icon={<AlignLeftOutlined />}>
					<NavLink to="/traces" style={{ fontSize: 12, textDecoration: "none" }}>
						Traces
					</NavLink>
				</Menu.Item>
				<Menu.Item key="/spans" icon={<ExperimentOutlined />}>
					<NavLink to="/spans" style={{ fontSize: 12, textDecoration: "none" }}>
						Spans
					</NavLink>
				</Menu.Item>
			</Menu>
		</Sider>
	);
};
