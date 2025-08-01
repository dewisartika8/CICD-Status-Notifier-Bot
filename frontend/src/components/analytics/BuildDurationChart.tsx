import React, { useEffect, useState } from 'react';
import { metricsApi } from '../../services/api';
import { Line } from 'react-chartjs-2';

const BuildDurationChart: React.FC = () => {
    const [buildDurations, setBuildDurations] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadBuildDurations = async () => {
            try {
                const response = await metricsApi.getBuildTrends();
                const durationsData = response.data?.data || response.data || [];
                setBuildDurations(durationsData);
            } catch (err: any) {
                setError(err?.message || 'Failed to load build durations');
            } finally {
                setLoading(false);
            }
        };

        loadBuildDurations();
    }, []);

    if (loading) return <div>Loading build durations...</div>;
    if (error) return <div>Error: {error}</div>;
    
    // Use mock data if no real data
    const mockData = [
        { date: '2024-01', average: 120 },
        { date: '2024-02', average: 110 },
        { date: '2024-03', average: 105 },
    ];
    
    const dataToUse = buildDurations.length > 0 ? buildDurations : mockData;
    
    const data = {
        labels: dataToUse.map((duration: any) => duration.date),
        datasets: [
            {
                label: 'Average Build Duration (seconds)',
                data: dataToUse.map((duration: any) => duration.average),
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                fill: true,
            }
        ],
    };

    const options = {
        responsive: true,
        plugins: {
            legend: {
                position: 'top' as const,
            },
            title: {
                display: true,
                text: 'Build Duration Trends',
            },
        },
    };

    return (
        <div>
            <Line data={data} options={options} />
        </div>
    );
};

export default BuildDurationChart;