// Affiche les détails d'un personnage dans un modal
function showDetails(characterId) {
    const modal = document.getElementById('character-modal');
    const content = document.getElementById('modal-content');
    
    // Affichage du loader
    content.innerHTML = '<div class="loading">Chargement...</div>';
    modal.style.display = 'block';
    document.body.classList.add('modal-open');

    fetch(`/api/character/${characterId}`)
        .then(response => {
            if (!response.ok) throw new Error('Erreur réseau');
            return response.json();
        })
        .then(data => {
            if (!data.data || !data.data.results || !data.data.results[0]) {
                throw new Error('Données invalides');
            }
            const character = data.data.results[0];
            content.innerHTML = `
                <h2>${escapeHtml(character.name)}</h2>
                <img src="${character.thumbnail.path}.${character.thumbnail.extension}" 
                     alt="${escapeHtml(character.name)}"
                     onerror="this.src='/static/img/placeholder.jpg'">
                <p>${character.description || 'Aucune description disponible.'}</p>
                <div class="character-stats">
                    <p>Comics: ${character.comics.available}</p>
                    <p>Series: ${character.series.available}</p>
                    <p>Stories: ${character.stories.available}</p>
                </div>
            `;
        })
        .catch(error => {
            console.error('Error:', error);
            content.innerHTML = `
                <div class="error-message">
                    Une erreur est survenue lors du chargement des détails.
                    <button onclick="closeModal()" class="retry-btn">Fermer</button>
                </div>
            `;
        });
}

// Sécurise le contenu HTML
function escapeHtml(unsafe) {
    if (!unsafe) return '';
    return unsafe
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}

// Fermeture du modal
function closeModal() {
    const modal = document.getElementById('character-modal');
    if (modal) {
        modal.style.display = 'none';
        document.body.classList.remove('modal-open');
    }
}

// Fermeture avec Escape
document.addEventListener('keydown', event => {
    if (event.key === 'Escape') closeModal();
});

// Fermeture en cliquant en dehors du modal
document.addEventListener('click', event => {
    const modal = document.getElementById('character-modal');
    if (event.target === modal) {
        closeModal();
    }
}); 