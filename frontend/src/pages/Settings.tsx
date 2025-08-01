import React from 'react';
import Layout from '../components/common/Layout';

const Settings: React.FC = () => {
    return (
        <Layout>
            <div className="settings-container">
                <h1>Settings</h1>
                <p>Configure your preferences and settings for the dashboard.</p>
                {/* Add form elements for user settings here */}
            </div>
        </Layout>
    );
};

export default Settings;