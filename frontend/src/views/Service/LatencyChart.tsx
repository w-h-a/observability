import React from "react";
import { Line } from "react-chartjs-2";
import { ChartOptions } from "chart.js";
import { ServiceMetricItem } from "../../updaters/domain";

interface LatencyChartProps {
	data: ServiceMetricItem[];
}

export const LatencyChart = (props: LatencyChartProps) => {
	const ref = React.createRef() as any;

	const data = (c: HTMLElement) => {
		const canvas = c as HTMLCanvasElement;

		const ctx = canvas.getContext("2d");

		const gradient = ctx?.createLinearGradient(0, 0, 0, 100);

		gradient?.addColorStop(0, "rgba(250,174,50,1)");
		gradient?.addColorStop(1, "rgba(250,174,50,1)");

		return {
			labels: props.data.map((item) => new Date(item.timestamp / 1000000)),
			datasets: [
				{
					label: "p99 Latency",
					data: props.data.map((item) => item.p99 / 1000000), // converting latency from nano sec to ms
					pointRadius: 0.5,
					borderColor: "rgba(250,174,50,1)",
					borderWidth: 2,
				},
				{
					label: "p95 Latency",
					data: props.data.map((item) => item.p95 / 1000000), // converting latency from nano sec to ms
					pointRadius: 0.5,
					borderColor: "rgba(227, 74, 51, 1.0)",
					borderWidth: 2,
				},
				{
					label: "p50 Latency",
					data: props.data.map((item) => item.p50 / 1000000), // converting latency from nano sec to ms
					pointRadius: 0.5,
					borderColor: "rgba(57, 255, 20, 1.0)",
					borderWidth: 2,
				},
			],
		};
	};

	const options: ChartOptions = {
		maintainAspectRatio: true,
		responsive: true,
		title: {
			display: true,
			text: "",
			fontSize: 20,
			position: "top",
			padding: 8,
			fontFamily: "Arial",
			fontStyle: "regular",
			fontColor: "rgb(200, 200, 200)",
		},
		legend: {
			display: true,
			position: "bottom",
			align: "center",
			labels: {
				fontColor: "rgb(200, 200, 200)",
				fontSize: 10,
				boxWidth: 10,
				usePointStyle: true,
			},
		},
		tooltips: {
			mode: "label",
			bodyFontSize: 12,
			titleFontSize: 12,

			callbacks: {
				label: function (tooltipItem, data) {
					if (typeof tooltipItem.yLabel === "number") {
						return (
							data.datasets![tooltipItem.datasetIndex!].label +
							" : " +
							tooltipItem.yLabel.toFixed(2)
						);
					} else {
						return "";
					}
				},
			},
		},
		scales: {
			yAxes: [
				{
					stacked: false,
					ticks: {
						beginAtZero: false,
						fontSize: 10,
						autoSkip: true,
						maxTicksLimit: 6,
					},
					gridLines: {
						borderDash: [1, 4],
						color: "#D3D3D3",
						lineWidth: 0.25,
					},
				},
			],
			xAxes: [
				{
					type: "time",
					distribution: "linear",
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

	return (
		<div>
			{/* popup */}
			<div style={{ textAlign: "center" }}>Latency in ms</div>
			<Line ref={ref} data={data} options={options} />
		</div>
	);
};
