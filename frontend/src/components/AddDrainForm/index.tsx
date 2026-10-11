import { MinusCircleOutlined, PlusOutlined } from "@ant-design/icons";
import { Button, Form, Input, Popconfirm, Space } from "antd";
import type React from "react";
import { workerIdPatternFromFields } from "@/components/DrainsTable/utils";

interface Field {
  key: string;
  value: string;
}

interface Props {
  onSubmit: (
    workerIdPattern: Record<string, string>,
    onSuccess: () => void,
  ) => void;
  loading: boolean;
}

const AddDrainForm: React.FC<Props> = ({ onSubmit, loading }) => {
  const [form] = Form.useForm<{ fields: Field[] }>();
  const fields = Form.useWatch("fields", form);

  return (
    <Form
      form={form}
      onFinish={(values) =>
        onSubmit(workerIdPatternFromFields(values.fields ?? []), () =>
          form.resetFields(),
        )
      }
    >
      <Form.List name="fields" initialValue={[{ key: "", value: "" }]}>
        {(listFields, { add, remove }) => (
          <>
            {listFields.map(({ key, name }) => (
              <Space key={key} align="baseline">
                <Form.Item
                  name={[name, "key"]}
                  dependencies={listFields
                    .filter((other) => other.name !== name)
                    .map((other) => ["fields", other.name, "key"])}
                  rules={[
                    {
                      required: true,
                      whitespace: true,
                      message: "Missing key",
                    },
                    {
                      validator: (_, value: string | undefined) =>
                        value &&
                        (form.getFieldValue("fields") as Field[]).filter(
                          (field) => field?.key === value,
                        ).length > 1
                          ? Promise.reject(new Error("Duplicate key"))
                          : Promise.resolve(),
                    },
                  ]}
                >
                  <Input placeholder="hostname" />
                </Form.Item>
                <Form.Item name={[name, "value"]}>
                  <Input placeholder="server123.example.com" />
                </Form.Item>
                <MinusCircleOutlined onClick={() => remove(name)} />
              </Space>
            ))}
            <Form.Item>
              <Button
                type="dashed"
                onClick={() => add({ key: "", value: "" })}
                icon={<PlusOutlined />}
              >
                Add field
              </Button>
            </Form.Item>
          </>
        )}
      </Form.List>
      <Form.Item>
        <Popconfirm
          title="Create drain?"
          description={
            fields?.length
              ? "Matching workers finish their current operation and receive no new ones."
              : "This pattern is empty and drains every worker in this size class queue."
          }
          onConfirm={() => form.submit()}
        >
          <Button type="primary" loading={loading}>
            Create drain
          </Button>
        </Popconfirm>
      </Form.Item>
    </Form>
  );
};

export default AddDrainForm;
