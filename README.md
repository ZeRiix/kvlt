<p align="center">
  <img src="./docs/assets/logo.png" alt="logo" width="300" height="auto" />
</p>
<p align="center">
  <span style="font-size: 24px; font-weight: bold;">Kvlt Database</span>
</p>
<p align="center">
  <a href="#">
    <img src='https://img.shields.io/badge/lang-Go-00ADD8?logo=go&style=plastic' alt='Go' />
  </a>
  <a href="#">
    <img src="https://img.shields.io/badge/coverage-0%25-red?style=plastic" alt="lang">
  </a>
  <a href="https://hub.docker.com/r/z4riix/kvlt">
    <img src="https://img.shields.io/docker/pulls/z4riix/kvlt?style=plastic&logo=docker" alt="Docker Hub">
  </a>
</p>

## Présentation

**KVLT** est une base de données clé-valeur en mémoire inspirée de Redis, conçue pour offrir un stockage de données simple, rapide et performant. Elle propose une solution légère pour les applications nécessitant un accès rapide aux données sans la complexité des systèmes de base de données traditionnels.

## Caractéristiques principales

- **Stockage en mémoire** : Accès ultra-rapide aux données
- **Architecture simple** : API REST intuitive et facile à utiliser
- **Flexibilité** : Stockage de différents types de données (chaînes, nombres, objets JSON)
- **Développé en Go** : Performances optimales et faible empreinte mémoire
- **Conteneurisé** : Déploiement facile via Docker

## Cas d'utilisation

**KVLT** est particulièrement adapté pour :

- Cache d'application
- Gestion des sessions utilisateurs
- Stockage temporaire de résultats de requêtes fréquentes
- Partage de données entre différents services
- Environnements d'apprentissage et de test

## Installation

```bash
docker run -d -v kvlt-data:/root/data -p 8080:8080 ---name kvlt z4riix/kvlt
```

### Compose configuration

```yaml
services:
  kvlt:
    image: z4riix/kvlt
    container_name: kvlt
    ports:
      - "8080:8080"
    environment:
      - CLEANER_TIME="@every 3s" # default 3s
    volumes:
      - kvlt-data:/root/data
volumes:
  kvlt-data:
    driver: local
```

## Ressources

- [Documentation du client KVLT](./client/README.md)
- [Documentation KVLT](./docs/README.md)
- [Cahier des charges](./docs/specifications.md)

## :warning: Avertissement

**Ce projet est créé dans le cadre d'un cours de développement et n'est pas destiné à être utilisé en production**. Merci de ne pas l'utiliser ni de prendre exemple sur ce projet pour vos propres projets.
