package database

import (
	"food-recipe/models"
	"fmt"
	"context"
	"github.com/machinebox/graphql"
	"time"
)

func SavePurchase(purchase models.Purchase) error {
    req := graphql.NewRequest(`
        mutation($user_id: String!, $amount: String!, $status: String!, $recipe_id: String, $tx_ref: String!, $created_at: timestamptz!) {
            insert_purchases_one(object: {
                user_id: $user_id,
                amount: $amount,
                status: $status,
                recipe_id: $recipe_id,
                tx_ref: $tx_ref,
                created_at: $created_at
            }) {
                id
            }
        }
    `)

    req.Var("user_id", purchase.UserID)
    req.Var("amount", purchase.Amount)
    req.Var("status", purchase.Status)
    req.Var("recipe_id", purchase.RecipeID)
    req.Var("tx_ref", purchase.TxRef)
    req.Var("created_at", time.Now().UTC())

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := client.Run(ctx, req, nil); err != nil {
        return fmt.Errorf("failed to save purchase: %w", err)
    }

    return nil
}
