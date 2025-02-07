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
import { Store } from "antd/es/form/interface";
import FormItem from "antd/es/form/FormItem";
import { DurationModelForm } from "./DurationModalForm";
import { GenericVisualizations, GraphType } from "./GenericVisualizations";
import { AppDispatch, RootState } from "../../updaters/store";
import { SpansUpdater } from "../../updaters/spans/spans";
import { ServiceUpdater } from "../../updaters/service/service";
import { ServicesUpdater } from "../../updaters/services/services";
import { ClientContext } from "../../clients/query/clientCtx";
import {
	FilteredQuery,
	Operator,
	SpanKind,
} from "../../clients/query/filteredQuery";

// TODO: status code
enum TagValue {
	service = "service",
	operation = "operation",
	kind = "kind",
	minDuration = "minDuration",
	maxDuration = "maxDuration",
}

enum CustomVisualizationField {
	dimension = "dimension",
	aggregation = "aggregation",
	interval = "interval",
	graphType = "graph_type",
}

// TODO: status code?
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
];

export const SpanFilter = () => {
	// clients
	const { queryClient } = useContext(ClientContext);

	// local state
	const [modalVisible, setModalVisible] = useState(false);
	const [spanFilters, setSpanFilters] = useState<FilteredQuery>({
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
	const operationNames = useSelector((state: RootState) => state.operationNames);
	const tags = useSelector((state: RootState) => state.tags);
	const customMetrics = useSelector((state: RootState) => state.customMetrics);

	const dispatch: AppDispatch = useDispatch();

	// retrieve filtered spans
	useEffect(() => {
		if (
			spanFilters.service ||
			spanFilters.kind ||
			spanFilters.operation ||
			(spanFilters.duration &&
				(spanFilters.duration.min || spanFilters.duration.max)) ||
			spanFilters.tags.length !== 0
		) {
			dispatch(SpansUpdater.Spans(queryClient, maxMinTime, spanFilters));
		}
	}, [dispatch, queryClient, maxMinTime, spanFilters]);

	// retrieve service names
	useEffect(() => {
		dispatch(ServicesUpdater.ServiceNames(queryClient));
	}, [dispatch, queryClient]);

	const [initialForm] = Form.useForm();

	// filter on service
	const onChangeService = (value: string) => {
		setSpanFilters({ ...spanFilters, service: value });
	};

	useEffect(() => {
		if (spanFilters.service) {
			dispatch(
				ServiceUpdater.OperationNames(queryClient, spanFilters.service ?? ""),
			);
			dispatch(ServiceUpdater.Tags(queryClient, spanFilters.service ?? ""));
			initialForm.setFieldsValue({ service: spanFilters.service });
		}
	}, [dispatch, initialForm, queryClient, spanFilters.service]);

	// filter on operation
	const onChangeOperation = (value: string) => {
		setSpanFilters({ ...spanFilters, operation: value });
	};

	useEffect(() => {
		if (spanFilters.operation) {
			initialForm.setFieldsValue({ operation: spanFilters.operation });
		}
	}, [initialForm, spanFilters.operation]);

	// filter on span kind
	const onChangeKind = (value: SpanKind) => {
		setSpanFilters({ ...spanFilters, kind: value });
	};

	useEffect(() => {
		if (spanFilters.kind) {
			initialForm.setFieldsValue({ kind: spanFilters.kind });
		}
	}, [initialForm, spanFilters.kind]);

	// filter on duration
	const onClickDuration = () => {
		setModalVisible(true);
	};

	const onCreateDuration = (values: Store) => {
		setModalVisible(false);

		const { min, max } = values;

		setSpanFilters({
			...spanFilters,
			duration: {
				min: min ? (Number(min) * 1000000).toString() : "",
				max: max ? (Number(max) * 1000000).toString() : "",
			},
		});
	};

	useEffect(() => {
		let durationButtonText = "Duration";

		if (spanFilters.duration) {
			if (spanFilters.duration.min === "" && spanFilters.duration.max !== "") {
				durationButtonText = `Duration< ${Number(spanFilters.duration.max) / 1000000}ms`;
			} else if (
				spanFilters.duration.min !== "" &&
				spanFilters.duration.max === ""
			) {
				durationButtonText = `Duration> ${Number(spanFilters.duration.min) / 1000000}ms`;
			} else if (
				spanFilters.duration.min !== "" &&
				spanFilters.duration.max !== ""
			) {
				durationButtonText = `${Number(spanFilters.duration.min) / 1000000}ms <Duration< ${Number(spanFilters.duration.max) / 1000000}ms`;
			}
		}

		initialForm.setFieldsValue({ duration: durationButtonText });
	}, [initialForm, spanFilters.duration]);

	const [tagForm] = Form.useForm();

	// filter on tags
	const onTagFormSubmit = (values: any) => {
		setSpanFilters({
			...spanFilters,
			tags: [
				...spanFilters.tags,
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
				setSpanFilters({ ...spanFilters, service: "" });
				break;
			case TagValue.operation:
				setSpanFilters({ ...spanFilters, operation: "" });
				break;
			case TagValue.kind:
				setSpanFilters({ ...spanFilters, kind: SpanKind.default });
				break;
			case TagValue.minDuration:
				setSpanFilters({
					...spanFilters,
					duration: { min: spanFilters.duration?.min ?? "", max: "" },
				});
				break;
			case TagValue.maxDuration:
				setSpanFilters({
					...spanFilters,
					duration: { min: "", max: spanFilters.duration?.max ?? "" },
				});
				break;
		}
	};

	const onCloseTagTag = (tag: {
		key: string;
		value: string;
		operator: Operator;
	}) => {
		setSpanFilters({
			...spanFilters,
			tags: spanFilters.tags?.filter((t) => {
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
			SpansUpdater.CustomMetrics(
				queryClient,
				{
					minTime: maxMinTime.minTime - 15 * 60,
					maxTime: maxMinTime.maxTime + 15 * 60,
				},
				dimension,
				aggregation,
				spanFilters,
			),
		);
	}, [dispatch, queryClient, maxMinTime, dimension, aggregation, spanFilters]);

	return (
		<div>
			<Card>
				<div>Filter Spans</div>
				<Form
					form={initialForm}
					layout="inline"
					initialValues={{ service: "", operation: "", duration: "Duration" }}
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
					<FormItem name="operation">
						<Select
							showSearch
							style={{ width: 180 }}
							onChange={onChangeOperation}
							placeholder="Select Operation"
							allowClear
						>
							{operationNames.map((name: string, idx: number) => {
								return (
									<Select.Option value={name} key={idx}>
										{name}
									</Select.Option>
								);
							})}
						</Select>
					</FormItem>
					<FormItem name="kind">
						<Select
							showSearch
							style={{ width: 180 }}
							onChange={onChangeKind}
							placeholder="Select Span Kind"
							allowClear
						>
							{[SpanKind.server, SpanKind.client].map(
								(name: SpanKind, idx: number) => {
									return (
										<Select.Option value={name} key={idx}>
											{name}
										</Select.Option>
									);
								},
							)}
						</Select>
					</FormItem>
					<FormItem name="duration">
						<Input style={{ width: 180 }} type="button" onClick={onClickDuration} />
					</FormItem>
				</Form>
				<Card
					style={{ padding: 6, marginTop: 10, marginBottom: 10 }}
					bodyStyle={{ padding: 6 }}
				>
					{!spanFilters.service ? null : (
						<Tag
							style={{ fontSize: 14, padding: 8 }}
							closable
							onClose={() => onCloseTag(TagValue.service)}
						>
							service:{spanFilters.service}
						</Tag>
					)}
					{!spanFilters.operation ? null : (
						<Tag
							style={{ fontSize: 14, padding: 8 }}
							closable
							onClose={() => onCloseTag(TagValue.operation)}
						>
							operation:{spanFilters.operation}
						</Tag>
					)}
					{!spanFilters.kind ? null : (
						<Tag
							style={{ fontSize: 14, padding: 8 }}
							closable
							onClose={() => onCloseTag(TagValue.kind)}
						>
							kind:{spanFilters.kind}
						</Tag>
					)}
					{!spanFilters.duration || !spanFilters.duration.min ? null : (
						<Tag
							style={{ fontSize: 14, padding: 8 }}
							closable
							onClose={() => onCloseTag(TagValue.minDuration)}
						>
							minDuration:
							{(Number(spanFilters.duration.min) / 1000000).toString()}ms
						</Tag>
					)}
					{!spanFilters.duration || !spanFilters.duration.max ? null : (
						<Tag
							style={{ fontSize: 14, padding: 8 }}
							closable
							onClose={() => onCloseTag(TagValue.maxDuration)}
						>
							maxDuration:
							{(Number(spanFilters.duration.max) / 1000000).toString()}ms
						</Tag>
					)}
					{!spanFilters.tags
						? null
						: spanFilters.tags
								.filter((t, i) => {
									const found = spanFilters.tags.findIndex((e) => {
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
				{modalVisible && (
					<DurationModelForm
						onCreate={onCreateDuration}
						durationFilterValues={
							spanFilters.duration &&
							(spanFilters.duration.min || spanFilters.duration.max)
								? {
										min: (Number(spanFilters.duration.min) / 1000000).toString() || "",
										max: (Number(spanFilters.duration.max) / 1000000).toString() || "",
									}
								: { min: "", max: "" }
						}
						onCancel={() => {
							setModalVisible(false);
						}}
					/>
				)}
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
