import { useEffect, useState } from 'react';
import { fetchProjects } from '../services/api';
import { Project } from '../types';

const useProjects = () => {
    const [projects, setProjects] = useState<Project[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const loadProjects = async () => {
            try {
                const response = await fetchProjects();
                const projectsData = response.data.projects || [];
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
    }, []);

    return { projects, loading, error };
};

export default useProjects;