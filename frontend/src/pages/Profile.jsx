
import React, { useState, useEffect, useContext } from 'react';
import styles from "../styles/Profile.module.css";
import '../App.css';
import Navbar from "../components/navbar/Navbar";
import Footer from "../components/footer/Footer";
import GithubAuthButton from "../components/GitHubLogin";
import AuthContext from "../contexts/AuthContext";

const Profile = () => {
    const { user } = useContext(AuthContext);
    const [isGithubLinked, setIsGithubLinked] = useState(true);
    const [profileData, setProfileData] = useState(null);
    const [error, setError] = useState(null);
    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const response = await fetch(`http://localhost:8081/v1/profile/${String(user.id)}`);
                if (!response.ok) {
                    throw new Error(`Ошибка сети: ${response.status}`);
                }
                const data = await response.json();
                setProfileData(data);
                setIsGithubLinked(data.isGithubLinked);
            } catch (err) {
                console.error("Ошибка получения профиля:", err);
                setError(err.message);
            }
        };

        fetchProfile();
    }, [user.id]); // Перезапуск запроса при изменении ID

    if (error) return <div>Ошибка: {error}</div>;
    if (!profileData) return <div>Загрузка...</div>;

    return (
        <>
            <Navbar />
            <div className={styles.profileContainer}>
                {/* Левая часть профиля */}
                <div className={styles.profileLeft}>
                    <div className={styles.profileInfo}>
                        <div className={styles.avatarContainer}>
                            <img
                                src={profileData.avatarUrl || "/default-avatar.png"}
                                alt="Profile"
                                className={styles.avatarImage}
                            />
                        </div>
                        <h1 className={styles.profileName}>{profileData.username}</h1>
                        <a href={`mailto:${profileData.email}`} className={styles.profileEmail}>
                            {profileData.email}
                        </a>
                        <p className={styles.profileDescription}>{profileData.bio || "Нет описания"}</p>

                        {profileData.cvUrl && (
                            <a href={profileData.cvUrl} target="_blank" rel="noopener noreferrer" className={styles.cvLink}>
                                📄 Открыть резюме
                            </a>
                        )}
                    </div>
                </div>

                {/* Правая часть профиля */}
                <div className={styles.profileRight}>
                    <div className={styles.topPanel}>
                        <div className={styles.panelContent}>
                            <h2 className={styles.panelTitle}>Профиль разработчика</h2>
                            {!isGithubLinked ? (
                                <GithubAuthButton />
                            ) : (
                                <div className={styles.githubLinked}>
                                    <svg className={styles.checkIcon} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                    </svg>
                                    GitHub привязан
                                    <a
                                        href={`https://github.com/${profileData.githubUsername}`}
                                        target="_blank"
                                        rel="noopener noreferrer"
                                        className={styles.githubLink}
                                    >
                                         @{profileData.githubUsername}
                                    </a>
                                </div>
                            )}
                        </div>
                    </div>

                {/* Content that appears after GitHub is linked */}
                {isGithubLinked && (
                    <div className={styles.fadeIn}>
                        {/* Experience - Full Width */}
                        <div className={styles.fullWidthSection}>
                            <h2 className={styles.sectionTitle}>Опыт работы</h2>
                            <div className={styles.contentItem}>
                                <p className={styles.contentTitle}>Senior Frontend Developer</p>
                                <p className={styles.contentSubtitle}>Tech Company</p>
                                <p className={styles.contentDate}>2020 - настоящее время</p>
                            </div>
                        </div>

                        {/* Three Equal Columns */}
                        <div className={styles.threeColumnGrid}>
                            <div className={styles.gridCard}>
                                <h2 className={styles.sectionTitle}>Образование</h2>
                                <div className={styles.contentItem}>
                                    <p className={styles.contentTitle}>МГУ им. Ломоносова</p>
                                    <p className={styles.contentSubtitle}>Факультет ВМК</p>
                                    <p className={styles.contentDate}>2015 - 2019</p>
                                </div>
                            </div>

                            <div className={styles.gridCard}>
                                <h2 className={styles.sectionTitle}>Языки</h2>
                                <div className={styles.contentItem}>
                                    <p className={styles.contentTitle}>Английский</p>
                                    <p className={styles.contentSubtitle}>C1 Advanced</p>
                                </div>
                                <div className={styles.contentItem}>
                                    <p className={styles.contentTitle}>Русский</p>
                                    <p className={styles.contentSubtitle}>Родной</p>
                                </div>
                            </div>

                            <div className={styles.gridCard}>
                                <h2 className={styles.sectionTitle}>Сертификаты</h2>
                                <div className={styles.contentItem}>
                                    <p className={styles.contentTitle}>AWS Solutions Architect</p>
                                    <p className={styles.contentDate}>2023</p>
                                </div>
                            </div>
                        </div>

                        <div className={styles.twoColumnGrid}>
                            <div className={styles.gridCard}>
                                <h2 className={styles.sectionTitle}>Образование</h2>
                                <div className={styles.contentItem}>
                                    <p className={styles.contentTitle}>МГУ им. Ломоносова</p>
                                    <p className={styles.contentSubtitle}>Факультет ВМК</p>
                                    <p className={styles.contentDate}>2015 - 2019</p>
                                </div>
                            </div>

                            <div className={styles.gridCard}>
                                <h2 className={styles.sectionTitle}>Языки</h2>
                                <div className={styles.contentItem}>
                                    <p className={styles.contentTitle}>Английский</p>
                                    <p className={styles.contentSubtitle}>C1 Advanced</p>
                                </div>
                                <div className={styles.contentItem}>
                                    <p className={styles.contentTitle}>Русский</p>
                                    <p className={styles.contentSubtitle}>Родной</p>
                                </div>
                            </div>
                        </div>
                    </div>
                )}
            </div>
            </div>
            <Footer />
        </>
    );
};

export default Profile;
