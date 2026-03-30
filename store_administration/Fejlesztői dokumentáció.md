# Netstore webalkalmazás - üzlet adminisztráció fejlesztői dokumentáció

## Backend

### Felhasznált technológiák

1. Go programozási nyelv
2. go-sql-driver/mysql - az adatbázis kapcsolathoz használt sql driver
3. joho/godotenv - a .env fájlban tárolt adatok beolvasására
4. labstack/echo/v4 - a webszervernek használt keretrendszer
5. golang.org/x/crypto - jelszó kezeléshez (bcrypt, 12-es cost)
6. Docker

### Package-k

#### Config

##### Config struktúrái

1. Config:
   - dbConfig: DBConfig (az adatbázis csatlakozásához használt konfiguráció) - környezeti változók a db package-ben kifejtve
   - Ip: string (az ip cím, amire a webszerver hallgat, alapértelmezetten: 0.0.0.0) - környezeti változó: IP_ADDRESS
   - Port: string (a port, amire a webszerver hallgat, alapértelmezetten: 8000) - környezeti változó: PORT
   - CentralServerAddress: string (a központi szerver ip címe, alapértelmezetten localhost) - környezeti változó: CENTRAL_SERVER_ADDRESS
   - CentralServerPort: string (a központi szerver portja, alapértelmezetten 42069) - környezeti változó: CENTRAL_SERVER_PORT
   - Psk: string (a pre-shared key, ami a központi szerverre való csatlakozáskor az authentikációhoz szükséges) - környezeti változó: PSK

   - Apply(): error függvény: a konfigurációban található adatokat alkalmazza
   - setupDb(): error függvény: a konfigurációban található adatbázisra csatlakozik (az Apply metódus használja)
   - ToAddress(): string függvény: visszaadja a konfiguráció alapján a címet, amire a webszerver hallgat

##### Config függvényei

1. CreateApplicationConfig(): Config függvény: beolvassa a környezeti változókat és létrehozza a Config struktúrát, majd visszaadja azt

##### Config használata

A program indulásakor létre kell hozni a konfigurációt a CreateApplicationConfig() segítségével, majd alkalmazni kell azt az Apply() segítségével.

#### Db

##### Db konstansai

1. ENV_DB_USERNAME - annak a környezeti változónak a neve, amiben az adatbázis csatlakozásra használt felhasználónév van
2. ENV_DB_PASSWORD - annak a környezeti változónak a neve, amiben az adatbázis csatlakozásra használt jelszó van
3. ENV_DB_PROTOCOL - annak a környezeti változónak a neve, amiben az adatbázis csatlakozásra használt protokoll van
4. ENV_DB_HOST - annak a környezeti változónak a neve, amiben az adatbázis ip címe van
5. ENV_DB_PORT - annak a környezeti változónak a neve, amiben az adatbázis portja van
6. ENV_DB_NAME - annak a környezeti változónak a neve, amiben az adatbázis neve van

##### Db struktúrái

1. DBConfig
   - username: string (az adatbázis csatlakozáshoz használt felhasználónév)
   - password: string (az adatbázis csatlakozáshoz használt jelszó)
   - protocol: string (az adatbázis csatlakozáshoz használt protokoll)
   - host: string (az adatbázis ip címe)
   - port: string (az adatbázis portja)
   - dbName: string (az adatbázis neve)

   - ToConnectionString(): string függvény: a DBConfig struktúrában tárolt adatok alapján az sql drivernek megfelelő csatlakozási stringet adja vissza

##### Db függvényei

1. CreateDatabaseConfig(): DBConfig függvény: kiolvassa a környezeti változókból az adatokat, majd visszaadja a struktúrát, amit előállít belőlük

2. ConnectToConfig(DBConfig): error függvény: paraméterként átvesz db konfigurációt és csatlakozik az adatbázishoz a benne található adatokkal, visszaad bármilyen esetlegesen felmerülő errort

3. Disconnect() függvény: zárja az adatbázis kapcsolatot

##### Db használata

A db package-t a config package használja. A CreateDatabaseConfig segítségével létrehoz egy DBConfig struktúrát, majd az alapján csatlakozik az adatbázisra.

#### Model

##### Model általános feltételezései

A modellek lekérdezésekkor általában csak a neveket adják vissza a joinolt táblákból, az azonosítókat nem feltétlen (pl.: nyitvatartási időnél a hét napjai, felhasználó rangja).
Feltöltéskor viszont pont az azonosítókat várják, a nevek megadása nem szükséges.
A Validate kezdetű függvények úgynevezett validátor függvények, adatellenőrzést végeznek, és felhasználó számára megjeleníthető hibaüzenetet adnak vissza.
Amikor egy függvény az összes adatbázisban található adatot visszaadja, a töröltek nem számítanak bele.
A kapcsolótáblák töltéséhez, vagy az azokban található adatok lekérdezéséhez a modellek segédfüggvényeket használnak.
A kapcsolótáblák töltése tranzakciókon belül történik, hogy minden adat biztosan mentésre kerüljön, hiba esetén pedig semmi.

##### Model konstansai

1. DAYS_A_WEEK: a hét napjainak száma
2. MAX_WORK_HOURS_A_DAY: egy nap alatt maximálisan dolgozható órák száma
3. MAX_WEEKLY_HOURS: egy héten maximálisan dolgozható órák száma
4. PASSWORD_HASH_COST: a bcrypt algoritmus costja
5. HOURS_A_DAY: egy nap óráinak száma
6. EXPIRES_IN_DAYS: ennyi nap az alapértelmezett lejárati ideje a sessionnek

##### Model struktúrái

1. WeekDay:
   - Id: int (a hét napjának azonosítója)
   - Name: string (a hét napja)

2. Role:
   - Id: int (a rang azonosítója)
   - Name: string (a rang neve)

3. StoreType:
   - Id: int (az üzlettípus azonosítója)
   - Name: string (az üzlettípus neve)

4. StoreDetail:
   - Address: string (az üzlet címe)
   - CentralServerAddress: string (a központi szerver címe)
   - CentralServerPort: uint16 (a központi szerver portja)
   - StoreTypeId: int (az aktuális üzlettípus azonosítója)
   - StoreTypeName: string (az aktuális ülettípus neve)

5. OpenHour:
   - Id: int (a nyitvatartási idő azonosítója)
   - OpensAt: string (a nyitvatartási idő kezdete HH:mm formátumban) -- ha ez a kettő 00:00 --
   - ClosesAt: string (a nyitvatartási idő vége HH:mm formátumban) -- akkor aznap az üzlet zárva --
   - WeekDayIds: []int (azon napok azonosítója, amikor érvényes a nyitvatartás, feltöltéskor használatos)
   - WeekDays: []string (azon napok neve, amikor érvényes a nyitvatartás, lekérdezéskor használatos)
   - DeletedAt: sql.NullTime (a nyitvatartási idő törlésének időpontja, vagy null a nem töröltek esetén)

   - InsertNewOpenHour(): error függvény: menti az adatbázisba az új nyitvatartási időt, visszaadja az esetlegesen felmerülő hibát
   - UpdateOpenHour(): error függvény: módosítja az adatbázisban a nyitvatartási időt, visszaadja az esetlegesen felmerülő hibát
   - DeleteOpenHour(): error függvény: törli a nyitvatartási időt, visszaadja az esetlegesen felmerülő hibát
   - ValidateInsert(): error függvény: ellenőrzi, hogy a tárolt adatok alkalmasak-e az adatbázisba mentésre
   - ValidateUpdate(): error függvény: ellenőrzi, hogy a tárolt adatokra módosítható-e az adatbázisban tárolt nyitvatartási idő
   - ValidateDelete(): error függvény: ellenőrzi, hogy minden szükséges adat megérkezett-e a nyitvatartási idő törléséhez

   Azt, hogy melyik nyitvatartási idő melyik napokon érvényes, az open_day nevű kapcsolótábla tárolja.

6. ContractType:
   - Id: int (a szerződéstípus azonosítója)
   - Name: string (a szerződéstípus neve)
   - WeeklyHours: int (a szerződéstípushoz tartozó heti munkaórák száma)
   - DeletedAt: sql.NullTime (a szerződéstípus törlésének időpontja, vagy null a nem töröltek esetén)

   - InsertNewContractType(): error függvény: beszúrja az adatbázisba az új szerződéstípust, visszaadja a felmerülő esetleges hibát
   - UpdateContractType(): error függvény: frissíti a meglévő szerződésttípust, visszaadja a felmerülő esetleges hibát
   - DeleteContractType(): error függvény: törli a meglévő szerződéstípust, visszadja a felmerülő esetleges hibát
   - ValidateInsert(): error függvény: ellenőrzi, hogy a tárolt adatok alkalmasak-e az adatbázisba mentésre
   - ValidateUpdate(): error függvény: ellenőrzi, hogy a tárolt adatok alkalmasak-e, hogy arra legyen frissítve az adatbázis
   - ValidateDelete(): error függvény: ellenőrzi, hogy megtalálható-e az összes szükséges adat a szerződéstípus törléséhez

7. User:
   - Id: int (a felhasználó azonosítója)
   - Firstname: string (a felhasználó keresztneve)
   - Lastname: string (a felhasználó vezetékneve)
   - Username: string (a felhasználónév, amit bejelentkezéshez szükséges)
   - Password: string (a felhasználó jelszava, csak létrehozáskor és jelszómódosításkor használt mező)
   - PasswordChanged: bool (azt tárolja, hogy a felhasználó megváltoztatta-e már a jelszavát, ami egyszer kötelező)
   - PhoneNumber: string (a felhasználó telefonszáma, opcionális)
   - Email: string (a felhasználó email címe, opcionális)
   - Role: string (a felhasználó rangja)
   - RoleId: int (a felhasználó rangjának azonosítója)
   - DeletedAt: sql.NullTime (a felhasználó törlésének időpontja, vagy null a nem töröltek esetén)

   - InsertNewUser(): error függvény: menti az adatbázisba az új felhasználót (fontos: feltételezi, hogy a felhasználó jelszava már hashelve van), visszaadja az esetlegesen felmerülő hibát
   - UpdateUser(): error függvény: frissíti a felhasználó adatait, visszaadja az esetlegesen felmerülő hibát
   - DeleteUser(): error függvény: törli a felhasználót, visszaadja az esetlegesen felmerülő hibát
   - ValidateInsert(): error függvény: ellenőrzi az adatok helyességét az adatbázisba mentés előtt
   - ValidateUpdate(): error függvény: ellenőrzi az adatok helyességét az adatbázis frissítése előtt
   - ValidateDelete(): error függvény: ellenőrzi, hogy minden adat megérkezett-e a felhasználó törléséhez, és törölhető-e a felhasználó
   - UpdatePassword(): error függvény: frissíti a felhasználó jelszavát az adatbázisban (fontos: feltételezi, hogy a jelszó hashelve van), visszaadja az esetlegesen felmerülő hibát
   - EncryptPassword(): error függvény: hasheli a felhasználó jelszavát és visszaadja az esetlegesen felmerülő hibát
   - LogoutUser(): error függvény: kijelentkezteti a felhasználót (lejártá teszi a sessionjét) és visszaadja az esetlegesen felmerülő hibát

8. Session:
   - Id: int (a session azonosítója)
   - UserId: int (azon felhasználó azonosítója, akihez a session tartozik)
   - Token: string (a sessiont azonosító token, a kliens oldalon cookie-ban tárolódik)
   - ExpiresAt: time.Time (a session lejárati ideje, minden felhasználói interakció után frissül, 1 hét inaktivitás után jár le)

   - InsertNewSession(): error függvény: az adatbázisba menti az új sessiont, visszaadja az esetlegesen felmerülő hibát
   - UpdateExpiry(): error függvény: frissíti a session lejárati idejét 1 héttel a hívás utánra
   - ChangeExpiredToNew(): error függvény: a lejárt sessiont cseréli egy újra (új tokennel és frissített lejárati idővel)
   - setNewExpiresAt() függvény: segédfüggvény, a többi ennek segítségével állít be új lejárati időt

9. ContractDay:
   - Id: int (a szerződésnap azonosítója)
   - StartingTime: string (a szerződéshez tartozó, a hét adott napján teljesítendő munkaidő kezdete)
   - EndingTime: string (a szerződéshez tartozó, a hét adott napján teljesítendő munkaidő vége)
   - WeekDayId: int (a hét adott napjának azonosítója)
   - WeekDay: string (a hét adott napja)
   - DeletedAt: sql.NullTime (a szerződésnap törlésének időpontja, vagy null a nem töröltek esetén)

10. Contract:
    - Id: int (a szerződés azonosítója)
    - UserName: string (annak a felhasználó neve, akié a szerződés)
    - UserId: int (annak a felhasználó azonosítója, akié a szerződés)
    - ContractType: string (a szerződéshez tartozó szerződéstípus)
    - ContractTypeId: int (a szerződéshez tartozó szerződéstípus azonosítója)
    - Salary: int (a szerződéshez jára fizetés)
    - Filename: sql.Nullstring (a szerződéshez esetlegesen feltöltött PDF fájl elérési útja, vagy null nem feltöltött PDF esetén)
    - StartsAt: time.Time (a szerződés kezdete)
    - EndsAt: sql.Nulltime (a szerződés vége határozott idejű szerződésnél, null ha határozatlan)
    - ContractDays: []ContractDay (a szerződéshez tartozó szerződésnapok)
    - DeletedAt: sql.NullTime (a szerződés törlésének időpontja, vagy null a nem töröltek esetén)

    - InsertNewContract(): error függvény: menti az adatbázisba az új szerződés adatait
    - UpdateContract(): error függvény: frissíti az adatbázisban szereplő szerződés adatait
    - DeleteContract(): error függvény: törli a szerződést
    - ValidateInsert(): error függvény: ellenőrzi, hogy a tárolt adatok menthetők-e az adatbázisba
    - ValidateUpdate(): error függvény: ellenőrzi, hogy lehet-e az új adatokra frissíteni az adatbázist
    - ValidateDelete(): error függvény: ellenőrzi, hogy az öszes adat megérkezett-e a szerződés törléséhez
    - DeleteContractFileFromDB(): error függvény: a szerződés PDF verzójának törlése esetén kiveszi a fájl elérési útját az adatbázisból, visszaadja az esetlegesen felmerülő hibát

    A szerződéshez tartozó szerződésnapokat a contract_day kapcsolótábla tárolja.

##### Model függvényei

1. GetAllWeekDay(): []WeekDay, error függvény: visszaadja az adatbázisban található hét napjait és az esetlegesen felmerülő hibát

2. GetAllRole(): []Role, error függvény: visszaadja az adatbázisban található rangokat és az esetlegesen felmerülő hibát

3. GetAllStoreType: []StoreType, error függvény: visszaadja az adatbázisban található összes üzlettípust és az esetlegesen felmerülő hibát

4. GetStoreDetail: StoreDetail, error függvény: visszaadja az adatbázisban található adatait a üzletnek (csak 1 sort tárol az adatbázis, ezért az mindig létezik alapértelmezett értékekkel, ha még nincs beállítva)

5. GetOpenHours(bool): []OpenHour, error függvény: paraméterként átveszi, hogy csak az adott évi és következő évi nyitvatartási időt adja-e vissza (ha nem, akkor az összeset), és visszaadja a nyitvatartási időket, vagy bármilyen menet közben felmerülő hibát

6. getWeekDaysForOpenHour(int): []string, error függvény: segédfüggvény, paraméterként átveszi egy nyitvatartási idő azonosítóját, majd visszaadja a hozzá tartozó hét napjait és az esetlegesen felmerülő hibát

7. insertWeekDaysForOpenHour(int, []int, \*sql.Tx): error függvény: segédfüggvény, paraméterként átveszi a nyitvatartási idő azonosítóját és a hozzá tartozó hét napjai azonosítókat, valamint a tranzakciót, aminek részeként ezek mentésre kerülnek, menti őket az open_day kapcsolótáblába és visszaadja a menet közben esetlegesen felmerülő errort

8. GetAllContractType(): []ContractType, error függvény: visszaadja az adatbázisban tárolt szerződéstípusokat és az esetlegesen felmerülő hibát

9. GetUserByUsername(string): User, error függvény: paraméterként átveszi egy felhasználó nevét és visszaadja a felhasználót, vagy a menet közben felmerült hibát (bejelentkezéskor használatos)

10. GetUserByUserId(int): User, error függvény: paraméterként átveszi egy felhasználó azonosítóját és visszaadja a felhasználót, vagy a menet közben felmerült hibát

11. GetAllUser: []User, error függvény: visszaadja az adatbázisban található összes felhasználót és az esetlegesen felmerült hibát

12. CheckPassword(string, string): bool függvény: paraméterként átveszi a begépelt és a hashelt jelszót, visszaadja, hogy a jelszavak egyeznek-e

13. GetSessionByUserId(int): Session, error függvény: paraméterként átvesz egy felhasználói azonosítót és visszaadja a felhasználóhoz tartozó sessiont (akkor is ha lejárt, hogy frissítve legyen ha kell), és az esetlegesen felmerülő hibát

14. GetSessionByToken(string): Session, error függvény: paraméterként átvesz egy session tokent és visszaadja a hozzá tartozó sessiont (csak ha érvényes), vagy az esetlegesen felmerülő hibát

15. GetContractByUserId(int): Contract, error függvény: paraméterként átvesz egy felhasználói azonosítót, majd visszaadja a hozzá tartozó szerződést és az esetlegesen felmerülő hibát

16. getContractDaysForContract(id): []ContractDay, error függvény: segédfüggvény, paraméterként átveszi a szerződés azonosítóját, majd visszaadja a szerződéshez tartozó szerződésnapokat és az esetlegesen felmerülő hibát

17. insertContractDaysForContract(int, []ContractDay, \*sql.Tx): error függvény: segédfüggvény, paraméterként átveszi a szerződés azonosítóját, a szerződéshez tartozó szerződésnapokat és a tranzakciót, majd meni a contract_day kapcsolótáblába az adatokat és visszaadja az esetlegesen felmerülő hibát

18. validaContractDays([]ContractDay): error függvény: paraméterként átveszi a szerződésnapokat, és egyesével ellenőrzi, hogy az adatok megfelelnek-e a követelményeknek és menthetők-e az adatbázisba

##### Model használata

A model package az adatbázisban található adatok go-ban leképzése, ezért alkalmazásszerte használatos.
Elsősorban a handler-ökben és middleware-kben.

#### Auth

##### Auth konstansai

1. A [role.go] fájlban megtalálható az összes rang neve és azonosítója konstansként.
2. ErrUserNotFound: hibás felhasználónév esetén használt error
3. ErrBadPassword: hibás jelszó esetén használt error

##### Auth függvényei

1. generateToken(int): string függvény: segédfüggvény: paraméterként átvesz egy hosszt, létrehoz egy random stringet azzal a hosszal, base64 urlencoding segítségével átalakítja azt, majd ezt visszaadja (session tokenek generálása)

2. LoginUser(string, string, \*context.Context): error függvény: paraméterként átveszi a felhasználónevet és a begépelt jelszót, és egy kontextust, ellenőrzi az adatok helyességét és annak megfelelően elutasítja a bejelentkezést (hibát dob), vagy kezeli a session, visszaadja az esetlegesen felmerült hibát (fontos: visszaadhat ErrUserNotFound és ErrBadPassword errort, amire külön ellenőrizni kell, mert nem valódi hibák)

3. CreateOrUpdateSessionForUser(int, \*context.Context): error függvény: paraméterként átveszi a felhasználó azonosítóját és egy kontextust, majd vagy létrehoz egy sessiont, vagy ha már volt, akkor frissíti azt, visszaadja a menet közben felmerülő hibát (a LoginUser függvény használja)

4. createSessionForUser(int, \*context.Context): error függvény: segédfüggvény, paraméterként átveszi a felhasználó azonosítóját és egy kontextust, majd létrehoz egy sessiont (a CreateOrUpdateSessionForUser hívja meg kizárolag akkor, ha adott felhasználónak még sosem volt sessionje, egyébként azt frissíti, vagy cseréli a tokenét a lejártsági állapottól függően)

5. CanUserSetRole(model.User, int): bool függvény: paraméterként átvesz egy felhasználót és egy rang azonosítót, és visszaadja, hogy adott felhasználónak van-e joga a rangot beállítani (üzletvezető rangot csak üzletvezető állíthat be)

6. CanUserDisablePasswordChange(model.User): bool függvény: paraméterként átvesz egy felhasználót és visszaadja, hogy van-e joga kikapcsolni a kötelező jelszó módosítást (csak az üzletvezető teheti meg)

##### Auth működése

Az authentikáció session-ökkel történik, amik az adatbázisban tárolódnak. A kliens cookie-kban kapja meg a session tokeneket.

A package egy bejelentkezni kívánó felhasználó adatait kezeli. A LoginUser meghívása indítja a folyamatot, ami ellenőrzi, hogy adott felhasználónévvel létezik-e felhasználó, hibát dob ha nem. Majd ellenőrzi, hogy helyes-e a jelszó, hibát dob, ha nem. Utána a session kezelése történik. Ha van érvényes session, akkor frissül a lejárati ideje. Ha van session, de már lejárt, akkor új session tokent kap és frissül a lejárati ideje. Ha sosem volt sessionje a felhasználónak, akkor új session jön létre.

A kontextus átadása azért szükséges, mert abba kerül mentésre a felhasználó, valamint a session, amit külsős packagek majd kiolvasnak.

#### Middleware

##### Middleware függvényei

1. AuthenticateUser: ez egy az echo keretrendszernek megfelelő middleware függvény. Ellenőrzi az auth_token cookie alapján, hogy a felhasználó be van-e jelentkezve (van-e érvényes session-je). Ha nem, akkor átirányítja a kérést a bejelentkező oldalra. Ha be van jelentkezve a felhasználó, akkor frissíti a sessionjét, lementi a felhasználót a kontextusba és tovább engedi a kérést.

2. AuthorizeStoreLeaderOrHr: ez egy az echo keretrendszernek megfeleő middleware függvény. Csak az AuthenticateUser után fut, és ellenőrzi, hogy a kontextusba mentett felhasználó üzletvezető, vagy hr ranggal rendelkezik-e. Ha nem, a kérést elutasítja, ha igen, akkor továbbengedi.

3. AuthorizeStoreLeader: ez egy az echo keretrendszernek megfeleő middleware függvény. Csak az AuthenticateUser után fut, és ellenőrzi, hogy a kontextusba mentett felhasználó üzletvezető ranggal rendelkezik-e. Ha nem, a kérést elutasítja, ha igen, akkor továbbengedi.

##### Middleware használata

A middleware függvényeket az API végpontok regisztrálásakor használja az alkalmazás, hogy minden műveletet csak az arra jogosulak hajthassanak végre.

#### Route

##### Route konstansai

1. CONTRACT_FOLDER: a mappa neve, ahova a feltöltött szerződés PDF fájlok kerülnek

##### Route struktúrái

1. LoginRequest (bejelentkezéskor érkező adatok fogadására alkalmas struktúra)
   - Username: string
   - Password: string

2. ConnectInfo (központi szerverre csatlakozáskor használt adatok fogadására alkalmas struktúra)
   - Ip: string (a központi szerver címe, lehet domain is)
   - Port: string (a központi szerver portja)
   - Psk: string (a központi szerveren használt authentikációhoz szükséges PSK)

3. PasswordUpdate (jelszó módosításhoz szükséges adatok fogadására alkalmas struktúra)
   - OldPassword: string
   - NewPassword: string

4. PasswordReset (jelszó visszaállításához szükséges adatok fogadására alkalmas struktúra)
   - UserId: int
   - Password: string
   - PasswordConfirm: string

##### Route függvényei

Alapértelmezetten az echo keretrendszernek megfelelő handler függvények találatók a package-ben, ezért a function signature csak akkor lesz kiírva, ha attól eltér.

1. HandleLogin: kezeli a bejelentkezési kérelmet, a sessiont és a cookiet is (API: POST /api/login)

2. createLoginErrorCodeAndMessage(err): int string függvény: segédfüggvény, paraméterként átvesz egy errort (amit a auth.LoginUser függvény ad vissza) és előállítja az az alapján a kliensnek visszaküldendő HTTP válasz kódot és hibaüzenetet

3. SetSessionCookie(echo.Context, model.Session) függvény: paraméterként átveszi az echo kontextusát és a sessiont, és beállítja a cookiet

4. HandleLogout: kijelentkezteti a felhasználót (API: GET /api/logout)

5. HandleGetEcho: a frontend használja a bejelentkezés ellenőrzésére (ha visszaad 204 NOCONTENT-et, akkor a felhasználó be van jelentkezve, egyébként át kell irányítani a felhasználót a bejelentkezésre) (API: GET /api/echo)

6. HandleGetContractByUserId: visszaadja a felhasználó szerződését (API: GET /api/contract)

7. HandlePostContract: menti az új szerződést az adatbázisba (API: POST /api/contract)

8. saveContractFile(\*multipart.FileHeader): error függvény: segédfüggvény, menti az esetlegesen feltöltött PDF fájlt

9. deleteContractFile(model.Contract): error függvény: segédfüggvény, paraméterként átvesz egy szerződést és törli a hozzá tartozó PDF fájlt

10. HandleUpdateContract: frissíti egy szerződés adatait (API: PUT /api/contract)

11. HandleDeleteContract: töröl egy szerződést (API: DELETE /api/contract)

12. HandleGetContractFile: visszaküldi a szerződéses fájlt (API: GET /api/contract-file)

13. HandleDeleteContractFile: tölri a szerződéses fájlt (API: DELETE /api/contract-file)

14. HandleGetAllContractType: visszaadja az összes szerződéstípust (API: GET /api/contract-type)

15. HandlePostContractType: ment egy új szerződéstípust (API: POST /api/contract-type)

16. HandleUpdateContractType: frissít egy szerződéstípust (API: PUT /api/contract-type)

17. HandleDeleteContractType: töröl egy szerződéstípust (API: DELETE /api/contract-type)

18. HandleGetOpenHours: visszaadja a nyitvatartási időket (API: GET /api/open-hour)

19. HandlePostOpenHour: menti az új nyitvatartási időt (API: POST /api/open-hour)

20. HandleUpdateOpenHour: frissít egy nyitvatartási időt (API: PUT /api/open-hour)

21. HandleDeleteOpenHour: töröl egy nyitvatartási időt (API: DELETE /api/open-hour)

22. HandleGetAllRole: visszaadja az összes rangot (API: GET /api/role)

23. HandleGetStoreDetail: visszaadja az üzlet adatait (API: GET /api/store-detail)

24. HandleUpdateStoreDetail: frissíti az üzlet adatait (API: PUT /api/store-detail)

25. HandleGetAllStoreType: visszaadja az összes üzlettípust (API: GET /api/store-type)

26. HandlePostUser: új felhasználót ment (API: POST /api/user)

27. HandleUpdateUser: frissít egy felhasználót (API: PUT /api/user)

28. HandleDeleteUser: töröl egy felhasználót (API: DELETE /api/user)

29. HandleGetUser: visszaadja az aktuális felhasználót (API: GET /api/user)

30. HandleGetAllUser: visszaadja az összes felhasználót (API: GET /api/all-user)

31. HandleUpdateUserPassword: kezeli a jelszófrissítést (API: POST /api/password-change)

32. HandlePasswordReset: kezeli a jelszó visszaállítást (API: POST /api/password-reset)

33. HandleGetWeekDays: visszaadja a hét napjait (API: GET /api/weekdays)

34. HandleGetConnect: visszaad egy üzenet, ami tartalmazza, hogy a rendszer éppen csatlakoztatva van-e a központi szerverhez, vagy nem

35. HandlePostConnect: beérkező csatlakozási adatok alapján kapcsolódik a központi szerverre (új NetworkManager létrehozása)

36. HandlePostNetworkSearch: amennyiben a rendszer csatlakozik a központi szerverre, keresési kérést indít, és visszaadja annak válaszát (vagy időtúllépés hibát)

#### Network

##### Network konstansai

1. VERSION: a protokoll szerinti verzió (a protokoll verziója, amit a kliens használ)
2. VERSION_SIZE: a verziót jelentő bytok száma
3. MESSAGE_LENGTH_SIZE: az üzenet hosszát jelentő byteok száma
4. HEADER_SIZE: a protokoll szerinti fejléc hossza
5. TIMEOUT_IN_SECONDS: az az időmennyiség, amennyit a rendszer vár a hálózatos keresés eredményére
6. STATUS_NOT_CONNECTED: a kliens nincs csatlakozva a központi szerverre
7. STATUS_CAN_SEARCH: a kliens csatlakozva van a központi szerverhez, és indíthat keresési kérést
8. STATUS_WAITING_FOR_ANSER: a kliens éppen egy hálózatos keresés eredményére vár (nem indíthat újat)
9. A protokoll szerinti üzenettípusok konstansai is megtalálhatók a package-ben
10. MSG_EOF: a protokoll szerinti EOF byte
11. UUID_LENGTH: a string formátumú UUID hossza

##### Network struktúrái

1. A protokollban meghatározott üzenettípusok és azok típusai megtalálhatók a központi szerverhez hasonló módon

2. TcpHeader: a központi szerverhez hasonló módon

3. Connection:
   - Conn: net.Conn (a mögöttes beépített kapcsolat a központi szerverhez)
   - ServerAnswerChan: chan Message (a csatorna, ahova az API végpont által küldött keresésre beérkezett üzenetet küldeni kell)
   - mutex: \*sync.Mutex (a megosztott erőforássok használata során ne keletkezzen race condition az egyes szálak között)

   - Authenticate(string): error függvény: paraméterként átveszi a PSK-t, és lebonyolítja a protokoll szerinti authentikációt a központi szerverrel, visszaküldi az esetlegesen felmerülő hibát
   - ReadLoop() függvény: külön szálon elkezd hallgatni a szervertől érkező üzenetekre
   - ReadHeader(): \*TcpHeader, error, error: kiolvas egy protokoll szerinti fejlécet, amit visszaad, vagy egy protokoll hibát, vagy egy hálózati hibát, ha történt
   - ReadPayload(uint32): \*TcpMessage függvény: paraméterként átveszi az üzenet hosszát, majd visszaadja az abból előállított TcpMessage struktúra pointerért, vagy nilt, hiba esetén
   - HandleMessage(\*TcpMessage) függvény: paraméterként átvesz egy üzenetet, átalakítja a megfelelő specifikus üzenettípussá, majd továbbküldi feldolgozásra
   - GetSearchResults(\*ClientSearchMessage) függvény: paraméterként átveszi a szervertől kapott keresési kérést, lekéri az eredményt, létrehozza belőle az AnswerMessage struktúrát és elküldi az üzenetet a központi szervernek (fontos: külön szálon hívandó)
   - GiveServerAnser(Message) függvény: paraméterként átvesz egy üzenetet, amit az API végpontnak vissza kell adnia
   - SendSearchRequest(\*SearchMessage, chan Message) függvény: paraméterként átvesz egy keresési kérés üzenetet, valamint a csatornát, ahova az arra érkezett választ küldeni kell (GiveAnswer függvény)
   - SendMessage(Message): error függvény: paraméterként átvesz egy üzenetet, és elküldi a központi szervernek, visszaadja a menet közben felmerülő hibát
   - write([]byte): error függvény: segédfüggvény, paraméterként átveszi a byteokat, amik az üzenetet képzik, elküldi a központi szervernek és visszaadja a menet közben felmerülő errort (SendMessage használja)
   - Close() függvény: bezárja a kapcsolatot

4. SearchResult (a hálózatos keresési kérésre kapott válasz struktúrája):
   - OpenHours: []model.OpenHour (a választ adó üzlet nyitvatartási ideje)
   - StoreDetail: model.StoreDetail (a választ adó üzlet adatai)
   - Products: bármi (a komplex szűrés végpont által visszaküldött válasz)

5. NetworkManager:
   - Connection: \*Connection (a manager által használt kapcsolat)
   - Status: int (a manager státusza)
   - psk: string (a manager által az authentikációhoz használt PSK)
   - mutex: \*sync.Mutex

   - StartReadLoop() függvény: elkezdi a hálózati hallgatást
   - IsConnected(): bool függvény: visszaadja thread-safe módon, hogy a manager csatlakozik-e a központi szerverre
   - SearchNetwork([]byte): []byte, error függvény: paraméterként átveszi a komplex szűrés paramétereit (JSON blob), elküldi a kérést a szervernek, majd vár meghatározott időt a válasza, amit ha megkap, visszaküldi az API válaszaként, egyébként időtúllépés hibaüzenetet ad (API hívás használja)
   - GetSearchResults([]byte): []byte, error függvény: paraméterként átveszi a komplex szűrés paramétereit (JSON blob), lekéri a nyitvatartási időt és az üzlet adatait, majd visszaküldi a központi szervernek az eredményeket (központi szerver kérése hívja meg)
   - CallApi([]byte): any, error függvény: paraméterként átveszi a komplex szűrés paramétereit (JSON blob), lekérdezi a termék adminisztráció felület végpontjától az adatot (auth_token cookienak használva a psk-t), majd visszaadja a választ vagy a menet közben felmerülő errort
   - Disconnect() függvény: a manager státuszát átállítja nem csatlakoztatottra (a valós kapcsolat bontása esetén hívódik meg, pl.: hálózati hiba esetén)

##### Network használata

Létre kell hozni egy managert a NewNetworkManager() használatával, majd elkezdeni a readloopot. Az IsConnected függvényt egy API hívás kezelője használja, hogy megállapítsa, csatlakozva van-e a központi szerverre a rendszer. Ha az API hívás hálózatos keresést indít, a SearchNetwork() használandó, aminek paraméterként át kell adni a komplex szűrés paramétereit képző JSON blobot. Ha a központi szerver küld megválaszolandó keresési kérést, azt kezeli a GetSearchResults() függvény.

### Általános API szabályok

Új erőforrás létrehozásánal az azonosító értelemszerűen nem kerül átadásra.
Meglévő erőforrás módosításánál az azonosító és az összes többi adat átadásra kell kerüljön.
Erőforrás törléséhez elég az azonosító, de átadható a többi adat is.
Minden erőforrás soft-delete formájában törlődik, kivéve a kapcsolótáblák általt tárolt adatokat, ha azok módosításra kerültek. Ha törölték a kapcsolódó adatot, akkor soft-delete történik a kapcsolótáblában is, ha viszont csak módosították, akkor az előző verzióhoz tartozó adatok hard-delete módon törlődnek, míg az újak bekerülnek.
Soft-delete esetén az esetlegesen kapcsolódó fájlok nem törlődnek. Fájlokat csak az explicit fájl törlésére szolgáló végpont törli.

Az AuthenticateUser middleware az összes végpont előtt lefut, kivéve a logint, és kiolvassa az auth_token cookiet, ezért a táblázatba nem került be.
Ennek megfelelően bármelyik végpont adhat vissza 401-es státuszkódot, vagy átirányítást a /login oldalra.
Az összes végpont előtt szerepel a /api prefix, ez a táblázatba nem került be.
A JSON üzenet a következő formában van:

- hiba esetén: {"error": "<hibaüzenet>"}
- siker esetén : {"message": "<siker üzenet>"}

Egyéb JSON body esetén a modellben definiált adatok kerülnek átadásra.

### API végpontok táblázat

| Metódus | Végpont                            | Request body                          | Response kód       | Response body                       | Middleware               |
| ------- | ---------------------------------- | ------------------------------------- | ------------------ | ----------------------------------- | ------------------------ |
| POST    | /login                             | JSON LoginRequest                     | 400, 500, 200      | JSON üzenet                         | -                        |
| GET     | /logout                            | -                                     | 500, 200           | JSON üzenet                         | -                        |
| GET     | /contract-type                     | -                                     | 500, 200           | JSON üzenet / szerződéstípusok      | AuthorizeStoreLeaderOrHr |
| POST    | /contract-type                     | JSON szerződéstípus                   | 400, 500, 201      | JSON üzenet                         | AuthorizeStoreLeader     |
| PUT     | /contract-type                     | JSON szerződéstípus                   | 400, 500, 200      | JSON üzenet                         | AuthorizeStoreLeader     |
| DELETE  | /contract-type                     | JSON szerződéstípus                   | 400, 500, 204      | JSON üzenet / NoContent             | AuthorizeStoreLeader     |
| GET     | /store-type                        | -                                     | 500, 200           | JSON üzenet / üzlettípusok          | AuthorizeStoreLeader     |
| GET     | /store-detail                      | -                                     | 500, 200           | JSON üzenet / üzlet adatok          | AuthorizeStoreLeader     |
| PUT     | /store-detail                      | JSON üzletadatok                      | 400, 500, 200      | JSON üzenet                         | AuthorizeStoreLeader     |
| GET     | /user                              | -                                     | 200                | JSON felhasználó                    | -                        |
| GET     | /all-user                          | -                                     | 500, 200           | JSON üzenet / felhasználók          | AuthorizeStoreLeaderOrHr |
| POST    | /user                              | JSON felhasználó                      | 400, 500, 201      | JSON üzenet                         | AuthorizeStoreLeaderOrHr |
| PUT     | /user                              | JSON felhasználó                      | 400, 500, 200      | JSON üzenet                         | AuthorizeStoreLeaderOrHr |
| DELETE  | /user                              | JSON felhasználó                      | 400, 500, 204      | JSON üzenet / NoContent             | AuthorizeStoreLeaderOrHr |
| GET     | /role                              | -                                     | 500, 200           | JSON rangok                         | AuthorizeStoreLeaderOrHr |
| POST    | /password-change                   | JSON PasswordChange                   | 400, 500, 200      | JSON üzenet                         | -                        |
| POST    | /password-reset                    | JSON PasswordReset                    | 400, 500, 200      | JSON üzenet                         | AuthorizeStoreLeader     |
| GET     | /open-hour                         | -                                     | 500, 200           | JSON üzenet / nyitvatartási idők    | AuthorizeStoreLeader     |
| POST    | /open-hour                         | JSON nyitvatartási idő                | 400, 500, 201      | JSON üzenet                         | AuthorizeStoreLeader     |
| PUT     | /open-hour                         | JSON nyitvatartási idő                | 400, 500, 200      | JSON üzenet                         | AuthorizeStoreLeader     |
| DELETE  | /open-hour                         | JSON nyitvatartási idő                | 400, 500, 204      | JSON üzenet / NoContent             | AuthorizeStoreLeader     |
| GET     | /weekdays                          | -                                     | 500, 200           | JSON üzenet / hét napjai            | AuthorizeStoreLeader     |
| GET     | /contract?userId=<userId>          | -                                     | 400, 500, 204, 200 | JSON üzenet / szerződés / NoContent | AuthorizeStoreLeaderOrHr |
| POST    | /contract                          | Form: JSON szerződés + opcionális PDF | 400, 500, 201      | JSON üzenet                         | AuthorizeStoreLeaderOrHr |
| PUT     | /contract                          | Form: JSON szerződés + opcionális PDF | 400, 500, 200      | JSON üzenet                         | AuthorizeStoreLeaderOrHr |
| DELETE  | /contract                          | JSON szerződés                        | 400, 500, 204      | JSON üzenet / NoContent             | AuthorizeStoreLeaderOrHr |
| GET     | /contract-file?filename=<filename> | -                                     | 400, 404, 200      | PDF fájl                            | AuthorizeStoreLeaderOrHr |
| DELETE  | /contract-file                     | JSON szerződés                        | 500, 204           | JSON üzenet / NoContent             | AuthorizeStoreLeaderOrHr |
| GET     | /echo                              | -                                     | 204                | NoContent                           | -                        |
| GET     | /connect                           | -                                     | 200                | JSON üzenet                         | AuthorizeStoreLeader     |
| POST    | /connect                           | JSON ConnectInfo                      | 400, 500, 200      | JSON üzenet                         | AuthorizeStoreLeader     |
| POST    | /network-search                    | JSON komplex szűrés paraméter         | 400, 500, 200      | JSON üzenet / keresési eredmény     | -                        |

## Frontend

### Felhasznált technológiák

1. Vuejs keretrendszer - a teljes UI-hoz
2. Vue-router - a navigációhoz
3. TailwindCSS - a dizájnhoz
4. Typescript programozási nyelv

### Projekt struktúra

Az összes forráskód egy vue projektnek megfelelő módon a src/ mappában van.

#### Types

Az alkalmazás által használt typescript típusok itt találhatók. Ezek megfelelnek a backenden található struktúráknak (a modelleknél), ezért a dokumentáció nem ismétli őket. A csak a frontenden használt típusok itt megtalálhatók.
A modellek esetén a típus mellett található egy osztály is, aminek konstruktora átvehet egy adott típusú objektumot. Ennek szerepe az adatok bekötése a UI-ba, ha kapnak paramétert, akkor módosítás céljából, egyébként új létrehozás céljából.
Ezen osztályokon létezik a to<Type>(): type függvény (pl.: toUser(): User), ami az osztályban található adatokból csinál egy adott típusú objektumot, ami JSON formátumban továbbítható a végpontok felé.
Az osztályokon szintén létezik a compare(object: type): boolean metódus (pl.: compare(user: User): boolean), ami megállapítja, hogy a kapott paraméter adatai megegyeznek-e az osztály adataival. Ezt a rendszer arra használja, hogy történt-e tényleges módosítás (tárolja az eredeti adatokat old<Változónév> néven, pl.: oldUser).

1. FeedbackType: a visszajelzés lehetséges típusai

2. Feedback:
   - type: FeedbackType
   - message: string

### Router

A router a vue-routernek megfelelő módon regisztrálja az utakat és a hozzájuk tartozó nézeteket.
Az App.vue megjeleníti a fejlécet (ami navigációs célt szolgál), valamint az éppen aktuális utat.

Tartalmaz egy middleware-t, ami a frontenden ellenőrzi a navigálás előtt, hogy a felhasználó be van-e jelentkezve, valamint cserélt-e már jelszót. Ha nincs bejelentkezve, a login oldalra viszi, ha nem cserélt jelszót, akkor a jelszó cseréhez, ha pedig minden rendben, továbbengedi.

### Components

A Views mappa elemei által használt komponensek.
Több helyen is megtalálható a következő függvény:

- validate(): {message: string, valid: boolean}
  Ez a függvény az adott komponens adait validálja az elküldésük előtt, a hibaüzenetet és a validálás sikerességét pedig visszaadja.

Az adatok kezelése általában két komponenst igényel:

1. Megjelnítő komponens (ez lehet kártya, vagy táblázat sor)
2. Adatbeállító (data) komponens (ez kezeli az adatokat módosításkor vagy mentéskor)

Ezen komponenseket a View kezeli.
Ezek a komponensek általános komponenseket (pl.: Feedback, SearchBar, Modal) használhatnak, akár többet is.

#### Feedback

Propként átvesz egy Feedback-et, és megjeleníti azt.
A láthatóságát a szülő komponens kezeli.

#### Header

Megjeleníti a navigációs menüt

#### Modal

Megjelenít egy felugró ablakot a propként átvett címmel, üzenettel, valamint a megerősítést jelentő gomb szövegével.
Emiteli, hogy Mégse vagy Megerősítés gombra kattintott a felhasználó.
A láthatóságát a szülő komponens kezeli.

#### SearchBar

Propsként átveszi a keresendő itemet, és biztosít egy input mezőt és egy gombot, aminek lenyomására emiteli a beírt keresési szöveget.
A keresés végrehajtása a szülő komponens dolga.

#### ContractTypes

##### ContractTypeCard

Propsként átvesz egy szerződéstípust, és megjelenít egy kártyát az adataival.
Emiteli a módosítási vagy törlési szándékot, aminek kezelése a szülőkomponens feladata.

##### ContractTypeData

Propsként átvesz egy szerződéstípus (módosítás) vagy nullt (új létrehozása), amihez biztosítja a szükséges felületet.
Az adatok mentését a komponens végzi.

#### OpenHour

##### OpenHourCard

Propsként átvesz egy nyitvatartási időt és megjeleníti az adatait.
Emiteli a módosítási vagy törlési szándékot, aminek kezelése a szülőkomponens feladata.

##### OpenHourData

Propsként átvesz egy nyitvatartási időt (módosítás) vagy nullt (új létrehozása), amihez biztosítja a felületet.
Az adatok mentését a komponens végzi.

#### Profile

##### PasswordChange

Felületet biztosít egy felhasználónak a saját jelszava cserélésére.
Az adatok mentését a komponens végzi.

#### Users

##### UserTable

Propsként átvesz egy felhasználó tömböt, amit táblázatosan megjelenít.
Emitel módosítási és törlési szándékot, de azokat ő maga is egy gyerek komponenstől kapja.

##### UserRow

A táblázat egy sorát képviseli, megjelníti a felhasználó adatait, amit propsként kap.
Emiteli a módosítási és törlési szándékot.

##### UserData

Propsként átvesz egy felhasználót (módosítás) vagy nullt (új létrehozása), amihez biztosítja a felületet.
Az adatok mentését a komponens végzi.

##### Contract

Propsként átvetsz egy userId-t, lekéri a hozzá tartozó szerződés adatait, amit a ContractData és ContractDay komponensek segítségével megjelenít. Ezek kezelését ez a komponens végzi, valamint a szerződéssel kapcsolatos adatkezelési műveleteket (mentés, törlés, fájl lekérdezése, törlése) is ez a komponens végzi.

##### ContractData

Propsként átveszi a szerződést, a szerződéstípusokat és a szerződés fájlt, ha létezik.
A szerződés összes adatának kezelését biztosítja, kivéve a szerződés napjait.
Az összes elvégezhető műveletet emiteli, amit aztán a Contract komponens kezel.

##### ContractDay

Propsként átveszi a hét napjait és a szerződés napjait, és felületet biztosít a szerződés napjainak beéllítására.
A lehetséges műveleteket emiteli, amiket a Contract komponens kezel.

### Views

#### ContractType

Lekéri a szerződéstípusokat, kezeli a Contract komponensek emitjeit, a feedbacket és a modalt is.
Ez a view végzi a szerződéstípushoz köthető adatkezelő műveleteket.

#### Home

A főoldal, csak a termék kezelés oldalra való átjutáshoz biztosít egy linket.

#### Login

A bejelentkezést biztosítja, ide kerül átirányításra a még nem bejelentkezett felhasználó.

#### OpenHour

Lekéri a nyitvatartási időket, kezeli az OpenHour komponensek emitjeit, a feedbacket és a modalt is.
Ez a view végzi a nyitvatartási időhöz köthető adatkezelő műveleteket.

#### Profile

Ez a nézet lekéri a felhasználó adatait, amit nem módosítható módon megjelenít.
Biztosít egy linket a jelszó cseréléshez.

#### StoreDetail

Lekéri, megjeleníti és kezelési felületet biztosít az üzlet adataihoz.
Az összes ehhez köthető műveletet, a központi szerverre való csatlakozással együtt ez a view intézi.

#### Users

Lekéri a felhasználókat, kezeli a feedbacket és modalt (a felhasználóhoz tartozót).
A Users komponenseket használja, vezérli mikor melyik látható (UserTable, UserData, vagy Contract).
Kezeli ezek emitjeit és a felhasználókhoz köthető adatkezelési műveleteket.
A szerződés modalját, feedbackjét és műveleteit a Contract komponens végzi.

## A docker konténer

Az alkalmazás minden része docker konténerben futásra van tervezve.
Ennek megfelelően megtalálható a projekt gyökérmappájában egy Dockerfile és .dockerignore.

### Frontend build

A docker először a frontendet építi ki:

1. Egy node:24-es docker image-ből indul ki.
2. Argumentumként átveszi azt a VITE\_ környezeti változót, amit a frontend használ a termék adminisztráció oldalának linkjeként.
3. A dockerben megszokott módon a /app könyvtárat jelöli meg kiindulóként, felmásolja a projektet, létrehozza a build tárolására szolgáló public mappát, majd belép a client mappába, amit buildel (npm i + npm run build), ennek az eredménye kerül a public mappába.
4. A build eredményére később frontend-build néven lehet hivatkozni.

### Backend build

A docker a backend buildelésével folytatja:

1. A golang:1.24-es docker image-ből indul ki.
2. A dockerben megszokott módon a /app könyvtárat jelöli meg kiindulóként, felmásolja a projektet, telepíti a függőségeket (go mod tidy), majd egy statikusan linkelt futtatható go binary-t buildel (CGO_ENABLED=0 go build -o main). Erre azért van szükség, hogy egy minimális alpine linuxon (ahol nincs libc) is fusson a projekt.
3. A build eredményére később backend-build néven lehet hivatkozni.

### A teljes összerakott projekt

Végül a két buildelt eredményt összesíti egy konténerben a docker:

1. Az alpine:latest image-ből indul ki.
2. A dockerben megszokott módon a /app könyvtárat jelöli meg kiindulóként.
3. Átmásolja a backend-buildből a go binary-t, a frontend-buildből a public mappát (amit a binary mellé kell helyezni, ugyanis a go backend azt a mappát szolgálja ki statikusan a webszerver részeként).
4. Létrehozza a contracts mappát, ahova a feltöltött szerződések fognak kerülni (szintén a binary mellett van, ugyanis a go backend ott keresi).
5. A konténer indító parancsaként pedig a binary-t indítja.
