import React, { forwardRef, useLayoutEffect, useRef, useState } from 'react';
import { Button, Card, CardProps, List, Popover, Space, theme } from 'antd';
import { MenuFoldOutlined, MenuUnfoldOutlined } from '@ant-design/icons';
import styles from './index.module.css';
import themeStyles from '@/theme/theme.module.css';

const { useToken } = theme;

interface HeaderProps {
  headerBits: React.ReactNode[];
  className?: string;
}

const Header = forwardRef<HTMLDivElement, HeaderProps>(({ headerBits, className }, ref) => {
  const actualClassName = [styles.header, className].join(' ');
  const actualTitleBits = headerBits.filter(titleBit => titleBit);
  return (
    <Space ref={ref} size="middle" className={actualClassName}>
      {actualTitleBits.map(
        (titleBit, index) =>
          titleBit && (
            <div key={index} className={actualClassName}>
              {titleBit}
            </div>
          ),
      )}
    </Space>
  );
});
Header.displayName = 'Header';

interface ExtraMenuProps {
  extraBits: React.ReactNode[];
  isExtraMenuOpen: boolean;
  setIsExtraMenuOpened: React.Dispatch<React.SetStateAction<boolean>>;
}

const ExtraMenu: React.FC<ExtraMenuProps> = ({ extraBits, isExtraMenuOpen, setIsExtraMenuOpened }) => {
  return (
    <Popover
      content={
        <List
          dataSource={extraBits}
          renderItem={item => <List.Item className={styles.item}>{item}</List.Item>}
          size="small"
          className={styles.list}
        />
      }
      trigger="click"
      onOpenChange={() => setIsExtraMenuOpened(!isExtraMenuOpen)}
      open={isExtraMenuOpen}
      placement="leftTop"
    >
      <Button className={styles.button}>
        {isExtraMenuOpen ? <MenuFoldOutlined rotate={180} /> : <MenuUnfoldOutlined rotate={180} />}
      </Button>
    </Popover>
  );
};

interface Props extends Omit<CardProps, 'title' | 'extra'> {
  icon: React.ReactNode;
  titleBits: React.ReactNode[];
  extraBits?: React.ReactNode[];
  className?: string;
}

export const PortalCard: React.FC<Props> = ({ icon, titleBits, extraBits, className, ...cardProps }) => {
  const { token } = useToken();
  const [isExtraMenuOpen, setIsExtraMenuOpened] = useState<boolean>(false);
  const [isExtraMenuDisplayed, setIsExtraMenuDisplayed] = useState<boolean>(true);
  const cardRef = useRef<HTMLDivElement>(null);
  const titleRef = useRef<HTMLDivElement>(null);
  const extraRef = useRef<HTMLDivElement>(null);
  const minimumSpaceBetweenTitleAndExtra = token.paddingLG * 2;
  useLayoutEffect(() => {
    const getShouldExtraMenuBeDisplayed = () => {
      if (cardRef.current && titleRef.current && extraRef.current) {
        return (
          cardRef.current.clientWidth <
          titleRef.current.clientWidth + extraRef.current.clientWidth + minimumSpaceBetweenTitleAndExtra
        );
      }
      return true;
    };
    const handleResize = () => {
      setIsExtraMenuDisplayed(getShouldExtraMenuBeDisplayed());
      if (!isExtraMenuDisplayed) {
        setIsExtraMenuOpened(false);
      }
    };
    handleResize();
    window.addEventListener('resize', handleResize);
    return () => {
      window.removeEventListener('resize', handleResize);
    };
  }, [isExtraMenuDisplayed, minimumSpaceBetweenTitleAndExtra]);
  const title = <Header ref={titleRef} headerBits={[icon, ...titleBits]} className={styles.title} />;
  const extra = extraBits?.length && (
    <div className={styles.extra}>
      <span className={isExtraMenuDisplayed ? styles.hidden : styles.visible}>
        <Header ref={extraRef} headerBits={extraBits ?? []} />
      </span>
      <span className={isExtraMenuDisplayed ? styles.visible : styles.hidden}>
        <ExtraMenu
          extraBits={extraBits}
          isExtraMenuOpen={isExtraMenuOpen}
          setIsExtraMenuOpened={setIsExtraMenuOpened}
        />
      </span>
    </div>
  );
  const extendedClassName = cardProps.type === 'inner' ? className : [className, styles.outer].join(' ');
  return (
    <Card title={title} extra={extra} bordered={false} className={extendedClassName} {...cardProps}>
      <Space ref={cardRef} direction="vertical" size="middle" className={themeStyles.space}>
        {cardProps.children}
      </Space>
    </Card>
  );
};

export default PortalCard;
