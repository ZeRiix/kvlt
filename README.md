<p align="center">
  <img src="./logo.png" alt="logo" width="300" height="auto" />
</p>
<p align="center">
  <span style="font-size: 24px; font-weight: bold;">Kvlt Database</span>
</p>
<p align="center">
  <a href="#">
    <img src='https://img.shields.io/badge/lang-Go-00ADD8?logo=go&style=plastic' alt='Go' />
  </a>
</p>

## Présentation Générale

**KVLT** est une base de données clé-valeur en mémoire inspirée de Redis, développée en Go. Elle offre un système de stockage rapide avec plusieurs plugins pour étendre ses fonctionnalités :

- **Persistance** via AOF (Append-Only File)
- **Expiration automatique** des données
- **Indexation** pour des requêtes rapides
- **Architecture modulaire** avec système de hooks

---

## Architecture par Couches

```mermaid
graph TD
    A[Application Layer<br/>main.go] --> B[Plugin Layer]
    
    B --> B1[AOF<br/>Persistance<br/>sur disque]
    B --> B2[Expiration<br/>TTL auto]
    B --> B3[Indexation<br/>Recherche]
    
    B1 --> C[Hook System<br/>Système d'événements pour plugins]
    B2 --> C
    B3 --> C
    
    C --> C1[OnSet<br/>Hooks]
    C --> C2[OnGet<br/>Hooks]
    C --> C3[OnDrop<br/>Hooks]
    
    C1 --> D["Core Store<br/>Stockage en mémoire<br/>map string Item + Méthodes<br/>Get, Set, Drop"]
    C2 --> D
    C3 --> D
    
    style A fill:#e1f5fe
    style B fill:#f3e5f5
    style C fill:#fff3e0
    style D fill:#e8f5e8
```

### Principe des Couches

1. **Core Store** : Stockage de base en mémoire avec les opérations CRUD
2. **Hook System** : Système d'événements permettant aux plugins de réagir aux opérations
3. **Plugin Layer** : Plugins modulaires qui étendent les fonctionnalités
4. **Application Layer** : Point d'entrée de l'application

---

## Le Store Principal

### Structure de Données

```go
type Item struct {
    Value interface{}  // Valeur stockée (peut être de n'importe quel type)
    Key   string      // Clé unique d'identification
}

type Store struct {
    data        map[string]Item  // Stockage principal en mémoire
    actionHooks ActionHooks     // Hooks pour les plugins
}
```

### Opérations de Base

- **`Get(key string)`** : Récupère un élément par sa clé
- **`Set(item Item)`** : Stocke un nouvel élément
- **`Drop(key string)`** : Supprime un élément

### Flux d'Exécution

```mermaid
flowchart LR
    A[Set Item] --> B[Update Map] --> C[Trigger Hooks]
    C --> D[Plugin Reactions]
    
    D --> E[AOF<br/>Logging]
    D --> F[Index<br/>Update]
    D --> G[Expiration<br/>Check]
    
    style A fill:#ffebee
    style B fill:#e8f5e8
    style C fill:#fff3e0
    style D fill:#f3e5f5
    style E fill:#e1f5fe
    style F fill:#e1f5fe
    style G fill:#e1f5fe
```

---

## Système de Hooks

Le système de hooks permet aux plugins de s'abonner aux événements du store :

```go
type ActionHooks struct {
    get  []func(item *Item)  // Déclenchés lors d'un Get
    set  []func(item *Item)  // Déclenchés lors d'un Set
    drop []func(item *Item)  // Déclenchés lors d'un Drop
}
```

### Avantages
- **Modularité** : Chaque plugin peut s'enregistrer indépendamment
- **Asynchrone** : Les hooks sont exécutés en goroutines
- **Extensibilité** : Facile d'ajouter de nouveaux plugins

---

## Plugin AOF (Append-Only File)

### Principe

L'AOF assure la **persistance des données** en journalisant toutes les opérations sur disque.

### Architecture AOF

```mermaid
graph LR
    subgraph "Memory Store"
        A[Set/Drop<br/>Operations]
    end
    
    subgraph "Disk Storage"
        subgraph "AOF System"
            B[Buffer Rotation<br/>QuantityBuffer]
            B --> B1[Buffer0]
            B --> B2[Buffer1]
            B --> B3[Buffer...]
            
            B1 --> C[AOF Files<br/>./buffer/timestamp]
            B2 --> C
            B3 --> C
            
            C --> D[Snapshot Files<br/>./data/key_files]
        end
    end
    
    A --> B
    
    style A fill:#ffebee
    style B fill:#e8f5e8
    style C fill:#fff3e0
    style D fill:#e1f5fe
```

### Configuration

```go
type OptionsAOF struct {
    IntervalAnalyzeBuffer time.Duration  // Fréquence de vidage des buffers
    IntervalSnapshot      time.Duration  // Fréquence de création des snapshots
    QuantityBuffer        int           // Nombre de buffers rotatifs
    AOFFolderPath         string        // Dossier des fichiers AOF
    SnapshotFolderPath    string        // Dossier des snapshots
    SplitChar             string        // Caractère de séparation
}
```

### Processus de Fonctionnement

```mermaid
graph TD
    A[Opérations Set/Drop] --> B[Buffering<br/>Stockage dans buffers rotatifs]
    B --> C[Export AOF<br/>Vidage périodique vers fichiers]
    C --> D[Snapshot<br/>Consolidation des AOF]
    D --> E[Recovery<br/>Rechargement au démarrage]
    
    B -.-> F[IntervalAnalyzeBuffer<br/>1 seconde]
    C -.-> G[Fichiers timestamp<br/>./buffer/]
    D -.-> H[IntervalSnapshot<br/>10 secondes]
    E -.-> I[Fichiers clé<br/>./data/]
    
    style A fill:#ffebee
    style B fill:#e8f5e8
    style C fill:#fff3e0
    style D fill:#e1f5fe
    style E fill:#f3e5f5
```

1. **Buffering** : Les opérations sont stockées dans des buffers rotatifs
2. **Export AOF** : Périodiquement, les buffers sont vidés vers des fichiers AOF
3. **Snapshot** : Les fichiers AOF sont consolidés en snapshots
4. **Recovery** : Au démarrage, les snapshots sont rechargés en mémoire

### Format des Fichiers AOF

```
action|\\|\\|key|\\|\\|{"Key":"test","Value":{"name":"john"}}
set|\\|\\|user1|\\|\\|{"Key":"user1","Value":{"name":"Alice","age":30}}
drop|\\|\\|user2|\\|\\|{"Key":"user2","Value":null}
```

### Avantages
- **Durabilité** : Aucune perte de données en cas de crash
- **Performance** : Système de buffers pour éviter les I/O trop fréquentes
- **Récupération** : Restauration automatique au démarrage

---

## Plugin d'Expiration

### Principe

Le plugin d'expiration permet de définir une **durée de vie** aux données stockées.

### Architecture d'Expiration

```mermaid
graph TB
    subgraph "Expiration System"
        A["expirationStore<br/>map timestamp → items"]
        
        A --> B1[timestamp1 → item1, item2, ...]
        A --> B2[timestamp2 → item3, item4, ...]
        A --> B3[timestamp3 → item5, ...]
        
        C[Background Goroutine<br/>Vérification continue]
        
        C --> D[Boucle temporelle<br/>currentTime]
        D --> E{"items trouvés pour<br/>currentTime ?"}
        E -->|Yes| F[Supprimer tous<br/>les items expirés]
        E -->|No| D
        F --> D
    end
    
    style A fill:#e8f5e8
    style C fill:#fff3e0
    style D fill:#e1f5fe
    style E fill:#ffebee
    style F fill:#f3e5f5
```

### Fonctionnement

1. **Détection** : Lors d'un `Set`, recherche de la propriété `expireAt`
2. **Stockage** : Les items avec expiration sont indexés par timestamp
3. **Surveillance** : Une goroutine vérifie en continu les expirations
4. **Suppression** : Les items expirés sont automatiquement supprimés

### Format des Données avec Expiration

```go
store.Set(Item{
    Key: "session_token",
    Value: map[string]interface{}{
        "token": "abc123",
        "userId": 42,
        "expireAt": time.Now().Unix() + 3600, // Expire dans 1 heure
    },
})
```

### Avantages
- **Automatique** : Pas besoin de gérer manuellement les suppressions
- **Flexible** : Expiration par item individuel
- **Performant** : Vérification optimisée par timestamp

---

## Plugin d'Indexation

### Principe

Le plugin d'indexation permet de **rechercher rapidement** des données par leurs propriétés internes.

### Architecture d'Indexation

```mermaid
graph TB
    subgraph "Indexation System"
        subgraph "Flatten Process"
            A["Original:<br/>{<br/>  user: {<br/>    profile: {<br/>      name: 'John',<br/>      age: 30<br/>    }<br/>  }<br/>}"]
            A --> B["Flattened:<br/>{<br/>  'user.profile.name': 'John',<br/>  'user.profile.age': 30<br/>}"]
        end
        
        subgraph "Index Storage"
            C["indexes['user.profile.name'] = {<br/>  stringStore: {<br/>    'John': {'item1': *Item, 'item2': *Item}<br/>    'Alice': {'item3': *Item}<br/>  }<br/>}"]
            
            D["indexes['user.profile.age'] = {<br/>  intStore: {<br/>    30: {'item1': *Item, 'item4': *Item}<br/>    25: {'item5': *Item}<br/>  }<br/>}"]
        end
        
        B --> C
        B --> D
    end
    
    style A fill:#ffebee
    style B fill:#e8f5e8
    style C fill:#fff3e0
    style D fill:#e1f5fe
```

### Structure des Index

```go
type Index struct {
    intStore     Record[int64, RecordItem]   // Index pour les entiers
    stringStore  Record[string, RecordItem]  // Index pour les chaînes
    nullStore    RecordItem                  // Index pour les valeurs null
    booleanStore Record[bool, RecordItem]    // Index pour les booléens
}

type Indexes Record[string, Index] // Index par chemin de propriété
```

### Processus de Recherche

```mermaid
graph LR
    A["Query: finder avec<br/>user.profile.name = John"] --> B[Lookup dans les index<br/>user.profile.name]
    B --> C[Vérifier stringStore<br/>pour 'John']
    C --> D[Retourner tous<br/>les Items trouvés]
    
    style A fill:#ffebee
    style B fill:#e8f5e8
    style C fill:#fff3e0
    style D fill:#e1f5fe
```

### Exemple d'Utilisation

```go
// Stockage
store.Set(Item{
    Key: "user1",
    Value: map[string]interface{}{
        "profile": map[string]interface{}{
            "name": "John",
            "age": int64(30),
        },
        "active": true,
    },
})

// Recherche
finder := InitIndexes(store)
results := finder("profile.name", "John")      // Trouve user1
results = finder("profile.age", int64(30))     // Trouve user1
results = finder("active", true)               // Trouve user1
```

### Avantages
- **Performance** : Recherche O(1) au lieu de O(n)
- **Flexibilité** : Index automatique sur toutes les propriétés
- **Types multiples** : Support des strings, int64, bool, null

---

## Exemple d'Utilisation

Voici un exemple complet d'utilisation de KVLT avec tous ses plugins :

```go
func main() {
    // Configuration AOF
    optionsAOF := store.OptionsAOF{
        IntervalAnalyzeBuffer: 1 * time.Second,
        IntervalSnapshot:      10 * time.Second,
        QuantityBuffer:        10,
        AOFFolderPath:         "./buffer",
        SnapshotFolderPath:    "./data",
        SplitChar:             "|\\|\\|",
    }

    // Initialisation du store
    storeInstance := store.NewStore()

    // Activation des plugins
    store.InitExpiration(storeInstance)
    store.InitAOF(storeInstance, optionsAOF)
    finder := store.InitIndexes(storeInstance)

    // Stockage d'un objet complexe avec expiration
    storeInstance.Set(store.Item{
        Key: "user_session",
        Value: map[string]interface{}{
            "userId": int64(123),
            "profile": map[string]interface{}{
                "name": "John Doe",
                "role": "admin",
            },
            "expireAt": time.Now().Unix() + 3600, // Expire dans 1h
        },
    })

    // Recherche par propriété
    adminUsers := finder("profile.role", "admin")
    fmt.Printf("Utilisateurs admin: %v\n", adminUsers)

    // Les données sont automatiquement:
    // - Persistées sur disque (AOF)
    // - Indexées pour la recherche
    // - Expirées après 1 heure
}
```

### Flux Complet

```mermaid
graph LR
    A[Set Item] --> B[Store in<br/>Memory] --> C[Trigger Hooks]
    
    C --> D[AOF Plugin<br/>• Log to buffer<br/>• Schedule export<br/>• Create snapshots]
    C --> E[Expiration Plugin<br/>• Check expireAt<br/>• Schedule cleanup<br/>• Remove expired]
    C --> F[Index Plugin<br/>• Flatten object<br/>• Update indexes<br/>• Enable search]
    
    style A fill:#ffebee
    style B fill:#e8f5e8
    style C fill:#fff3e0
    style D fill:#e1f5fe
    style E fill:#f3e5f5
    style F fill:#e8f5e8
```

---

## Avantages de l'Architecture

1. **Modularité** : Chaque plugin est indépendant et peut être activé/désactivé
2. **Performance** : 
   - Stockage en mémoire pour la vitesse
   - Index pour les recherches rapides
   - Buffers AOF pour optimiser les I/O
3. **Robustesse** :
   - Persistance via AOF
   - Récupération automatique au démarrage
   - Gestion automatique des expirations
4. **Extensibilité** : Système de hooks permet d'ajouter facilement de nouveaux plugins

Cette architecture fait de KVLT une base de données clé-valeur complète et performante, adaptée aux besoins modernes de stockage temporaire et de cache.

