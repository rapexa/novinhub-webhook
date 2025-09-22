# تست SMS Webhook

فایل‌های تست برای بررسی عملکرد webhook و ارسال SMS.

## 📁 فایل‌های تست

### 1. `test_simple_sms.py` - تست ساده
تست سریع و ساده برای بررسی عملکرد اولیه.

```bash
python test_simple_sms.py
```

**ویژگی‌ها:**
- ✅ Health Check
- ✅ تست SMS با شماره معتبر
- ✅ تست User ID خالی
- ✅ تست Admin ID

### 2. `test_sms_webhook.py` - تست جامع
تست کامل و جامع برای بررسی تمام قابلیت‌ها.

```bash
python test_sms_webhook.py
```

**ویژگی‌ها:**
- ✅ Health Check
- ✅ تست SMS تک
- ✅ تست SMS چندگانه
- ✅ تست جلوگیری از Duplicate
- ✅ تست شماره نامعتبر
- ✅ تست User ID خالی
- ✅ گزارش‌گیری کامل

## 🚀 نحوه اجرا

### 1. اجرای سرور
```bash
# در terminal اول
go run cmd/server/main.go
```

### 2. اجرای تست
```bash
# در terminal دوم
python test_simple_sms.py
# یا
python test_sms_webhook.py
```

## 📊 نتایج مورد انتظار

### ✅ تست‌های موفق:
- Health Check: 200 OK
- SMS معتبر: 200 OK + SMS ارسال می‌شود
- User ID خالی: 200 OK + SMS با "کاربر گرامی"

### ❌ تست‌های ناموفق:
- شماره نامعتبر: 200 OK + SMS ارسال نمی‌شود
- Duplicate SMS: 200 OK + SMS block می‌شود

## 🔍 بررسی Logs

برای بررسی جزئیات، logs سرور را چک کنید:

```bash
# در terminal سرور
# logs JSON format نمایش داده می‌شود
```

**Logs مهم:**
- `📲 SMS SENDING INITIATED` - شروع ارسال SMS
- `✅ SMS SENT SUCCESSFULLY` - SMS موفق
- `❌ SMS SEND FAILED` - SMS ناموفق
- `🚫 SMS BLOCKED - ALREADY SENT TODAY` - SMS block شده

## 🛠️ تنظیمات تست

### تغییر URL:
```python
WEBHOOK_URL = "http://localhost:8080/webhook"
HEALTH_URL = "http://localhost:8080/health"
```

### تغییر شماره‌های تست:
```python
TEST_PHONE_NUMBERS = [
    "09123456789",
    "09987654321", 
    "09155520952",
    # شماره‌های خود را اضافه کنید
]
```

### تغییر User ID های تست:
```python
TEST_USER_IDS = [
    "76599340",  # Admin ID
    "123456789",
    "987654321",
    # User ID های خود را اضافه کنید
]
```

## 📱 تست Pattern Management

### 1. تست Pattern فعلی:
```bash
# در Telegram
/start
📱 پترن امروز
```

### 2. تغییر Pattern:
```bash
# در Telegram
➡️ برو به پترن بعدی
```

### 3. لیست Pattern ها:
```bash
# در Telegram
📋 لیست پترن‌ها
```

## 🔧 عیب‌یابی

### مشکل: Connection Error
```
❌ خطای اتصال - مطمئن شوید سرور روی http://localhost:8080 در حال اجرا است
```
**راه حل:** سرور را اجرا کنید

### مشکل: SMS ارسال نمی‌شود
**بررسی کنید:**
- SMS service فعال است (`enabled: true`)
- API key صحیح است
- Pattern تنظیم شده است
- Originator تنظیم شده است

### مشکل: Pattern تغییر نمی‌کند
**بررسی کنید:**
- ربات تلگرام اجرا شده است
- Admin ID صحیح است
- Config درست load شده است

## 📈 مثال خروجی

```
🧪 شروع تست جامع SMS Webhook
============================================================

🏥 تست Health Check...
✅ Health Check موفق!
📊 Status: healthy
⏰ Timestamp: 2025-01-17T15:30:45Z

============================================================
1️⃣ تست SMS تک
🚀 تست SMS - 09123456789 (User: 76599340)
📱 شماره: 09123456789
👤 User ID: 76599340
📥 Response Status: 200
⏱️ Response Time: 45.23ms
✅ Webhook موفق!
📋 Response: {
  "status": "success",
  "message": "Webhook processed successfully"
}
✅ Webhook موفق!
```

## 🎯 نکات مهم

1. **سرور باید اجرا باشد** قبل از اجرای تست
2. **SMS service باید فعال باشد** در config
3. **Pattern باید تنظیم شده باشد** در config
4. **API key باید صحیح باشد** برای IPPanel
5. **Logs را بررسی کنید** برای جزئیات بیشتر

## 📞 پشتیبانی

اگر مشکلی داشتید:
1. Logs سرور را چک کنید
2. Config را بررسی کنید
3. Health check را تست کنید
4. Pattern management را چک کنید
