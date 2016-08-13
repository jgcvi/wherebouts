#! /bin/bash
GROUP=frontend-group
TEMPLATE=$GROUP-tmpl
IMAGE=debian-8
MACHINE_TYPE=f1-micro
MIN_INSTANCES=1
MAX_INSTANCES=1
TARGET_UTILIZATION=0.6
STARTUP_SCRIPT=startup-script.sh
STARTUP_SCRIPT=startup-script.sh
SCOPES="userinfo-email,cloud-platform"
ZONE=us-central1-a
SERVICE=frontend-web-service

# [START create_template]
gcloud compute instance-templates create $TEMPLATE \
  --image $IMAGE \
  --machine-type $MACHINE_TYPE \
  --scopes $SCOPES \
  #--metadata-from-file startup-script=$STARTUP_SCRIPT \
  # --metadata app-location=$WHEREBOUTS_DEPLOY_LOCATION \
  #--tags $TAGS
# [END create_template]

# [START create_group]
gcloud compute instance-groups managed \
  create $GROUP \
  --base-instance-name $GROUP \
  --size $MIN_INSTANCES \
  --template $TEMPLATE \
  --zone $ZONE
# [END create_group]

# [START create_named_port
gcloud compute instance-groups managed set-named-ports \
  $GROUP \
  --named-port http:8080 \
  --zone $ZONE
# [END create_named_port]

gcloud compute http-health-checks create ah-health-check \
  --request-path /_ah/health \
  --port 8080

gcloud compute backend-services create $SERVICE \
  --http-health-check ah-health-check

gcloud compute backend-services add-backend $SERVICE \
  --instance-group $GROUP \
  --zone $ZONE

gcloud compute url-maps create $SERVICE-map \
  --default-service $SERVICE

gcloud compute target-http-proxies create $SERVICE-proxy \
  --url-map $SERVICE-map

gcloud compute forwarding-rules create $SERVICE-http-rule \
  --global \
  --target-http-proxy $SERVICE-proxy \
  --port-range 80

gcloud compute instance-groups managed set-autoscaling \
  $GROUP \
  --max-num-replicas $MAX_INSTANCES \
  --target-load-balancing-utilization $TARGET_UTILIZATION \
  --zone $ZONE

gcloud compute firewall-rules create default-allow-http-8080 \
  --allow tcp:8080 \
  --source-ranges 0.0.0.0/0 \
  --target-tags http-server \
  --description "Allow port 8080 access to http-server"
