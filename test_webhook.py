#!/usr/bin/env python3
"""
NovinHub Webhook Test Script
تست کامل برای webhook های NovinHub

استفاده:
    python test_webhook.py
    python test_webhook.py --event message_created
    python test_webhook.py --url https://asllmarket.org/webhook
"""

import requests
import json
import time
import argparse
from typing import Dict, Any

# URL پیش‌فرض webhook
DEFAULT_WEBHOOK_URL = "https://asllmarket.org/webhook"

class WebhookTester:
    def __init__(self, webhook_url: str = DEFAULT_WEBHOOK_URL):
        self.webhook_url = webhook_url
        self.session = requests.Session()
        self.session.headers.update({
            'Content-Type': 'application/json',
            'User-Agent': 'NovinHub-Webhook-Tester/1.0'
        })

    def send_webhook(self, event_data: Dict[str, Any]) -> bool:
        """ارسال webhook event و بررسی پاسخ"""
        try:
            print(f"\n🚀 ارسال event: {event_data['type']}")
            print(f"📤 Data: {json.dumps(event_data, ensure_ascii=False, indent=2)}")
            
            start_time = time.time()
            response = self.session.post(
                self.webhook_url,
                json=event_data,
                timeout=5
            )
            end_time = time.time()
            
            response_time = (end_time - start_time) * 1000  # میلی‌ثانیه
            
            print(f"📥 Response Status: {response.status_code}")
            print(f"⏱️  Response Time: {response_time:.2f}ms")
            
            if response.status_code == 200:
                print("✅ موفقیت‌آمیز!")
                try:
                    response_json = response.json()
                    print(f"📋 Response: {json.dumps(response_json, ensure_ascii=False, indent=2)}")
                except:
                    print(f"📋 Response Text: {response.text}")
            else:
                print(f"❌ خطا: {response.status_code}")
                print(f"📋 Response: {response.text}")
                
            return response.status_code == 200
            
        except requests.exceptions.RequestException as e:
            print(f"❌ خطای شبکه: {e}")
            return False
        except Exception as e:
            print(f"❌ خطای غیرمنتظره: {e}")
            return False

    def test_message_created(self) -> bool:
        """تست event پیغام جدید"""
        event_data = {
            "type": "message_created",
            "user_id": 123456,
            "payload": {
                "id": "msg_789",
                "content": "سلام! این یک پیغام تستی است 🚀",
                "account": {
                    "id": "acc_456",
                    "name": "تست اکانت",
                    "platform": "instagram"
                },
                "socialUser": {
                    "id": "social_123",
                    "username": "@test_user",
                    "full_name": "کاربر تست",
                    "profile_pic": "https://example.com/pic.jpg"
                }
            }
        }
        return self.send_webhook(event_data)

    def test_comment_created(self) -> bool:
        """تست event کامنت جدید"""
        event_data = {
            "type": "comment_created",
            "user_id": "789012",
            "payload": {
                "id": "comment_456",
                "content": "این یک کامنت تستی است! 👍",
                "account": {
                    "id": "acc_789",
                    "name": "اکانت تست",
                    "platform": "telegram"
                },
                "socialUser": {
                    "id": "social_456",
                    "username": "@commenter",
                    "full_name": "کامنت گذار تست"
                },
                "accountPost": {
                    "id": "post_123",
                    "title": "پست تستی",
                    "url": "https://example.com/post/123"
                }
            }
        }
        return self.send_webhook(event_data)

    def test_autoform_completed(self) -> bool:
        """تست event تکمیل فرم هوشمند"""
        event_data = {
            "type": "autoform_completed",
            "user_id": 345678,
            "payload": {
                "id": "form_789",
                "messages": [
                    {
                        "question": "نام شما چیست؟",
                        "answer": "علی احمدی"
                    },
                    {
                        "question": "شماره تماس؟",
                        "answer": "09123456789"
                    },
                    {
                        "question": "محصول مورد علاقه؟",
                        "answer": "لپ‌تاپ"
                    }
                ],
                "socialUser": {
                    "id": "social_789",
                    "username": "@form_user",
                    "full_name": "علی احمدی",
                    "phone": "09123456789"
                }
            }
        }
        return self.send_webhook(event_data)

    def test_leed_created(self) -> bool:
        """تست event ایجاد لید جدید"""
        event_data = {
            "type": "leed_created",
            "user_id": "456789",
            "payload": {
                "id": "lead_123",
                "phone": "09987654321",
                "messages": [
                    {
                        "content": "سلام، من به محصولاتتون علاقه‌مندم",
                        "timestamp": "2025-09-17T22:00:00Z"
                    },
                    {
                        "content": "شماره من: 09987654321",
                        "timestamp": "2025-09-17T22:01:00Z"
                    }
                ],
                "socialUser": {
                    "id": "social_lead_456",
                    "username": "@potential_customer",
                    "full_name": "مشتری احتمالی",
                    "phone": "09987654321"
                }
            }
        }
        return self.send_webhook(event_data)

    def test_revalidate(self) -> bool:
        """تست event احراز هویت مجدد"""
        event_data = {
            "type": "revalidate",
            "user_id": 999999,
            "payload": {
                "timestamp": int(time.time()),
                "challenge": "revalidate_webhook_test_2025",
                "app_id": "novinhub_webhook_test"
            }
        }
        return self.send_webhook(event_data)

    def test_health_check(self) -> bool:
        """تست health check endpoint"""
        try:
            health_url = self.webhook_url.replace('/webhook', '/health')
            print(f"\n🏥 تست Health Check: {health_url}")
            
            response = self.session.get(health_url, timeout=5)
            print(f"📥 Status: {response.status_code}")
            
            if response.status_code == 200:
                print("✅ سرویس سالم است!")
                try:
                    health_data = response.json()
                    print(f"📋 Health Data: {json.dumps(health_data, ensure_ascii=False, indent=2)}")
                except:
                    print(f"📋 Response: {response.text}")
                return True
            else:
                print(f"❌ مشکل در سرویس: {response.status_code}")
                return False
                
        except Exception as e:
            print(f"❌ خطا در health check: {e}")
            return False

    def run_all_tests(self) -> Dict[str, bool]:
        """اجرای تمام تست‌ها"""
        print("🧪 شروع تست‌های جامع NovinHub Webhook")
        print(f"🎯 Target URL: {self.webhook_url}")
        print("=" * 50)
        
        results = {}
        
        # Health Check
        results['health'] = self.test_health_check()
        time.sleep(1)
        
        # Test all event types
        results['message_created'] = self.test_message_created()
        time.sleep(1)
        
        results['comment_created'] = self.test_comment_created()
        time.sleep(1)
        
        results['autoform_completed'] = self.test_autoform_completed()
        time.sleep(1)
        
        results['leed_created'] = self.test_leed_created()
        time.sleep(1)
        
        results['revalidate'] = self.test_revalidate()
        
        return results

    def print_summary(self, results: Dict[str, bool]):
        """خلاصه نتایج"""
        print("\n" + "=" * 50)
        print("📊 خلاصه نتایج تست:")
        
        passed = 0
        total = len(results)
        
        for test_name, success in results.items():
            status = "✅ موفق" if success else "❌ ناموفق"
            print(f"  {test_name}: {status}")
            if success:
                passed += 1
        
        print(f"\n📈 نتیجه کلی: {passed}/{total} تست موفق")
        
        if passed == total:
            print("🎉 همه تست‌ها موفقیت‌آمیز بودند!")
        else:
            print(f"⚠️  {total - passed} تست ناموفق بود.")


def main():
    parser = argparse.ArgumentParser(description='NovinHub Webhook Tester')
    parser.add_argument('--url', default=DEFAULT_WEBHOOK_URL, help='Webhook URL')
    parser.add_argument('--event', choices=[
        'message_created', 'comment_created', 'autoform_completed', 
        'leed_created', 'revalidate', 'health'
    ], help='تست یک event خاص')
    
    args = parser.parse_args()
    
    tester = WebhookTester(args.url)
    
    if args.event:
        # تست یک event خاص
        if args.event == 'health':
            success = tester.test_health_check()
        elif args.event == 'message_created':
            success = tester.test_message_created()
        elif args.event == 'comment_created':
            success = tester.test_comment_created()
        elif args.event == 'autoform_completed':
            success = tester.test_autoform_completed()
        elif args.event == 'leed_created':
            success = tester.test_leed_created()
        elif args.event == 'revalidate':
            success = tester.test_revalidate()
        
        print(f"\n🏁 نتیجه: {'موفق' if success else 'ناموفق'}")
    else:
        # تست همه
        results = tester.run_all_tests()
        tester.print_summary(results)


if __name__ == "__main__":
    main()
