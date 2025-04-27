# Guide d'utilisation de KVLT

**KVLT** est une base de données clé-valeur en mémoire légère et rapide, inspirée de Redis. Voici comment l'utiliser efficacement.

## Installation

- Via Docker (recommandé)

```bash
docker run -d -p 8080:8080 --name kvlt z4riix/kvlt
```

Docker compose:

```yaml
services:
  kvlt:
    image: z4riix/kvlt:latest
    container_name: kvlt
    ports:
      - "8080:8080"
    environment:
      - CLEANER_TIME=@every 5s
    volumes:
      - kvlt-data:/root/data

volumes:
  kvlt-data:
    driver: local
```

Puis lancez avec :

```bash
docker-compose up -d
```

- Via Go (non recommandé)

**KVLT** peut également être installé localement en clonant le dépôt et en le construisant :

```bash
git clone git@github.com:ZeRiix/kvlt.git
cd kvlt
go build -o kvlt
./kvlt
```

## Configuration

**KVLT** peut être configuré via des variables d'environnement :

| Variable     | Description	                             | Valeur par défaut |
|--------------|---------------------------------------------|-------------------|
| PORT         | Port d'écoute du serveur                    | 8080              |
| HOST         | Hôte d'écoute                               | 0.0.0.0           |
| CLEANER_TIME | Fréquence de nettoyage des données expirées | @every 3s         |
| DB_PATH      | Chemin de stockage des données persistantes | store.json        |

## Utilisation de l'API REST

**Définir une valeur**

Pour stocker une clé avec une valeur et éventuellement une durée d'expiration :

```bash
curl -X PUT "http://localhost:8080/value" \
  -H "Content-Type: application/json" \
  -d '{
      "key": "user:1", 
      "value": {"name": "John", "email": "john@example.com"}, 
      "duration": 3600
    }'
```

**Récupérer une valeur**

Pour récupérer la valeur associée à une clé :

```bash
curl "http://localhost:8080/value/user:1"
```

Exemple de réponse de succès :

```json
{
  "data": {
    "key": "user:1",
    "value": {
      "name": "John",
      "email": "john@example.com"
    }
  }
}
```
