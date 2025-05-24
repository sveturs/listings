build_n_deploy_front:
	gh workflow run "Frontend Build & Deploy" --ref main --repo sveturs/svetu
	sleep 10
	gh run list --workflow=frontend-build-n-deploy.yml --repo sveturs/svetu

build_n_deploy_back:
	gh workflow run "Backend Build & Deploy" --ref main --repo sveturs/svetu
	sleep 10
	gh run list --workflow=backend-build-n-deploy.yml --repo sveturs/svetu

