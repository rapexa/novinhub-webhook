#!/usr/bin/env python3
"""
تست debug برای مشکل SMS
"""

import requests
import json
import time

def test_webhook_with_debug():
    """تست webhook با debug بیشتر"""
    
    webhook_data = {
        "type": "leed_created",
        "user_id": "76599340",
        "payload": {
            "id": "debug_lead_123",
            "type": "number",
            "value": "09155520952",
            "message_id": "debug_msg_456",
            "created_at": "2025-01-17T15:30:45Z",
            "message": {
                "id": 123456,
                "text": "09155520952",
                "date": "2025-01-17 15:30:45",
                "type": "text"
            },
            "socialUser": {
                "id": 789012,
                "username": "debug_user",
                "name": "کاربر Debug",
                "social": "Instagram"
            }
        }
    }
    
    print("🚀 ارسال webhook برای debug SMS...")
    print(f"📱 شماره: {webhook_data['payload']['value']}")
    print(f"👤 User ID: {webhook_data['user_id']}")
    
    try:
        response = requests.post(
            "http://localhost:8080/webhook",
            json=webhook_data,
            headers={'Content-Type': 'application/json'},
            timeout=15
        )
        
        print(f"📥 Response Status: {response.status_code}")
        print(f"📋 Response: {response.text}")
        
        if response.status_code == 200:
            print("✅ Webhook موفق!")
            print("\n💡 حالا لاگ‌های سرور را بررسی کنید:")
            print("   - آیا retry mechanism کار می‌کند؟")
            print("   - آیا error message واضح‌تر شده؟")
            print("   - آیا API connection بهتر شده؟")
        else:
            print("❌ Webhook ناموفق!")
            
    except Exception as e:
        print(f"❌ خطا: {e}")

def test_ippanel_direct():
    """تست مستقیم IPPanel API"""
    print("\n🔍 تست مستقیم IPPanel API...")
    
    api_key = "OWY1ZTAyMjEtMThmNi00NzRiLWFhOTItZTEwMmFhNDQzZTliZTcwM2EzODg5NzUzNWMwOWE3ZDliYWUyYTExMWZlMzY="
    
    # تست دریافت اعتبار
    try:
        response = requests.get(
            "https://api2.ippanel.com/api/v1/sms/accounting/credit/show",
            headers={
                'Apikey': api_key,
                'Content-Type': 'application/json'
            },
            timeout=10
        )
        
        print(f"📊 Credit Check - Status: {response.status_code}")
        if response.status_code == 200:
            print(f"✅ API Key معتبر! Response: {response.text}")
        else:
            print(f"❌ API Key مشکل دارد! Response: {response.text}")
            
    except Exception as e:
        print(f"❌ خطا در تست API: {e}")

def main():
    print("🧪 تست Debug SMS Webhook")
    print("=" * 50)
    
    # تست مستقیم API
    test_ippanel_direct()
    
    # تست webhook
    test_webhook_with_debug()
    
    print("\n💡 راه‌حل‌های پیشنهادی:")
    print("1. API Key را بررسی کنید")
    print("2. اتصال اینترنت را چک کنید") 
    print("3. IPPanel service status را بررسی کنید")
    print("4. از API Key معتبر استفاده کنید")
    print("5. Rate limit را چک کنید")

if __name__ == "__main__":
    main()
