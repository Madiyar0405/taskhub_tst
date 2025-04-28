import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import Navbar from "../components/navbar/Navbar";
import Footer from "../components/footer/Footer";
import styles from "../styles/CompanyIssue.module.css";

const CompanyIssue = () => {
    const { issueId } = useParams(); // Получаем issueId из URL
    const [issue, setIssue] = useState(null);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchIssue = async () => {
            try {
                const response = await fetch(`http://localhost:8082/v1/issues/${issueId}`);
                if (!response.ok) {
                    throw new Error(`Ошибка загрузки: ${response.status}`);
                }
                const data = await response.json();
                setIssue(data); // Сохраняем данные в state
            } catch (err) {
                console.error("Ошибка получения issue:", err);
                setError(err.message);
            }
        };
        fetchIssue();
    }, [issueId]);

    if (error) return <div className={styles.error}>Ошибка: {error}</div>;
    if (!issue) return <div className={styles.loading}>Загрузка...</div>;

    return (
        <>
            <Navbar />
            <div className={styles.issueContainer}>
                <h2 className={styles.issueTitle}>{issue.title}</h2>
                <p className={styles.issueBody}>{issue.body}</p>
            </div>
            <Footer />
        </>
    );
};

export default CompanyIssue;
