#!/usr/bin/env sh
az login --service-principal --username "${AZURE_CLIENT_ID}" --password "${AZURE_CLIENT_SECRET}" --tenant "${AZURE_TENANT_ID}"
secretValue=$(az keyvault secret show --vault-name gola --name "DB-PASSWORD" --query "value" --output tsv)
if [ $? -eq 0 ]; then
    export DB_PASSWORD="${secretValue}"
    migrate -path /database -database postgres://"${DB_USER}":"${DB_PASSWORD}"@"${DB_HOST}":"${DB_PORT}"/"${DB_NAME}"?sslmode=disable -verbose up
else
    echo "Unable to retrieve password"
fi
