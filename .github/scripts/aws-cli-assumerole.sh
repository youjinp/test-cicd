unset  AWS_SESSION_TOKEN

temp_role=$(aws sts assume-role \
                    --role-arn "${TEST_ROLE}" \
                    --role-session-name "vgaltes-prod" \
                    --profile vgaltes-serverless)

AWS_ACCESS_KEY_ID=$(echo $temp_role | jq .Credentials.AccessKeyId | xargs)
AWS_SECRET_ACCESS_KEY=$(echo $temp_role | jq .Credentials.SecretAccessKey | xargs)
AWS_SESSION_TOKEN=$(echo $temp_role | jq .Credentials.SessionToken | xargs)

echo "::set-env name=VERSION::19.2.5"

echo "::set-env name=AWS_ACCESS_KEY_ID::${AWS_ACCESS_KEY_ID}"
echo "::set-env name=AWS_SECRET_ACCESS_KEY::${AWS_SECRET_ACCESS_KEY}"
echo "::set-env name=AWS_SESSION_TOKEN::${AWS_SESSION_TOKEN}"