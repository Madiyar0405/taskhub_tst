import React from 'react';
import PropTypes from 'prop-types';

const Input = ({ className = '', ...props }) => {
    return (
        <input
            className={`px-4 py-2 border rounded-lg focus:outline-none focus:ring focus:border-blue-300 ${className}`}
            {...props}
        />
    );
};

Input.propTypes = {
    className: PropTypes.string,
};

export default Input;
