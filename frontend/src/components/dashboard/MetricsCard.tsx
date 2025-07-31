import React from 'react';

const MetricsCard = ({ title, value, description }) => {
    return (
        <div className="metrics-card">
            <h3 className="metrics-card-title">{title}</h3>
            <p className="metrics-card-value">{value}</p>
            <p className="metrics-card-description">{description}</p>
        </div>
    );
};

export default MetricsCard;