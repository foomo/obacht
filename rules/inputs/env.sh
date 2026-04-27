#!/bin/sh
# Scan environment variables for suspicious patterns.
# NEVER outputs values, only variable names and matched patterns.

exact_matches="AWS_SECRET_ACCESS_KEY AWS_ACCESS_KEY_ID GITHUB_TOKEN GITLAB_TOKEN NPM_TOKEN DOCKER_PASSWORD SLACK_TOKEN SLACK_WEBHOOK_URL DATABASE_URL MYSQL_PASSWORD POSTGRES_PASSWORD REDIS_URL MONGODB_URI MONGO_URL AMQP_URL RABBITMQ_URL CELERY_BROKER_URL GCP_SA_KEY GITHUB_PAT DD_APP_KEY VAULT_DEV_ROOT_TOKEN_ID SUPABASE_SERVICE_ROLE_KEY"
suffix_patterns="_SECRET_KEY _ACCESS_KEY _LICENSE_KEY _PRIVATE_KEY _API_KEY _PASSWORD _CREDENTIALS _CREDENTIAL _ACCESS_TOKEN _AUTH_TOKEN _TOKEN _SECRET _AUTH _DSN"

suspicious="["
first=true

env | while IFS='=' read -r name rest; do
  matched=""
  pattern=""

  # Check exact matches.
  for exact in $exact_matches; do
    if [ "$name" = "$exact" ]; then
      matched="$name"
      pattern="exact:$exact"
      break
    fi
  done

  # Check suffix patterns.
  if [ -z "$matched" ]; then
    for suffix in $suffix_patterns; do
      case "$name" in
        *"$suffix")
          matched="$name"
          pattern="*$suffix"
          break
          ;;
      esac
    done
  fi

  if [ -n "$matched" ]; then
    if [ "$first" = true ]; then
      first=false
    else
      printf ','
    fi
    printf '{"name":"%s","pattern":"%s"}' "$matched" "$pattern"
  fi
done | {
  items=$(cat)
  printf '{"suspicious_vars": [%s]}' "$items"
}
