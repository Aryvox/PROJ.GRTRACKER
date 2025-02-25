function toggleFavorite(characterId) {
    console.log('Toggle favorite for:', characterId); // Debug
    
    fetch('/favoris', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: `characterId=${characterId}`
    })
    .then(response => {
        if (!response.ok) {
            console.error('Server response:', response.status);
            throw new Error('Network response was not ok');
        }
        const btn = document.querySelector(`[data-character-id="${characterId}"]`);
        if (btn) {
            const isActive = btn.classList.toggle('active');
            btn.innerHTML = isActive ? '❤️' : '🤍';
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('Erreur lors de la mise à jour des favoris');
    });
} 