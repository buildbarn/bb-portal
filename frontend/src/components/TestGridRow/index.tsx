
import React, { useState } from 'react';
import { useQuery } from '@apollo/client';
import FIND_TESTS from '@/app/tests/index.graphql';
import { FindTestsQueryVariables, OrderDirection, TestCollection, TestCollectionOrderField, TestSummary } from '@/graphql/__generated__/graphql';
import { TestStatusType } from '../TestGrid'
import { Space, Row, Statistic, Table } from 'antd';
import TestGridBtn from '../TestGridBtn';
import { TestStatusEnum } from '../TestStatusTag';


interface Props {
    rowLabel: string;
    first: number,
    reverseOrder: boolean
}

const TestGridRow: React.FC<Props> = ({ rowLabel, first, reverseOrder }) => {

    var { loading, data, previousData, error } = useQuery(FIND_TESTS, {
        variables: {
            first: first,
            where: { label: rowLabel },
            orderBy: {
                direction: OrderDirection.Desc,
                field: TestCollectionOrderField.FirstSeen
            }
        }, fetchPolicy: 'cache-and-network'
    });

    var activeData = loading ? previousData : data;
    let rowDataSrc: TestCollection[] = []

    if (error) {
        rowDataSrc = [];
    } else {
        const rowTestData = activeData?.findTests.edges?.flatMap(edge => edge?.node) ?? [];
        rowDataSrc = rowTestData.filter((x): x is TestCollection => !!x);
        if (reverseOrder) {
            rowDataSrc.reverse();
        }
    }
    return (
        <Space size="middle">
            {rowDataSrc.map((item) => (
                <TestGridBtn key={"test-grid-btn" + item.id} invocationId={item.bazelInvocation?.invocationID} status={item.overallStatus as TestStatusEnum} />
            ))}
        </Space>
    );
};

export default TestGridRow;