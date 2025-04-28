
import React, {useContext, useEffect, useState} from 'react';
import { useParams } from "react-router-dom";
import styles from "../styles/Company.module.css";
import Navbar from "../components/navbar/Navbar";
import Footer from "../components/footer/Footer";
import AuthContext from "../contexts/AuthContext";

const CompanyIssues = () => {
    const { companyId } = useParams();
    const [companyData, setCompanyData] = useState(null);
    const [error, setError] = useState(null);
    const { user } = useContext(AuthContext);


    useEffect(() => {
        const fetchCompanyData = async () => {
            try {
                const response = await fetch(`http://localhost:8082/v1/companies/${companyId}`);

                console.log(user.token+" token "+user)
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
                        <p>
                            <strong>Website:</strong> {companyData.websiteUrl}
                        </p>
                    </div>
                    <div className={styles.right}>
                        <p>
                            <strong>Website:</strong> {companyData.websiteUrl}
                        </p>
                        <p>
                            <strong>Work with us since:</strong> {formattedDate}
                        </p>
                        <p>
                            <strong>About us:</strong> {companyData.description}
                        </p>
                    </div>
                </div>
            </div>
            <Footer/>
        </>
    );
};

export default CompanyIssues;
