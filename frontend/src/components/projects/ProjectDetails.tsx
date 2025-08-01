import React, { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { fetchProjectDetails } from '../../services/api';
import { ProjectResponse } from '../../types';
import BuildHistory from './BuildHistory';
import StatusBadge from './StatusBadge';

const ProjectDetails: React.FC = () => {
    const { id } = useParams<{ id: string }>();
    const [project, setProject] = useState<ProjectResponse | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchProject = async () => {
            if (!id) return;
            
            try {
                const response = await fetchProjectDetails(id);
                setProject(response.data.data);
            } catch (err: any) {
                setError(err?.message || 'Failed to fetch project');
            } finally {
                setLoading(false);
            }
        };

        fetchProject();
    }, [id]);

    if (loading) {
        return <div>Loading...</div>;
    }

    if (error) {
        return <div>Error: {error}</div>;
    }

    if (!project) {
        return <div>Project not found</div>;
    }

    return (
        <div>
            <h1>{project.name}</h1>
            <StatusBadge status={project.status} />
            <p>Repository: {project.repository_url}</p>
            
            <BuildHistory projectId={project.id} />
        </div>
    );
};

export default ProjectDetails;
