    import React from "react";
    import { useNavigate } from "react-router-dom";
    import styles from "./CompanyCard.module.css";

    export default function CompanyCard({ company }) {
        const navigate = useNavigate();
        const handleClick = () => {
            navigate(`/companies/${company.companyId}`); // Переход на страницу компании
        };

        return (
            <div className={styles.companyCard} onClick={handleClick}>
                <img
                    src={company.logo} // Используем logo
                    alt={company.companyName} // Используем companyName
                    className={styles.companyImage}
                />
                <div className={styles.companyContent}>
                    <div className={styles.companyHeader}>
                        <img
                            src={company.logo} // Используем logo как иконку
                            alt={company.companyName} // Используем companyName
                            className={styles.companyIcon}
                        />
                        <h3 className={styles.companyName}>{company.companyName}</h3>
                    </div>
                    {/* Если нет описания, можно просто не выводить этот элемент */}
                    {company.description && (
                        <p className={styles.companyDescription}>
                            {company.description}
                        </p>
                    )}
                </div>
            </div>
        );
    }
