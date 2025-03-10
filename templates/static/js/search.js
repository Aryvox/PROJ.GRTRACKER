document.addEventListener('DOMContentLoaded', function() {
    const sortSelect = document.getElementById('sort');
    const comicsSelect = document.getElementById('comics');
    const seriesSelect = document.getElementById('series');
    
    // Appliquer les filtres quand ils changent
    sortSelect.addEventListener('change', applyFilters);
    comicsSelect.addEventListener('change', applyFilters);
    seriesSelect.addEventListener('change', applyFilters);
    
    // Restaurer les valeurs des filtres depuis l'URL
    restoreFiltersFromURL();
});

function restoreFiltersFromURL() {
    const urlParams = new URLSearchParams(window.location.search);
    
    if (urlParams.has('sort')) {
        document.getElementById('sort').value = urlParams.get('sort');
    }
    if (urlParams.has('comics')) {
        document.getElementById('comics').value = urlParams.get('comics');
    }
    if (urlParams.has('series')) {
        document.getElementById('series').value = urlParams.get('series');
    }
}

function applyFilters() {
    const cards = document.querySelectorAll('.character-card');
    const sortValue = document.getElementById('sort').value;
    const comicsValue = document.getElementById('comics').value;
    const seriesValue = document.getElementById('series').value;
    
    // Convertir les cartes en tableau pour le tri
    let cardsArray = Array.from(cards);
    
    // Filtrer par nombre de comics
    if (comicsValue) {
        cardsArray = cardsArray.filter(card => {
            const comicsCount = parseInt(card.dataset.comics);
            const [min, max] = comicsValue.split('-').map(n => n === '+' ? Infinity : parseInt(n));
            return comicsCount >= min && (max === Infinity ? true : comicsCount <= max);
        });
    }
    
    // Filtrer par nombre de séries
    if (seriesValue) {
        cardsArray = cardsArray.filter(card => {
            const seriesCount = parseInt(card.dataset.series);
            const [min, max] = seriesValue.split('-').map(n => n === '+' ? Infinity : parseInt(n));
            return seriesCount >= min && (max === Infinity ? true : seriesCount <= max);
        });
    }
    
    // Trier les cartes
    if (sortValue) {
        cardsArray.sort((a, b) => {
            const nameA = a.querySelector('h2').textContent;
            const nameB = b.querySelector('h2').textContent;
            return sortValue === 'name_asc' ? 
                nameA.localeCompare(nameB) : 
                nameB.localeCompare(nameA);
        });
    }
    
    // Masquer toutes les cartes
    cards.forEach(card => card.style.display = 'none');
    
    // Afficher les cartes filtrées
    cardsArray.forEach(card => card.style.display = '');
    
    // Mettre à jour l'URL avec les filtres
    updateURL();
}

function updateURL() {
    const sortValue = document.getElementById('sort').value;
    const comicsValue = document.getElementById('comics').value;
    const seriesValue = document.getElementById('series').value;
    const searchQuery = document.querySelector('.search-input').value;
    
    const params = new URLSearchParams(window.location.search);
    
    if (searchQuery) params.set('q', searchQuery);
    if (sortValue) params.set('sort', sortValue);
    if (comicsValue) params.set('comics', comicsValue);
    if (seriesValue) params.set('series', seriesValue);
    
    const newURL = `${window.location.pathname}?${params.toString()}`;
    window.history.pushState({}, '', newURL);
} 