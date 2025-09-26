#!/usr/bin/env python3
"""
تست ساده SMS Webhook
تست سریع برای بررسی ارسال SMS
"""

import requests
import json
import time

# تنظیمات
WEBHOOK_URL = "https://asllmarket.org/webhook"

def test_sms_webhook(phone_number, user_id):
    """تست ارسال SMS با webhook"""
    
    # ایجاد webhook event
    webhook_data = {
        "type": "leed_created",
        "user_id": user_id,
        "payload": {
            "id": "test_lead_123",
            "type": "number",
            "value": phone_number,
            "message_id": "test_msg_456",
            "created_at": "2025-01-17T15:30:45Z",
            "message": {
                "id": 123456,
                "text": phone_number,
                "date": "2025-01-17 15:30:45",
                "type": "text"
            },
            "socialUser": {
                "id": 789012,
                "username": f"test_user_{user_id}",
                "name": f"کاربر تست {user_id}",
                "social": "Instagram"
            }
        }
    }
    
    print(f"🚀 ارسال webhook برای SMS...")
    print(f"📱 شماره: {phone_number}")
    print(f"👤 User ID: {user_id}")
    print(f"📤 Data: {json.dumps(webhook_data, ensure_ascii=False, indent=2)}")
    
    try:
        response = requests.post(
            WEBHOOK_URL,
            json=webhook_data,
            headers={'Content-Type': 'application/json'},
            timeout=10
        )
        
        print(f"📥 Response Status: {response.status_code}")
        print(f"📋 Response: {response.text}")
        
        if response.status_code == 200:
            print("✅ Webhook موفق!")
            return True
        else:
            print("❌ Webhook ناموفق!")
            return False
            
    except requests.exceptions.ConnectionError:
        print("❌ خطای اتصال - مطمئن شوید سرور در حال اجرا است")
        return False
    except Exception as e:
        print(f"❌ خطا: {e}")
        return False

def test_health_check():
    """تست health check"""
    try:
        response = requests.get("https://asllmarket.org/health", timeout=5)
        if response.status_code == 200:
            print("✅ Health Check موفق!")
            print(f"📊 {response.json()}")
            return True
        else:
            print(f"❌ Health Check ناموفق: {response.status_code}")
            return False
    except Exception as e:
        print(f"❌ خطا در Health Check: {e}")
        return False

def main():
    print("🧪 تست ساده SMS Webhook")
    print("=" * 40)
    
    # تست Health Check
    print("\n1️⃣ تست Health Check...")
    if not test_health_check():
        print("❌ سرور در حال اجرا نیست!")
        return
    
    # تست SMS
    print("\n2️⃣ تست ارسال SMS...")
    
    # شماره‌های تست
    test_cases = [
        ("09123456789", "76599340"),  # شماره معتبر + Admin ID
        ("09987654321", "123456789"),  # شماره معتبر + User ID
        ("09115392188", ""),  # شماره معتبر + User ID خالی
    ]
    
    for phone, user_id in test_cases:
        print(f"\n📱 تست: {phone} (User: {user_id or 'خالی'})")
        success = test_sms_webhook(phone, user_id)
        print(f"نتیجه: {'✅ موفق' if success else '❌ ناموفق'}")
        time.sleep(2)  # تاخیر بین تست‌ها
    
    print("\n🏁 تست کامل شد!")
    print("💡 برای تست بیشتر، فایل test_sms_webhook.py را اجرا کنید")

if __name__ == "__main__":
    main()
