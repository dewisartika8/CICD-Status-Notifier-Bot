import React from 'react';
import Layout from '../components/common/Layout';
import SuccessRateChart from '../components/analytics/SuccessRateChart';
import BuildDurationChart from '../components/analytics/BuildDurationChart';
import TrendsChart from '../components/analytics/TrendsChart';

const Analytics = () => {
    return (
        <Layout>
            <h1>Analytics Dashboard</h1>
            <div className="analytics-charts">
                <SuccessRateChart />
                <BuildDurationChart />
                <TrendsChart />
            </div>
        </Layout>
    );
};

export default Analytics;