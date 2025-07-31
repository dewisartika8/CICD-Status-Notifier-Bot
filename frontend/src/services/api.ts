import axios from 'axios';

const API_BASE_URL = 'http://localhost:5000/api/v1'; // Update with your backend URL

// Function to get the dashboard overview
export const getDashboardOverview = async () => {
    const response = await axios.get(`${API_BASE_URL}/dashboard/overview`);
    return response.data;
};

// Function to get the list of projects
export const getProjects = async () => {
    const response = await axios.get(`${API_BASE_URL}/dashboard/projects`);
    return response.data;
};

// Function to get project details by ID
export const getProjectDetails = async (projectId) => {
    const response = await axios.get(`${API_BASE_URL}/dashboard/projects/${projectId}/details`);
    return response.data;
};

// Function to get build history for a specific project
export const getBuildHistory = async (projectId) => {
    const response = await axios.get(`${API_BASE_URL}/projects/${projectId}/builds`);
    return response.data;
};

// Function to get real-time project status
export const getProjectStatus = async (projectId) => {
    const response = await axios.get(`${API_BASE_URL}/projects/${projectId}/status`);
    return response.data;
};

// Function to get success rate metrics
export const getSuccessRateMetrics = async () => {
    const response = await axios.get(`${API_BASE_URL}/metrics/success-rate`);
    return response.data;
};

// Function to get average build duration metrics
export const getAverageBuildDuration = async () => {
    const response = await axios.get(`${API_BASE_URL}/metrics/average-build-duration`);
    return response.data;
};

// Function to get build trends
export const getBuildTrends = async () => {
    const response = await axios.get(`${API_BASE_URL}/metrics/build-trends`);
    return response.data;
};