import React, {useEffect, useRef, useState} from "react";
import { useNavigate } from "react-router-dom";
import axios from "axios";

const GitHubCallback = () => {
    const navigate = useNavigate();
    const isRequestSent = useRef(false); // Флаг для предотвращения двойного запроса

    useEffect(() => {
        if (isRequestSent.current) return; // Если запрос уже отправлялся, выходим
        isRequestSent.current = true; // Устанавливаем флаг

        const query = new URLSearchParams(window.location.search);
        const githubCode = query.get("code");

        console.log("🔹 GitHub Code:", githubCode);

        if (githubCode) {
            const authToken = localStorage.getItem("authToken");
            console.log("🔹 Токен перед отправкой запроса:", authToken);

            axios
                .post(
                    "http://localhost:8081/v1/profile/github",
                    { githubCode },
                    {
                        headers: {
                            "Content-Type": "application/json",
                            "Accept": "application/json",
                            Authorization: authToken ? `Bearer ${authToken}` : "",
                        },
                    }
                )
                .then((response) => {
                    console.log("✅ Успешный ответ:", response.data);
                    alert(response.data.message);
                    navigate("/profile");
                })
                .catch((error) => {
                    if (error.response) {
                        console.error("❌ Ошибка от сервера:", error.response.status, error.response.data);
                        alert(`Ошибка: ${error.response.status} - ${error.response.data.message || "Неизвестная ошибка"}`);
                    } else if (error.request) {
                        console.error("⚠️ Нет ответа от сервера:", error.request);
                        alert("Сервер не отвечает. Попробуйте позже.");
                    } else {
                        console.error("🚨 Ошибка при отправке запроса:", error.message);
                        alert("Ошибка при отправке запроса.");
                    }
                    navigate("/");
                });
        } else {
            alert("GitHub authorization failed.");
            navigate("/");
        }
    }, [navigate]);

    return <p>Processing GitHub authorization...</p>;
};

export default GitHubCallback;
