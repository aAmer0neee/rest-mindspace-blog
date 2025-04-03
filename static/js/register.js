document.getElementById("registerForm").addEventListener("submit", async function (e) {
    e.preventDefault();

    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
        const response = await fetch("/auth/register", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ username, password })
        });

        if (response.ok) {
            showToast("Регистрация успешна", "success");
        } else {
            if (response.status === 400) {
                showToast("Неверные данные", "error");
            } else {
                showToast("Неизвестная ошибка, попробуйте снова", "error");
            }
        }

    } catch (error) {
        console.error("Ошибка:", error);
        showToast("Произошла ошибка при регистрации. Попробуйте снова.", "error");
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
