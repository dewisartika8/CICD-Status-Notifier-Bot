import React, { useEffect } from 'react';
import Layout from '../components/common/Layout';
import Overview from '../components/dashboard/Overview';
import ProjectList from '../components/dashboard/ProjectList';
import MetricsCard from '../components/dashboard/MetricsCard';
import { useMetrics } from '../hooks/useMetrics';
import { useProjects } from '../hooks/useProjects';
import { fetchDashboardData } from '../services/api';

const Dashboard = () => {
  const { metrics, fetchMetrics } = useMetrics();
  const { projects, fetchProjects } = useProjects();

  useEffect(() => {
    fetchMetrics();
    fetchProjects();
  }, [fetchMetrics, fetchProjects]);

  return (
    <Layout>
      <h1>Dashboard</h1>
      <Overview metrics={metrics} />
      <MetricsCard metrics={metrics} />
      <ProjectList projects={projects} />
    </Layout>
  );
};

export default Dashboard;