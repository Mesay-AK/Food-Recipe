package models


type PaymentRequest struct {
    Amount      string `json:"amount"`
    Currency    string `json:"currency"`
    Email       string `json:"email"`
    FirstName   string `json:"first_name"`
    LastName    string `json:"last_name"`
    TxRef       string `json:"tx_ref"`
    CallbackURL string `json:"callback_url"`
    ReturnURL   string `json:"return_url"`
}


type PaymentResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
    Data    struct {
        CheckoutURL string `json:"checkout_url"`
    } `json:"data"`
}


type VerifyResponse struct {
    Status string `json:"status"`
    Data   struct {
        TxRef  string `json:"tx_ref"`
        Status string `json:"status"`
    } `json:"data"`
}


type Purchase struct {
    ID        string  `json:"id"`
    UserID    string  `json:"user_id"`   
    Amount    string `json:"amount"`    
    Status    string `json:"status"`    
    RecipeID  int64  `json:"recipe_id"` 
    TxRef     string `json:"tx_ref"`    
    CreatedAt string `json:"created_at"`
}
