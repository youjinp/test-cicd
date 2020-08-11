.PHONY: build

export DEPLOYMENT_BUCKET_NAME = test-cicd-20200811-deployment
export STACK_NAME = test-cicd-20200811
export AWS_REGION = ap-southeast-2
export AWS_BRANCH = develop

deploy:
	$(info [*] Packaging and deploying services...)
	$(MAKE) deploy.createDeploymentBucket
	$(MAKE) deploy.api

deploy.createDeploymentBucket:
	$(info [*] Checking if deployment bucket exists...)
	@if aws s3 ls "s3://$${DEPLOYMENT_BUCKET_NAME}-$${AWS_REGION}" 2>&1 | grep -q 'NoSuchBucket'; then \
		echo Creating a bucket; \
		aws s3 mb s3://$${DEPLOYMENT_BUCKET_NAME}-$${AWS_REGION} \
			--region $${AWS_REGION}; \
		aws s3api put-public-access-block \
			--bucket $${DEPLOYMENT_BUCKET_NAME}-$${AWS_REGION} \
			--public-access-block-configuration BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true; \
	else \
		echo Bucket exists; \
	fi

deploy.api:
	sam package \
		--s3-bucket $${DEPLOYMENT_BUCKET_NAME}-$${AWS_REGION} \
		--region $${AWS_REGION} \
		--output-template-file packaged.yaml && \

	sam deploy \
		--s3-bucket $${DEPLOYMENT_BUCKET_NAME}-$${AWS_REGION} \
		--template-file packaged.yaml \
		--stack-name $${STACK_NAME}-api-$${AWS_BRANCH} \
		--capabilities CAPABILITY_IAM CAPABILITY_NAMED_IAM \
		--region $${AWS_REGION} \
		--no-fail-on-empty-changeset
