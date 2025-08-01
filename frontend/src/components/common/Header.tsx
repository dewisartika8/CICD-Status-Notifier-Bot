import React from 'react';

const Header: React.FC = () => {
    return (
        <header className="header">
            <h1 className="header-title">CI/CD Status Notifier Dashboard</h1>
            <nav className="header-nav">
                <ul>
                    <li><a href="/">Overview</a></li>
                    <li><a href="/projects">Projects</a></li>
                    <li><a href="/analytics">Analytics</a></li>
                    <li><a href="/settings">Settings</a></li>
                </ul>
            </nav>
        </header>
    );
};

export default Header;