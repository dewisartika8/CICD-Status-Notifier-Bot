import React, { useEffect, useState } from 'react';
import { fetchProjects } from '../../services/api';
import ProjectItem from './ProjectItem';

const ProjectList: React.FC = () => {
    const [projects, setProjects] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadProjects = async () => {
            try {
                const data = await fetchProjects();
                setProjects(data);
            } catch (err) {
                setError('Failed to load projects');
            } finally {
                setLoading(false);
            }
        };

        loadProjects();
    }, []);

    if (loading) {
        return <div>Loading projects...</div>;
    }

    if (error) {
        return <div>{error}</div>;
    }

    return (
        <div>
            <h2>Project List</h2>
            <ul>
                {projects.map((project) => (
                    <ProjectItem key={project.id} project={project} />
                ))}
            </ul>
        </div>
    );
};

export default ProjectList;