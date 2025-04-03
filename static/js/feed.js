let currentPage = 1;
let totalPosts = 0;
const pageSize = 3;

async function fetchFeed(page) {
    try {
        const response = await fetch(`/data?page=${page}`);
        const data = await response.json();
        totalPosts = data.total || 0;
        renderFeed(data.posts);
        updatePagination();
    } catch (error) {
        console.error("Ошибка при получении фида:", error);
    }
}

function renderFeed(posts) {
    const container = document.getElementById('feed-container');
    container.innerHTML = '';

    if (!posts || posts.length === 0) {
        container.innerHTML = '<p>Записей еще нет</p>';
        return;
    }

    posts.forEach(post => {
        const cleanPreview = document.createElement("div");
        cleanPreview.innerHTML = post.preview;
        const previewText = cleanPreview.textContent.trim();

        const postElement = `
            <a href="/article?id=${post.id}" class="feed-item">
                <h2>${post.title}</h2>
                <div class="content-preview">
                    <p>${previewText}</p>
                </div>
            </a>
        `;
        container.insertAdjacentHTML("beforeend", postElement);
    });
}

function updatePagination() {
    const totalPages = Math.max(1, Math.ceil(totalPosts / pageSize));
    document.getElementById('page-info').textContent = `Страница ${currentPage} из ${totalPages}`;

    document.getElementById('prev-btn').disabled = currentPage <= 1;
    document.getElementById('next-btn').disabled = currentPage >= totalPages;
}

function changePage(direction) {
    const totalPages = Math.max(1, Math.ceil(totalPosts / pageSize));
    
    if ((direction === -1 && currentPage > 1) || (direction === 1 && currentPage < totalPages)) {
        currentPage += direction;
        fetchFeed(currentPage);
    }
}

document.getElementById('prev-btn').addEventListener('click', () => changePage(-1));
document.getElementById('next-btn').addEventListener('click', () => changePage(1));

fetchFeed(currentPage);
