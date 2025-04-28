import React, { useContext, useState } from 'react';
import AuthContext from "../contexts/AuthContext";
import { useLocation, useNavigate } from "react-router-dom";
import { CheckCircle, XCircle } from "lucide-react"; // оставляем иконки
import style from "../styles/Verification.module.css";
import Navbar from "../components/navbar/Navbar";
import Footer from "../components/footer/Footer";


const Verification = () => {
    const location = useLocation();
    const email = location.state?.email;
    const [code, setCode] = useState('');
    const { verifyCode } = useContext(AuthContext);
    const navigate = useNavigate();
    const [error, setError] = useState('');
    const [status, setStatus] = useState('waiting'); // 'waiting' | 'success' | 'error'

    const handleSubmit = async () => {
        if (!email) {
            console.error("No email found.");
            setStatus('error');
            setError('Email не найден');
            return;
        }
        try {
            await verifyCode(code, email);
            setStatus('success');
            setTimeout(() => {
                navigate('/signIn');
            }, 1500);
        } catch (err) {
            setStatus('error');
            setError('Ошибка верификации! Пожалуйста, попробуйте снова.');
            console.error('Verification failed:', err);
        }
    };

    const handleResend = () => {
        setStatus('waiting');
        setError('');
        console.log('Resending verification code');
        // Здесь должна быть логика повторной отправки кода
    };

    return (
        <>
            <Navbar />
        <div className={style.container}>
            <div className={style.card}> {/* Используем div вместо Card */}
                <div className={style.cardHeader}> {/* Используем div вместо CardHeader */}
                    <h2 className={style.title}>Подтверждение электронной почты</h2> {/* CardTitle */}
                </div>
                <div className={style.cardContent}> {/* Используем div вместо CardContent */}
                    <div className={style.content}>
                        <p className={style.text}>
                            Мы отправили код подтверждения на вашу электронную почту.<br/>
                            Пожалуйста, введите его ниже.
                        </p>
                        {status === 'success' && (
                            <div className={style.success}>
                                <CheckCircle className={style.icon} />
                                <span>Почта успешно подтверждена!</span>
                            </div>
                        )}
                        {status === 'error' && (
                            <div className={style.error}>
                                <XCircle className={style.icon} />
                                <span>{error || 'Неверный код. Попробуйте снова.'}</span>
                            </div>
                        )}
                        <input
                            type="text"
                            placeholder="Введите код"
                            maxLength={6}
                            value={code}
                            onChange={(e) => setCode(e.target.value)}
                            className={style.input}
                        />
                        <button onClick={handleSubmit} className={style.button}>Подтвердить</button>
                        <div className={style.resendContainer}>
                            <button onClick={handleResend} className={style.resendButton}>
                                Отправить код повторно
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>
            <Footer/>
        </>
    );
};

export default Verification;
