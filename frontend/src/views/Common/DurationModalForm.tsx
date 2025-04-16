import { Col, Form, InputNumber, Modal, Row } from "antd";
import { Store } from "antd/es/form/interface";

interface DurationModalFormProps {
	onCreate: (values: Store) => void;
	onCancel: () => void;
	durationFilterValues: { min: string; max: string };
}

export const DurationModelForm = (props: DurationModalFormProps) => {
	const [form] = Form.useForm();

	return (
		<Modal
			visible={true}
			title="Choose min and max values of duration"
			okText="Apply"
			cancelText="Cancel"
			onCancel={props.onCancel}
			onOk={() => {
				form
					.validateFields()
					.then((values) => {
						props.onCreate(values);
						form.resetFields();
					})
					.catch((err) => {
						console.log("Failed to parse duration modal form fields: ", err);
					});
			}}
		>
			<Form
				form={form}
				layout="horizontal"
				name="form_in_modal"
				initialValues={props.durationFilterValues}
			>
				<Row>
					<Col span={12}>
						<Form.Item name="min" label="Min (in ms)">
							<InputNumber />
						</Form.Item>
					</Col>
					<Col span={12}>
						<Form.Item name="max" label="Max (in ms)">
							<InputNumber />
						</Form.Item>
					</Col>
				</Row>
			</Form>
		</Modal>
	);
};
