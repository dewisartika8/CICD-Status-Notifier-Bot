import React, { useEffect, useState } from 'react';
import { Line } from 'react-chartjs-2';
import { fetchSuccessRate } from '../../services/api';

const SuccessRateChart = () => {
    const [data, setData] = useState({
        labels: [],
        datasets: [
            {
                label: 'Success Rate',
                data: [],
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                borderWidth: 2,
            },
        ],
    });

    useEffect(() => {
        const getSuccessRateData = async () => {
            try {
                const response = await fetchSuccessRate();
                const successRates = response.data; // Assuming the response structure
                const labels = successRates.map(rate => rate.date); // Adjust based on actual data structure
                const values = successRates.map(rate => rate.successRate); // Adjust based on actual data structure

                setData({
                    labels,
                    datasets: [
                        {
                            ...data.datasets[0],
                            data: values,
                        },
                    ],
                });
            } catch (error) {
                console.error('Error fetching success rate data:', error);
            }
        };

        getSuccessRateData();
    }, []);

    return (
        <div>
            <h2>Build Success Rate Over Time</h2>
            <Line data={data} />
        </div>
    );
};

export default SuccessRateChart;