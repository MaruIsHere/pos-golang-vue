#!/bin/bash
echo "Starting Concurrency Test: 10 Cashiers firing concurrently..."

# First, get a token. Assuming 'admin' 'admin123' works.
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}' | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo "Failed to get token!"
  exit 1
fi

echo "Got token, starting 10 background cashiers..."

for i in {1..10}; do
  (
    for j in {1..5}; do
      curl -s -X POST http://localhost:8080/api/pos/orders \
        -H "Authorization: Bearer $TOKEN" \
        -H "Content-Type: application/json" \
        -d '{
          "customer_name": "Cashier '$i' Order '$j'",
          "payment_method": "cash",
          "paid_amount": 50000,
          "items": [
            { "product_id": "c6233b66-831e-4509-9f44-884ecba3108c", "quantity": 1, "subtotal": 15000 }
          ]
        }' > /dev/null
    done
    echo "Cashier $i finished."
  ) &
done

wait
echo "All cashiers finished. Checking total orders."
curl -s -X GET http://localhost:8080/api/reports/dashboard \
  -H "Authorization: Bearer $TOKEN" | grep -o '"total_orders":[^,]*'

