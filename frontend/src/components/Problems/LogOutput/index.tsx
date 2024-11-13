import React, {useEffect, useState} from 'react';
import {util} from 'zod';
import {BlobReference} from '@/graphql/__generated__/graphql';
import {LogViewerCard} from "@/components/LogViewer";
import { env } from 'next-runtime-env';
import { Button, Card, Space, Tooltip } from 'antd';
import ButtonGroup from 'antd/es/button/button-group';
import { DownloadOutlined, WarningOutlined } from '@ant-design/icons';
import Link from 'next/link';
import styles from "@/components/LogViewer/index.module.css"
import DownloadButton from '@/components/DownloadButton';

interface Props {
    blobReference: BlobReference;
    tailBytes?: number;
}

const DEFAULT_TAIL_BYTES = 10_000;



/* eslint-disable consistent-return */
const LogOutput: React.FC<Props> = ({blobReference, tailBytes = DEFAULT_TAIL_BYTES}) => {
    // if (blobReference.ephemeralURL != "" && env('NEXT_PUBLIC_BROWSER_URL')) {
    //     const url = new URL(blobReference.ephemeralURL, env('NEXT_PUBLIC_BROWSER_URL'))
    //     var logDownload = <Space direction="horizontal" size="small">
    //       <Card>
    //         <ButtonGroup>
    //           <Button>
    //             <DownloadOutlined/>
    //             <Link href={url.toString()} download="test.log" target="_self">Download Log File</Link>
    //           </Button>
    //           <Tooltip title="Depending on the configuration of your remote cache, this link may point to ephemeral content and it could disappear.">
    //             <Button icon={<WarningOutlined/>} danger/>
    //           </Tooltip>
    //         </ButtonGroup>
    //       </Card>
    //     </Space>
    //   }
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

    if (blobReference.ephemeralURL != "" && env('NEXT_PUBLIC_BROWSER_URL')) {
        console.log("blobReference.ephemeralURL=", blobReference.ephemeralURL)
        const url = new URL(blobReference.ephemeralURL, env('NEXT_PUBLIC_BROWSER_URL'))
        return <Space direction="horizontal" size="small">
        <Card>
          <ButtonGroup>
            <DownloadButton url={url.toString()} enabled={true} buttonLabel="Download Log File" fileName="output.log" />
            <Tooltip title="Depending on the configuration of your remote cache, this link may point to ephemeral content and it could disappear.">
              <Button icon={<WarningOutlined/>} danger />
            </Tooltip>
          </ButtonGroup>
        </Card>
      </Space>
      }

    return <LogViewerCard log={validContent}/>;
};

export default LogOutput;
