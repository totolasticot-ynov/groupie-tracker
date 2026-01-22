# Groupie Tracker

Une application web Go pour explorer et tracker les artistes musicaux, leurs concerts et leurs relations géographiques.

## 🎯 Fonctionnalités

- 🎵 **Exploration d'artistes**: Parcourez une base de données complète d'artistes
- 🔍 **Recherche avancée**: Trouvez des artistes par nom, année de création, nombre de membres, etc.
- 🗺️ **Localisation des concerts**: Visualisez les dates et lieux des concerts par artiste
- 💳 **Paiements sécurisés**: Intégration PayPal pour les transactions
- 📊 **API RESTful**: Accédez aux données via des endpoints JSON

## 📋 Prérequis

- Go 1.21+
- Un navigateur web moderne
- Compte PayPal Sandbox (pour les tests de paiement)

## 🚀 Installation

1. Clonez le repository:
```bash
git clone https://github.com/yourusername/groupie-tracker.git
cd groupie-tracker
```

2. Installez les dépendances:
```bash
go mod download
```

3. Lancez l'application:
```bash
go run main.go
```

L'application démarre sur `http://localhost:8080`

## 📁 Structure du projet

```
groupie-tracker/
├── main.go                 # Point d'entrée de l'application
├── internal/
│   └── server/
│       ├── server.go      # Serveur HTTP et routes
│       ├── page.go        # Rendu des pages
│       └── paypal.go      # Intégration PayPal
├── templates/             # Fichiers HTML
├── static/                # Fichiers CSS, JS, images
└── README.md
```

## 🔗 Routes principales

| Route | Description |
|-------|-------------|
| `/` | Page d'accueil |
| `/search` | Recherche d'artistes |
| `/explore` | Page d'exploration avec intégration PayPal |
| `/artist?id=<id>` | Détails d'un artiste spécifique |
| `/api/artists` | API JSON - Liste des artistes |
| `/api/search?q=<query>` | API JSON - Recherche d'artistes |

## 💻 Utilisation

### Rechercher un artiste
1. Allez sur la page `/search`
2. Entrez le nom ou les critères de recherche
3. Cliquez sur "Rechercher"

### Consulter les concerts
Visitez la page d'un artiste pour voir toutes les dates et localisations de ses concerts.

### Effectuer un paiement
Accédez à `/explore` et cliquez sur "Payer avec PayPal" pour tester les paiements en mode Sandbox.

## 🔧 Configuration

Les variables d'environnement importantes:
- `PORT` - Port du serveur (défaut: 8080)
- `PAYPAL_CLIENT_ID` - ID client PayPal
- `PAYPAL_CLIENT_SECRET` - Secret client PayPal

## 📝 Licence

Ce projet est sous licence MIT.

## 👤 Auteur

Dhordain Thomas, Benoit Augustin, Klapczynski Esteban

## 🤝 Contribution

Les contributions sont bienvenues ! Veuillez créer une pull request avec vos modifications.

## 📞 Support

Pour toute question ou problème, veuillez ouvrir une issue sur GitHub.