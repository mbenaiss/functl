build:
	cd cli && \
	statik -src=../api && \
	go install 
