document.addEventListener('DOMContentLoaded', function() {
    console.log('Script chargé');
    
    // Initialiser les favoris
    if (!localStorage.getItem('favoriteIds')) {
        localStorage.setItem('favoriteIds', '[]');
    }

    // Marquer les boutons favoris
    markFavoriteButtons();

    // Si on est sur la page favoris, afficher les favoris
    if (window.location.pathname === '/favoris') {
        displayFavorites();
    }
});

function toggleFavorite(button, characterId) {
    console.log('Toggle favori pour:', characterId);
    
    let favoriteIds = JSON.parse(localStorage.getItem('favoriteIds') || '[]');
    const index = favoriteIds.indexOf(characterId);
    
    if (index === -1) {
        // Ajouter aux favoris
        favoriteIds.push(characterId);
        button.textContent = '❤️';
    } else {
        // Retirer des favoris
        favoriteIds.splice(index, 1);
        button.textContent = '🤍';
    }
    
    localStorage.setItem('favoriteIds', JSON.stringify(favoriteIds));
    console.log('Favoris mis à jour:', favoriteIds);

    // Actualiser la page des favoris si nécessaire
    if (window.location.pathname === '/favoris') {
        displayFavorites();
    }
}

function displayFavorites() {
    console.log('Affichage des favoris');
    const container = document.getElementById('favorites-grid');
    if (!container) {
        console.error('Container des favoris non trouvé');
        return;
    }

    const favoriteIds = JSON.parse(localStorage.getItem('favoriteIds') || '[]');
    console.log('IDs des favoris:', favoriteIds);

    if (favoriteIds.length === 0) {
        container.innerHTML = '<p class="no-favorites">Aucun favori pour le moment</p>';
        return;
    }

    // Charger les détails de chaque personnage
    Promise.all(favoriteIds.map(id => 
        fetch(`/api/character/${id}`)
            .then(response => {
                if (!response.ok) throw new Error('Network response was not ok');
                return response.json();
            })
            .then(data => data.data.results[0])
            .catch(error => {
                console.error('Error fetching character:', id, error);
                return null;
            })
    ))
    .then(characters => {
        const validCharacters = characters.filter(char => char !== null);
        if (validCharacters.length === 0) {
            container.innerHTML = '<p class="error">Aucun personnage trouvé</p>';
            return;
        }

        container.innerHTML = validCharacters.map(char => `
            <div class="character-card">
                <img src="${char.thumbnail.path}.${char.thumbnail.extension}" alt="${char.name}">
                <h2>${char.name}</h2>
                <p>${char.description || 'Aucune description disponible.'}</p>
                <div class="card-actions">
                    <button onclick="toggleFavorite(this, '${char.id}')" class="favorite-btn">❤️</button>
                    <a href="/character/details/${char.id}" class="details-btn">📖 Détails</a>
                </div>
            </div>
        `).join('');
    })
    .catch(error => {
        console.error('Error loading favorites:', error);
        container.innerHTML = '<p class="error">Erreur lors du chargement des favoris</p>';
    });
}

function markFavoriteButtons() {
    const favoriteIds = JSON.parse(localStorage.getItem('favoriteIds') || '[]');
    document.querySelectorAll('.favorite-btn').forEach(button => {
        const onclick = button.getAttribute('onclick');
        if (onclick) {
            const match = onclick.match(/toggleFavorite\(this,\s*'(\d+)'\)/);
            if (match && favoriteIds.includes(match[1])) {
                button.textContent = '❤️';
            }
        }
    });
}

// Charger les favoris si on est sur la page des favoris
if (window.location.pathname === '/favoris') {
    console.log('On favorites page, setting up...'); // Debug log
    document.addEventListener('DOMContentLoaded', displayFavorites);
}

// S'assurer que le script est chargé
console.log('Favorites script loaded');

// Ajouter un style pour le bouton désactivé
document.head.insertAdjacentHTML('beforeend', `
    <style>
        .favorite-btn:disabled {
            opacity: 0.5;
            cursor: not-allowed;
        }
    </style>
`);

// Fonction pour ajouter un favori
function addToFavorites(characterId, name, description, image) {
    let favorites = JSON.parse(localStorage.getItem('favorites') || '[]');
    
    // Vérifier si déjà dans les favoris
    if (!favorites.some(f => f.id === characterId)) {
        favorites.push({
            id: characterId,
            name: name,
            description: description,
            image: image
        });
        localStorage.setItem('favorites', JSON.stringify(favorites));
    }
    return true;
}

// Fonction pour retirer un favori
function removeFromFavorites(characterId) {
    let favorites = JSON.parse(localStorage.getItem('favorites') || '[]');
    favorites = favorites.filter(f => f.id !== characterId);
    localStorage.setItem('favorites', JSON.stringify(favorites));
    
    // Si on est sur la page favoris, actualiser l'affichage
    if (window.location.pathname === '/favoris') {
        displayFavorites();
    }
}

// Fonction pour basculer l'état favori
function toggleFavorite(button, characterId, name, description, image) {
    let favorites = JSON.parse(localStorage.getItem('favorites') || '[]');
    const isFavorite = favorites.some(f => f.id === characterId);
    
    if (isFavorite) {
        removeFromFavorites(characterId);
        button.textContent = '🤍';
    } else {
        addToFavorites(characterId, name, description, image);
        button.textContent = '❤️';
    }
}

// Initialiser l'affichage des favoris au chargement de la page
document.addEventListener('DOMContentLoaded', function() {
    if (window.location.pathname === '/favoris') {
        displayFavorites();
    }
}); 