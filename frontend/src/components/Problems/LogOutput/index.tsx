import DownloadButton from '@/components/DownloadButton';
import { LogViewerCard } from "@/components/LogViewer";
import { BlobReference } from '@/graphql/__generated__/graphql';
import { generateUrlFromEphemeralUrl } from '@/utils/urlGenerator';
import { WarningOutlined } from '@ant-design/icons';
import { Button, Card, Space, Tooltip } from 'antd';
import ButtonGroup from 'antd/es/button/button-group';
import React, { useEffect, useState } from 'react';

interface Props {
    blobReference: BlobReference;
    tailBytes?: number;
}

const DEFAULT_TAIL_BYTES = 10_000;



/* eslint-disable consistent-return */
const LogOutput: React.FC<Props> = ({blobReference, tailBytes = DEFAULT_TAIL_BYTES}) => {
    const [contents, setContents] = useState<string>("");
    useEffect(() => {
        fetch(blobReference.downloadURL, {
            headers: {"Range": `bytes=-${tailBytes}`},
        })
            .then(response => response.text())
            .then(data => {
                setContents(data)
            })
            .catch(error => {
                // Handle any errors
            });
    }, [blobReference.downloadURL, setContents, tailBytes])

    let validContent = contents;

    if (validContent.length >= DEFAULT_TAIL_BYTES) {
        // If tail is partial, we strip the first line since it could be incomplete.
        // 100 lines matches previous backend behavior.
        const lines = contents.split('\n');

        // We know tail is incomplete content, so we request one more and always drop the first.
        const selectedLines = lines.slice(-101);
        selectedLines.shift();
        validContent = selectedLines.join('\n')
    }

    if (blobReference.ephemeralURL != "") {
        const url = generateUrlFromEphemeralUrl("", blobReference.ephemeralURL);
        return <Space direction="horizontal" size="small">
            <Card>
                <ButtonGroup>
                    <DownloadButton url={url} enabled={true} buttonLabel="Download Log File" fileName="output.log" />
                    <Tooltip title="Depending on the configuration of your remote cache, this link may point to ephemeral content and it could disappear.">
                        <Button icon={<WarningOutlined />} danger />
                    </Tooltip>
                </ButtonGroup>
            </Card>
        </Space>
    }

    return <LogViewerCard log={validContent}/>;
};

export default LogOutput;
