#!/bin/bash
# API Test Examples - Identity & Organization Service
# Tập hợp các ví dụ curl để test các API endpoints

API_URL="http://localhost:3001/api/v1"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Identity & Organization Service - API Examples ===${NC}\n"

# ==================== COMPANIES ====================

echo -e "${GREEN}1. CREATE COMPANY${NC}"
echo "POST $API_URL/companies"
curl -X POST $API_URL/companies \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Tech",
    "address": "123 Nguyen Trai, HCMC",
    "phone": "+84-28-0000-0001",
    "email": "contact@acmetech.com",
    "website": "https://acmetech.example",
    "status": 1,
    "subscription_plan": 1,
    "subscription_start_date": "2025-01-01",
    "subscription_end_date": "2025-12-31",
    "max_employees": 500,
    "max_devices": 50
  }' | jq .
echo -e "\n---\n"

# Get COMPANY_ID from previous response (you need to copy it)
# For demo purposes, using a placeholder
COMPANY_ID="00000000-0000-0000-0000-000000000001"

echo -e "${GREEN}2. GET ALL COMPANIES${NC}"
echo "GET $API_URL/companies"
curl -X GET $API_URL/companies | jq .
echo -e "\n---\n"

echo -e "${GREEN}3. GET COMPANY BY ID${NC}"
echo "GET $API_URL/companies/{company_id}"
curl -X GET $API_URL/companies/$COMPANY_ID | jq .
echo -e "\n---\n"

echo -e "${GREEN}4. UPDATE COMPANY${NC}"
echo "PUT $API_URL/companies/{company_id}"
curl -X PUT $API_URL/companies/$COMPANY_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Acme Tech Updated",
    "status": 1,
    "max_employees": 1000
  }' | jq .
echo -e "\n---\n"

# ==================== USERS ====================

echo -e "${GREEN}5. CREATE USER (Employee)${NC}"
echo "POST $API_URL/users"
curl -X POST $API_URL/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "alice.acme@example.com",
    "phone": "0900001001",
    "password": "SecurePass123",
    "full_name": "Alice Nguyen",
    "role": 2,
    "company_id": "'$COMPANY_ID'",
    "employee_code": "EMP-001",
    "department": "Engineering",
    "position": "Senior Engineer",
    "hire_date": "2025-01-15",
    "salary": 50000
  }' | jq .
echo -e "\n---\n"

# Get USER_ID from previous response
USER_ID="00000000-0000-0000-0000-000000000011"

echo -e "${GREEN}6. GET ALL USERS${NC}"
echo "GET $API_URL/users"
curl -X GET $API_URL/users | jq .
echo -e "\n---\n"

echo -e "${GREEN}7. GET USERS BY COMPANY${NC}"
echo "GET $API_URL/users?company_id={company_id}"
curl -X GET "$API_URL/users?company_id=$COMPANY_ID" | jq .
echo -e "\n---\n"

echo -e "${GREEN}8. GET USER BY ID${NC}"
echo "GET $API_URL/users/{user_id}"
curl -X GET $API_URL/users/$USER_ID | jq .
echo -e "\n---\n"

echo -e "${GREEN}9. UPDATE USER${NC}"
echo "PUT $API_URL/users/{user_id}"
curl -X PUT $API_URL/users/$USER_ID \
  -H "Content-Type: application/json" \
  -d '{
    "full_name": "Alice Nguyen Updated",
    "position": "Lead Engineer",
    "salary": 60000
  }' | jq .
echo -e "\n---\n"

# ==================== FACE DATA ====================

echo -e "${GREEN}10. UPLOAD FACE DATA${NC}"
echo "POST $API_URL/users/{user_id}/face-data"
curl -X POST $API_URL/users/$USER_ID/face-data \
  -H "Content-Type: application/json" \
  -d '{
    "image_url": "https://example.com/face-001.jpg",
    "quality_score": 0.95
  }' | jq .
echo -e "\n---\n"

echo -e "${GREEN}11. GET FACE DATA LIST${NC}"
echo "GET $API_URL/users/{user_id}/face-data"
curl -X GET $API_URL/users/$USER_ID/face-data | jq .
echo -e "\n---\n"

# Get FID from previous response
FID="00000000-0000-0000-0000-000000000021"

echo -e "${GREEN}12. DELETE FACE DATA${NC}"
echo "DELETE $API_URL/users/{user_id}/face-data/{fid}"
curl -X DELETE $API_URL/users/$USER_ID/face-data/$FID | jq .
echo -e "\n---\n"

echo -e "${GREEN}13. DELETE USER${NC}"
echo "DELETE $API_URL/users/{user_id}"
curl -X DELETE $API_URL/users/$USER_ID | jq .
echo -e "\n---\n"

echo -e "${GREEN}14. DELETE COMPANY${NC}"
echo "DELETE $API_URL/companies/{company_id}"
curl -X DELETE $API_URL/companies/$COMPANY_ID | jq .
echo -e "\n---\n"

echo -e "${BLUE}=== All tests completed ===${NC}\n"

# Note: Make sure jq is installed: sudo apt-get install jq
