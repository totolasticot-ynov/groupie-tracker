document.addEventListener("DOMContentLoaded", () => {
    const searchInput = document.getElementById("search-input");
    const suggestionsBox = document.getElementById("suggestions-box");

    if (searchInput) {
        searchInput.addEventListener("input", async (e) => {
            const query = e.target.value.trim();
            if (query.length < 1) { suggestionsBox.style.display = "none"; return; }

            try {
                const response = await fetch(`/search?q=${encodeURIComponent(query)}`);
                const results = await response.json();
                
                suggestionsBox.innerHTML = "";
                
                if (results.length > 0) {
                    suggestionsBox.style.display = "block";
                    results.forEach(item => {
                        const div = document.createElement("div");
                        div.className = "suggestion-item";
                        // Affichage Type en couleur Accent
                        div.innerHTML = `
                            <span style="font-weight:600; color:white;">${item.text}</span> 
                            <span style="color:#1db954; font-size:0.75em; border:1px solid #1db954; padding:2px 6px; border-radius:4px;">${item.type}</span>
                        `;
                        div.addEventListener("click", () => window.location.href = `/artist?id=${item.artistId}`);
                        suggestionsBox.appendChild(div);
                    });
                } else { suggestionsBox.style.display = "none"; }
            } catch (error) { console.error(error); }
        });

        // Fermer si on clique ailleurs
        document.addEventListener("click", (e) => {
            if (!searchInput.contains(e.target) && !suggestionsBox.contains(e.target)) {
                suggestionsBox.style.display = "none";
            }
        });
    }
});
