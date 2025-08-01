// src/services/auth.ts

import axios from 'axios';

const API_URL = '/api/v1/auth';

// Function to handle user login
export const login = async (username: string, password: string) => {
    try {
        const response = await axios.post(`${API_URL}/login`, { username, password });
        if (response.data.token) {
            localStorage.setItem('user', JSON.stringify(response.data));
        }
        return response.data;
    } catch (error: any) {
        throw error?.response?.data || error?.message || 'Login failed';
    }
};

// Function to handle user logout
export const logout = () => {
    localStorage.removeItem('user');
};

// Function to check if the user is logged in
export const isLoggedIn = () => {
    const user = localStorage.getItem('user');
    return user !== null;
};

// Function to get the current user
export const getCurrentUser = () => {
    const userData = localStorage.getItem('user');
    return userData ? JSON.parse(userData) : null;
};