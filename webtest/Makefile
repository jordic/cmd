

TVERSION = r2
GKE_PROJECT=balmy-cab-106909

compile:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

build:
	docker build  -t transtic_tester .
	docker tag -f transtic_tester eu.gcr.io/$(GKE_PROJECT)/transtic-tester:$(TVERSION)
	gcloud docker push eu.gcr.io/$(GKE_PROJECT)/transtic-tester:$(TVERSION)
	

