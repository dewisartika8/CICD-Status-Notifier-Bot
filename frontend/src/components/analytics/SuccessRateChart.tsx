import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import { dashboardApi } from '../../services/api';

const SuccessRateChart: React.FC = () => {
    const [successRates, setSuccessRates] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadSuccessRates = async () => {
            try {
                const response = await dashboardApi.getOverview();
                const ratesData = response.data.data.success_rate;
                setSuccessRates([{ name: 'Success Rate', value: ratesData }]);
            } catch (err: any) {
                setError(err?.message || 'Failed to load success rates');
            } finally {
                setLoading(false);
            }
        };

        loadSuccessRates();
    }, []);

    if (loading) return <div>Loading success rates...</div>;
    if (error) return <div>Error: {error}</div>;

    // Use mock data if no real data or ensure it's an array
    const mockData = [
        { date: '2024-01', successRate: 85 },
        { date: '2024-02', successRate: 90 },
        { date: '2024-03', successRate: 88 },
    ];
    
    const dataToUse = Array.isArray(successRates) && successRates.length > 0 ? successRates : mockData;

    const data = {
        labels: dataToUse.map((rate: any) => rate.date),
        datasets: [
            {
                label: 'Success Rate (%)',
                data: dataToUse.map((rate: any) => rate.successRate),
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
                text: 'Build Success Rate Trends',
            },
        },
        scales: {
            y: {
                beginAtZero: true,
                max: 100,
                ticks: {
                    callback: function(value: any) {
                        return value + '%';
                    }
                }
            }
        }
    };

    return (
        <div>
            <Line data={data} options={options} />
        </div>
    );
};

export default SuccessRateChart;