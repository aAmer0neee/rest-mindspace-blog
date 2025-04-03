document.getElementById("postForm").addEventListener("submit", async function(e) {
    e.preventDefault();  

    const content = document.getElementById("content").value;

    const formData = { content };

    try {
        const response = await fetch("/admin", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(formData)
        });

        if (!response.ok) throw new Error("Ошибка при отправке данных");

        showToast("Статья успешно добавлена!", "success");

        document.getElementById("content").value = "";
    } catch (error) {
        console.error("Ошибка:", error);
        showToast("Ошибка при отправке статьи.", "error");
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