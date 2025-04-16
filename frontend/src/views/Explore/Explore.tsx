import { useContext, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import {
	AutoComplete,
	Button,
	Card,
	Form,
	Input,
	Select,
	Space,
	Tag,
} from "antd";
import FormItem from "antd/es/form/FormItem";
import { GenericVisualizations, GraphType } from "./GenericVisualizations";
import { AppDispatch, RootState } from "../../updaters/store";
import { MetricsUpdater } from "../../updaters/metrics/metrics";
import { SpansUpdater } from "../../updaters/spans/spans";
import { ServiceUpdater } from "../../updaters/service/service";
import { ServicesUpdater } from "../../updaters/services/services";
import { ClientContext } from "../../clients/query/clientCtx";
import {
	FilteredQuery,
	Operator,
	SpanKind,
} from "../../clients/query/filteredQuery";

enum TagValue {
	service = "service",
}

enum CustomVisualizationField {
	dimension = "dimension",
	aggregation = "aggregation",
	interval = "interval",
	graphType = "graph_type",
}

const dimensions = [
	{
		title: "Calls",
		key: "calls",
		value: "calls",
	},
	{
		title: "Duration",
		key: "duration",
		value: "duration",
	},
	{
		title: "CPU",
		key: "cpu",
		value: "cpu",
	},
];

const aggregations = [
	{
		dimension: "calls",
		defaultSelected: { title: "Count", key: "count", value: "count" },
		optionsAvailable: [
			{ title: "Count", key: "count", value: "count" },
			{ title: "Rate (per sec)", key: "rate_per_sec", value: "rate_per_sec" },
		],
	},
	{
		dimension: "duration",
		defaultSelected: { title: "p99", key: "p99", value: "p99" },
		optionsAvailable: [
			{ title: "p99", key: "p99", value: "p99" },
			{ title: "p95", key: "p95", value: "p95" },
			{ title: "p50", key: "p50", value: "p50" },
			{ title: "avg", key: "avg", value: "avg" },
		],
	},
	{
		dimension: "cpu",
		defaultSelected: { title: "Count", key: "count", value: "count" },
		optionsAvailable: [{ title: "Count", key: "count", value: "count" }],
	},
];

export const Explore = () => {
	// clients
	const { queryClient } = useContext(ClientContext);

	// local state
	const [filters, setFilters] = useState<FilteredQuery>({
		service: "",
		operation: "",
		kind: SpanKind.default,
		duration: { min: "", max: "" },
		tags: [],
	});
	const [dimension, setDimension] = useState("calls");
	const [aggregation, setAggregation] = useState("count");
	// const [step, setStep] = useState("60");

	// store state
	const maxMinTime = useSelector((state: RootState) => state.maxMinTime);
	const serviceNames = useSelector((state: RootState) => state.serviceNames);
	const tags = useSelector((state: RootState) => state.tags);
	const customMetrics = useSelector((state: RootState) => state.customMetrics);

	const dispatch: AppDispatch = useDispatch();

	// retrieve filtered spans
	useEffect(() => {
		if (
			filters.service ||
			filters.kind ||
			filters.operation ||
			(filters.duration && (filters.duration.min || filters.duration.max)) ||
			filters.tags.length !== 0
		) {
			dispatch(SpansUpdater.Spans(queryClient, maxMinTime, filters));
		}
	}, [dispatch, queryClient, maxMinTime, filters]);

	// retrieve service names
	useEffect(() => {
		dispatch(ServicesUpdater.ServiceNames(queryClient));
	}, [dispatch, queryClient]);

	const [initialForm] = Form.useForm();

	// filter on service
	const onChangeService = (value: string) => {
		setFilters({ ...filters, service: value });
	};

	useEffect(() => {
		if (filters.service) {
			dispatch(ServiceUpdater.OperationNames(queryClient, filters.service ?? ""));
			dispatch(ServiceUpdater.Tags(queryClient, filters.service ?? ""));
			initialForm.setFieldsValue({ service: filters.service });
		}
	}, [dispatch, initialForm, queryClient, filters.service]);

	// filter on tags
	const [tagForm] = Form.useForm();

	const onTagFormSubmit = (values: any) => {
		setFilters({
			...filters,
			tags: [
				...filters.tags,
				{
					key: values.tag_key,
					operator: values.operator,
					value: values.tag_value,
				},
			],
		});
	};

	const onChangeTagKey = (value: string) => {
		tagForm.setFieldsValue({ tag_key: value });
	};

	// remove filters
	const onCloseTag = (value: TagValue) => {
		switch (value) {
			case TagValue.service:
				setFilters({ ...filters, service: "" });
				break;
		}
	};

	const onCloseTagTag = (tag: {
		key: string;
		value: string;
		operator: Operator;
	}) => {
		setFilters({
			...filters,
			tags: filters.tags?.filter((t) => {
				return (
					t.key !== tag.key && t.value !== tag.value && t.operator !== tag.operator
				);
			}),
		});
	};

	// custom visualizations stuff
	const [customVizForm] = Form.useForm();

	const onCustomVizValuesChange = (changedValues: any) => {
		const field = Object.keys(changedValues)[0];

		switch (field) {
			case CustomVisualizationField.dimension:
				const tempAgg = aggregations.filter((a) => {
					return a.dimension === changedValues[field];
				})[0];

				customVizForm.setFieldsValue({
					aggregation: tempAgg.defaultSelected.value,
				});

				const values = customVizForm.getFieldsValue([
					CustomVisualizationField.dimension,
					CustomVisualizationField.aggregation,
				]);

				setDimension(values[CustomVisualizationField.dimension]);
				setAggregation(values[CustomVisualizationField.aggregation]);

				break;
			case CustomVisualizationField.aggregation:
				setAggregation(changedValues[field]);
				break;
			case CustomVisualizationField.interval:
				break;
			case CustomVisualizationField.graphType:
				break;
		}
	};

	useEffect(() => {
		dispatch(
			MetricsUpdater.CustomMetrics(
				queryClient,
				{
					minTime: maxMinTime.minTime - 15 * 60,
					maxTime: maxMinTime.maxTime + 15 * 60,
				},
				dimension,
				aggregation,
				filters,
			),
		);
	}, [dispatch, queryClient, maxMinTime, dimension, aggregation, filters]);

	return (
		<div>
			<Card>
				<div>Filter Data</div>
				<Form
					form={initialForm}
					layout="inline"
					initialValues={{ service: "" }}
					style={{ marginTop: 10, marginBottom: 10 }}
				>
					<FormItem rules={[{ required: true }]} name="service">
						<Select
							showSearch
							style={{ width: 180 }}
							onChange={onChangeService}
							placeholder="Select Service"
							allowClear
						>
							{serviceNames.map((name: string, idx: number) => {
								return (
									<Select.Option value={name} key={idx}>
										{name}
									</Select.Option>
								);
							})}
						</Select>
					</FormItem>
				</Form>
				<Card
					style={{ padding: 6, marginTop: 10, marginBottom: 10 }}
					bodyStyle={{ padding: 6 }}
				>
					{!filters.service ? null : (
						<Tag
							style={{ fontSize: 14, padding: 8 }}
							closable
							onClose={() => onCloseTag(TagValue.service)}
						>
							service:{filters.service}
						</Tag>
					)}
					{!filters.tags
						? null
						: filters.tags
								.filter((t, i) => {
									const found = filters.tags.findIndex((e) => {
										return (
											t.key === e.key && t.operator === e.operator && t.value === e.value
										);
									});
									return i === found;
								})
								.map((t) => {
									return (
										<Tag
											style={{ fontSize: 14, padding: 8 }}
											closable
											key={`${t.key}-${t.operator}-${t.value}`}
											onClose={() => onCloseTagTag(t)}
										>
											{t.key} {t.operator} {t.value}
										</Tag>
									);
								})}
				</Card>
				<div>Select service to get tag suggestions</div>
				<Form
					form={tagForm}
					layout="inline"
					onFinish={onTagFormSubmit}
					initialValues={{ operator: "equals" }}
					style={{ marginTop: 10, marginBottom: 10 }}
				>
					<FormItem rules={[{ required: true }]} name="tag_key">
						<AutoComplete
							options={tags.map((key: string) => {
								return { value: key };
							})}
							style={{ width: 200, textAlign: "center" }}
							onChange={onChangeTagKey}
							filterOption={(input: string, option: { value: string } | undefined) => {
								return !!(
									option && option.value.toUpperCase().includes(input.toUpperCase())
								);
							}}
							placeholder="Tag Key"
						/>
					</FormItem>
					<FormItem name="operator">
						<Select style={{ width: 120, textAlign: "center" }}>
							<Select.Option value="equals">EQUAL</Select.Option>
							<Select.Option value="contains">CONTAINS</Select.Option>
						</Select>
					</FormItem>
					<FormItem rules={[{ required: true }]} name="tag_value">
						<Input
							style={{ width: 160, textAlign: "center" }}
							placeholder="Tag Value"
						/>
					</FormItem>
					<FormItem>
						<Button type="primary" htmlType="submit">
							{" "}
							Apply Tag Filter{" "}
						</Button>
					</FormItem>
				</Form>
			</Card>
			<Card>
				<div>Custom Visualizations</div>
				<Form
					form={customVizForm}
					onValuesChange={onCustomVizValuesChange}
					initialValues={{
						[CustomVisualizationField.dimension]: dimension,
						[CustomVisualizationField.aggregation]: "Count",
						[CustomVisualizationField.interval]: "1m",
						[CustomVisualizationField.graphType]: "line",
					}}
				>
					<Space>
						<Form.Item name={CustomVisualizationField.dimension}>
							<Select style={{ width: 120 }}>
								{dimensions.map((d) => {
									return (
										<Select.Option key={d.key} value={d.value}>
											{d.title}
										</Select.Option>
									);
								})}
							</Select>
						</Form.Item>
						<Form.Item name={CustomVisualizationField.aggregation}>
							<Select style={{ width: 120 }}>
								{aggregations
									.filter((a) => {
										return a.dimension === dimension;
									})[0]
									.optionsAvailable.map((a) => {
										return (
											<Select.Option key={a.key} value={a.value}>
												{a.title}
											</Select.Option>
										);
									})}
							</Select>
						</Form.Item>
						<Form.Item name={CustomVisualizationField.interval}>
							<Select style={{ width: 120 }} allowClear>
								<Select.Option value="1m">1 min</Select.Option>
								<Select.Option value="5m">5 min</Select.Option>
								<Select.Option value="30m">30 min</Select.Option>
							</Select>
						</Form.Item>
						<Form.Item name={CustomVisualizationField.graphType}>
							<Select style={{ width: 120 }} allowClear>
								<Select.Option value={GraphType.line}>Line</Select.Option>
								<Select.Option value={GraphType.bar}>Bar</Select.Option>
							</Select>
						</Form.Item>
					</Space>
				</Form>
				<GenericVisualizations graphType={GraphType.line} data={customMetrics} />
			</Card>
		</div>
	);
};
