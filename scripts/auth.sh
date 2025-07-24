#!/bin/sh

gcloud auth application-default login \
    --account "${LOGIN_ACCOUNT}" \
    --project gcas-form-lg-dev \
    --scopes='https://www.googleapis.com/auth/admin.directory.user,https://www.googleapis.com/auth/iam.test,https://www.googleapis.com/auth/cloud-platform'
