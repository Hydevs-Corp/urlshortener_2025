# Raccourcisseur d'URL en Go

Un service web performant de raccourcissement et de gestion d'URLs développé en Go. Cette application transforme des URL longues en codes courts et uniques. Lorsqu'une URL courte est consultée, le système redirige instantanément l'utilisateur vers l'URL d'origine tout en enregistrant le clic de manière asynchrone afin d'assurer une latence de redirection nulle.

Le service comprend aussi un composant de surveillance qui vérifie périodiquement la disponibilité des URL longues et consigne tout changement d'état. L'interaction se fait via une API RESTful et une interface en ligne de commande (CLI) complète.

## ✨ Fonctionnalités

- **Raccourcissement d'URL** : Génère des codes courts uniques de 6 caractères alphanumériques. Gère les collisions via un mécanisme de retry.
- **Redirection instantanée** : Redirige les utilisateurs vers l'URL originale en utilisant un code de statut `302 Found` pour une rapidité maximale.
- **Analytics asynchrone** : Le suivi des clics est traité en arrière-plan avec des Goroutines et des channels bufferisés, garantissant que la redirection utilisateur n'est jamais bloquée.
- **Surveillance de la santé des URLs** : Vérifie périodiquement si les URL longues sont encore accessibles (réponses HTTP 200/3xx). En cas de changement d'état, une notification factice est écrite dans les logs du serveur.
- **API RESTful** : API claire pour créer, gérer et récupérer les statistiques des liens.
- **Interface en ligne de commande (CLI)** : Une CLI complète pour interagir avec le service sans interface graphique.

## 🛠️ Stack technique

- **Go** : Langage principal.
- **Gin** : Framework HTTP performant pour construire l'API REST.
- **GORM** : ORM pour la persistance avec SQLite.
- **Cobra** : Bibliothèque pour créer une CLI moderne.
- **Viper** : Gestion de configuration.
- **SQLite** : Base de données embarquée, sans serveur.

## 🚀 Pour commencer

Suivez ces étapes pour configurer le projet et exécuter l'application.

### 1. Prérequis

- [Go](https://golang.org/doc/install) (version 1.21 ou supérieure)
- [Git](https://git-scm.com/)

### 2. Installation

1.  **Clonez le dépôt :**

    ```sh
    git clone https://github.com/axellelanca/urlshortener.git
    cd urlshortener
    ```

2.  **Téléchargez les dépendances :**

    ```sh
    go mod tidy
    ```

3.  **Construisez l'exécutable :**

    Cette commande compile l'application et crée un binaire `url-shortener` à la racine du projet.

    ```sh
    go build -o url-shortener
    ```

### 3. Initialisation de la base de données

Avant de démarrer le serveur, créez le fichier de base de données SQLite et ses tables en exécutant les migrations GORM :

```sh
./url-shortener migrate
```

Vous devriez voir un message de succès confirmant la création des tables. Un fichier `url_shortener.db` sera créé à la racine du projet.

## utilisation

### 1. Démarrer le serveur

Cette commande démarre le serveur web Gin, les workers de traitement des clics asynchrones, et le moniteur de santé des URLs.

```sh
./url-shortener run-server
```

Gardez cette fenêtre de terminal ouverte. Elle affichera les logs des requêtes HTTP, du traitement des clics, et des notifications de surveillance des URLs.

### 2. Interagir avec le service (dans un nouveau terminal)

Ouvrez une **nouvelle fenêtre de terminal** pour utiliser la CLI ou tester l'API pendant que le serveur est en cours d'exécution.

#### Créer une URL courte (CLI)

Raccourcissez une URL longue en utilisant la commande `create` :

```sh
./url-shortener create --url="https://www.youtube.com/watch?v=dQw4w9WgXcQ"
```

La sortie sera similaire à ceci :

```
URL courte créée avec succès:
Code: XYZ123
URL complète: http://localhost:8080/XYZ123
```

#### Accéder à l'URL courte

1.  Ouvrez votre navigateur web et accédez à l'URL courte fournie (par exemple, `http://localhost:8080/XYZ123`).
2.  Vous serez redirigé instantanément vers l'URL longue d'origine.
3.  Dans le terminal du serveur, vous verrez des logs indiquant qu'un clic a été enregistré.

#### Voir les statistiques du lien (CLI)

Vérifiez combien de fois votre URL courte a été visitée :

```sh
./url-shortener stats --code="XYZ123"
```

La sortie affichera le nombre total de clics :

```
Statistiques pour le code court: XYZ123
URL longue: https://www.youtube.com/watch?v=dQw4w9WgXcQ
Total de clics: 1
```

## 🌐 Points de terminaison de l'API

| Méthode | Point de terminaison              | Description                                                              |
| :------ | :-------------------------------- | :----------------------------------------------------------------------- |
| `GET`   | `/health`                         | Vérifie la santé du service.                                             |
| `POST`  | `/api/v1/links`                   | Crée une nouvelle URL courte. Attend `{"long_url": "..."}`.              |
| `GET`   | `/{shortCode}`                    | Redirige vers l'URL d'origine et enregistre le clic.                     |
| `GET`   | `/api/v1/links/{shortCode}/stats` | Récupère les statistiques (clics totaux) pour une URL courte spécifique. |

#### Exemple avec `curl`

**Créer une URL courte :**

```sh
curl -X POST http://localhost:8080/api/v1/links
     -H "Content-Type: application/json"
     -d '{"long_url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ"}'
```

**Vérifier la santé du service :**

```sh
curl http://localhost:8080/health
```

## 🤝 Contribuer

Les contributions sont les bienvenues ! N'hésitez pas à soumettre une demande de tirage ou à ouvrir un problème pour tout bogue ou demande de fonctionnalité.

## Architecture du Projet

```
url-shortener/
├── cmd/
│   ├── root.go             # Initialise la commande racine Cobra et ses sous-commandes
│   ├── server/
│   │   └── server.go       # Logique pour la commande 'run-server' (lance le serveur Gin, les workers de clics, le moniteur)
│   └── cli/
│       ├── create.go       # Logique pour la commande 'create' (crée un lien via CLI)
│       ├── stats.go        # Logique pour la commande 'stats' (affiche les statistiques d'un lien via CLI)
│       └── migrate.go      # Logique pour la commande 'migrate' (exécute les migrations GORM)
├── internal/
│   ├── api/
│   │   └── handlers.go     # Fonctions de gestion des requêtes HTTP (handlers Gin pour les routes API)
│   ├── models/
│   │   ├── link.go         # Définition de la structure GORM 'Link'
│   │   └── click.go        # Définition de la structure GORM 'Click'
│   ├── services/
│   │   ├── link_service.go # Logique métier pour les liens (ex: génération de code, validation)
│   │   └── click_service.go # Logique métier pour les clics (optionnel, peut être directement dans le worker si simple)
│   ├── workers/
│   │   └── click_worker.go # Goroutine et logique pour l'enregistrement asynchrone des clics
│   ├── monitor/
│   │   └── url_monitor.go  # Logique pour la surveillance périodique de l'état des URLs
│   ├── config/
│   │   └── config.go       # Chargement et structure de la configuration de l'application (Viper)
│   └── repository/
│       ├── link_repository.go # Interface et implémentation GORM pour les opérations CRUD sur 'Link'
│       └── click_repository.go # Interface et implémentation GORM pour les opérations CRUD sur 'Click'
├── configs/
│   └── config.yaml         # Fichier de configuration par défaut pour Viper
├── go.mod                  # Fichier de module Go (liste des dépendances du projet)
├── go.sum                  # Sommes de contrôle pour la sécurité des dépendances
└── README.md               # Documentation du projet (installation, utilisation, etc.)

```
