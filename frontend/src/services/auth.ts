// src/services/auth.ts

import axios from 'axios';

const API_URL = '/api/v1/auth';

// Function to handle user login
export const login = async (username, password) => {
    try {
        const response = await axios.post(`${API_URL}/login`, { username, password });
        if (response.data.token) {
            localStorage.setItem('user', JSON.stringify(response.data));
        }
        return response.data;
    } catch (error) {
        throw error.response.data;
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
    return JSON.parse(localStorage.getItem('user'));
};