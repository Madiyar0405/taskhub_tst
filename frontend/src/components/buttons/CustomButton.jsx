import React from 'react';
import Button from '../ui/button';
import { Link } from 'lucide-react';

const CustomButton = ({
                          text = '',
                          textColor = 'text-white',
                          backgroundColor = 'bg-primary',
                          href,
                          onClick
                      }) => {
    const buttonClasses = `flex items-center gap-2 ${textColor} ${backgroundColor}`;

    if (href) {
        return (
            <Button
                asChild
                className={buttonClasses}
                onClick={onClick}
            >
                <a href={href} target="_blank" rel="noopener noreferrer">
                    {text}
                    <Link className="w-4 h-4" />
                </a>
            </Button>
        );
    }

    return (
        <Button
            className={buttonClasses}
            onClick={onClick}
        >
            {text}
        </Button>
    );
};

export default CustomButton;