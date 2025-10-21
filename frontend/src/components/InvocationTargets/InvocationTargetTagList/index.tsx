import { Popover, Space, Tag } from "antd";

const TAGS_TO_SHOW = 3;

interface Props {
  tags: string[] | undefined | null;
}

export const InvocationTargetTagList: React.FC<Props> = ({ tags }) => {
  if (!tags || tags.length === 0) {
    return;
  }
  const sortedTags = tags.slice().sort();
  const splitIndex =
    sortedTags.length <= TAGS_TO_SHOW ? sortedTags.length : TAGS_TO_SHOW - 1;
  const trailingTagsLength = sortedTags.length - splitIndex;

  const tagList = (
    <Space direction="horizontal" size={0}>
      {sortedTags.slice(0, splitIndex).map((tag) => (
        <Tag key={tag} color="blue">
          {tag}
        </Tag>
      ))}
      {trailingTagsLength > 0 && <Tag color="blue">+{trailingTagsLength}</Tag>}
    </Space>
  );

  if (trailingTagsLength > 0) {
    return (
      <Popover
        placement="top"
        title={
          <Space direction="vertical" align="center" size={0}>
            {sortedTags.map((tag) => (
              <Tag key={tag} color="blue">
                {tag}
              </Tag>
            ))}
          </Space>
        }
      >
        {tagList}
      </Popover>
    );
  } else {
    return tagList;
  }
};
