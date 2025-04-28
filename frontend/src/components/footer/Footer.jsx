import React from 'react';
import '../../App.css'

const Footer = () => {
    return (
        <footer className="footer body" style={{
            textAlign: "center",
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            width: "100%",
        }}>
            <p style={{ color: "var(--primary-color)" }}>
                &copy; 2025 Task Hub. All Rights Reserved.
            </p>
        </footer>

    );
}

export default Footer;