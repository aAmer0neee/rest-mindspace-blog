document.addEventListener("DOMContentLoaded", function () {
    const form = document.getElementById("authForm");

    if (form) {
        form.addEventListener("submit", async function (e) {
            e.preventDefault();

            const username = document.getElementById("username").value;
            const password = document.getElementById("password").value;

            try {
                const response = await fetch("/auth", {
                    method: "POST",
                    credentials: "include",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ username, password })
                });

                if (response.ok) {
                    window.location.href = "/admin";
                } else {
                    if (response.status === 401) {
                        showToast("Неверный логин или пароль", "error");
                    } else if (response.status === 400) {
                        showToast("Неверные данные", "error");
                    } else {
                        showToast("Неизвестная ошибка, попробуйте снова", "error");
                    }
                }

            } catch (error) {
                console.error("Ошибка авторизации:", error);
                showToast("Ошибка авторизации, попробуйте снова", "error");
            }
        });
    } else {
        console.error('Форма с id="authForm" не найдена');
    }
});

function showToast(message, type) {
    const toast = document.getElementById("toast");
    toast.innerText = message;
    toast.className = `toast show ${type}`;

    setTimeout(() => {
        toast.className = "toast";
    }, 3000);
}