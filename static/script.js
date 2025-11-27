// [cite: 52] Les variables passées dans le HTML sont maintenant accessibles ici.

console.log("Le script.js est chargé !");

// Vérification que le transfert a fonctionné
if (typeof tableauAvecTousLesArtistes !== 'undefined') {
    console.log("J'ai accès aux artistes dans script.js ! Nombre d'artistes : " + tableauAvecTousLesArtistes.length);
    
    // Exemple : C'est ici que tu coderas ta barre de recherche
    // en filtrant "tableauAvecTousLesArtistes"
}
