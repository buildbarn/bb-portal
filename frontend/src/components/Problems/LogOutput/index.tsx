import React, {useEffect, useState} from 'react';
import {util} from 'zod';
import {BlobReference} from '@/graphql/__generated__/graphql';
import {LogViewerCard} from "@/components/LogViewer";

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

    return <LogViewerCard log={validContent}/>;
};

export default LogOutput;
