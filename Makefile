.PHONY: deploy-backend deploy-frontend

deploy-backend:
	gcloud run deploy colorboxd-backend --image europe-southwest1-docker.pkg.dev/colorboxd/colorboxd-repo/colorboxd-backend:latest --region us-central1

deploy-frontend:
	gcloud run deploy colorboxd --image europe-southwest1-docker.pkg.dev/colorboxd/colorboxd-repo/colorboxd:latest --region us-central1
