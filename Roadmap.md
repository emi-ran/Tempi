## 🧭 **Prompt: Tempi CLI Roadmap**

**Amaç:**
Windows üzerinde çalışan, `.exe` olarak derlenen bir **CLI aracı** geliştir.
Aracın adı **Tempi** olacak.
Tempi’nin görevi, geçici işler için rastgele klasörler oluşturmak, bunları belirli bir süre (default 4 saat) aktif tutmak ve süre dolduğunda otomatik olarak silmek.
Ayrıca sistem yeniden başlasa bile zamanlama korunmalı.

---

### 🎯 **Temel Özellikler**

1. **CLI Komutları**

   - `tempi`
     → Varsayılan ayarlarla rastgele isimli bir klasör oluşturur.
     → Örneğin: `C:\Temp\tempi_20251106_123456`
     → Deadtime (ömrü): 4 saat (değiştirilebilir).
   - `tempi new --deadtime 2h`
     → Ömrü 2 saat olan geçici klasör oluşturur.
   - `tempi deletenow`
     → Tüm Tempi klasörlerini ve ilgili işlemleri hemen siler.
   - `tempi list`
     → Aktif geçici klasörleri ve kalan sürelerini gösterir.
   - `tempi clean`
     → Süresi dolmuş klasörleri ve içeriğini temizler.
   - `tempi --help`
     → Yardım ve komut listesini gösterir.

---

### ⚙️ **Davranış ve İç Mantık**

- Tempi her oluşturduğu klasör için bir metadata kaydı tutar:

  ```json
  {
    "path": "C:\\Temp\\tempi_20251106_123456",
    "created_at": "2025-11-06T12:34:56Z",
    "deadtime": "4h",
    "expires_at": "2025-11-06T16:34:56Z"
  }
  ```

- Bu kayıtlar örneğin `C:\Users\<user>\AppData\Local\Tempi\registry.json` dosyasında tutulur.
- Sistem yeniden başlasa bile bu kayıtlar korunur.
- Tempi başlatıldığında veya `clean` komutu çalıştığında bu kayıtlar kontrol edilir:

  - Süresi dolan klasörlerin silinmesi planlanır.
  - Silinmeden önce:

    - Klasörde aktif işlemler (örneğin `Process Explorer` tarzı API’lerle) tespit edilir.
    - Bu işlemler sonlandırılmaya çalışılır.

- Klasörün içindeki **en son değiştirilen dosyanın zamanı** `deadtime`’ı aşarsa klasör “aktif değil” sayılır ve temizlenir.

---

### 🧩 **Teknik Gereksinimler**

- **Dil:** Go
- **CLI Framework:** `cobra` (veya `urfave/cli/v2`)
- **Binary:** Tek dosya `.exe`
- **Zamanlama:**

  - Windows Task Scheduler veya arka planda çalışan hafif bir servis tarzı mekanizma kullanılabilir.
  - Alternatif olarak Tempi her çalıştığında kayıtları kontrol edip expired klasörleri temizleyebilir.

- **Veri Kaydı:** JSON dosyası (AppData altında).
- **Random Folder Name:** UUID veya zaman damgası ile oluşturulacak.

---

### 🧠 **Ekstra Özellikler (opsiyonel)**

- `tempi config set default_deadtime 2h`
  → Varsayılan ömrü ayarla.
- `tempi status`
  → Aktif görevlerin log’unu göster.
- `tempi update`
  → Yeni sürüm kontrolü.
- `tempi purge`
  → Tüm metadata’yı ve geçmiş klasörleri sil.

---

### 💾 **Örnek Kullanımlar**

```bash
tempi
# => Created folder C:\Temp\tempi_20251106_123456
#    This folder will expire in 4h.

tempi new --deadtime 1h
# => Created folder C:\Temp\tempi_20251106_133000
#    This folder will expire in 1h.

tempi list
# => [1] C:\Temp\tempi_20251106_123456 (expires in 3h12m)
# => [2] C:\Temp\tempi_20251106_133000 (expires in 58m)

tempi deletenow
# => Deleted 2 expired folders.
```

---

### 🧩 **Geliştirme Aşamaları (Roadmap)**

1. **Temel CLI yapısı oluştur (cobra ile)**
2. **`new` komutu:** klasör oluşturma + metadata kaydı
3. **`clean` komutu:** expired klasörleri bulma ve silme
4. **`deletenow` komutu:** hepsini anında silme
5. **Zamanlama mekanizması:** deadtime takibi (Timer veya yeniden başlatmada kontrol)
6. **İşlem takibi:** klasördeki dosyaları kullanan process’leri bulup sonlandırma
7. **Config dosyası ve varsayılan ayarlar**
8. **--help çıktısı ve CLI ergonomisi**
9. **Test + `.exe` build + PATH entegrasyonu**
