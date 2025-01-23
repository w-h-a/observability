import { useContext, useEffect, useState } from "react";
import { useDispatch, useSelector } from "react-redux";
import { AutoComplete, Button, Card, Form, Input, Select, Tag } from "antd";
import { Store } from "antd/es/form/interface";
import FormItem from "antd/es/form/FormItem";
import { DurationModelForm } from "./DurationModalForm";
import { AppDispatch, RootState } from "../../updaters/store";
import { SpansUpdater } from "../../updaters/spans/spans";
import { ServiceUpdater } from "../../updaters/service/service";
import { ServicesUpdater } from "../../updaters/services/services";
import { ClientContext } from "../../clients/query/clientCtx";
import { FilteredQuery, Operator } from "../../clients/query/filteredQuery";

// TODO: status code
enum TagValue {
	service = "service",
	operation = "operation",
	minDuration = "minDuration",
	maxDuration = "maxDuration",
}

export const SpanFilter = () => {
	// clients
	const { queryClient } = useContext(ClientContext);

	// local state
	const [modalVisible, setModalVisible] = useState(false);
	const [spanFilters, setSpanFilters] = useState<FilteredQuery>({
		service: "",
		operation: "",
		duration: { min: "", max: "" },
		tags: [],
	});

	// store state
	const maxMinTime = useSelector((state: RootState) => state.maxMinTime);
	const serviceNames = useSelector((state: RootState) => state.serviceNames);
	const operationNames = useSelector((state: RootState) => state.operationNames);
	const tags = useSelector((state: RootState) => state.tags);

	const dispatch: AppDispatch = useDispatch();

	// retrieve filtered spans
	useEffect(() => {
		if (
			spanFilters.service ||
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

	return (
		<div>
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
		</div>
	);
};
