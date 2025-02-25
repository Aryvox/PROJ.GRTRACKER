// Fonction de recherche
async function searchCharacters(query) {
    try {
        const response = await fetch(`/api/search?q=${encodeURIComponent(query)}`);
        const data = await response.json();
        return data;
    } catch (error) {
        console.error('Error searching characters:', error);
        return null;
    }
}

// Gestion des favoris

// Gestion de la pagination
function updatePagination(currentPage, totalPages) {
    const prevButton = document.getElementById('prev-page');
    const nextButton = document.getElementById('next-page');
    
    prevButton.disabled = currentPage <= 1;
    nextButton.disabled = currentPage >= totalPages;
    
    prevButton.onclick = () => {
        if (currentPage > 1) {
            window.location.href = `/collection?page=${currentPage - 1}`;
        }
    };
    
    nextButton.onclick = () => {
        if (currentPage < totalPages) {
            window.location.href = `/collection?page=${currentPage + 1}`;
        }
    };
} 