import React from 'react';

const StatusBadge = ({ status }) => {
    const getStatusClass = (status) => {
        switch (status) {
            case 'success':
                return 'bg-green-500 text-white';
            case 'failure':
                return 'bg-red-500 text-white';
            case 'pending':
                return 'bg-yellow-500 text-white';
            default:
                return 'bg-gray-500 text-white';
        }
    };

    return (
        <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold ${getStatusClass(status)}`}>
            {status.charAt(0).toUpperCase() + status.slice(1)}
        </span>
    );
};

export default StatusBadge;