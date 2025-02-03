import { useContext, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { ForceGraph2D } from "react-force-graph";
import { Spin } from "antd";
import { AppDispatch, RootState } from "../../updaters/store";
import { ServicesUpdater } from "../../updaters/services/services";
import { ClientContext } from "../../clients/query/clientCtx";

export const ServiceMap = () => {
	const { queryClient } = useContext(ClientContext);

	const maxMinTime = useSelector((state: RootState) => state.maxMinTime);
	const services = useSelector((state: RootState) => state.services);
	const serviceDependencies = useSelector(
		(state: RootState) => state.serviceDependencies,
	);

	const dispatch: AppDispatch = useDispatch();

	useEffect(() => {
		dispatch(ServicesUpdater.Services(queryClient, maxMinTime));
		dispatch(ServicesUpdater.ServiceDependency(queryClient, maxMinTime));
	}, [dispatch, queryClient, maxMinTime]);

	if (
		services[0].serviceName.length === 0 ||
		(serviceDependencies[0].parent.length === 0 &&
			serviceDependencies[0].child.length === 0)
	) {
		return <Spin />;
	}

	const graphData = ServicesUpdater.Graph(services, serviceDependencies);

	return (
		<div>
			<ForceGraph2D
				cooldownTicks={100}
				graphData={graphData}
				nodeLabel="id"
				linkAutoColorBy={(d) => d.target}
				linkDirectionalParticles="value"
				linkDirectionalParticleSpeed={(d) => d.value}
				nodeCanvasObject={(node, ctx) => {
					let color = "#88CEA5";
					if (Number(node.errorRate) > 0) {
						color = "#F98989";
					}

					ctx.fillStyle = color;
					ctx.beginPath();
					if (node.x && node.y) {
						ctx.arc(node.x, node.y, 16, 0, 2 * Math.PI, false);
					}
					ctx.fill();

					const fontSize = 6;
					ctx.font = `${fontSize}px Roboto`;

					const label = node.id;

					ctx.textAlign = "center";
					ctx.textBaseline = "middle";
					ctx.fillStyle = "#333333";
					if (node.x && node.y) {
						ctx.fillText(label, node.x, node.y);
					}
				}}
				nodePointerAreaPaint={(node, color, ctx) => {
					ctx.fillStyle = color;
					ctx.beginPath();
					if (node.x && node.y) {
						ctx.arc(node.x, node.y, 5, 0, 2 * Math.PI, false);
					}
					ctx.fill();
				}}
				onNodeDragEnd={(node) => {
					node.fx = node.x;
					node.fy = node.y;
				}}
			/>
		</div>
	);
};
