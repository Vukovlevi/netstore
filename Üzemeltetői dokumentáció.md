# A netstore projekt üzemeltetői dokumentációja

## Az üzemeltetés módja

A projekt legkönnyebben docker segítségével üzemeltethető. Az alkalmazás egyes részeit külön-külön konténerbe kell helyezni. Ezen konténerek buildelésének módjai a fejlesztői dokumentációkban kaptak helyet.  
Ennek megfelelően a telepítő számítógépen szükséges a docker jelenléte, egyéb technikai követelménye az alkalmazásnak nincs, viszont internet kapcsolat szükséges.

## Gyors indítás és szimuláció

Amennyiben szimulálni kéne egy teljes hálózatot, a projekt gyökérmappájában helyet kapó docker compose ezt biztosítja. Tartalmazza a szükséges központi szervert, valamint 3 üzletet, amiből az első kettő teszt adatokkal rendelkezik, a harmadik pedig üres.

Az indításhoz a következő szabad portok szükségesek: 3307, 3308, 3309, 8000, 10000, 11000, 12000, 13000, 14000, 42069. Ezen portok átállíthatók a compose fájlban, erről a dokumentáció későbbi részében lesz szó.

Amennyiben ezek megfelelők, a compose fájlt tartalmazó mappában kiadható a: docker compose up -d parancs. Ez első indításkor létrehozza az adatbázisokat (DB_TYPE környezeti változó alapján szerkezetet, vagy teszt adatokat), kiépíti a konténereket és elindítja őket. A futás után ellenőrizendő, hogy minden konténer elindult-e, az üzlet adminisztrációk folyamatosan újraindulhatnak, amíg nincs kész az adatbázis. Amint mindegyik fut, a teljes szimuláció elindult.

A használathoz a 3 üzletre lehet belépni, az elsőre a 8000-es, a másodikra a 11000-es, a harmadikra a 13000-es porton keresztül.
Az ez utáni teendők a felhasználói dokumentációban kerültek kifejtésre.

Amennyiben az adatbázisokat újra akarjuk telepíteni (mondjuk mert megváltoztattuk a DB_TYPE környezeti változó értékét, és szeretnénk ha lenne hatása), akkor a: docker compose down -v paranccsal állíthatók le a konténerek, hogy az adatbázisok törlődjenek. Ha adatvesztés nélkül kell leállítani, akkor: docker compose down parancs.

Ha mondjuk frissült valamelyik alkalmazásrész, ezért újra kéne építeni a konténert, akkor a: docker compose up -d --build paranccsal tehetjük meg (csak leállított konténerek esetén). Ha nem szükséges újra felépíteni őket, csak el szeretnénk indítani, akkor a: docker compose up -d parancs használható.

### A központi szerver

A szimulációban a központi szerver egy belső docker hálózaton fut az üzletekkel. Ennek megfelelően van konfigurálva a fájlban, környezeti változókon keresztül az összes üzlet adminisztráció rész. A valóságban ez egy külön szerveren kapna helyet, míg az üzlet rész minden üzletben lenne telepítve.

### Az üzlet rész

Az üzlet rész három konténerből tevedők össze: az adatbázis, a termék adminisztráció, valamint az üzlet adminisztráció. A szimulációban található 3 üzlet miatt ezen konténerek 3-szor szerepelnek. A hozzájuk rendelt környezeti változókon keresztül szabadon konfigurálható a működésük.

## A valós üzemeltetés

Éles körülmények között is dockerben érdemes üzemeltetni a rendszereket. A központi szerver konténerét felállítani egy nyilvános ip címmel rendelkező szerveren. Az üzlethez tartozó 3 konténert pedig az üzlet belső szerverén érdemes futtatni. A környezeti változókon keresztül konfigurálható az összes elérés. A szimulációhoz használt docker compose alapján érdemes elkészíteni a valós helyzetben használtakat, jellemzően csak kitörlendő némi rész:

1. A központi szerver telepítéséhez a compose fájl central_server service-ére van szükség, valamint a central_server mappára. A compose és a mappa legyen egy szinten. A többi törölhető.
2. Egy üzlet telepítéséhez a compose fájl mysql, prod_admin és store_admin service-eire van szükség, valamint a product_administraion, store_administration és db mappákra. A többi törölhető.

A db mappa az adatbázis konténer által lesz használva, ahol az init mappában található scriptek felmountolódnak volume-ként oda, ahol automatikusan lefutásra kerülnek. A script pedig vagy a db/adatok, vagy a db/szerkezet mappában található sql fájlt futtatja le, a DB_TYPE környezeti változó alapján. Ez vagy egy pusztán szerkezeti, vagy teszt adatokkal feltöltött adatbázist csinál, valós helyzetben a szerkezeti használata ajánlott (ennek hogyanja később kifejtve).

## A docker compose

A docker compose egy technológia, ami lehetőve teszi, hogy egy compose fájl alapján (docker-compose.yml) több konténer kezelhető legyen (service-ek formájában kerülnek megadásra). A konténerek között egy belső virtuális hálózat is lérejön, aminek köszönhetően a benne helyet kapó konténerek a service nevét használva domain/ip címként tudnak egymásra hivatkozni.

### A projektben található compose

A projektnek része a gyökérmappában található docker-compose.yml. Ez a compose fájl egy teljes hálózatot szimulál egy központi szerverrel, valamint 3 üzlettel.

### Service-ek

#### Központi szerver

A központi szerver a fájlban megtalálható első service.

1. A service neve: central_server
2. A build paraméter azt mondja meg, hol található a kiépítendő és használandó docker konténer ehhez a service-hez (ez a ./central_server mappa a projekt kiosztásában).
3. A konténer neve central_server lesz.
4. A konténer portjainak megfeleltetése a valós számítógép portjaihoz. A két port közül az első a valós számítógépé, a második pedig a konténer belső portja. A központi szerver alapértelmezetten a 42069-es protot használja és a 0.0.0.0-ás ip címen hallgat.
5. A konténernek átadott környezeti változók segítségével konfigurálható a konténer működése. Az IP_ADDRESS-szel módosítható az ip cím, amin hallgat (nem javasolt), a PORT segítségével a port (ezen esetben a store_admin service változóit is módosítani érdemes, de nem kötelező), valamint a PSK, ami az előre megbeszélt kulcs az authentikáláshoz (ezt is érdemes a store_admin service-nél is átírni, de nem kötelező).

#### Az adatbázis

A mysql adatbázist két alkalmazásrész is használja: az üzlet adminisztráció és a termék adminisztráció. A compose fájlban megtalálható második service.

1. A service neve: mysql
2. Ez nem helyi konténert buildel, hanem a mysql:8.0-ás imaget használja.
3. A konténer neve: mysql.
4. A restart: always sor azt mondja meg, hogy ha nem sikerült a konténer indítása, próbálja újra.
5. Környezeti változóknak az image által használtak kerülnek megadásra, a root jelszó, valamint egy alapértelmezett adatbázis (nem azt használja az alkalmazás, hanem egy netstore nevű adatbázist, amit az init script és a volume-ként használt db mappa hoz létre). Ezek mellett a DB_TYPE környezeti változó is átadásra kerül. Ennek két lehetséges értéke van: adatok és szerkezet. Az itt megadott értékeket kell majd használni a többi konténer adatbázis csatlakozásával kapcsolatos környezeti változóinál.
6. A konténer portjai szintén itt kerülnek megadása. A külső port a 3307-esre van állítva, míg a belső 3306. Mivel az az alapértelmezett mysql port, ezért a gazda gépen már foglalt lehet, ezért van 3307-re irányítva, de ez a fájlban módosítható. A többi konténer a belső hálózatnak köszönhetően úgyis tudja használni a 3306-os portot.
7. Volumes-nak egyszer fel van csatolva a mysql_data, ahova az adatbázik kerülnek, hogy konténer újraindítás esetén ne vesszenek el az adatok. Ezek mellett fel van csatolva a db mappa több ízben is, egyrészt az image által használt init.d mappaként, mert automatikusan lefuttatja az ott található scripteket a konténer. A db/init mappában található rövid shell script pedig lefuttatja a DB_TYPE környezeti változó alapján a pusztán szerkezetet létrehozó sql-t, vagy a teszt adatokat is tartalmazó sql-t (amennyiben még nem létezik az adatbázis). Ezek a fájlok a db/szerkezet és db/adatok mappában találhatók, és a konténer a /db volume-ról éri el őket.
8. Ezen kívül van egy healthcheck is, ami ellenőrzi, hogy készen áll-e a konténer a használatra (ennek egyéb konténereknél van jelentősége).

#### Termék adminisztráció

A termék adminisztráció alkalmazásrész konténere. A compose fájlban megtalálható harmadik service.

1. A service neve: prod_admin
2. A build paraméter azt mondja meg, hol található a kiépítendő és használandó docker konténer ehhez a service-hez (ez a ./product_administration mappa a projekt kiosztásában). Valamint átadja paraméterként az alkalmazásrész párját képző üzlet adminisztráció linkjét (amit egy bármely gépen elindított böngésző használni tud, szóval nem a belső docker hálózatot használja).
3. A konténer neve: prod_admin.
4. A port megfeleltések kerülnek meghatározásra, a konténer belső portja a 80, ezt a gazdagép bármely szabad projtára be lehet állítani. A fájlban ez a 10000-es port.
5. A konfigurációhoz használt környezeti változók kerülnek átadásra, a PSK amit authentikációhoz használ (egy hálózatos keresés által indított api hívásnál), valamint az adatbázis kapcsolathoz szükséges adatok (itt már használható a belső docker hálózat hostnévként), amik az adatbázis konténernél kerülnek megadásra.
6. Ez a konténer az adatbázis készenlététől függ.

#### Üzlet adminisztráció

Az üzlet adminisztráció alakalmazásrész konténere. A compose fájlban megtalálható negyedik service.

1. A service neve: store_admin
2. A build paraméter azt mondja meg, hol található a kiépítendő és használandó docker konténer ehhez a service-hez (ez a ./store_administration mappa a projekt kiosztásában). Valamint átadja paraméterként az alkalmazásrész párját képző termék adminisztráció linkjét (amit egy bármely gépen elindított böngésző használni tud, szóval nem a belső docker hálózatot használja).
3. A konténer neve: store_admin.
4. A port megfeleltések kerülnek meghatározásra, a konténer belső portja a 8000 (ez konfigurálható), ezt a gazdagép bármely szabad projtára be lehet állítani. A fájlban ez a 8000-es port.
5. A konfigurációhoz használt környezeti változók kerülnek átadásra, a PSK amit a központi szerverre való csatlakozáskor authentikációra használ, valamint az adatbázis kapcsolathoz szükséges adatok (itt már használható a belső docker hálózat hostnévként), amik az adatbázis konténernél kerülnek megadásra. Ezen túl konfigurálható az IP_ADDRESS és PORT környezeti változókkal, hogy mire hallgat a szerver. Az ip alapértelmezetten 0.0.0.0, ezt nem érdemes változtatni, a port viszont konfigurálható (viszont akkor a párt képző termék adminisztráció résznél átírandó a VITE változó). Ezek mellett itt kerülnek beállításra a központi szerverre csatlakozáshoz szükséges adatok, amit a központi szerver konténernél lehet konfigurálni (a portot, a címe vagy a docker belső hálózaté, ha szimulációban futnak, vagy a szerver valós címe egy éles környezetben).
6. Ez a konténer az adatbázis készenlététől függ, de a restart: always itt is fut, mivel az adatbázisnak idő kell amíg elindul, így addig újra kell indítani a konténert, amíg az nincs kész.

Nem kötelező konfigurálni a központi szerverre csatlakozáshoz szükséges adatokat, mivel az az oldalon is megadható, de ajánlott a compose fájlban.
