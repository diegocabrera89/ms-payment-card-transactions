package utils

import (
	"github.com/diegocabrera89/ms-payment-card-transactions/constantsmicro"
	"github.com/diegocabrera89/ms-payment-card-transactions/models"
	"github.com/google/uuid"
	"time"
)

// BuildCreateToken build processor object.
func BuildCreateToken(token *models.Token) {
	token.TokenID = uuid.New().String()       // Generate a unique ID for the client
	token.CreatedAt = time.Now().UTC().Unix() //Date in UTC
	token.Status = constantsmicro.StatusTokenActive
}
