import "./TraceGraph.css";
import { useContext, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { useParams } from "react-router-dom";
import { Affix, Card, Col, Row, Space, Spin } from "antd";
import { TraceWaterfall } from "./TraceWaterfall";
import { SelectedSpanDetails } from "../Spans/SelectedSpanDetails";
import { AppDispatch, RootState } from "../../updaters/store";
import { TracesUpdater } from "../../updaters/traces/traces";
import { ClientContext } from "../../clients/query/clientCtx";
import { Tree } from "../../updaters/domain";

export const TraceGraph = () => {
	const { queryClient } = useContext(ClientContext);

	const { id } = useParams<{ id: string }>();

	const [clickedSpanTags, setClickedSpanTags] = useState<Tree>();
	const [sortedTree, setSortedTree] = useState<Tree[]>([]);

	const maxMinTime = useSelector((state: RootState) => state.maxMinTime);
	const traces = useSelector((state: RootState) => state.traces);

	const dispatch: AppDispatch = useDispatch();

	useEffect(() => {
		dispatch(TracesUpdater.Traces(queryClient, maxMinTime, "", id));
	}, [dispatch, queryClient, maxMinTime, id]);

	useEffect(() => {
		if (traces[0].events.length !== 0) {
			const tree = TracesUpdater.SpansToTree(structuredClone(traces[0].events));
			const sorted = TracesUpdater.SortTree([structuredClone(tree)]);
			setSortedTree(sorted);
		}
	}, [traces]);

	if (traces[0].events.length === 0 || sortedTree.length === 0) {
		return <Spin />;
	}

	return (
		<Row gutter={{ xs: 8, sm: 16, md: 24, lg: 32 }}>
			<Col md={18} sm={18}>
				<Space direction="vertical" size="middle" style={{ width: "100%" }}>
					<Affix offsetTop={24}>
						<Card
							id="collapsable"
							style={{ background: "#333333", borderRadius: "5px" }}
						>
							<TraceWaterfall
								tree={sortedTree}
								setSpanTagsInfo={(z: any) => setClickedSpanTags(z.data)}
							/>
						</Card>
					</Affix>
				</Space>
			</Col>
			<Col md={6} sm={6}>
				<Affix offsetTop={24}>
					<SelectedSpanDetails data={clickedSpanTags} />
				</Affix>
			</Col>
		</Row>
	);
};
