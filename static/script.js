(function() {
    'use strict';
    
    let debounceTimer;
    const DEBOUNCE_DELAY = 300;
    
    document.addEventListener('DOMContentLoaded', function() {
        const input = document.getElementById('search-input');
        const box = document.getElementById('suggestions');
        
        if (!input || !box) {
            console.warn('Éléments de recherche non trouvés');
            return;
        }
        
        // Debounced search function
        function performSearch(query) {
            if (query.length < 1) {
                box.style.display = 'none';
                return;
            }
            
            const url = `/api/search?q=${encodeURIComponent(query)}`;
            
            fetch(url)
                .then(res => {
                    if (!res.ok) throw new Error(`HTTP ${res.status}`);
                    return res.json();
                })
                .then(data => {
                    box.innerHTML = '';
                    
                    if (Array.isArray(data) && data.length > 0) {
                        box.style.display = 'block';
                        const fragment = document.createDocumentFragment();
                        
                        data.forEach(item => {
                            const div = document.createElement('div');
                            div.className = 'sugg-item';
                            div.setAttribute('role', 'option');
                            div.innerHTML = `<span>${escapeHtml(item.text)}</span> <span class="sugg-tag">${escapeHtml(item.type)}</span>`;
                            div.style.cursor = 'pointer';
                            
                            div.addEventListener('click', () => {
                                window.location.href = `/artist?id=${item.artistId}`;
                            });
                            
                            fragment.appendChild(div);
                        });
                        
                        box.appendChild(fragment);
                    } else {
                        box.style.display = 'none';
                    }
                })
                .catch(error => {
                    console.error('Erreur recherche:', error);
                    box.style.display = 'none';
                });
        }
        
        input.addEventListener('input', (e) => {
            const query = e.target.value.trim();
            clearTimeout(debounceTimer);
            debounceTimer = setTimeout(() => performSearch(query), DEBOUNCE_DELAY);
        });
        
        document.addEventListener('click', (e) => {
            if (!input.contains(e.target) && !box.contains(e.target)) {
                box.style.display = 'none';
            }
        });
        
        // Close on Escape key
        document.addEventListener('keydown', (e) => {
            if (e.key === 'Escape') {
                box.style.display = 'none';
            }
        });
    });
    
    // XSS Prevention
    function escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }
})();



