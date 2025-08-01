import React from 'react';
import { Line } from 'react-chartjs-2';
import useMetrics from '../../hooks/useMetrics';

const TrendsChart = () => {
    const { metrics, loading, error } = useMetrics();

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error loading trends data</div>;
    }

    // Mock data for demo since metrics doesn't have trendsData
    const trendsData = [
        { date: '2024-01', successRate: 85, averageDuration: 120 },
        { date: '2024-02', successRate: 90, averageDuration: 110 },
        { date: '2024-03', successRate: 88, averageDuration: 105 },
    ];

    const data = {
        labels: trendsData.map((item: any) => item.date),
        datasets: [
            {
                label: 'Build Success Rate',
                data: trendsData.map((item: any) => item.successRate),
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                fill: true,
            },
            {
                label: 'Average Build Duration',
                data: trendsData.map((item: any) => item.averageDuration),
                borderColor: 'rgba(255, 99, 132, 1)',
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                fill: true,
            },
        ],
    };

    const options = {
        responsive: true,
        scales: {
            y: {
                beginAtZero: true,
            },
        },
    };

    return (
        <div>
            <h2>Trends in Project Metrics</h2>
            <Line data={data} options={options} />
        </div>
    );
};

export default TrendsChart;