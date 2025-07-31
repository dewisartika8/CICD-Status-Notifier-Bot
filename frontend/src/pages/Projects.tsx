import React, { useEffect } from 'react';
import { useProjects } from '../hooks/useProjects';
import Layout from '../components/common/Layout';
import ProjectList from '../components/dashboard/ProjectList';
import Overview from '../components/dashboard/Overview';

const Projects = () => {
    const { projects, fetchProjects } = useProjects();

    useEffect(() => {
        fetchProjects();
    }, [fetchProjects]);

    return (
        <Layout>
            <h1>Projects Dashboard</h1>
            <Overview />
            <ProjectList projects={projects} />
        </Layout>
    );
};

export default Projects;