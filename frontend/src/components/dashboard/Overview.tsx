import React, { useEffect, useState } from 'react';
import { fetchDashboardOverview } from '../../services/api';
import MetricsCard from './MetricsCard';
import ProjectList from './ProjectList';
import { DashboardOverview } from '../../types';

const Overview: React.FC = () => {
    const [overviewData, setOverviewData] = useState<DashboardOverview | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadOverview = async () => {
            try {
                const response = await fetchDashboardOverview();
                const data = response.data?.data || response.data;
                setOverviewData(data);
            } catch (err: any) {
                setError(err?.message || 'Failed to load overview');
            } finally {
                setLoading(false);
            }
        };

        loadOverview();
    }, []);

    if (loading) return <div>Loading overview...</div>;
    if (error) return <div>Error: {error}</div>;
    if (!overviewData) return <div>No overview data available</div>;

    const successRate = overviewData.total_builds > 0
        ? ((overviewData.successful_builds / overviewData.total_builds) * 100).toFixed(1)
        : '0';

    return (
        <div>
            <h1>Dashboard Overview</h1>
            <div className="metrics-grid">
                <MetricsCard 
                    title="Total Projects" 
                    value={overviewData.total_projects} 
                    description="Total number of projects"
                />
                <MetricsCard 
                    title="Active Builds" 
                    value={overviewData.active_projects} 
                    description="Currently running builds"
                />
                <MetricsCard 
                    title="Success Rate" 
                    value={`${successRate}%`}
                    description="Overall build success rate"
                />
            </div>
            {/* Project list will be loaded separately */}
            <div style={{ marginTop: '2rem' }}>
                <h3>Project Overview</h3>
                <p>Projects section - data loaded separately via ProjectList component</p>
            </div>
        </div>
    );
};

export default Overview;