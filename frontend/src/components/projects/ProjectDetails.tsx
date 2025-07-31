import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { fetchProjectDetails } from '../../services/api';
import BuildHistory from './BuildHistory';
import StatusBadge from './StatusBadge';

const ProjectDetails = () => {
    const { id } = useParams();
    const [project, setProject] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const getProjectDetails = async () => {
            try {
                const data = await fetchProjectDetails(id);
                setProject(data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        };

        getProjectDetails();
    }, [id]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    return (
        <div>
            <h1>{project.name}</h1>
            <StatusBadge status={project.status} />
            <p>{project.description}</p>
            <h2>Build History</h2>
            <BuildHistory projectId={project.id} />
        </div>
    );
};

export default ProjectDetails;