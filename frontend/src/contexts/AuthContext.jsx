import React, { createContext, useEffect, useState } from "react";
import axios from "axios";

const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [user, setUser] = useState(null);
    const [isLoading, setIsLoading] = useState(true);
    const [verificationCodeSent, setVerificationCodeSent] = useState(false);
    const API_BASE_URL = "http://localhost:8081";

    useEffect(() => {
        const token = localStorage.getItem("authToken");
        const email = localStorage.getItem("userEmail");
        if (token) {
            const decoded = decodeJWT(token);
            if (decoded) {
                const userId = decoded.uid;
                setUser({ token, id: userId, email });
                console.log("User loaded:", { id: userId, email });
            }
        }
        setIsLoading(false); // ✅ Завершаем загрузку
    }, []);

    const logIn = async (email, password) => {
        try {
            const response = await axios.post(
                `${API_BASE_URL}/v1/login`,
                { email, password }
            );
            const { token } = response.data;

            if (token) {
                const decoded = decodeJWT(token);
                localStorage.setItem("authToken", token);
                localStorage.setItem("userEmail", email);
                setUser({ token, id: decoded?.uid, email });

                axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;


                return { success: true };
            } else {
                const errorMsg = "В ответе от сервера нет токена";
                console.error(`❌ ${errorMsg}:`, response.data);
                return { success: false, message: errorMsg };
            }
        } catch (error) {
            let errorMessage = "Ошибка при входе в систему";

            console.log(error.response.data.message +" status amodus");
            // Детальная обработка ошибок
            if (error.response) {
                // Ответ от сервера был получен, но статус не 2xx
                if(error.response.status === 400){
                    errorMessage = "Error occurred, try again";
                }
                else if (error.response.status === 401) {
                    errorMessage = "Неверный email или пароль";
                } else if (error.response.status === 404) {
                    errorMessage = "Сервер не найден. Проверьте подключение к интернету";
                } else if (error.response.status === 500) {
                    errorMessage = "Ошибка сервера. Попробуйте позже";
                } else if (error.response.data && error.response.data.message) {
                    errorMessage = error.response.data.message;
                }
            } else if (error.request) {
                // Запрос был отправлен, но ответ не получен
                errorMessage = "Сервер не отвечает. Проверьте подключение к интернету";
            } else {
                // Что-то пошло не так при настройке запроса
                errorMessage = error.message || "Неизвестная ошибка";
            }

            console.error("SignIn failed. Error details:", error.response?.data || error.message);

            // Возвращаем объект с деталями ошибки
            return { success: false, message: errorMessage };
        }
    };

    const signUp = async (email, password, username) => {
        let errorMessage = "some Error";

        try {
            await axios.post(`${API_BASE_URL}/v1/signup`, { email, password, username });
            setVerificationCodeSent(true);
        } catch (error) {
            errorMessage = error.response.data.message
            console.error("SignUp failed. Error details:", error.response?.data || error.message);
        }
        return { success: false, message: errorMessage };

    };

    const verifyCode = async (OTP, email) => {
        try {
            const response = await axios.post(`${API_BASE_URL}/v1/signup/confirm`, { email, otp: OTP });
            const { token } = response.data;
            if (token) {
                localStorage.setItem("authToken", token);
                localStorage.setItem("userEmail", email);
                setUser({ token, email });
            }
            alert("Verification successful. You are now logged in.");
        } catch (error) {
            console.error("Verification failed. Error details:", error.response?.data || error.message);
            alert("Verification failed. Please try again.");
        }
    };

    const logOut = () => {
        setUser(null);
        localStorage.removeItem("authToken");
        localStorage.removeItem("userEmail");
        console.log("User logged out and token cleared.");
    };

    const decodeJWT = (token) => {
        try {
            const payload = token.split(".")[1];
            const decoded = JSON.parse(atob(payload));
            return decoded;
        } catch (error) {
            console.error("Ошибка при декодировании JWT:", error);
            return null;
        }
    };

    return (
        <AuthContext.Provider value={{ user, isLoading, logIn, logOut, signUp, verifyCode, verificationCodeSent }}>
            {children}
        </AuthContext.Provider>
    );
};

export default AuthContext;
