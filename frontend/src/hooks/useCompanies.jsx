import { useState, useEffect } from "react";

export default function useCompanies() {
    const [companies, setCompanies] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        console.log("Fetching companies...");

        fetch("http://localhost:8082/v1/companies", {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })
            .then(response => {
                console.log("Response: ", response); // Выводим сам объект Response

                if (!response.ok) {
                    throw new Error("Ошибка загрузки данных");
                }

                return response.json();
            })
            .then(data => {
                console.log("data: ", data); // Выводим данные ответа

                // Теперь мы ожидаем, что объект будет содержать ключ `companies`, который является массивом
                if (data && Array.isArray(data.companies)) {
                    setCompanies(data.companies); // Извлекаем массив из объекта
                } else {
                    console.error("Полученные данные не содержат массив companies:", data);
                    setError("Полученные данные не содержат массив companies.");
                }

                setLoading(false);
            })
            .catch(error => {
                console.error("Ошибка запроса:", error);
                setError(error.message);
                setLoading(false);
            });
    }, []);

    return { companies, loading, error };
}
