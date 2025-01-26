import { useEffect, useState } from "react";
import { useDispatch } from "react-redux";
import { Button, Form, Select, Space } from "antd";
import FormItem from "antd/es/form/FormItem";
import { AppDispatch } from "../../updaters/store";
import { TimeUpdater } from "../../updaters/time/time";

const options = [
	// { value: "custom", key: "custom", label: "Custom" },
	{ value: "15min", key: "15min", label: "Last 15 min" },
	{ value: "30min", key: "30min", label: "Last 30 min" },
	{ value: "1hr", key: "1hr", label: "Last 1 hour" },
	{ value: "6hr", key: "6hr", label: "Last 6 hour" },
	{ value: "1day", key: "1day", label: "Last 1 day" },
	{ value: "1week", key: "1week", label: "Last 1 week" },
];

export const DateTimeSelector = () => {
	const [timeInterval, setTimeInterval] = useState("15min");

	const dispatch: AppDispatch = useDispatch();

	const [form] = Form.useForm();

	const onIntervalSelect = (value: string) => {
		setTimeInterval(value);
	};

	const onRefreshClick = () => {
		setTimeInterval(timeInterval);
	};

	useEffect(() => {
		dispatch(TimeUpdater.Time(timeInterval));
	}, [dispatch, timeInterval]);

	return (
		<div>
			<Space style={{ float: "right", display: "block" }}>
				<Space>
					<Form
						form={form}
						layout="inline"
						initialValues={{ interval: "15min" }}
						style={{ marginTop: 10, marginBottom: 10 }}
					>
						<FormItem name="interval">
							<Select onSelect={onIntervalSelect} value={timeInterval}>
								{options.map((opt) => {
									return (
										<Select.Option value={opt.value} key={opt.key}>
											{opt.label}
										</Select.Option>
									);
								})}
							</Select>
						</FormItem>
						<FormItem name="refresh">
							<Button type="primary" onClick={onRefreshClick}>
								refresh
							</Button>
						</FormItem>
					</Form>
				</Space>
			</Space>
		</div>
	);
};
