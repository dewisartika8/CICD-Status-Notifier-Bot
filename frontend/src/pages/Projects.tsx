import React, { useEffect } from 'react';
import useProjects from '../hooks/useProjects';
import Layout from '../components/common/Layout';
import ProjectList from '../components/dashboard/ProjectList';
// import Overview from '../components/dashboard/Overview';

const Projects = () => {
    const { projects, loading, error } = useProjects();

    // Projects are automatically loaded by the hook
    // No need for manual effect

    return (
        <Layout>
            <h1>Projects Dashboard</h1>
            {/* <Overview /> */}
            <ProjectList projects={projects} />
        </Layout>
    );
};

export default Projects;