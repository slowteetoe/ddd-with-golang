package main

import (
	"context"
	"log"

	"github.com/Rhymond/go-money"
	"github.com/google/uuid"

	coffeeco "slowteetoe.com/coffeeco/internal"
	"slowteetoe.com/coffeeco/internal/payment"
	"slowteetoe.com/coffeeco/internal/purchase"
	"slowteetoe.com/coffeeco/internal/store"
)

func main() {

	ctx := context.Background()

	// This is the test key from Stripe's documentation. Feel free to use it, no charges will actually be made.
	stripeTestAPIKey := "sk_test_4eC39HqLyjWDarjtT1zdp7dc"

	// This is a test token from Stripe's documentation. Feel free to use it, no charges will actually be made.
	cardToken := "tok_visa"

	// This is the credentials for mongo if you run docker-compose up in this repo.
	mongoConString := "mongodb://root:example@localhost:27017"
	csvc, err := payment.NewStripeService(stripeTestAPIKey)
	if err != nil {
		log.Fatal(err)
	}

	purchaseRepo, err := purchase.NewMongoRepo(ctx, mongoConString)
	if err != nil {
		log.Fatal(err)
	}
	if err := purchaseRepo.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	storeRepo, err := store.NewMongoRepo(ctx, mongoConString)
	if err != nil {
		log.Fatal(err)
	}
	if err := storeRepo.Ping(ctx); err != nil {
		log.Fatal(err)
	}

	sSvc := store.NewService(storeRepo)

	svc := purchase.NewService(csvc, purchaseRepo, sSvc)

	someStoreID := uuid.New()

	pur := &purchase.Purchase{
		CardToken: &cardToken,
		Store: store.Store{
			ID: someStoreID,
		},
		ProductsToPurchase: []coffeeco.Product{{
			ItemName:  "item1",
			BasePrice: *money.New(3300, "USD"),
		}},
		PaymentMeans: payment.MEANS_CARD,
	}
	if err := svc.CompletePurchase(ctx, someStoreID, pur, nil); err != nil {
		log.Fatal(err)
	}

	log.Println("purchase was successful")
}
