import { useEffect, useState } from 'react';
import { fetchMetrics } from '../services/api';

const useMetrics = () => {
    const [metrics, setMetrics] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const getMetrics = async () => {
            try {
                const data = await fetchMetrics();
                setMetrics(data);
            } catch (err) {
                setError(err);
            } finally {
                setLoading(false);
            }
        };

        getMetrics();
    }, []);

    return { metrics, loading, error };
};

export default useMetrics;