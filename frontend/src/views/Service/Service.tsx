import { useContext, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router-dom";
import { Card, Col, Row, Tabs } from "antd";
import { EndpointsTable } from "./EndpointsTable";
import { RequestRateChart } from "./RequestRateChart";
import { AppDispatch, RootState } from "../../updaters/store";
import { ServiceUpdater } from "../../updaters/service/service";
import { ClientContext } from "../../clients/query/clientCtx";

export const Service = () => {
	const { queryClient } = useContext(ClientContext);

	const { service } = useParams<{ service: string }>();

	const maxMinTime = useSelector((state: RootState) => state.maxMinTime);
	const endpoints = useSelector((state: RootState) => state.endpoints);
	const serviceMetrics = useSelector((state: RootState) => state.serviceMetrics);

	const dispatch: AppDispatch = useDispatch();

	useEffect(() => {
		dispatch(ServiceUpdater.Endpoints(queryClient, maxMinTime, service));
		dispatch(ServiceUpdater.ServiceMetrics(queryClient, maxMinTime, service));
	}, [dispatch, queryClient, maxMinTime, service]);

	return (
		<Tabs defaultActiveKey="1">
			<Tabs.TabPane tab="Service Metrics" key="1">
				<Row gutter={32} style={{ margin: 20 }}>
					<Col span={12}>
						{/* <Card bodyStyle={{ padding: 10 }}>
							<LatencyLineChart
								data={serviceMetrics}
								// popupClickHandler={onTracePopupClick}
							/>
						</Card> */}
					</Col>
					<Col span={12}>
						<Card bodyStyle={{ padding: 10 }}>
							<RequestRateChart data={serviceMetrics} />
						</Card>
					</Col>
				</Row>
				<Row gutter={32} style={{ margin: 20 }}>
					<Col span={12}>
						{/* <Card bodyStyle={{ padding: 10 }}>
							<ErrorRateChart
								data={serviceMetrics}
								// popupClickHandler={onErrTracePopupClick}
							/>
						</Card> */}
					</Col>
					<Col span={12}>
						<EndpointsTable data={endpoints} />
					</Col>
				</Row>
			</Tabs.TabPane>
		</Tabs>
	);
};
