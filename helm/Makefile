NAMESPACE ?= default
RELEASE ?= erida
VALUES ?= ./helm/values.yaml
STRESS_DEPLOY ?= false

validate:
	@helm template \
		--namespace $(NAMESPACE) \
		--values $(VALUES) \
		$(RELEASE) ./helm | kubeval --strict --ignore-missing-schemas

view:
	@helm template \
		--namespace $(NAMESPACE) \
		--values $(VALUES) \
		--set erida.env.SMTP_USER="$(SMTP_USER)" \
		--set erida.env.SMTP_PASS="$(SMTP_PASS)" \
		--set erida.env.SLACK_TOKEN="$(SLACK_TOKEN)" \
		--set stress.enabled=$(STRESS_DEPLOY) \
		$(RELEASE) ./helm

up:
	@helm upgrade --install \
		--namespace $(NAMESPACE) \
		--values ./helm/values.yaml \
		--set erida.env.SMTP_USER="$(SMTP_USER)" \
		--set erida.env.SMTP_PASS="$(SMTP_PASS)" \
		--set erida.env.SLACK_TOKEN="$(SLACK_TOKEN)" \
		--set stress.enabled=$(STRESS_DEPLOY) \
		$(RELEASE) ./helm

down:
	@helm uninstall \
		--namespace $(NAMESPACE) \
		--wait \
		--timeout 5m \
		$(RELEASE)