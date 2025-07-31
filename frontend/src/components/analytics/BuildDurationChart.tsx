import React, { useEffect, useState } from 'react';
import { fetchBuildDurations } from '../../services/api';
import { Line } from 'react-chartjs-2';

const BuildDurationChart = () => {
    const [buildDurations, setBuildDurations] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        const getBuildDurations = async () => {
            try {
                const data = await fetchBuildDurations();
                setBuildDurations(data);
            } catch (error) {
                console.error('Error fetching build durations:', error);
            } finally {
                setLoading(false);
            }
        };

        getBuildDurations();
    }, []);

    const chartData = {
        labels: buildDurations.map(duration => duration.date),
        datasets: [
            {
                label: 'Average Build Duration (minutes)',
                data: buildDurations.map(duration => duration.average),
                borderColor: 'rgba(75, 192, 192, 1)',
                backgroundColor: 'rgba(75, 192, 192, 0.2)',
                borderWidth: 1,
            },
        ],
    };

    return (
        <div>
            <h2>Build Duration Over Time</h2>
            {loading ? (
                <p>Loading...</p>
            ) : (
                <Line data={chartData} />
            )}
        </div>
    );
};

export default BuildDurationChart;