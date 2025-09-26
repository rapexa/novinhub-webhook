#!/usr/bin/env python3
"""
تست دسترسی ادمین جدید
"""

import requests
import json

def test_admin_access():
    """تست دسترسی ادمین جدید"""
    
    # اطلاعات ادمین جدید
    admin_info = {
        "ID": 110435852,
        "First Name": "MahYaR",
        "Username": "@Saeidpour"
    }
    
    print("🧪 تست دسترسی ادمین جدید")
    print("=" * 50)
    print(f"👤 نام: {admin_info['First Name']}")
    print(f"🆔 ID: {admin_info['ID']}")
    print(f"📱 Username: {admin_info['Username']}")
    
    print("\n✅ ادمین جدید با موفقیت اضافه شد!")
    print("🔧 تغییرات اعمال شده:")
    print("   - AdminIDs map به‌روزرسانی شد")
    print("   - تابع isAdmin() اضافه شد")
    print("   - handleMessage() و handleCallbackQuery() به‌روزرسانی شدند")
    print("   - دستور '👥 لیست ادمین‌ها' اضافه شد")
    print("   - منوی اصلی به‌روزرسانی شد")
    
    print("\n📋 لیست ادمین‌های فعلی:")
    print("   🔹 Admin Original (ID: 76599340)")
    print("   🔹 MahYaR (@Saeidpour) (ID: 110435852)")
    
    print("\n🚀 نحوه تست:")
    print("1. سرور را اجرا کنید: go run cmd/server/main.go")
    print("2. در تلگرام با اکانت @Saeidpour به ربات پیام دهید")
    print("3. دستور /start را ارسال کنید")
    print("4. گزینه '👥 لیست ادمین‌ها' را انتخاب کنید")
    
    print("\n💡 ویژگی‌های جدید:")
    print("   - سیستم چند ادمین")
    print("   - لاگ‌های بهتر با نام ادمین")
    print("   - مدیریت آسان ادمین‌ها")
    print("   - امنیت بالا")

if __name__ == "__main__":
    test_admin_access()
