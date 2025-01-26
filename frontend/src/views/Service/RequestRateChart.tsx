import React from "react";
import { Line } from "react-chartjs-2";
import { ChartOptions } from "chart.js";
import { ServiceMetricItem } from "../../updaters/domain";

interface RequestRateChartProps {
	data: ServiceMetricItem[];
}

export const RequestRateChart = (props: RequestRateChartProps) => {
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
					label: "Requests per sec",
					data: props.data.map((item) => item.callRate),
					pointRadius: 0.5,
					borderColor: "rgba(250,174,50,1)",
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
			padding: 2,
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
			<div style={{ textAlign: "center" }}>Requests per sec</div>
			<Line ref={ref} data={data} options={options} />
		</div>
	);
};
