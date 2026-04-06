# SonarQube Local

## Demarrer SonarQube

Depuis la racine du projet :

```bash
docker compose -f infra/docker/docker-compose.sonar.yml up -d
```

Interface locale :

```text
http://localhost:9000
```

Identifiants par defaut au premier login :

```text
admin / admin
```

## Generer un token

Dans SonarQube :

1. Se connecter
2. Aller dans `My Account`
3. `Security`
4. Generer un token

## Lancer une analyse locale

```bash
export SONAR_TOKEN="<ton_token>"
bash scripts/sonar-local.sh
```

Options :

```bash
export SONAR_HOST_URL="http://localhost:9000"
export SONAR_PROJECT_KEY="storm"
```

## Arreter SonarQube

```bash
docker compose -f infra/docker/docker-compose.sonar.yml down
```

Pour supprimer aussi les donnees :

```bash
docker compose -f infra/docker/docker-compose.sonar.yml down -v
```
