import React, { useEffect, useState } from 'react';
import { webhookApi } from '../../services/api';
import { WebhookEventResponse } from '../../types';

interface BuildHistoryProps {
    projectId: string;
}

const BuildHistory: React.FC<BuildHistoryProps> = ({ projectId }) => {
    const [builds, setBuilds] = useState<WebhookEventResponse[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchBuilds = async () => {
            try {
                const response = await webhookApi.getWebhookEventsByProject(projectId);
                setBuilds(response.data.data);
            } catch (err: any) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        fetchBuilds();
    }, [projectId]);

    if (loading) {
        return <div>Loading build history...</div>;
    }

    if (error) {
        return <div>Error fetching build history: {error}</div>;
    }

    return (
        <div>
            <h2>Build History for Project {projectId}</h2>
            <table>
                <thead>
                    <tr>
                        <th>Build Number</th>
                        <th>Status</th>
                        <th>Duration</th>
                        <th>Date</th>
                    </tr>
                </thead>
                <tbody>
                    {builds.map(build => (
                        <tr key={build.id}>
                            <td>{build.delivery_id}</td>
                            <td>{build.event_type}</td>
                            <td>N/A</td>
                            <td>{new Date(build.created_at).toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default BuildHistory;