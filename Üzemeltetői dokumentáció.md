# A netstore projekt üzemeltetői dokumentációja

## Az üzemeltetés módja

A projekt legkönnyebben docker segítségével üzemeltethető. Az alkalmazás egyes részeit külön-külön konténerbe kell helyezni. Ezen konténerek buildelése a fejlesztői dokumentációkban kaptak helyet.

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
5. Környezeti változóknak az image által használtak kerülnek megadásra, a root jelszó, valamint egy alapértelmezett adatbázis (nem azt használja az alkalmazás, erről később). Ezek mellett a DB_TYPE környezeti változó is átadásra kerül. Ennek két lehetséges értéke van: adatok és szerkezet. Az itt megadott értékeket kell majd használni a többi konténer adatbázis csatlakozásával kapcsolatos környezeti változóinál.
6. A konténer portjai szintén itt kerülnek megadása. A külső port a 3307-esre van állítva, míg a belső 3306. Mivel az az alapértelmezett mysql port, ezért a gazda gépen már foglalt lehet, ezért van 3307-re irányítva, de ez a fájlban módosítható. A többi konténer a belső hálózatnak köszönhetően úgyis tudja használni a 3306-os portot.
7. Volumes-nak egyszer fel van csatolva a mysql_data, ahova az adatbázik kerülnek, hogy konténer újraindítás esetén ne vesszenek el az adatok. Ezek mellett fel van csatolva a db mappa több ízben is, egyrészt az image által használt init.d mappaként, mert automatikusan lefuttatja az ott található scripteket a konténer. A db/init mappában található rövid shell script pedig lefuttatja a DB_TYPE környezeti változó alapján a pusztán szerkezetet létrehozó sql-t, vagy a teszt adatokat is tartalmazó sql-t (amennyiben még nem létezik az adatbázis). Ezek a fájlok a db/szerkezet és db/adatok mappában találhatók, és a konténer a /db volume-ról éri el őket.
8. Ezen kívül van egy healthcheck is, ami ellenőrzi, hogy készen áll-e a konténer a használatra (ennek egyéb konténereknél van jelentősége).

### A központi szerver üzemeltetése

A központi szerver elkülönülten a többi konténertől is futtatható. Éles környezetben ez az alkalmazásrész kifejezetten külön, egy szerveren lenne futtatva, ami az egyes üzletek számára elérhető.  
A teljes projekt gyökérmappájában található egy docker-compose.yml fájl, ebben egy teljes rendszer szimulációja van kiépítve (3 üzlet + központi szerver). Amennyiben csak a központi szerverre van szükség, a compose fájl central_server service-ére, valamint a central_server mappában helyet kapó kódra van szükség.

1. Töltsük le a teljes projektet.
2. A compose fájlt helyezzük a central_server mappával egy szintre. (Ez a teljes projekt letöltésekor teljesül)
3. Ha szükséges, vegyük ki a compose fájl többi service-ét, és töröljük a nem szükséges egyéb mappákat (db, store_administration, product_administration). (Ezt csak akkor, ha kizárólag központi szervert telepítünk).
