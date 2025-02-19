import { Button, Col, Divider, Row } from "antd";
import { usePathname, useRouter } from "next/navigation";
import styles from "./index.module.css";

interface Props {
  filterInvocationId: string | null;
}

const OperationsInvocationFilter: React.FC<Props> = ({
  filterInvocationId,
}) => {
  const router = useRouter();
  const pathname = usePathname();

  const handleFilterClear = () => router.replace(pathname);

  if (filterInvocationId) {
    return (
      <Row>
        <Col span={4} className={styles.alignLeft}>
          <h3>Invocation ID:</h3>
        </Col>
        <Col span={16} className={styles.alignCenter}>
          <pre>{decodeURIComponent(filterInvocationId)}</pre>
        </Col>
        <Col span={4} className={styles.alignRight}>
          <Button type="primary" onClick={handleFilterClear}>
            Clear Invocation ID Filter
          </Button>
        </Col>
        <Divider />
      </Row>
    );
  }
};

export default OperationsInvocationFilter;
