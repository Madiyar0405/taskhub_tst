import React, {useContext, useEffect, useState} from 'react';
import {useParams,Link} from "react-router-dom";
import styles from "../styles/Company.module.css";
import Navbar from "../components/navbar/Navbar";
import Footer from "../components/footer/Footer";

const CompanyInfo = () => {
    const { companyId } = useParams();
    const [companyData, setCompanyData] = useState(null);
    const [issues, setIssues] = useState([]); // Состояние для хранения issues
    const [error, setError] = useState(null);



    useEffect(() => {
        const fetchCompanyData = async () => {
            try {
                const response = await fetch(`http://localhost:8082/v1/companies/${companyId}`);
                if (!response.ok) {
                    throw new Error(`Ошибка сети: ${response.status}`);
                }
                const data = await response.json();
                setCompanyData(data);
            } catch (err) {
                console.error("Ошибка получения данных компании:", err);
                setError(err.message);
            }
        };

        fetchCompanyData();
    }, [companyId]);

    useEffect(() => {
        const fetchIssues = async () => {
            try {
                const response = await fetch(`http://localhost:8082/v1/companies/${companyId}/issues`);

                if (!response.ok) {
                    throw new Error(`Ошибка сети: ${response.status}`);
                }
                const data = await response.json();
                setIssues(data.issues);
                console.log(data);
            } catch (err) {
                console.error("Ошибка получения issues:", err);
                setError(err.message);
            }
        };

        fetchIssues();
    }, [companyId]);

    if (error) return <div>Ошибка: {error}</div>;
    if (!companyData) return <div>Загрузка...</div>;

    const createdAt = new Date(companyData.createdAt);
    const formattedDate = createdAt.toLocaleDateString("ru-RU", {
        year: "numeric",
        month: "long",
        day: "numeric",
    });

    return (
        <>
            <Navbar />
            <div className={styles.companyContainer}>
                <h2 className={styles.companyName}>{companyData.companyName}</h2>
                <div className={styles.contentWrapper}>
                    <div className={styles.left}>
                        <div className={styles.avatarContainer}>
                            <img
                                src={companyData.logo}
                                alt="Profile"
                                className={styles.avatarImage}
                            />
                        </div>
                        <p><strong>Website:</strong> {companyData.websiteUrl}</p>
                    </div>
                    <div className={styles.right}>
                        <p><strong>Work with us since:</strong> {formattedDate}</p>
                        <p><strong>About us:</strong> {companyData.description}</p>
                    </div>
                </div>

                {/* Секция Issues */}
                <div className={styles.issuesSection}>
                    <h3>Company Tasks</h3>
                    <ul className={styles.issuesList}>
                        {issues.map((issue, index) => (
                            <li key={index} className={styles.issueItem}>
                                <strong>{issue.title}</strong>
                                <Link to={`/issues/${issue.issueId}`} className={styles.issueButton}>            View Issue
                                </Link>
                            </li>
                        ))}
                    </ul>

                </div>
            </div>
            <Footer/>
        </>
    );
};

export default CompanyInfo;
