import React, { useEffect, useState } from 'react';
import { fetchDashboardOverview } from '../../services/api';
import MetricsCard from './MetricsCard';
import ProjectList from './ProjectList';
import './Overview.css'; // Assuming you have some CSS for styling

const Overview = () => {
    const [overviewData, setOverviewData] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const getOverviewData = async () => {
            try {
                const data = await fetchDashboardOverview();
                setOverviewData(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        getOverviewData();
    }, []);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div className="overview-container">
            <h2>Dashboard Overview</h2>
            <div className="metrics">
                <MetricsCard title="Total Projects" value={overviewData.totalProjects} />
                <MetricsCard title="Active Builds" value={overviewData.activeBuilds} />
                <MetricsCard title="Success Rate" value={`${overviewData.successRate}%`} />
            </div>
            <h3>Project List</h3>
            <ProjectList projects={overviewData.projects} />
        </div>
    );
};

export default Overview;