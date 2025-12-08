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
                        div.innerHTML = `<span style="font-weight:bold; color:white;">${item.text}</span> <span style="color:#e91e63; font-size:0.8em;">${item.type}</span>`;
                        div.addEventListener("click", () => window.location.href = `/artist?id=${item.artistId}`);
                        suggestionsBox.appendChild(div);
                    });
                } else { suggestionsBox.style.display = "none"; }
            } catch (error) { console.error(error); }
        });
        document.addEventListener("click", (e) => {
            if (!searchInput.contains(e.target) && !suggestionsBox.contains(e.target)) suggestionsBox.style.display = "none";
        });
    }
});
