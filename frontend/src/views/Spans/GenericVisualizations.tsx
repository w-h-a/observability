import { Bar, Line } from "react-chartjs-2";
import { CustomMetric } from "../../updaters/domain";

export enum GraphType {
	line = "line",
	bar = "bar",
}

interface GenericVisualizationsProps {
	graphType: GraphType;
	data: CustomMetric[];
}

export const GenericVisualizations = (props: GenericVisualizationsProps) => {
	const data = {
		labels: props.data.map((m) => new Date(m.timestamp / 1000000)),
		datasets: [
			{
				data: props.data.map((m) => m.value),
				borderColor: "rgba(250,174,50,1)",
				backgroundColor:
					props.graphType === GraphType.bar ? "rgba(250,174,50,1)" : "",
			},
		],
	};

	const options = {
		responsive: true,
		maintainAspectRatio: false,
		legend: {
			display: false,
		},
		scales: {
			yAxes: [
				{
					gridLines: {
						drawBorder: false,
					},
					ticks: {
						display: false,
					},
				},
			],
			xAxes: [
				{
					type: "time",
					ticks: {
						beginAtZero: false,
						fontSize: 10,
						autoSkip: true,
						maxTicksLimit: 10,
					},
				},
			],
		},
	};

	switch (props.graphType) {
		case GraphType.line:
			return (
				<div>
					<Line data={data} options={options} />
				</div>
			);
		case GraphType.bar:
			return (
				<div>
					<Bar data={data} options={options} />
				</div>
			);
	}
};
