# Raccourcisseur d'URL en Go

Un service web performant de raccourcissement et de gestion d'URLs développé en Go. Cette application transforme des URL longues en codes courts et uniques. Lorsqu'une URL courte est consultée, le système redirige instantanément l'utilisateur vers l'URL d'origine tout en enregistrant le clic de manière asynchrone afin d'assurer une latence de redirection nulle.

Le service comprend aussi un composant de surveillance qui vérifie périodiquement la disponibilité des URL longues et consigne tout changement d'état. L'interaction se fait via une API RESTful et une interface en ligne de commande (CLI) complète.

## ✨ Fonctionnalités

-   **Raccourcissement d'URL** : Génère des codes courts uniques de 6 caractères alphanumériques. Gère les collisions via un mécanisme de retry.
-   **Redirection instantanée** : Redirige les utilisateurs vers l'URL originale en utilisant un code de statut `302 Found` pour une rapidité maximale.
-   **Analytics asynchrone** : Le suivi des clics est traité en arrière-plan avec des Goroutines et des channels bufferisés, garantissant que la redirection utilisateur n'est jamais bloquée.
-   **Surveillance de la santé des URLs** : Vérifie périodiquement si les URL longues sont encore accessibles (réponses HTTP 200/3xx). En cas de changement d'état, une notification factice est écrite dans les logs du serveur.
-   **API RESTful** : API claire pour créer, gérer et récupérer les statistiques des liens.
-   **Interface en ligne de commande (CLI)** : Une CLI complète pour interagir avec le service sans interface graphique.

## 🛠️ Stack technique

-   **Go** : Langage principal.
-   **Gin** : Framework HTTP performant pour construire l'API REST.
-   **GORM** : ORM pour la persistance avec SQLite.
-   **Cobra** : Bibliothèque pour créer une CLI moderne.
-   **Viper** : Gestion de configuration.
-   **SQLite** : Base de données embarquée, sans serveur.

## 🚀 Pour commencer

Suivez ces étapes pour configurer le projet et exécuter l'application.

### 1. Prérequis

-   [Go](https://golang.org/doc/install) (version 1.21 ou supérieure)
-   [Git](https://git-scm.com/)

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

## 📁 Structure du projet

Le projet suit une structure modulaire qui sépare les préoccupations pour la maintenabilité et l'évolutivité.

```
url-shortener/
├── cmd/                  # Commandes Cobra pour la CLI et le serveur
├── internal/             # Logique principale de l'application (privée)
│   ├── api/              # Gestionnaires HTTP Gin
│   ├── config/           # Configuration Viper
│   ├── models/           # Modèles de données GORM (Link, Click)
│   ├── monitor/          # Logique de surveillance de la santé des URLs
│   ├── repository/       # Couche d'accès aux données (implémentation GORM)
│   ├── services/         # Logique métier
│   └── workers/          # Workers en arrière-plan asynchrones
├── configs/              # Fichiers de configuration (ex: config.yaml)
├── go.mod                # Dépendances du module Go
├── main.go               # Point d'entrée principal de l'application
└── url_shortener.db      # Fichier de base de données SQLite
```

## 🤝 Contribuer

Les contributions sont les bienvenues ! N'hésitez pas à soumettre une demande de tirage ou à ouvrir un problème pour tout bogue ou demande de fonctionnalité.

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier [LICENSE](LICENSE) pour plus de détails.

## Objectif du Projet

Ce TP vous met au défi de construire un service web performant de raccourcissement et de gestion d'URLs en Go. Votre application permettra de transformer une URL longue en une URL courte et unique. Chaque fois qu'une URL courte est visitée, le système redirigera instantanément l'utilisateur vers l'URL originale tout en enregistrant le clic de manière asynchrone, pour ne jamais ralentir la redirection.

Le service inclura également un moniteur pour vérifier périodiquement la disponibilité des URLs longues et notifier tout changement d'état. L'interaction se fera via une API RESTful et une interface en ligne de commande (CLI) complète.

## Connaissances Mobilisées

Ce projet est une synthèse complète et pratique de tous les concepts abordés durant ce module de Go (normalement il n'y aura pas trop de surprise) :

-   Syntaxe Go de base (structs, maps, boucles, conditions, etc.)
-   Concurrence (Goroutines, Channels) pour les tâches asynchrones et non-bloquantes
-   Interfaces CLI avec [Cobra](https://cobra.dev/)
-   Gestion des erreurs
-   Manipulation de données (JSON) pour les APIs
-   APIs RESTful avec le framework web [Gin](https://gin-gonic.com/)
-   Persistance des données avec l'ORM [GORM](https://gorm.io/) et SQLite
-   Gestion de configuration avec [Viper](https://github.com/spf13/viper)
-   Design patterns courants (Repository, Service) pour une architecture propre

## Fonctionnalités Attendues

### Core Features (Obligatoires)

1. **Raccourcissement d'URLs** :

-   Générer des codes courts uniques (6 caractères alphanumériques).
-   Gérer les collisions lors de la génération de codes via une logique de retry.

2. **Redirection instantanée** :

-   Rediriger les utilisateurs vers l'URL originale sans latence (code HTTP 302).
-   Analytics asynchrones :
-   Enregistrer les détails de chaque clic en arrière-plan via des Goroutines et un Channel bufferisé. La redirection ne doit jamais être bloquée par l'enregistrement du clic.

3. **Surveillance de l'état des URLs** :

-   Le service doit vérifier périodiquement (intervalle configurable via Viper) si les URLs longues sont toujours accessibles (réponse HTTP 200/3xx).
-   Si l'état d'une URL change (accessible leftrightarrow inaccessible), une fausse notification doit être générée dans les logs du serveur (ex: "[NOTIFICATION] L'URL ... est maintenant INACCESSIBLE.").

4. **APIs REST (via Gin)** :

-   `GET /health` : Vérifie l'état de santé du service.
-   `POST /api/v1/links` : Crée une nouvelle URL courte (attend un JSON {"long_url": "..."}).
-   `GET /{shortCode}` : Gère la redirection et déclenche l'analytics asynchrone.
-   `GET /api/v1/links/{shortCode}/stats` : Récupère les statistiques d'un lien (nombre total de clics).

5. **Interface CLI (via Cobra)** :

-   `./url-shortener run-server` : Lance le serveur API, les workers de clics et le moniteur d'URLs.
-   `./url-shortener create --url="https://..."` : Crée une URL courte depuis la ligne de commande.
-   `./url-shortener stats --code="xyz123"` : Affiche les statistiques d'un lien donné.
-   `./url-shortener migrate` : Exécute les migrations GORM pour la base de données.

6. **Features Avancées (Bonus - si le temps le permet)**

-   URLs personnalisées : Permettre aux utilisateurs de proposer leur propre alias (ex: /mon-alias-perso).
-   Expiration des liens : Les URLs courtes peuvent avoir une durée de vie limitée.
-   Rate limiting : Protection simple par IP pour les créations de liens.

## Architecture du Projet

Le projet suit une structure modulaire classique pour les applications Go, qui sépare bien les différences préoccupations du projet :

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

## Démarrage et Utilisation du Projet

Suivez ces étapes pour mettre en place le projet et tester votre application (quand elle fonctionnera, évidemment).

### 1. Préparation Initiale

1. **Clonez le dépôt :**

```bash
git clone https://github.com/axellelanca/urlshortener.git
cd urlshortener # Naviguez vers le dossier du projet cloné
```

2. **Téléchargez et nettoyez les dépendances :**

```bash
go mod tidy
```

## Pour tester votre projet :

### Construisez l'exécutable :

Ceci compile votre application et crée un fichier url-shortener à la racine du projet.

```bash
go build -o url-shortener
```

Désormais, toutes les commandes seront lancées avec ./url-shortener.

### Initialisation de la Base de Données

Avant de démarrer le serveur, créez le fichier de base de données SQLite et ses tables :

1.  **Exécutez les migrations :**

```bash
./url-shortener migrate
```

Un message de succès confirmera la création des tables. Un fichier url_shortener.db sera créé à la racine du projet.

### Lancer le Serveur et les Processus de Fond

C'est l'étape qui démarre le cœur de votre application. Elle démarre le serveur web, les workers qui enregistrent les clics, et le moniteur d'URLs.

Démarrez le service :

```bash
./url-shortener run-server
```

Laissez ce terminal ouvert et actif. Il affichera les logs du serveur HTTP, des workers de clics et du moniteur d'URLs.

### 4. Interagir avec le Service (Utilise un **Nouveau Terminal**)

Ouvre une **nouvelle fenêtre de terminal** pour exécuter les commandes CLI et tester les APIs pendant que le serveur est en cours d'exécution.

#### 4.1. Créer une URL courte (via la CLI)

Raccourcis une URL longue en utilisant la commande `create` :

```bash
./url-shortener create --url="https://www.example.com/ma-super-url-de-test-pour-le-tp-go-final"
```

Tu obtiendras un message similaire à :

```bash
URL courte créée avec succès:
Code: XYZ123
URL complète: http://localhost:8080/XYZ123
```

Note le Code (ex: XYZ123) et l'URL complète pour les étapes suivantes.

#### 4.2. Accéder à l'URL courte (via Navigateur)

1. Ouvre ton navigateur web et accède à l'URL complète que tu as obtenue (par exemple, http://localhost:8080/XYZ123).
2. Le navigateur devrait te rediriger instantanément vers l'URL longue originale. Dans le terminal où le serveur tourne (./url-shortener run-server), tu devrais voir des logs indiquant qu'un clic a été détecté et envoyé au worker asynchrone.

#### 4.3. Consulter les Statistiques (via la CLI)

Vérifie combien de fois ton URL courte a été visitée :

1. Affiche les statistiques :

```
./url-shortener stats --code="XYZ123"
```

Le terminal affichera :

```
Statistiques pour le code court: XYZ123
URL longue: [https://www.example.com/ma-super-url-de-test-pour-le-tp-go-final](https://www.example.com/ma-super-url-de-test-pour-le-tp-go-final)
Total de clics: 1
```

(Le nombre de clics augmentera à chaque fois que tu accèderas à l'URL courte via ton navigateur).

#### 4.4. Tester l'API de Santé (via curl)

Vérifie si ton serveur est bien opérationnel :

1. Exécute la commande curl :

```
curl http://localhost:8080/health
```

Tu devrais obtenir :

```
{"status":"ok"}
```

#### 4.5. Observer le Moniteur d'URLs

Le moniteur fonctionne en arrière-plan et vérifie la disponibilité des URLs longues toutes les 5 minutes (par défaut).

Observe les logs dans le terminal où run-server tourne. Si l'état d'une URL que tu as raccourcie change (par exemple, si le site devient inaccessible), tu verras un message [NOTIFICATION] similaire à :

```
[NOTIFICATION] Le lien XYZ123 ([https://url-hors-ligne.com](https://url-hors-ligne.com)) est passé de ACCESSIBLE à INACCESSIBLE !
```

(Pour tester cela, tu pourrais raccourcir une URL vers un site que tu sais hors ligne ou une adresse IP inexistante, et attendre l'intervalle de surveillance.)

### 5. Arrêter le Serveur

Quand tu as terminé tes tests et que tu souhaites arrêter le service :

1. Dans le terminal où ./url-shortener run-server tourne, appuie sur :

```
Ctrl + C
```

Tu verras des logs confirmant l'arrêt propre du serveur.

## Barème de Notation (/20)

### 1. Robustesse Technique & Fonctionnelle (12 points)

-   1 point : Le projet se lance via ./url-shortener run-server.
-   4 points : Implémentation correcte de la redirection non-bloquante (GET /{shortCode}) avec utilisation efficace des goroutines et channels pour les analytics.
-   2 points : Le moniteur d'URLs fonctionne correctement, vérifie les URLs périodiquement et génère des notifications logiques.
-   3 points : Toutes les APIs REST et commandes CLI obligatoires (create, stats, migrate) sont fonctionnelles et robustes.
-   2 points : Gestion des erreurs pertinentes.

### 2. Qualité du Code & Documentation (2 points)

-   2 points : Code propre, lisible, **bien commenté** et code respectant les conventions Go vu en cours, et README pertinent.
-   2 points : Organisation des commits Git avec des messages clairs et pertinents.

### 3. Entretien Technique (4 points)

-   2 points : En Groupe : Votre capacité à expliquer et à défendre votre code lors d'un entretien individuel/en groupe. Cela inclut la compréhension de l'architecture, l'explication du fonctionnement asynchrone (workers, moniteur), et votre capacité à répondre aux questions techniques sur votre code. Vous devrez être capables de naviguer dans votre projet et de justifier vos choix.
-   2 points : Questions individuelles

### 4. Points faciles

-   1 point si votre code compile
-   1 point si vous faites des erreurs personnalisées
