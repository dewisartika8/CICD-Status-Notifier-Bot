import { useEffect, useState } from 'react';
import { dashboardApi } from '../services/api';
import { Metrics } from '../types';

const useMetrics = () => {
    const [metrics, setMetrics] = useState<Metrics | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadMetrics = async () => {
            try {
                const response = await dashboardApi.getOverview();
                const metricsData = response.data.data;
                const mappedMetrics: Metrics = {
                    successRate: metricsData.success_rate,
                    averageDuration: metricsData.average_duration,
                    deploymentFrequency: 0,
                    buildsToday: 0,
                    buildsThisWeek: 0,
                    buildsThisMonth: metricsData.total_builds
                };
                setMetrics(mappedMetrics);
            } catch (err: any) {
                setError(err?.message || 'Failed to load metrics');
            } finally {
                setLoading(false);
            }
        };

        loadMetrics();
    }, []);

    return { metrics, loading, error };
};

export default useMetrics;