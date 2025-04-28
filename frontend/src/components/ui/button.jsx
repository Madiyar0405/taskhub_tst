import React from 'react';
import PropTypes from 'prop-types';

const Button = ({ children, className = '', onClick, asChild = false }) => {
    const Tag = asChild ? 'div' : 'button';
    return (
        <Tag
            className={`px-4 py-2 rounded-lg font-medium ${className}`}
            onClick={onClick}
        >
            {children}
        </Tag>
    );
};

Button.propTypes = {
    children: PropTypes.node.isRequired,
    className: PropTypes.string,
    onClick: PropTypes.func,
    asChild: PropTypes.bool,
};

export default Button;
