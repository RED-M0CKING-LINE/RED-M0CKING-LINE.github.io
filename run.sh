podman compose -f ./deploy/compose/prod.compose.yaml down
git pull
podman compose -f ./deploy/compose/prod.compose.yaml up -d --build