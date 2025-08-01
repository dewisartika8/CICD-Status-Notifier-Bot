import React, { useEffect, useState } from 'react';
import { fetchProjects } from '../../services/api';
import ProjectItem from './ProjectItem';
import { Project } from '../../types';

interface ProjectListProps {
    projects?: Project[];
}

const ProjectList: React.FC<ProjectListProps> = ({ projects: propProjects }) => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (propProjects) {
            setProjects(propProjects);
            setLoading(false);
            return;
        }

        const loadProjects = async () => {
            try {
                const response = await fetchProjects();
                const projectsData = response.data?.projects || [];
                const mappedProjects: Project[] = projectsData.map((p) => ({
                    ...p,
                    buildHistory: [],
                    isActive: p.status === 'active'
                }));
                setProjects(mappedProjects);
            } catch (err: any) {
                setError(err?.message || 'Failed to load projects');
            } finally {
                setLoading(false);
            }
        };

        loadProjects();
    }, [propProjects]);

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