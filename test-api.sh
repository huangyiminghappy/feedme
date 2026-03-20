#!/bin/bash

echo "🧪 McDonald's Bot System - API Test Script"
echo "=========================================="
echo ""

echo "1️⃣  Resetting system..."
curl -s -X POST http://localhost:8080/api/reset | python3 -m json.tool
echo ""
sleep 1

echo "2️⃣  Creating orders (Normal -> VIP -> Normal)..."
curl -s -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{"id":"order1","type":"Normal"}' | python3 -m json.tool
sleep 0.1
curl -s -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{"id":"vip1","type":"VIP"}' | python3 -m json.tool
sleep 0.1
curl -s -X POST http://localhost:8080/api/orders -H "Content-Type: application/json" -d '{"id":"order2","type":"Normal"}' | python3 -m json.tool
echo ""
sleep 0.5

echo "3️⃣  Checking pending orders (VIP should be first)..."
curl -s http://localhost:8080/api/orders/pending | python3 -m json.tool
echo ""

echo "4️⃣  Adding a bot..."
curl -s -X POST http://localhost:8080/api/bots -H "Content-Type: application/json" -d '{"id":"bot1"}' | python3 -m json.tool
echo ""
sleep 0.5

echo "5️⃣  Checking which order is being processed (should be VIP)..."
curl -s http://localhost:8080/api/bots | python3 -m json.tool
echo ""

echo "6️⃣  System statistics..."
curl -s http://localhost:8080/api/stats | python3 -m json.tool
echo ""

echo "✅ Test completed! Check the results above."
echo ""
echo "💡 Tips:"
echo "   - VIP order should be processed first"
echo "   - Open http://localhost:5173/ to see the UI"
echo "   - Backend is running on http://localhost:8080"
