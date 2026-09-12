import {
  FileOutlined,
  FolderOpenOutlined,
  FolderOutlined,
  LinkOutlined,
  StopOutlined,
} from "@ant-design/icons";
import { Button, theme } from "antd";
import { Spin } from "antd/lib";
import React, {
  useCallback,
  useContext,
  useEffect,
  useRef,
  useState,
} from "react";
import { LinkButton } from "@/components/LinkButton";
import { readableFileSize, readableFileSizeFromString } from "@/utils/filesize";
import { CompareContext } from ".";
import styles from "./directoryComponent.module.css";
import type { DirectoryRenderData, FileRenderData, NodeType } from "./types";
import { formattedFileName, styleMap } from "./utils";

const { useToken } = theme;

interface DirectoryComponentProps {
  baseName?: string;
  name: string;
  renderData?: FileRenderData;
}
const INDENT = 10;

export const TopDirectoryComponent: React.FC<DirectoryComponentProps> = ({
  name,
}) => {
  const [focusedNodeId, setFocusedNodeId] = React.useState<string | null>(null);
  const rootReference = useRef<HTMLDivElement>(null);
  const { getData } = useContext(CompareContext);
  const fullPath = name;
  const [data, setData] = useState<DirectoryRenderData | null>();

  useEffect(() => {
    getData(fullPath).then((data) => {
      setData(data);
    });
  }, [fullPath, getData]);

  const handleKeyDown = (event: React.KeyboardEvent<HTMLDivElement>) => {
    const root = rootReference.current;
    if (!root) return;

    const current = event.target as HTMLElement;
    const treeitem = current.closest('[role="treeitem"]') as HTMLElement | null;
    if (!treeitem) return;

    switch (event.key) {
      case "ArrowDown": // Intentional fall-though
      case "ArrowUp": {
        event.preventDefault();
        const items = Array.from(
          root.querySelectorAll<HTMLElement>('[role="treeitem"]'),
        );
        const index = items.indexOf(treeitem);
        if (index !== -1) {
          const offset = event.key === "ArrowDown" ? 1 : -1;
          items[index + offset]?.focus();
        }
        break;
      }

      case "ArrowRight": {
        event.preventDefault();
        const isExpanded = treeitem.getAttribute("aria-expanded");
        if (isExpanded === "false") {
          treeitem.click();
        } else if (isExpanded === "true") {
          const nodeOuterDiv = treeitem.parentElement;
          const group = nodeOuterDiv?.nextElementSibling;
          const firstChild =
            group?.querySelector<HTMLElement>('[role="treeitem"]');
          firstChild?.focus();
        }
        break;
      }

      case "ArrowLeft": {
        event.preventDefault();
        const isExpanded = treeitem.getAttribute("aria-expanded");
        if (isExpanded === "true") {
          treeitem.click();
        } else {
          const group = treeitem.closest('[role="group"]');
          const parentOuterDiv = group?.previousElementSibling;
          const parentTreeitem =
            parentOuterDiv?.querySelector<HTMLElement>('[role="treeitem"]');
          parentTreeitem?.focus();
        }
        break;
      }
      case " ": {
        event.preventDefault();
        treeitem.click();
        break;
      }
    }
  };

  const depth = fullPath.split("/").length;
  if (!data) {
    return (
      <div style={{ marginLeft: depth * INDENT }}>
        <Button type="text" disabled={true}>
          {name}
          <Spin />
        </Button>
      </div>
    );
  }

  return (
    <div
      role="tree"
      aria-label="Browser directory tree"
      ref={rootReference}
      onKeyDown={handleKeyDown}
    >
      {/* biome-ignore lint/a11y/useSemanticElements: role="group" is required for ARIA tree structure. Fieldset has default styling we don't want to use. */}
      <div role="group">
        {data.directories.map((entry) => (
          <DirectoryComponent
            key={entry.name}
            baseName={fullPath}
            name={entry.name}
            renderData={entry}
          />
        ))}
        {data.files.map((value) => {
          const path = `${fullPath}/${value.name}`;
          return (
            <DirectoryNode
              key={value.name}
              file={value}
              depth={depth + 1}
              role="treeitem"
              tabIndex={focusedNodeId === path ? 0 : -1}
              onFocus={() => setFocusedNodeId(path)}
            />
          );
        })}
        {data.symlinks.map((value) => {
          const path = `${fullPath}/${value.name}`;
          return (
            <DirectoryNode
              key={value.name}
              file={value}
              depth={depth + 1}
              role="treeitem"
              tabIndex={focusedNodeId === path ? 0 : -1}
              onFocus={() => setFocusedNodeId(path)}
            />
          );
        })}
      </div>
    </div>
  );
};

const DirectoryComponent: React.FC<DirectoryComponentProps> = ({
  baseName,
  name,
  renderData,
}) => {
  const { closeDir, openDir, getData, openDirs } = useContext(CompareContext);
  const fullPath = baseName ? `${baseName}/${name}` : name;
  const expanded = openDirs.has(fullPath);
  const [data, setData] = useState<DirectoryRenderData | null>();
  const [focusedNodeId, setFocusedNodeId] = useState<string | null>(null);
  const toggleExpanded = useCallback(
    (event: React.MouseEvent) => {
      if (event.shiftKey || event.ctrlKey || event.metaKey || event.altKey) {
        return; // The user is trying to open the link to the directory
      }
      if (expanded) {
        closeDir(fullPath);
      } else {
        openDir(fullPath);
      }
    },
    [fullPath, closeDir, openDir, expanded],
  );

  useEffect(() => {
    getData(fullPath).then((data) => {
      setData(data);
    });
  }, [fullPath, getData]);

  const depth = fullPath.split("/").length;
  if (!data) {
    return (
      <div style={{ marginLeft: depth * INDENT }}>
        <Button type="text" disabled={true}>
          {name}
          <Spin />
        </Button>
      </div>
    );
  }
  if (renderData && renderData?.nodeType !== "aborted") {
    renderData.nodeType = expanded ? "directory_open" : "directory";
  }

  return (
    <div>
      {renderData && (
        <DirectoryNode
          file={renderData}
          depth={depth}
          onClick={toggleExpanded}
          role="treeitem"
          aria-expanded={expanded}
          tabIndex={focusedNodeId === fullPath ? 0 : -1}
          onFocus={() => setFocusedNodeId(fullPath)}
        />
      )}
      {expanded && data && (
        // biome-ignore lint/a11y/useSemanticElements: role="group" is required for ARIA tree structure. Fieldset has default styling we don't want to use.
        <div role="group">
          {data.directories.map((entry) => (
            <DirectoryComponent
              key={entry.name}
              baseName={fullPath}
              name={entry.name}
              renderData={entry}
            />
          ))}
          {data.files.map((value) => {
            const path = `${fullPath}/${value.name}`;
            return (
              <DirectoryNode
                key={value.name}
                file={value}
                depth={depth + 1}
                role="treeitem"
                tabIndex={focusedNodeId === path ? 0 : -1}
                onFocus={() => setFocusedNodeId(path)}
              />
            );
          })}
          {data.symlinks.map((value) => {
            const path = `${fullPath}/${value.name}`;
            return (
              <DirectoryNode
                key={value.name}
                file={value}
                depth={depth + 1}
                role="treeitem"
                tabIndex={focusedNodeId === path ? 0 : -1}
                onFocus={() => setFocusedNodeId(path)}
              />
            );
          })}
        </div>
      )}
    </div>
  );
};

const iconMap: Record<NodeType, React.ReactNode> = {
  file: <FileOutlined />,
  symlink: <LinkOutlined />,
  directory: <FolderOutlined />,
  directory_open: <FolderOpenOutlined />,
  aborted: <StopOutlined />,
};

interface FileProps {
  file: FileRenderData;
  depth: number;
  role?: React.AriaRole;
  tabIndex?: number;
  "aria-expanded"?: boolean;
  onClick?: React.MouseEventHandler<HTMLButtonElement | HTMLAnchorElement>;
  onFocus?: React.FocusEventHandler<HTMLButtonElement | HTMLAnchorElement>;
}

const DirectoryNode = React.forwardRef<
  HTMLButtonElement | HTMLAnchorElement,
  FileProps
>(
  (
    {
      file,
      depth,
      role,
      tabIndex,
      "aria-expanded": ariaExpanded,
      onClick,
      onFocus,
    },
    ref,
  ) => {
    const { token } = useToken();
    const { mergedMode } = useContext(CompareContext);

    const commonProps = {
      ref,
      role: role,
      tabIndex: tabIndex,
      "aria-expanded": ariaExpanded,
      onFocus: onFocus,
      onClick: (
        event: React.MouseEvent<HTMLButtonElement | HTMLAnchorElement>,
      ) => {
        onClick?.(event);

        if (
          !event.shiftKey &&
          !event.ctrlKey &&
          !event.metaKey &&
          !event.altKey &&
          (file.nodeType === "directory" || file.nodeType === "directory_open")
        ) {
          event.preventDefault(); // Prevent the user from following the link when they try to open/close a folder
        }
      },
      style: styleMap(file.compareState, token),
      className: styles.directoryComponent,
    };
    const innerContent = [
      <div className={styles.fileName} key="name">
        {iconMap[file.nodeType] || null}
        &nbsp;&nbsp;&nbsp;
        {formattedFileName(file.name, file.willBePrefetched)}
        {file.target && ` -> `}
        <span
          style={
            file.targetCompare === "different"
              ? styleMap("sub_different", token)
              : {}
          }
        >
          {file.target && file.target}
        </span>
      </div>,
      <div key="size">
        {!!file.hash || !!file.compareHash ? (
          <TotalSizeDisplay hash={file.hash} compareHash={file.compareHash} />
        ) : (
          <SizeDisplay
            sizeBytes={file.sizeBytes}
            compareSizeBytes={file.compareSizeBytes}
          />
        )}
        &nbsp;&nbsp;&nbsp;&nbsp;
        {(file.compareState !== "missing" || mergedMode) && (
          <code
            style={
              file.permissionsCompare === "different"
                ? styleMap("sub_different", token)
                : {}
            }
          >
            {file.permissions}
          </code>
        )}
      </div>,
    ];

    return (
      <div
        style={{
          marginLeft: (depth - 2) * INDENT,
        }}
      >
        {typeof file.linkOptions === "string" ? (
          <Button
            type="text"
            href={file.linkOptions}
            target="_blank"
            download={true}
            {...commonProps}
          >
            {innerContent}
          </Button>
        ) : (
          <LinkButton buttonType="text" {...file.linkOptions} {...commonProps}>
            {innerContent}
          </LinkButton>
        )}
      </div>
    );
  },
);
const SizeDisplay: React.FC<{
  sizeBytes?: string;
  compareSizeBytes?: string;
}> = ({ sizeBytes, compareSizeBytes }) => {
  const { mergedMode } = useContext(CompareContext);
  const different: boolean = sizeBytes !== compareSizeBytes;
  const displayBase: boolean = !!sizeBytes;
  const displayCompare: boolean = different && !!compareSizeBytes && mergedMode;
  return (
    <code>
      {displayBase && sizeBytes && readableFileSizeFromString(sizeBytes)}
      {displayBase && displayCompare && " / "}
      {displayCompare &&
        compareSizeBytes &&
        readableFileSizeFromString(compareSizeBytes)}
    </code>
  );
};

const TotalSizeDisplay: React.FC<{
  hash?: string;
  compareHash?: string;
}> = ({ hash, compareHash }) => {
  const { token } = theme.useToken();
  const { getTotalSizeData, mergedMode } = useContext(CompareContext);
  const totalSizeData = hash ? getTotalSizeData(hash) : undefined;
  const compareSizeData = compareHash
    ? getTotalSizeData(compareHash)
    : undefined;
  if (!totalSizeData?.totalSize && !compareSizeData?.totalSize) {
    return null;
  }
  const different: boolean =
    totalSizeData?.totalSize !== compareSizeData?.totalSize;
  const displayBase: boolean = !!totalSizeData?.totalSize;
  const displayCompare: boolean =
    different && !!compareSizeData?.totalSize && mergedMode;
  if (!displayBase && !displayCompare) {
    return null;
  }
  return (
    <code>
      {displayBase && (
        <span
          className={
            displayCompare && mergedMode ? styles.totalSizeBadge : undefined
          }
          style={displayCompare && mergedMode ? styleMap("unique", token) : {}}
        >
          <span title={"Total size may be larger than what is displayed here"}>
            {totalSizeData?.fullyResolved === false ? ">" : ""}
          </span>
          {totalSizeData && readableFileSize(totalSizeData.totalSize, 3)}
        </span>
      )}
      {displayBase && displayCompare && " | "}
      {displayCompare && (
        <span
          className={
            displayBase && mergedMode ? styles.totalSizeBadge : undefined
          }
          style={displayBase && mergedMode ? styleMap("missing", token) : {}}
        >
          <span title={"Total size may be larger than what is displayed here"}>
            {compareSizeData?.fullyResolved === false ? ">" : ""}
          </span>
          {compareSizeData && readableFileSize(compareSizeData.totalSize, 3)}
        </span>
      )}
    </code>
  );
};
