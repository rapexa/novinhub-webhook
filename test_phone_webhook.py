#!/usr/bin/env python3
"""
Test script to simulate webhook events with phone numbers
This script sends test webhook payloads to verify phone number extraction
"""

import requests
import json

# Webhook endpoint
WEBHOOK_URL = "http://localhost:8080/webhook"

# Test message with phone number (similar to the log you provided)
test_message_payload = {
    "type": "message_created",
    "user_id": 128659,
    "payload": {
        "id": "3939210205",
        "message_id": "test_message_id",
        "type": "text",
        "attachment": None,
        "text": "09155520952",  # Phone number from your log
        "date": "2025-09-18 01:31:06",
        "can_delete": 0,
        "can_reaction": "1",
        "is_auto_response": None,
        "social_user_id": 571859378,
        "from": None,
        "reactions": [],
        "account_id": 185315,
        "conversation_id": 3475816668,
        "account": {
            "id": 185315,
            "identifier": "17841457765637051",
            "name": "alirezaasll",
            "login_required": 0,
            "can_send_direct": 1,
            "can_send_comment": 1,
            "can_send_post": 1,
            "relogin_reason": None,
            "created_at": "2025-07-31T09:33:19.000000Z",
            "type": "Instagram",
            "profile_url": "https://instagram.com/alirezaasll",
            "info": None,
            "automation_status": True,
            "social_user_id": 173588463
        },
        "socialUser": {
            "id": 571859378,
            "user_id": "4082895595259389",
            "username": "mohammad_gharehbagh_",
            "name": "Inventor...",
            "image": "https://example.com/profile.jpg",
            "social": "Instagram",
            "updated_at": "2025-09-18T01:31:09.000000Z"
        }
    }
}

# Test message with multiple phone numbers
test_multiple_phones_payload = {
    "type": "message_created",
    "user_id": 128659,
    "payload": {
        "id": "test_multiple",
        "text": "سلام، شماره‌های من: 09155520952 و +989123456789 هستند",
        "account": {},
        "socialUser": {}
    }
}

# Test message with no phone number
test_no_phone_payload = {
    "type": "message_created",
    "user_id": 128659,
    "payload": {
        "id": "test_no_phone",
        "text": "سلام، چطوری؟ فالو دارم",
        "account": {},
        "socialUser": {}
    }
}

# Test lead created event (this is the main event we should focus on!)
# Note: user_id (128659) will be used as the 'code' variable in SMS pattern
test_lead_created_payload = {
    "type": "leed_created",
    "user_id": 128659,
    "payload": {
        "id": "26033286",
        "type": "number",
        "value": "09155520952",
        "message_id": "3939210734",
        "created_at": "2025-09-18T01:31:30.000000Z",
        "message": {
            "id": 3939210734,
            "message_id": "test_message_id",
            "type": "text",
            "attachment": None,
            "text": "09155520952",
            "date": "2025-09-18 01:31:29",
            "can_delete": 0,
            "can_reaction": 1,
            "is_auto_response": None,
            "social_user_id": 571859378,
            "from": None,
            "reactions": [],
            "account_id": 185315,
            "conversation_id": 3475816668
        },
        "socialUser": {
            "id": 571859378,
            "user_id": "4082895595259389",
            "username": "mohammad_gharehbagh_",
            "name": "Inventor...",
            "image": "https://example.com/profile.jpg",
            "social": "Instagram",
            "updated_at": "2025-09-18T01:31:30.000000Z"
        }
    }
}

# Test lead created event with empty user_id (should use "کاربر گرامی" as code)
test_lead_empty_user_payload = {
    "type": "leed_created",
    "user_id": "",  # Empty user_id
    "payload": {
        "id": "26033287",
        "type": "number",
        "value": "09123456789",
        "message_id": "3939210735",
        "created_at": "2025-09-18T01:31:30.000000Z",
        "message": {
            "id": 3939210735,
            "text": "09123456789",
        },
        "socialUser": {
            "id": 571859379,
            "username": "@test_user_empty",
            "name": "کاربر تست",
        }
    }
}

def send_test_webhook(payload, test_name):
    """Send a test webhook payload"""
    try:
        print(f"\n🧪 Testing: {test_name}")
        print(f"📤 Sending payload with text: '{payload['payload'].get('text', 'N/A')}'")
        
        response = requests.post(
            WEBHOOK_URL,
            json=payload,
            headers={'Content-Type': 'application/json'},
            timeout=10
        )
        
        print(f"✅ Response Status: {response.status_code}")
        if response.text:
            print(f"📥 Response Body: {response.text}")
        
    except requests.exceptions.ConnectionError:
        print(f"❌ Connection failed - Make sure webhook server is running on {WEBHOOK_URL}")
    except Exception as e:
        print(f"❌ Error: {e}")

def test_duplicate_sms_prevention():
    """Test duplicate SMS prevention by sending message + lead events"""
    print("\n🧪 Testing Daily SMS Limit System")
    print("-" * 40)
    
    # First send message_created
    send_test_webhook(test_message_payload, "1️⃣ MESSAGE CREATED (Just logged)")
    
    # Then send lead_created (simulating real scenario)
    send_test_webhook(test_lead_created_payload, "2️⃣ LEAD CREATED (First SMS today - should be sent)")
    
    # Send lead_created again (should be blocked - same day)
    send_test_webhook(test_lead_created_payload, "3️⃣ LEAD CREATED AGAIN (Should be blocked - same day)")
    
    print("\n💡 Note: Tomorrow this same user can receive SMS again!")

def main():
    print("🚀 Testing Phone Number Detection in Webhook")
    print("=" * 50)
    
    # Test cases
    test_cases = [
        (test_message_payload, "Single phone number (09155520952)"),
        (test_multiple_phones_payload, "Multiple phone numbers"),
        (test_no_phone_payload, "No phone numbers"),
        (test_lead_created_payload, "🎯 LEAD CREATED EVENT (MAIN TEST!) 🎯"),
        (test_lead_empty_user_payload, "🎯 LEAD WITH EMPTY USER_ID (Should use 'کاربر گرامی' as code) 🎯")
    ]
    
    for payload, name in test_cases:
        send_test_webhook(payload, name)
    
    # Test duplicate prevention
    test_duplicate_sms_prevention()
    
    print("\n" + "=" * 50)
    print("🏁 Test completed!")
    print("💡 Check your webhook server logs to see:")
    print("   - Phone number detection")
    print("   - Daily SMS limit system in action")
    print("   - Cache system working (24-hour cooldown)")
    print("   - Tomorrow reset functionality")

if __name__ == "__main__":
    main()
