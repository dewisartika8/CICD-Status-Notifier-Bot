import React, { useEffect, useState } from 'react';
import { fetchBuildHistory } from '../../services/api';

const BuildHistory = ({ projectId }) => {
    const [builds, setBuilds] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const getBuildHistory = async () => {
            try {
                const data = await fetchBuildHistory(projectId);
                setBuilds(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        getBuildHistory();
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
                            <td>{build.buildNumber}</td>
                            <td>{build.status}</td>
                            <td>{build.duration} seconds</td>
                            <td>{new Date(build.date).toLocaleString()}</td>
                        </tr>
                    ))}
                </tbody>
            </table>
        </div>
    );
};

export default BuildHistory;