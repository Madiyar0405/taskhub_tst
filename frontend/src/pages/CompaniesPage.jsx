import React, {useEffect, useState} from "react";
import CompanyCard from "../components/CompanyCard/CompanyCard";
import styles from "../styles/CompanyPage.module.css";
import Navbar from "../components/navbar/Navbar";
import Footer from "../components/footer/Footer";
import useCompanies from "../hooks/useCompanies";

const fallbackCompanies = [
    {
        id: 1,
        name: "Google",
        description: "Innovative technology solutions",
        image: "https://via.placeholder.com/150",
        icon: "https://via.placeholder.com/50",
    },
    {
        id: 2,
        name: "AgroMax",
        description: "Agricultural advancements",
        image: "https://via.placeholder.com/150",
        icon: "https://via.placeholder.com/50",
    },
    {
        id: 3,
        name: "MediCare",
        description: "Healthcare redefined",
        image: "https://via.placeholder.com/150",
        icon: "https://via.placeholder.com/50",
    },
];

export default function CompaniesPage() {
    const { companies, loading, error } = useCompanies();
    const dataToShow = error ? fallbackCompanies : companies;

    console.log("dataToShow:", dataToShow); // Добавить лог для проверки данных

    return (
        <>
            <Navbar />
            <div className={styles.companyPageContainer}>
                <h1 className={styles.pageTitle}>Наши компании</h1>
                {loading && <p>Загрузка...</p>}
                <div className={styles.companyGrid}>
                    {dataToShow.map((company, index) => (
                        <CompanyCard key={company.id || index} company={company} />
                        ))}
                </div>
            </div>
            <Footer />
        </>
    );
}


