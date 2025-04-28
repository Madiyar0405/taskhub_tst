const exchangeCodeForToken = async (code) => {
    try {
        const response = await fetch("http://localhost:8080/oauth/github/token", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ code }),
        });

        const data = await response.json();
        console.log("Access Token:", data.access_token);
    } catch (error) {
        console.error("Ошибка при обмене кода на токен:", error);
    }
};

// Получаем код из URL
const params = new URLSearchParams(window.location.search);
const code = params.get("code");

if (code) {
    exchangeCodeForToken(code);
}
