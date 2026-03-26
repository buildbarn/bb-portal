import { Button, Col, Divider, Row } from "antd";
import styles from "./index.module.css";
import type { OperationsFilterParams } from "@/routes/operations.index";
import { useNavigate } from "@tanstack/react-router";

interface Props {
  filter: OperationsFilterParams;
}

const OperationsInvocationFilter: React.FC<Props> = ({ filter }) => {
  const navigate = useNavigate()
  if (filter) {
    return (
      <Row>
        <Col span={4} className={styles.alignLeft}>
          <h3>Invocation ID:</h3>
        </Col>
        <Col span={16} className={styles.alignLeft}>
          <pre>{JSON.stringify(filter, null, 2)}</pre>
        </Col>
        <Col span={4} className={styles.alignRight}>
          <Button type="primary" onClick={() => navigate({ search: undefined })}>
            Clear Invocation ID Filter
          </Button>
        </Col>
        <Divider />
      </Row>
    );
  }
};

export default OperationsInvocationFilter;
