
import React, { useState } from 'react';
import { useQuery } from '@apollo/client';
import FIND_TESTS from '@/app/tests/index.graphql';
import { FindTestsQueryVariables, OrderDirection, TargetPair, TargetPairOrderField, TestCollection, TestCollectionOrderField, TestSummary } from '@/graphql/__generated__/graphql';
import { TargetStatusType } from '../TargetGrid'
import { Space, Row, Statistic, Table } from 'antd';
import { TestStatusEnum } from '../../TestStatusTag';
import TargetGridBtn from '../TargetGridBtn';
import { FIND_TARGETS } from '@/app/targets/graphql';


interface Props {
    rowLabel: string;
    first: number,
    reverseOrder: boolean
}

const TargetGridRow: React.FC<Props> = ({ rowLabel, first, reverseOrder }) => {

    var { loading, data, previousData, error } = useQuery(FIND_TARGETS, {
        variables: {
            first: first,
            where: { label: rowLabel },
            orderBy: {
                direction: OrderDirection.Desc,
                field: TargetPairOrderField.Duration
            }
        }, fetchPolicy: 'cache-and-network'
    });

    var activeData = loading ? previousData : data;
    let rowDataSrc: TargetPair[] = []

    if (error) {
        rowDataSrc = [];
    } else {
        const rowTestData = activeData?.findTargets.edges?.flatMap(edge => edge?.node) ?? [];
        rowDataSrc = rowTestData.filter((x): x is TargetPair => !!x);
        if (reverseOrder) {
            rowDataSrc.reverse();
        }
    }
    return (
        <Space size="middle">
            {rowDataSrc.map((item) => (
                <TargetGridBtn key={"target-grid-btn" + item.id} invocationId={item.bazelInvocation?.invocationID} status={item.success ?? null} />
            ))}
        </Space>
    );
};

export default TargetGridRow;