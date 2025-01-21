import "./TraceGraph.css";
import { useContext, useEffect } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router-dom";
import { Card, Col, Row, Space } from "antd";
import { flamegraph } from "d3-flame-graph";
import * as d3 from "d3";
// @ts-ignore
import * as d3Tip from "d3-tip";
import { AppDispatch, RootState } from "../../updaters/store";
import { SpansUpdater } from "../../updaters/spans/spans";
import { ClientContext } from "../../clients/query/clientCtx";
import { Config, EnvVar } from "../../config/config";

export const TraceGraph = () => {
	const { queryClient } = useContext(ClientContext);

	const { id } = useParams<{ id: string }>();

	const spanMatrixForATrace = useSelector(
		(state: RootState) => state.spanMatrixForATrace,
	);

	const dispatch: AppDispatch = useDispatch();

	useEffect(() => {
		dispatch(SpansUpdater.SpansByTraceId(queryClient, id));
	}, [dispatch, queryClient, id]);

	const tip = d3Tip
		.default()
		.attr("class", "d3-tip")
		.html(function (d: any) {
			let display =
				d.data.name +
				"<br>duration: " +
				d.data.value / 1000000 +
				"ms" +
				"<br>status: " +
				d.data.code;

			for (const tag of d.data.tags) {
				display += `<br>${tag.key}: ${tag.value}`;
			}

			return display;
		});

	const graph = flamegraph()
		.cellHeight(20)
		.transitionDuration(500)
		.inverted(true)
		.minFrameSize(10)
		.elided(false)
		.differential(false)
		.selfValue(true)
		.sort(true);

	// hack for now
	if (Config.GetInstance().get(EnvVar.ENVIRONMENT) !== "test") {
		graph.tooltip(tip);
	}

	useEffect(() => {
		if (spanMatrixForATrace[0].events.length !== 0) {
			const tree = SpansUpdater.SpansToTree(
				structuredClone(spanMatrixForATrace[0].events),
			);
			d3.select("#graph").datum(tree).call(graph);
		}
	}, [spanMatrixForATrace, graph]);

	return (
		<Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
			<Col md={24} sm={24}>
				<Space direction="vertical" size="middle" style={{ width: "100%" }}>
					<Card bodyStyle={{ padding: 80 }} style={{ height: 500 }}>
						<div
							style={{
								display: "flex",
								justifyContent: "center",
								flexDirection: "column",
								alignItems: "center",
							}}
						>
							<div style={{ textAlign: "center", marginTop: 40 }}>
								Trace Graph For Trace {id}{" "}
							</div>
							<div id="graph" style={{ fontSize: 12, marginTop: 80 }}></div>
						</div>
					</Card>
				</Space>
			</Col>
		</Row>
	);
};
