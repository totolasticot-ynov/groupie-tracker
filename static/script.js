document.addEventListener("DOMContentLoaded", () => {
    const input = document.getElementById("search-input");
    const box = document.getElementById("suggestions");

    if (input) {
        input.addEventListener("input", async (e) => {
            const q = e.target.value.trim();
            if (q.length < 1) { box.style.display = "none"; return; }

            try {
                const res = await fetch(`/search?q=${encodeURIComponent(q)}`);
                const data = await res.json();
                box.innerHTML = "";
                
                if (data.length > 0) {
                    box.style.display = "block";
                    data.forEach(item => {
                        const div = document.createElement("div");
                        div.className = "sugg-item";
                        div.innerHTML = `<span>${item.text}</span> <span class="sugg-tag">${item.type}</span>`;
                        div.onclick = () => window.location.href = `/artist?id=${item.artistId}`;
                        box.appendChild(div);
                    });
                } else { box.style.display = "none"; }
            } catch (e) { console.error(e); }
        });

        document.addEventListener("click", (e) => {
            if (!input.contains(e.target) && !box.contains(e.target)) box.style.display = "none";
        });
    }
});
