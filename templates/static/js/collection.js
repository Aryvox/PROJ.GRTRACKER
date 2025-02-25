function toggleFavorite(button, characterId) {
    const formData = new FormData();
    formData.append('characterId', characterId);
    
    fetch('/favoris', {
        method: 'POST',
        body: formData
    })
    .then(response => {
        if (response.ok) {
            if (button.classList.contains('active')) {
                button.classList.remove('active');
                button.innerHTML = '❤️ Favoris';
            } else {
                button.classList.add('active');
                button.innerHTML = '💖 Favori';
            }
        } else {
            console.error('Erreur lors de la mise à jour des favoris');
        }
    })
    .catch(error => {
        console.error('Erreur:', error);
    });
}