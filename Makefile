# if the KEY environment variable is not set to either stage or prod, makefile will fail
# KEY is confirmed below in the check_env directive
# example:
# for stage run: ENV=stage make
# for production run: ENV=prod make
include .env

# found yolo at: https://azer.bike/journal/a-good-makefile-for-go/

AWS_STACK_NAME ?= $(PROJECT_NAME)

deploy: check_env build awspackage awsdeploy

check_env:
	@echo -n "Your environment is $(ENV)? [y/N] " && read ans && [ $${ans:-N} = y ]

clean:
	@rm -rf dist
	@mkdir -p dist

build: clean
	@for dir in `ls handler`; do \
		GOOS=linux go build -o dist/$$dir github.com/pulpfree/gsales-pdf-reports/handler/$$dir; \
	done
	@cp ./config/defaults.yml dist/
	@cp -r ./image dist/
	@echo "build successful"

# "create-bucket-configuration" line only if not in us-east-1, or apparently some other regions as well...
configure:
	@ aws s3api create-bucket \
		--bucket $(AWS_BUCKET_NAME) \
		--region $(AWS_REGION) \
		--create-bucket-configuration LocationConstraint=$(AWS_REGION)

# watch: Run given command when code changes. e.g; make watch run="echo 'hey'"
# @yolo -i . -e vendor -e bin -e dist -c $(run)

watch:
	@yolo -i . -e vendor -e dist -c "make build"

run: build
	sam local start-api -n env.json

validate:
	sam validate

awspackage:
	@aws cloudformation package \
   --template-file ${FILE_TEMPLATE} \
   --output-template-file ${FILE_PACKAGE} \
   --s3-bucket $(AWS_LAMBDA_BUCKET) \
   --s3-prefix $(AWS_BUCKET_PREFIX) \
   --profile $(AWS_PROFILE) \
   --region $(AWS_REGION)

awsdeploy:
	@aws cloudformation deploy \
	--template-file ${FILE_PACKAGE} \
	--region $(AWS_REGION) \
	--stack-name $(AWS_STACK_NAME) \
	--capabilities CAPABILITY_IAM \
	--profile $(AWS_PROFILE) \
	--force-upload \
	--parameter-overrides \
		ParamCertificateArn=$(CERTIFICATE_ARN) \
		ParamCustomDomainName=$(CUSTOM_DOMAIN_NAME) \
		ParamENV=$(ENV) \
		ParamHostedZoneId=$(HOSTED_ZONE_ID) \
		ParamKMSKeyID=$(KMS_KEY_ID) \
		ParamProjectName=$(PROJECT_NAME) \
		ParamReportBucket=${AWS_REPORT_BUCKET} \
		ParamSSMPath=$(SSM_PARAM_PATH) \
		ParamSecurityGroupIds=$(SECURITY_GROUP_IDS) \
		ParamSubnetIds=$(SUBNET_IDS) \
		ParamUserPoolArn=$(USER_POOL_ARN)

describe:
	@aws cloudformation describe-stacks \
		--region $(AWS_REGION) \
		--stack-name $(AWS_STACK_NAME)

outputs:
	@ make describe \
		| jq -r '.Stacks[0].Outputs'
