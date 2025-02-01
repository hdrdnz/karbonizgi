-kategorilere göre soru hazırlanması
-sorulara göre karbon ayak izinin hesaplanması ve skora göre yorumlama


KATEGORİLER -Birey
Ev Enerji Kullanımı (Elektrik, ısınma, soğutma)
Ulaşım (Araç kullanımı, toplu taşıma, uçak seyahati)
Beslenme (Et tüketimi, süt ürünleri, hazır gıdalar)
Atık ve Geri Dönüşüm (Çöp miktarı, geri dönüşüm, kompost)
Alışveriş (Kıyafet, elektronik eşya tüketimi, plastik poşet kullanımı)
Yaşam Alanı (Konut türü, hane halkı sayısı, enerji verimliliği)


--skora göre yorumlama  //burada api de toplam skoru iste
if total_emissions <= 167:
    print("🌱 Muhteşem bir iş çıkarmışsınız! Yaşam tarzınız gerçekten çevre dostu. Gezegenimiz sizin gibi kahramanlara ihtiyaç duyuyor. Bu harika alışkanlıklarınızı başkalarına da anlatmaya ne dersiniz?")
elif total_emissions <= 333:
    print("✅ Harika gidiyorsunuz! Küresel ortalamaya yakınsınız ve bu hiç de kolay bir şey değil. Ufak tefek ayarlamalarla dünyaya olan katkınızı daha da artırabilirsiniz. Birlikte daha yeşil bir gelecek mümkün!")
elif total_emissions <= 584:
    print("⚠️ Yolun başındasınız ama endişelenmeyin! Her büyük değişim, küçük adımlarla başlar. Belki bir hafta et tüketimini azaltabilir veya enerji tasarruflu cihazlara geçiş yapabilirsiniz. Doğa bu çabalarınızı kesinlikle fark edecek!")
else:
    print("🛑 Zaman harekete geçme zamanı! Evet, karbon ayak iziniz yüksek görünüyor, ama hiçbir şey için geç değil. Enerji tasarrufu, toplu taşıma ve geri dönüşüm gibi adımlarla dünyaya nefes aldırabilirsiniz. Unutmayın, her bir bilinçli tercih gezegen için bir iyilik!")

ŞİRKET KATEGORİLER
Üretim ve Sanayi Sektörü
Hizmet Sektörü
Lojistik ve Taşımacılık Sektörü
Tarım ve Gıda Sektörü
İnşaat ve Altyapı Sektörü
Perakende ve Toptan Satış
Enerji ve Kamu Hizmetleri






Enerji Tüketimi (Scope 2)
Taşımacılık ve Lojistik (Scope 1 ve Scope 3)
Üretim ve Süreç Emisyonları (Scope 1)
Atık ve Geri Dönüşüm (Scope 3)
İş Seyahatleri ve Çalışan Hareketliliği (Scope 3)
Ofis ve Bina Yönetimi (Scope 1 ve Scope 2)


Bu veriler, farklı enerji kullanımı yöntemlerinin çevreye olan karbon salım etkisini temsil eder. Karbon salım değerleri, CO₂ biriminde (kilogram) yıllık toplam salınım miktarlarını göstermektedir.

 //Doğalgaz: Ortalama yıllık tüketim 1200 m³ → 5.5 kg CO₂/m³.
// Kömür: Ortalama yıllık tüketim 2.5 ton → 2.5 kg CO₂/kg.
// Elektrik: Ortalama yıllık tüketim 6000 kWh → 0.5 kg CO₂/kWh (şebeke elektriği).
// Güneş Enerjisi: 0 kg CO₂ (yenilenebilir enerji).

!!!!!!!!!!!!!!!!!! belirli bir aralık oluşturarak kullanıcı seçimlerini ona uygun yorumlama ayarla.







