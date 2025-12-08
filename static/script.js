document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("search-input");
    const suggestionsBox = document.getElementById("suggestions-box");

    if (searchInput) {
        searchInput.addEventListener("input", async (e) => {
            const query = e.target.value.trim();

            if (query.length < 1) {
                suggestionsBox.style.display = "none";
                return;
            }

            try {
                // Appel à l'API de recherche Go
                const response = await fetch(`/search?q=${encodeURIComponent(query)}`);
                const results = await response.json();

                // Vider les anciennes suggestions
                suggestionsBox.innerHTML = "";

                if (results.length > 0) {
                    suggestionsBox.style.display = "block";
                    results.forEach(item => {
                        // Création de l'élément visuel
                        const div = document.createElement("div");
                        div.className = "suggestion-item";
                        
                        // Affichage formaté : "Texte" à gauche, "TYPE" à droite
                        div.innerHTML = `
                            <span class="suggestion-text">${item.text}</span>
                            <span class="suggestion-type">${item.type}</span>
                        `;

                        // Clic = Redirection vers l'artiste
                        div.addEventListener("click", () => {
                            window.location.href = `/artist?id=${item.artistId}`;
                        });

                        suggestionsBox.appendChild(div);
                    });
                } else {
                    suggestionsBox.style.display = "none";
                }

            } catch (error) {
                console.error("Erreur recherche:", error);
            }
        });

        // Cacher les suggestions si on clique ailleurs
        document.addEventListener("click", (e) => {
            if (!searchInput.contains(e.target) && !suggestionsBox.contains(e.target)) {
                suggestionsBox.style.display = "none";
            }
        });
    }
});
