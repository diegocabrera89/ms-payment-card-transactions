compile = env GOOS=linux  GOARCH=arm64  go build -v -ldflags '-s -w -v' -tags lambda.norpc -o
zipper = zip -j -r
test_to_file = go test -coverprofile=coverage.out
percent = go tool cover -func=coverage.out | sed 's~\([^/]\{1,\}/\)\{3\}~~' | sed '$d' | sort -g -r -k 3
profile = pet

build: clean gomodgen import
	$(compile) bin/cmd/createTokenHandler/bootstrap cmd/createTokenHandler/create_token_handler.go
	$(compile) bin/cmd/chargeCardHandler/bootstrap cmd/chargeCardHandler/charge_card_handler.go
	$(compile) bin/cmd/getTokenHandler/bootstrap cmd/getTokenHandler/get_token_handler.go
	$(compile) bin/cmd/refundCardHandler/bootstrap cmd/refundCardHandler/refund_card_handler.go
	$(compile) bin/cmd/getCardTransactionHandler/bootstrap cmd/getCardTransactionHandler/get_card_transaction_handler.go

zip:
	$(zipper) bin/cmd/createTokenHandler/createTokenHandler.zip bin/cmd/createTokenHandler/bootstrap
	$(zipper) bin/cmd/chargeCardHandler/chargeCardHandler.zip bin/cmd/chargeCardHandler/bootstrap
	$(zipper) bin/cmd/getTokenHandler/getTokenHandler.zip bin/cmd/getTokenHandler/bootstrap
	$(zipper) bin/cmd/refundCardHandler/refundCardHandler.zip bin/cmd/refundCardHandler/bootstrap
	$(zipper) bin/cmd/getCardTransactionHandler/getCardTransactionHandler.zip bin/cmd/getCardTransactionHandler/bootstrap

clean:
	go clean
	rm -rf ./bin ./vendor go.sum

deploy: build zip
	sls deploy --aws-profile $(profile) --verbose

undeploy:
	sls remove --aws-profile $(profile) --verbose

import:
	go mod download github.com/aws/aws-lambda-go
	go mod download github.com/diegocabrera89/ms-payment-core

	go get github.com/diegocabrera89/ms-payment-core/dynamodbcore
	go get github.com/diegocabrera89/ms-payment-core/helpers
	go get github.com/diegocabrera89/ms-payment-card-transactions/utils

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
