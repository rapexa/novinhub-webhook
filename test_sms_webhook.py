#!/usr/bin/env python3
"""
تست کامل webhook برای ارسال SMS
شبیه‌سازی webhook event با شماره تلفن و user_id
"""

import requests
import json
import time
import random
from datetime import datetime

# تنظیمات
WEBHOOK_URL = "http://localhost:8080/webhook"
HEALTH_URL = "http://localhost:8080/health"

# شماره‌های تست (شماره‌های معتبر ایرانی)
TEST_PHONE_NUMBERS = [
    "09123456789",
    "09987654321", 
    "09155520952",
    "09301234567",
    "09111111111"
]

# User ID های تست
TEST_USER_IDS = [
    "76599340",  # Admin ID
    "123456789",
    "987654321",
    "555666777",
    "111222333"
]

class WebhookTester:
    def __init__(self, webhook_url=WEBHOOK_URL):
        self.webhook_url = webhook_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'User-Agent': 'SMS-Webhook-Tester/1.0'
        })

    def test_health_check(self):
        """تست health check"""
        print("\n🏥 تست Health Check...")
        try:
            response = self.session.get(HEALTH_URL, timeout=5)
            if response.status_code == 200:
                print("✅ Health Check موفق!")
                health_data = response.json()
                print(f"📊 Status: {health_data.get('status', 'unknown')}")
                print(f"⏰ Timestamp: {health_data.get('timestamp', 'unknown')}")
                return True
            else:
                print(f"❌ Health Check ناموفق: {response.status_code}")
                return False
        except Exception as e:
            print(f"❌ خطا در Health Check: {e}")
            return False

    def create_lead_webhook(self, phone_number, user_id):
        """ایجاد webhook event برای leed_created"""
        return {
            "type": "leed_created",
            "user_id": user_id,
            "payload": {
                "id": f"lead_{random.randint(100000, 999999)}",
                "type": "number",
                "value": phone_number,
                "message_id": f"msg_{random.randint(100000, 999999)}",
                "created_at": datetime.now().isoformat() + "Z",
                "message": {
                    "id": random.randint(100000, 999999),
                    "text": phone_number,
                    "date": datetime.now().strftime("%Y-%m-%d %H:%M:%S"),
                    "type": "text"
                },
                "socialUser": {
                    "id": random.randint(100000, 999999),
                    "username": f"test_user_{user_id}",
                    "name": f"کاربر تست {user_id}",
                    "social": "Instagram"
                }
            }
        }

    def send_webhook(self, event_data, test_name):
        """ارسال webhook event"""
        try:
            print(f"\n🚀 {test_name}")
            print(f"📱 شماره: {event_data['payload']['value']}")
            print(f"👤 User ID: {event_data['user_id']}")
            
            start_time = time.time()
            response = self.session.post(
                self.webhook_url,
                json=event_data,
                timeout=10
            )
            end_time = time.time()
            
            response_time = (end_time - start_time) * 1000
            
            print(f"📥 Response Status: {response.status_code}")
            print(f"⏱️ Response Time: {response_time:.2f}ms")
            
            if response.status_code == 200:
                print("✅ Webhook موفق!")
                try:
                    response_json = response.json()
                    print(f"📋 Response: {json.dumps(response_json, ensure_ascii=False, indent=2)}")
                except:
                    print(f"📋 Response Text: {response.text}")
                return True
            else:
                print(f"❌ Webhook ناموفق: {response.status_code}")
                print(f"📋 Response: {response.text}")
                return False
                
        except requests.exceptions.ConnectionError:
            print(f"❌ خطای اتصال - مطمئن شوید سرور روی {self.webhook_url} در حال اجرا است")
            return False
        except Exception as e:
            print(f"❌ خطای غیرمنتظره: {e}")
            return False

    def test_single_sms(self, phone_number, user_id):
        """تست ارسال SMS برای یک شماره"""
        event_data = self.create_lead_webhook(phone_number, user_id)
        test_name = f"تست SMS - {phone_number} (User: {user_id})"
        return self.send_webhook(event_data, test_name)

    def test_multiple_sms(self, count=3):
        """تست ارسال SMS برای چندین شماره"""
        print(f"\n📱 تست ارسال {count} SMS...")
        results = []
        
        for i in range(count):
            phone = random.choice(TEST_PHONE_NUMBERS)
            user_id = random.choice(TEST_USER_IDS)
            
            success = self.test_single_sms(phone, user_id)
            results.append(success)
            
            # تاخیر بین تست‌ها
            if i < count - 1:
                time.sleep(2)
        
        return results

    def test_duplicate_prevention(self):
        """تست جلوگیری از ارسال duplicate SMS"""
        print("\n🔄 تست جلوگیری از Duplicate SMS...")
        
        # انتخاب یک شماره و user_id ثابت
        phone = random.choice(TEST_PHONE_NUMBERS)
        user_id = random.choice(TEST_USER_IDS)
        
        print(f"📱 شماره ثابت: {phone}")
        print(f"👤 User ID ثابت: {user_id}")
        
        # ارسال اول (باید موفق باشد)
        print("\n1️⃣ ارسال اول (باید موفق باشد):")
        success1 = self.test_single_sms(phone, user_id)
        
        time.sleep(1)
        
        # ارسال دوم (باید block شود - same day)
        print("\n2️⃣ ارسال دوم (باید block شود - same day):")
        success2 = self.test_single_sms(phone, user_id)
        
        return success1, success2

    def test_invalid_phone(self):
        """تست شماره تلفن نامعتبر"""
        print("\n❌ تست شماره تلفن نامعتبر...")
        
        invalid_phones = [
            "123456789",  # شماره نامعتبر
            "0912345678",  # کوتاه
            "091234567890",  # بلند
            "08123456789",  # پیش‌شماره نامعتبر
        ]
        
        results = []
        for phone in invalid_phones:
            event_data = self.create_lead_webhook(phone, random.choice(TEST_USER_IDS))
            test_name = f"تست شماره نامعتبر: {phone}"
            success = self.send_webhook(event_data, test_name)
            results.append(success)
            time.sleep(1)
        
        return results

    def test_empty_user_id(self):
        """تست user_id خالی"""
        print("\n👤 تست User ID خالی...")
        
        event_data = self.create_lead_webhook(
            random.choice(TEST_PHONE_NUMBERS), 
            ""  # User ID خالی
        )
        
        test_name = "تست User ID خالی (باید از 'کاربر گرامی' استفاده کند)"
        return self.send_webhook(event_data, test_name)

    def run_comprehensive_test(self):
        """اجرای تست جامع"""
        print("🧪 شروع تست جامع SMS Webhook")
        print("=" * 60)
        
        # تست Health Check
        if not self.test_health_check():
            print("❌ Health Check ناموفق - تست متوقف شد")
            return
        
        # تست‌های مختلف
        test_results = {}
        
        # 1. تست SMS تک
        print("\n" + "="*60)
        print("1️⃣ تست SMS تک")
        test_results['single_sms'] = self.test_single_sms(
            random.choice(TEST_PHONE_NUMBERS),
            random.choice(TEST_USER_IDS)
        )
        
        # 2. تست SMS چندگانه
        print("\n" + "="*60)
        print("2️⃣ تست SMS چندگانه")
        test_results['multiple_sms'] = self.test_multiple_sms(3)
        
        # 3. تست جلوگیری از Duplicate
        print("\n" + "="*60)
        print("3️⃣ تست جلوگیری از Duplicate")
        success1, success2 = self.test_duplicate_prevention()
        test_results['duplicate_prevention'] = (success1, success2)
        
        # 4. تست شماره نامعتبر
        print("\n" + "="*60)
        print("4️⃣ تست شماره نامعتبر")
        test_results['invalid_phone'] = self.test_invalid_phone()
        
        # 5. تست User ID خالی
        print("\n" + "="*60)
        print("5️⃣ تست User ID خالی")
        test_results['empty_user_id'] = self.test_empty_user_id()
        
        # خلاصه نتایج
        self.print_summary(test_results)

    def print_summary(self, results):
        """چاپ خلاصه نتایج"""
        print("\n" + "="*60)
        print("📊 خلاصه نتایج تست:")
        print("="*60)
        
        # تست SMS تک
        print(f"1️⃣ SMS تک: {'✅ موفق' if results['single_sms'] else '❌ ناموفق'}")
        
        # تست SMS چندگانه
        multiple_success = sum(results['multiple_sms'])
        multiple_total = len(results['multiple_sms'])
        print(f"2️⃣ SMS چندگانه: {multiple_success}/{multiple_total} موفق")
        
        # تست Duplicate
        success1, success2 = results['duplicate_prevention']
        print(f"3️⃣ Duplicate Prevention:")
        print(f"   - ارسال اول: {'✅ موفق' if success1 else '❌ ناموفق'}")
        print(f"   - ارسال دوم: {'✅ Block شد' if not success2 else '❌ Block نشد'}")
        
        # تست شماره نامعتبر
        invalid_success = sum(results['invalid_phone'])
        invalid_total = len(results['invalid_phone'])
        print(f"4️⃣ شماره نامعتبر: {invalid_success}/{invalid_total} Block شد")
        
        # تست User ID خالی
        print(f"5️⃣ User ID خالی: {'✅ موفق' if results['empty_user_id'] else '❌ ناموفق'}")
        
        print("\n" + "="*60)
        print("💡 نکات مهم:")
        print("- بررسی logs سرور برای جزئیات بیشتر")
        print("- مطمئن شوید SMS service فعال است")
        print("- pattern management درست کار می‌کند")
        print("- daily limit system فعال است")
        print("="*60)

def main():
    print("🚀 SMS Webhook Tester")
    print("تست کامل webhook برای ارسال SMS")
    print("=" * 60)
    
    tester = WebhookTester()
    
    # اجرای تست جامع
    tester.run_comprehensive_test()
    
    print("\n🏁 تست کامل شد!")
    print("برای تست دستی، از تابع‌های زیر استفاده کنید:")
    print("- tester.test_single_sms('09123456789', '123456')")
    print("- tester.test_duplicate_prevention()")
    print("- tester.test_invalid_phone()")

if __name__ == "__main__":
    main()
