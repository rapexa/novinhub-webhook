#!/usr/bin/env python3
"""
تست اتصال به IPPanel API
بررسی وضعیت API Key و اتصال به سرویس
"""

import requests
import base64
import json

# تنظیمات IPPanel
API_KEY = "OWY1ZTAyMjEtMThmNi00NzRiLWFhOTItZTEwMmFhNDQzZTliZTcwM2EzODg5NzUzNWMwOWE3ZDliYWUyYTExMWZlMzY="
BASE_URL = "https://api2.ippanel.com/api/v1"

def test_api_connection():
    """تست اتصال به IPPanel API"""
    print("🔍 تست اتصال به IPPanel API...")
    
    # تست دریافت اعتبار (Credit)
    try:
        response = requests.get(
            f"{BASE_URL}/sms/accounting/credit/show",
            headers={
                'Apikey': API_KEY,
                'Content-Type': 'application/json'
            },
            timeout=10
        )
        
        print(f"📥 Response Status: {response.status_code}")
        print(f"📋 Response Headers: {dict(response.headers)}")
        
        if response.status_code == 200:
            data = response.json()
            print("✅ اتصال موفق!")
            print(f"📊 Response Data: {json.dumps(data, ensure_ascii=False, indent=2)}")
            return True
        else:
            print(f"❌ خطا در API: {response.status_code}")
            print(f"📋 Error Response: {response.text}")
            return False
            
    except requests.exceptions.ConnectionError:
        print("❌ خطای اتصال - احتمالاً API در دسترس نیست")
        return False
    except requests.exceptions.Timeout:
        print("❌ Timeout - API پاسخ نمی‌دهد")
        return False
    except Exception as e:
        print(f"❌ خطای غیرمنتظره: {e}")
        return False

def test_send_sms():
    """تست ارسال SMS"""
    print("\n📱 تست ارسال SMS...")
    
    test_data = {
        "code": "a2xjmxbszf27a7e",
        "sender": "+9850002040000000",
        "recipient": "09155520952",
        "variable": {
            "code": "تست"
        }
    }
    
    try:
        response = requests.post(
            f"{BASE_URL}/sms/pattern/normal/send",
            json=test_data,
            headers={
                'Apikey': API_KEY,
                'Content-Type': 'application/json'
            },
            timeout=15
        )
        
        print(f"📥 Response Status: {response.status_code}")
        print(f"📋 Response Data: {response.text}")
        
        if response.status_code == 200:
            print("✅ ارسال SMS موفق!")
            return True
        else:
            print(f"❌ خطا در ارسال SMS: {response.status_code}")
            return False
            
    except Exception as e:
        print(f"❌ خطا در ارسال SMS: {e}")
        return False

def decode_api_key():
    """رمزگشایی API Key برای بررسی"""
    print("\n🔐 بررسی API Key...")
    try:
        decoded = base64.b64decode(API_KEY).decode('utf-8')
        print(f"✅ API Key رمزگشایی شد: {decoded}")
    except Exception as e:
        print(f"❌ خطا در رمزگشایی API Key: {e}")

def main():
    print("🧪 تست اتصال IPPanel API")
    print("=" * 50)
    
    # بررسی API Key
    decode_api_key()
    
    # تست اتصال
    connection_ok = test_api_connection()
    
    if connection_ok:
        # تست ارسال SMS
        test_send_sms()
    else:
        print("\n💡 راه‌حل‌های پیشنهادی:")
        print("1. API Key را بررسی کنید")
        print("2. اتصال اینترنت را چک کنید")
        print("3. وضعیت سرویس IPPanel را بررسی کنید")
        print("4. از API Key معتبر استفاده کنید")

if __name__ == "__main__":
    main()
