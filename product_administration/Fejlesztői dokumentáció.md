# Netstore webalkalmazás - termék adminisztráció fejlesztői dokumentáció

## Backend

### Felhasznált technológiák

1. PHP programozási nyelv
2. MySQL adatbázis
3. Apache webszerver
4. Docker

### Modulok

#### Sql_functions

##### Sql_functions függvényei

1. getDbConnection() függvény: létrehozza és visszaadja a MySQL adatbázis kapcsolatot a környezeti változókban tárolt adatok alapján (DB_HOST, DB_USERNAME, DB_PASSWORD, DB_NAME). A kapcsolat UTF-8 karakterkódolást használ.

2. getData() függvény: paraméterként átvesz egy SELECT lekérdezést, opcionálisan a prepared statement típusokat és paramétereket, végrehajtja a lekérdezést és JSON-kompatibilis tömböt ad vissza az eredményekkel.

3. changeData() függvény: paraméterként átvesz egy INSERT, UPDATE vagy DELETE lekérdezést, opcionálisan a prepared statement típusokat és paramétereket, végrehajtja a lekérdezést. Sikeres végrehajtás esetén true-t ad vissza, hiba esetén a hibaüzenetet.

##### Sql_functions használata

Az sql_functions.php-t az összes CRUD művelet használja az adatbázis kapcsolat létrehozásához és a lekérdezések végrehajtásához.

#### Read_cookies (Authentikáció)

##### Read_cookies függvényei

1. authentication() függvény: ellenőrzi az auth_token cookie-t. Ha a cookie értéke megegyezik a PSK környezeti változóval, akkor a hozzáférés engedélyezve van (a hálózatos kereséshez szükséges). Egyébként a session táblában keresi a tokent, és ha talál érvényes sessiont, a felhasználó adatait (id, username, rang neve) a $_REQUEST['user'] változóba menti.

##### Read_cookies működése

Az authentikáció session-ökkel történik, amik az adatbázisban tárolódnak. A kliens cookie-kban kapja meg a session tokeneket. A PSK alapú authentikáció lehetővé teszi, hogy a hálózatos keresés során más üzletek lekérdezzék a termék adatokat anélkül, hogy felhasználói fiókkal rendelkeznének.

#### Middleware (Rangok)

##### Middleware konstansai

1. ROLE_UZLETVEZETO: "Üzletvezető" (üzletvezető rang neve)
2. ROLE_HR: "HR-es" (HR-es rang neve)
3. ROLE_RAKTARVEZETO: "Raktárvezető" (raktárvezető rang neve)
4. ROLE_RAKTARKEZELO: "Raktárkezelő" (raktárkezelő rang neve)
5. ROLE_PENZTAROS: "Pénztáros" (pénztáros rang neve)
6. ROLE_EGYEB: "Egyéb dolgozó" (egyéb dolgozó rang neve)

##### Middleware függvényei

1. requireRole(string $role) függvény: paraméterként átveszi a szükséges rang nevét, és ellenőrzi, hogy a bejelentkezett felhasználó rendelkezik-e azzal a ranggal. Visszaadja az ellenőrzés eredményét.

2. requireAdmin() függvény: ellenőrzi, hogy a bejelentkezett felhasználó üzletvezető ranggal rendelkezik-e.

3. requireHRAccess() függvény: ellenőrzi, hogy a bejelentkezett felhasználó HR-es vagy üzletvezető ranggal rendelkezik-e.

4. requireWarehouseManagerAccess() függvény: ellenőrzi, hogy a bejelentkezett felhasználó raktárvezető vagy üzletvezető ranggal rendelkezik-e.

5. requireProductManagementAccess() függvény: ellenőrzi, hogy a bejelentkezett felhasználó raktárkezelő, raktárvezető vagy üzletvezető ranggal rendelkezik-e.

6. requireSearchAccess() függvény: ellenőrzi, hogy a bejelentkezett felhasználó bármilyen ranggal rendelkezik-e (mindenki kereshet).

7. requireQuantityDecreaseAccess() függvény: ellenőrzi, hogy a bejelentkezett felhasználó pénztáros, raktárkezelő, raktárvezető vagy üzletvezető ranggal rendelkezik-e.

8. checkResourceAccess(string $method, string $resource) függvény: paraméterként átveszi a HTTP metódust és az erőforrás nevét, majd megállapítja, hogy a bejelentkezett felhasználónak van-e hozzáférése az adott erőforráshoz az adott művelettel. GET kérés esetén mindenki hozzáfér, írási műveletekhez (POST, PUT, DELETE) a canWrite() függvényt használja.

9. canWrite(string $resource) függvény: paraméterként átveszi az erőforrás nevét és megállapítja, hogy a bejelentkezett felhasználónak van-e írási joga az adott erőforráshoz. A kategóriák, alkategóriák, terméktípusok, márkák és tárolási körülmények kezeléséhez raktárvezető vagy üzletvezető rang szükséges. A termékek kezeléséhez raktárkezelő, raktárvezető vagy üzletvezető rang szükséges.

10. canDecreaseQuantity() függvény: ellenőrzi, hogy a bejelentkezett felhasználó csökkentheti-e a termékek mennyiségét (pénztáros, raktárkezelő, raktárvezető vagy üzletvezető).

##### Middleware használata

A middleware függvényeket a crud_operations.php használja az API végpontok kiszolgálásakor, hogy minden műveletet csak az arra jogosultak hajthassanak végre.

#### Crud_operations

##### Crud_operations működése

A crud_operations.php a backend fő belépési pontja. Feladata:

1. Beállítja a CORS fejléceket (Access-Control-Allow-Origin, Access-Control-Allow-Methods, Access-Control-Allow-Headers, Access-Control-Allow-Credentials), hogy a frontend hozzáférhessen az API-hoz.
2. Betölti az összes szükséges modult (loadenv, sql_functions, read_cookies, roles és az összes CRUD kezelő).
3. Elvégzi az authentikációt az authentication() függvény segítségével. Ha a felhasználó nincs bejelentkezve, 401-es státuszkódot ad vissza.
4. Az URL alapján meghatározza az erőforrás típusát és a megfelelő CRUD kezelőhöz irányítja a kérést.
5. Írási műveletek (POST, PUT, DELETE) esetén a checkResourceAccess() segítségével ellenőrzi a jogosultságot. Ha a felhasználónak nincs joga az adott művelethez, 403-as státuszkódot ad vissza.

Az auth végpont speciális: nem igényel jogosultság-ellenőrzést, és visszaadja a bejelentkezett felhasználó adatait (id, username, rang neve).

#### Model

##### Model struktúrái (adatbázis táblák)

1. Category (kategória):
   - id: int (a kategória azonosítója, elsődleges kulcs)
   - name: string (a kategória neve, egyedi)
   - deleted_at: date (a kategória törlésének időpontja, vagy null a nem töröltek esetén)

2. SubCategory (alkategória):
   - id: int (az alkategória azonosítója, elsődleges kulcs)
   - name: string (az alkategória neve, egyedi)
   - category_id: int (a szülő kategória azonosítója, külső kulcs)
   - deleted_at: date (az alkategória törlésének időpontja, vagy null a nem töröltek esetén)

3. StoringCondition (tárolási körülmény):
   - id: int (a tárolási körülmény azonosítója, elsődleges kulcs)
   - description: string (a tárolási körülmény leírása, egyedi)
   - deleted_at: date (a tárolási körülmény törlésének időpontja, vagy null a nem töröltek esetén)

4. ProductType (terméktípus):
   - id: int (a terméktípus azonosítója, elsődleges kulcs)
   - name: string (a terméktípus neve, egyedi)
   - description: text (a terméktípus leírása)
   - sub_id: int (a kapcsolódó alkategória azonosítója, külső kulcs)
   - storing_condition_id: int (a kapcsolódó tárolási körülmény azonosítója, külső kulcs)
   - deleted_at: date (a terméktípus törlésének időpontja, vagy null a nem töröltek esetén)

5. Brand (márka):
   - id: int (a márka azonosítója, elsődleges kulcs)
   - name: string (a márka neve, egyedi)
   - is_own: int (jelzi, hogy saját márkás-e a termék, 0 vagy 1)
   - is_temporary: int (jelzi, hogy ideiglenes márka-e, 0 vagy 1)
   - deleted_at: date (a márka törlésének időpontja, vagy null a nem töröltek esetén)

6. Product (termék):
   - id: int (a termék azonosítója, elsődleges kulcs)
   - name: string (a termék neve)
   - description: text (a termék leírása)
   - amount: int (a termék mennyisége raktáron)
   - size: string (a termék mérete/súlya)
   - size_type: string (a méret típusa, pl.: kg, liter, db)
   - expires_at: date (a termék lejárati dátuma, vagy null ha nem releváns)
   - price: int (a termék ára forintban)
   - discount: decimal (a termék kedvezménye, 0 és 1 közötti tizedestört, pl.: 0.2 = 20%)
   - warranty: date (a termék garanciájának lejárata, vagy null ha nincs)
   - type_id: int (a kapcsolódó terméktípus azonosítója, külső kulcs)
   - brand_id: int (a kapcsolódó márka azonosítója, külső kulcs)
   - deleted_at: date (a termék törlésének időpontja, vagy null a nem töröltek esetén)

#### CRUD kezelők

##### Category (crud/category.php)

1. GET: visszaadja az összes aktív kategóriát, vagy egy adott kategóriát az id query paraméter alapján. Csak azokat a kategóriákat adja vissza, amelyeknél a deleted_at null.

2. POST: létrehoz egy új kategóriát. A request body-ban a name mező szükséges. Ellenőrzi, hogy a megadott néven nem létezik-e már kategória. Sikeres létrehozás esetén 201-es státuszkódot ad vissza.

3. PUT: frissíti egy meglévő kategória nevét. A request body-ban az id és a name mező szükséges.

4. DELETE: törli (soft delete) a megadott azonosítójú kategóriát. Az id query paraméterként érkezik. Ellenőrzi, hogy nincs-e az adott kategóriához tartozó alkategória. Ha van, a törlés nem hajtódik végre. Sikeres törlés esetén 204-es státuszkódot ad vissza.

##### SubCategory (crud/sub_category.php)

1. GET: visszaadja az összes aktív alkategóriát a hozzájuk tartozó kategória nevével együtt vagy egy adott alkategóriát az id query paraméter alapján.

2. POST: létrehoz egy új alkategóriát. A request body-ban a name és a category_id mező szükséges. Ellenőrzi az egyediséget a category_id mező alapján. Sikeres létrehozás esetén 201-es státuszkódot ad vissza.

3. PUT: frissíti egy meglévő alkategória nevét és kategória azonosítóját. A request body-ban az id, name és category_id mező szükséges.

4. DELETE: törli (soft delete) a megadott azonosítójú alkategóriát. Ellenőrzi, hogy nincs-e az adott alkategóriához tartozó terméktípus. Ha van, a törlés nem hajtódik végre. Sikeres törlés esetén 204-es státuszkódot ad vissza.

##### StoringCondition (crud/storing_condition.php)

1. GET: visszaadja az összes aktív tárolási körülményt vagy egy adott tárolási körülményt az id query paraméter alapján.

2. POST: létrehoz egy új tárolási körülményt. A request body-ban a description mező szükséges. Ellenőrzi az egyediséget. Sikeres létrehozás esetén 201-es státuszkódot ad vissza.

3. PUT: frissíti egy meglévő tárolási körülmény leírását. A request body-ban az id és a description mező szükséges.

4. DELETE: törli (soft delete) a megadott azonosítójú tárolási körülményt. Ellenőrzi, hogy nincs-e az adott tárolási körülményhez kapcsolódó terméktípus. Ha van, a törlés nem hajtódik végre. Sikeres törlés esetén 204-es státuszkódot ad vissza.

##### ProductType (crud/product_type.php)

1. GET: visszaadja az összes aktív terméktípust, vagy egy adott terméktípust az id query paraméter alapján.

2. POST: létrehoz egy új terméktípust. A request body-ban a name, description, sub_id és storing_condition_id mező szükséges. Ellenőrzi az egyediséget. Sikeres létrehozás esetén 201-es státuszkódot ad vissza.

3. PUT: frissíti egy meglévő terméktípus adatait. A request body-ban az id, name, description, sub_id és storing_condition_id mező szükséges.

4. DELETE: törli (soft delete) a megadott azonosítójú terméktípust. Ellenőrzi, hogy nincs-e az adott terméktípushoz tartozó termék. Ha van, a törlés nem hajtódik végre. Sikeres törlés esetén 204-es státuszkódot ad vissza.

##### Brand (crud/brand.php)

1. GET: visszaadja az összes aktív márkát, vagy egy adott márkát az id query paraméter alapján.

2. POST: létrehoz egy új márkát. A request body-ban a name, is_own és is_temporary mező szükséges. Ellenőrzi az egyediséget. Sikeres létrehozás esetén 201-es státuszkódot ad vissza.

3. PUT: frissíti egy meglévő márka adatait. A request body-ban az id, name, is_own és is_temporary mező szükséges.

4. DELETE: törli (soft delete) a megadott azonosítójú márkát. Ellenőrzi, hogy nincs-e az adott márkához tartozó termék. Ha van, a törlés nem hajtódik végre. Sikeres törlés esetén 204-es státuszkódot ad vissza.

##### Product (crud/product.php)

1. GET: visszaadja az összes aktív terméket a kapcsolódó terméktípus és márka nevével együtt (JOIN), vagy egy adott terméket az id query paraméter alapján. Az egyedi termék lekérdezése a teljes kapcsolati láncolatot visszaadja (márka név, típus név, alkategória név, kategória név, tárolási körülmény leírása).

2. POST: létrehoz egy új terméket. A request body-ban a name, description, amount, size, size_type, expires_at, price, discount, warranty, type_id és brand_id mező szükséges. Tartalmaz dátum szanálási logikát: ha a dátum csak évszámból áll, YYYY-01-01 formátumra alakítja. Ha ISO formátumú (T-vel), levágja az idő részt, üres stringet null-ra cserél. Ellenőrzi a duplikátumokat a név és márka azonosító kombinációja alapján. Sikeres létrehozás esetén 201-es státuszkódot ad vissza.

3. PUT: frissíti egy meglévő termék adatait. Ugyanazt a dátum szanálási logikát használja, mint a POST.

4. DELETE: törli (soft delete) a megadott azonosítójú terméket. Sikeres törlés esetén 204-es státuszkódot ad vissza.

##### SearchProduct (crud/search_product.php)

1. POST: komplex keresési végpont, ami dinamikus lekérdezést végez a megadott szűrők alapján. A szűrhető mezők:
   - name: szöveges keresés (LIKE)
   - description: szöveges keresés (LIKE)
   - category_id: kategória azonosító szerinti szűrés
   - sub_category_id: alkategória azonosító szerinti szűrés
   - type_id: terméktípus azonosító szerinti szűrés
   - brand_id: márka azonosító szerinti szűrés
   - storing_condition_id: tárolási körülmény azonosító szerinti szűrés
   - amount_min, amount_max: mennyiség tartomány szűrés
   - price_min, price_max: ár tartomány szűrés
   - size_val: méret szerinti szűrés
   - size_type: mérettípus szerinti szűrés
   - is_discounted: csak kedvezményes termékek (discount > 0)
   - has_warranty: csak garanciával rendelkező termékek (warranty IS NOT NULL)
   - show_expired: ha false, kiszűri a lejárt termékeket (expires_at >= ma vagy NULL)
   - page: lapozás (alapértelmezetten 1)

   A lapozás 25 elemre van beállítva oldalanként. A válasz tartalmazza az adatokat (data), az összes találat számát (total), az aktuális oldalt (page) és a limitet (limit). Prepared statementeket használ az SQL injection megelőzésére.

### Általános API szabályok

Új erőforrás létrehozásánál az azonosító értelemszerűen nem kerül átadásra.
Meglévő erőforrás módosításánál az azonosító és az összes többi adat átadásra kell kerüljön.
Erőforrás törléséhez elég az azonosító.
Minden erőforrás soft-delete formájában törlődik (a deleted_at mező beállításra kerül).

Az authentication() függvény az összes végpont előtt lefut, és kiolvassa az auth_token cookiet, ezért a táblázatba nem került be.
Ennek megfelelően bármelyik végpont adhat vissza 401-es státuszkódot, ha a felhasználó nincs bejelentkezve.
Írási műveletek esetén a checkResourceAccess() ellenőrzi a jogosultságot, és 403-as státuszkódot adhat vissza, ha a felhasználónak nincs joga az adott művelethez.
Az összes végpont előtt szerepel az /api prefix, ez a táblázatba nem került be.
A JSON üzenet a következő formában van:

- hiba esetén: {"error": "\<hibaüzenet\>"}
- siker esetén: {"message": "\<siker üzenet\>"}

Egyéb JSON body esetén a modellben definiált adatok kerülnek átadásra.

## Frontend

### Felhasznált technológiák

1. React keretrendszer - a teljes UI-hoz
2. TailwindCSS - a dizájnhoz
3. Typescript programozási nyelv
4. Vite - a build eszköz

### Projekt struktúra

Az összes forráskód egy React projektnek megfelelő módon a src/ mappában van.

#### Típusok

Az alkalmazás által használt TypeScript típusok itt találhatók. Ezek megfelelnek a backenden található struktúráknak (lásd: modellek), ezért a dokumentáció nem ismétli őket. A csak a frontenden használt típusok itt megtalálhatók.

1. Category: id, name

2. SubCategory: id, name, category_id, category_name (opcionális, lekérdezéskor a backend JOIN-nal adja hozzá)

3. StoringCondition: id, description

4. ProductType: id, name, description, sub_id, storing_condition_id

5. Brand: id, name, is_own (0 vagy 1), is_temporary (0 vagy 1)

6. Product: id, name, description, amount, size, size_type, expires_at, price, discount, warranty, type_id, brand_id, brand_name (opcionális), type_name (opcionális), sub_category_name (opcionális), category_name (opcionális), storing_condition_name (opcionális)

7. ApiResponse: message (opcionális), id (opcionális), name (opcionális) - az API válaszainak fogadására alkalmas típus

8. NetworkStoreResult: store_detail (az üzlet adatai: address, storeTypeName), products (a keresési eredmény: data, total, page, limit) - a hálózatos keresés eredményének típusa

### Router

A router a React Router-nek megfelelő módon regisztrálja az utakat és a hozzájuk tartozó nézeteket.
Az App.tsx a BrowserRouter-t használja, az AuthProvider-t pedig az egész alkalmazás köré helyezi.

A RequireAuth komponens védi az összes belső útvonalat: ellenőrzi, hogy a felhasználó be van-e jelentkezve. Ha nincs bejelentkezve, a bejelentkezési oldalra irányítja át. A bejelentkezési oldal az üzlet adminisztráció alkalmazáson keresztül érhető el (a VITE_STORE_ADMIN_URL környezeti változó alapján).

Az alkalmazás útvonalai:

- / - Főoldal (Dashboard)
- /categories - Kategóriák kezelése
- /subcategories - Alkategóriák kezelése
- /storing-condition - Tárolási körülmények kezelése
- /product-types - Terméktípusok kezelése
- /products - Termékek kezelése
- /brands - Márkák kezelése
- /search - Részletes keresés
- /* - Átirányítás a főoldalra

### Context

#### AuthContext

Az alkalmazás (frontend) azon része, ami a felhasználó authentikációs állapotát és jogosultságait kezeli.

##### AuthContext adatai

1. user: a bejelentkezett felhasználó adatai (id, username, rang neve), vagy null ha nincs bejelentkezve
2. isAuthenticated: jelzi, hogy a felhasználó be van-e jelentkezve
3. loading: jelzi, hogy az authentikáció ellenőrzése folyamatban van-e

##### AuthContext függvényei

1. canWrite(resource: string): boolean függvény: paraméterként átveszi az erőforrás nevét és megállapítja, hogy a bejelentkezett felhasználó rendelkezik-e írási jogosultsággal az adott erőforráshoz. A kategóriák, alkategóriák, típusok, márkák és tárolási körülmények kezeléséhez raktárvezető vagy üzletvezető rang szükséges. A termékek kezeléséhez raktárkezelő, raktárvezető vagy üzletvezető rang szükséges. A felhasználók és szerződések kezeléséhez HR-es rang szükséges.

2. canAccess(resource: string): boolean függvény: paraméterként átveszi az erőforrás nevét és megállapítja, hogy a bejelentkezett felhasználónak van-e olvasási hozzáférése az adott erőforráshoz.

##### AuthContext működése

Az alkalmazás betöltésekor az AuthContext lekérdezi a ./api/auth végpontot, hogy ellenőrizze a felhasználó bejelentkezési állapotát. Ha a felhasználó nincs bejelentkezve, a VITE_STORE_ADMIN_URL környezeti változó által meghatározott üzlet adminisztrációs alkalmazás bejelentkezési oldalára irányítja át.

### Services

A services mappa az API hívásokat kezelő szolgáltatásokat tartalmazza. Minden szolgáltatás a fetch API-t használja JSON fejlécekkel és hibakezeléssel.

#### categoryService

1. getAll(): Promise\<Category[]\> függvény: visszaadja az összes kategóriát
2. getOne(id: number): Promise\<Category\> függvény: visszaad egy adott kategóriát
3. create(name: string): Promise\<ApiResponse\> függvény: létrehoz egy új kategóriát
4. update(id: number, name: string): Promise\<ApiResponse\> függvény: frissít egy kategóriát
5. delete(id: number): Promise\<void\> függvény: töröl egy kategóriát

#### subCategoryService

1. getAll(): Promise\<SubCategory[]\> függvény: visszaadja az összes alkategóriát
2. create(name: string, category_id: number): Promise\<ApiResponse\> függvény: létrehoz egy új alkategóriát
3. update(id: number, name: string, category_id: number): Promise\<ApiResponse\> függvény: frissít egy alkategóriát
4. delete(id: number): Promise\<void\> függvény: töröl egy alkategóriát

#### storingConditionService

1. getAll(): Promise\<StoringCondition[]\> függvény: visszaadja az összes tárolási körülményt (üres tömböt ad vissza 404-es válasz esetén)
2. create(description: string): Promise\<ApiResponse\> függvény: létrehoz egy új tárolási körülményt
3. update(id: number, description: string): Promise\<ApiResponse\> függvény: frissít egy tárolási körülményt
4. delete(id: number): Promise\<void\> függvény: töröl egy tárolási körülményt

#### productTypeService

1. getAll(): Promise\<ProductType[]\> függvény: visszaadja az összes terméktípust
2. create(name: string, description: string, sub_id: number, storing_condition_id: number): Promise\<ApiResponse\> függvény: létrehoz egy új terméktípust
3. update(id: number, name: string, description: string, sub_id: number, storing_condition_id: number): Promise\<ApiResponse\> függvény: frissít egy terméktípust
4. delete(id: number): Promise\<void\> függvény: töröl egy terméktípust

#### brandService

1. getAll(): Promise\<Brand[]\> függvény: visszaadja az összes márkát
2. create(name: string, is_own: boolean, is_temporary: boolean): Promise\<ApiResponse\> függvény: létrehoz egy új márkát (a boolean értékeket 1/0-ra konvertálja)
3. update(id: number, name: string, is_own: boolean, is_temporary: boolean): Promise\<ApiResponse\> függvény: frissít egy márkát
4. delete(id: number): Promise\<void\> függvény: töröl egy márkát

#### productService

1. getAll(): Promise\<Product[]\> függvény: visszaadja az összes terméket
2. create(data: ProductPayload): Promise\<ApiResponse\> függvény: létrehoz egy új terméket
3. update(id: number, data: ProductPayload): Promise\<ApiResponse\> függvény: frissít egy terméket
4. delete(id: number): Promise\<void\> függvény: töröl egy terméket

A ProductPayload interfész tartalmazza az összes termék mezőt: name, description, amount, size, size_type, expires_at, price, discount, warranty, type_id, brand_id.

#### searchFilters

1. SearchFilters interfész: a keresési szűrők típusa, amely tartalmazza az összes szűrhető mezőt (category_id, sub_category_id, type_id, brand_id, storing_condition_id, name, description, amount_min, amount_max, price_min, price_max, size_val, size_type, show_expired, has_warranty, is_discounted, page)

2. buildParams(filters: SearchFilters): Record\<string, any\> függvény: a szűrő objektumból csak a nem üres és nem nulla értékeket tartalmzó objektumot készít, ami JSON formátumban továbbítható az API felé

3. search(filters: SearchFilters): Promise függvény: a helyi keresési végpontra küld kérést (./api/search_product)

4. networkSearch(filters: SearchFilters): Promise függvény: hálózatos keresést indít a többi üzletben a VITE_STORE_ADMIN_URL környezeti változó alapján (\${VITE_STORE_ADMIN_URL}/api/network-search)

### Components

#### Layout komponensek

##### Layout

A fő elrendezési komponens, amely tartalmazza az oldalsávot és a fő tartalmi területet. Reszponzív kialakítású: asztali nézetben az oldalsáv mindig látható, mobil nézetben egy hamburger menü ikon segítségével nyitható és zárható. A fő tartalmi terület az Outlet komponenst használja a beágyazott útvonalak megjelenítéséhez.

##### Sidebar

A navigációs oldalsáv komponens. Tartalmazza a NetStore logót és az alkalmazás nevét, a navigációs menüpontokat, valamint a felhasználó profil szekcióját a kijelentkezés gombbal.
A menüpontok a felhasználó jogosultságaitól függően jelennek meg:
- Ha a felhasználó rendelkezik canWrite('category') jogosultsággal: Kategóriák, Alkategóriák, Tárolási körülmények, Terméktípusok, Márkák
- Ha a felhasználó rendelkezik canWrite('product') jogosultsággal: Termékek
- Mindenki számára elérhető: Részletes keresés
A navigációs linkek a NavLink komponenst használják, aktív állapotjelzéssel.

#### Auth komponensek

##### RequireAuth

Útvonal védő komponens. Ellenőrzi az isAuthenticated állapotot, és ha a felhasználó nincs bejelentkezve, átirányítja a bejelentkező oldalra. A betöltés alatt egy töltés animációt jelenít meg. Ha a felhasználó be van jelentkezve, az Outlet komponens segítségével megjeleníti a beágyazott útvonalakat.

##### AccessDenied

Jogosultság megtagadva komponens. Propsként átveszi az erőforrás nevét és a szükséges rangokat, majd megjeleníti a felhasználó aktuális rangját és a szükséges rangokat. Tartalmaz egy linket a főoldalra történő visszalépéshez.

#### UI komponensek

##### FeedbackMessage

Visszajelzés megjelenítő komponens. Propsként átvesz egy típust ('error' vagy 'success') és egy üzenetet. Hiba esetén piros, siker esetén zöld háttérrel jeleníti meg az üzenetet.

#### Form komponensek

Minden form komponens hasonló mintát követ:
1. Legördülő menü a meglévő elemek kiválasztásához vagy új létrehozásához
2. Input mezők az adatok megadásához
3. Mentés/Törlés gombok betöltési állapottal
4. Visszajelzés üzenetek (FeedbackMessage)
5. Szerkesztés módban egyes mezők (pl.: kapcsolódó entitások) le vannak tiltva

##### CategoryForm

Propsként átveszi a kategóriákat, és biztosítja a felületet kategóriák létrehozásához, módosításához és törléséhez. Egyetlen szerkeszthető mező: name. Tartalmaz egy "Új bejegyzés" gombot, amivel a szerkesztés módból új létrehozás módra válthat.

##### SubCategoryForm

Propsként átveszi az alkategóriákat és a kategóriákat. A szerkeszthető mezők: name és category_id (legördülő menü). Szerkesztés módban a kategória kiválasztása le van tiltva. Az alkategóriák listájában a kategória neve is megjelenik.

##### StoringConditionForm

Propsként átveszi a tárolási körülményeket. Egyetlen szerkeszthető mező: description.

##### ProductTypeForm

Propsként átveszi a terméktípusokat, kategóriákat, alkategóriákat és tárolási körülményeket. A szerkeszthető mezők: name, description (szövegterület), categoryId, subCategoryId, storingConditionId. Kategória-alkategória kaszkád logika: az alkategória kiválasztása a kategória alapján szűrődik. Szerkesztés módban a kategória és alkategória kiválasztása le van tiltva.

##### BrandForm

Propsként átveszi a márkákat. A szerkeszthető mezők: name, is_own (jelölőnégyzet), is_temporary (jelölőnégyzet). A jelölőnégyzetek két oszlopos elrendezésben jelennek meg.

##### ProductForm

A legösszetettebb form komponens. Propsként átveszi a termékeket, kategóriákat, alkategóriákat, terméktípusokat és márkákat. A szerkeszthető mezők: name, description, amount, size, sizeType, expiresAt, price, discount, warranty, categoryId, subCategoryId, typeId, brandId. Háromszintű kaszkád logika: kategória → alkategória → terméktípus. Szerkesztés módban a kapcsolati mezők (kategória, alkategória, típus, márka) le vannak tiltva. Tartalmaz dátum formázási logikát a expires_at és warranty mezőkhöz (formatDate() segédfüggvény, amely többféle dátum formátumot kezel).

##### SearchProductForm

Két fülű keresési felület. Az első fül (Kapcsolatok) tartalmazza a kategória, alkategória, típus, tárolási körülmény és márka szűrőket. A második fül (Tulajdonságok) tartalmazza a mennyiség tartomány, ár tartomány, méret, mérettípus, opciók (lejárt/garancia/kedvezményes) és leírás keresőt. A kiválasztott szűrők számát egy jelvény mutatja. Két keresési gomb érhető el: "Keresés" (helyi) és "Hálózatos keresés" (más üzletekben). A kiválasztott szűrők címkékként jelennek meg.

### Views

#### Dashboard

A főoldal, amely üdvözlő üzenetet jelenít meg a felhasználó nevével. Tartalmaz egy gradiens kártyát a felhasználó profil adataival, egy gyors elérési rácsot (2-4 oszlopos) ikonokkal és linkekkel a különböző funkciókhoz (a linkek a canWrite() jogosultság alapján szűrődnek), valamint egy jogosultsági listát (zöld pötty az engedélyezett műveleteknél). A keresés mindenki számára elérhető. Tartalmaz egy gombot a központi adminisztrációs panelre (VITE_STORE_ADMIN_URL).

#### CategoryManagement

Lekéri a kategóriákat a komponens betöltésekor, ha a felhasználónak van jogosultsága. Kezeli a CategoryForm komponens CRUD műveleteit, a visszajelzést és a törlés megerősítését. Ha a felhasználónak nincs jogosultsága, az AccessDenied komponenst jeleníti meg.

#### SubCategoryManagement

Lekéri az alkategóriákat és a kategóriákat párhuzamosan a komponens betöltésekor. Kezeli a SubCategoryForm komponens CRUD műveleteit. A form állapota tartalmazza a kategória függőséget.

#### StoringConditionManagement

Lekéri a tárolási körülményeket a komponens betöltésekor. Kezeli a StoringConditionForm komponens CRUD műveleteit a leírás mező alapján.

#### ProductTypeManagement

Párhuzamosan (Promise.all()) lekéri a terméktípusokat, kategóriákat, alkategóriákat és tárolási körülményeket. Kezeli a ProductTypeForm komponens CRUD műveleteit. A form állapota kezeli a kategória kaszkádolást.

#### BrandManagement

Lekéri a márkákat a komponens betöltésekor. Kezeli a BrandForm komponens CRUD műveleteit, beleértve az is_own és is_temporary jelzőket.

#### ProductManagement

A legösszetettebb nézet. Párhuzamosan lekéri a termékeket, kategóriákat, alkategóriákat, terméktípusokat és márkákat. Keresési funkcionalitással rendelkezik, ami legördülő eredménylistát jelenít meg (ikon + név + márka címke). Támogatja a sessionStorage-ot: a selectProductId segítségével előzetesen kiválasztható egy termék (a keresési oldalról navigálva). A formatDate() segédfüggvény többféle dátum formátumot kezel. Szerkesztés módban a kategória, alkategória, típus és márka kiválasztása le van tiltva.

#### SearchProductManagement

Két keresési móddal rendelkező nézet:
- Helyi keresés: eredménytáblázat 25 elemes lapozással
- Hálózatos keresés: üzlet kiválasztó legördülő menü, majd a kiválasztott üzlet termékeinek megjelenítése

Az eredménytáblázat oszlopai: név, márka, kategória/típus, ár, mennyiség/méret. A lapozás előre/hátra/számozott gombokkal történik. A termék sorokra kattintva a /products oldalra navigál a kiválasztott termékkel (sessionStorage-on keresztül). Mindkét keresési típushoz tartozik hibakezelés. Üres állapot esetén tájékoztató üzenet jelenik meg.

## A docker konténer

Az alkalmazás minden része docker konténerben futásra van tervezve.
Ennek megfelelően megtalálható a projektrész gyökérmappájában (product_administration mappában) egy Dockerfile és .dockerignore.

### Frontend build

A docker először a frontendet építi ki:

1. Egy node:24-es docker image-ből indul ki.
2. Argumentumként átveszi a VITE_STORE_ADMIN_URL környezeti változót, amit a frontend használ az üzlet adminisztráció oldalának linkjeként.
3. A dockerben megszokott módon a /app könyvtárat jelöli meg kiindulóként, felmásolja a projektet, majd belép a munkamappába. Az npm i segítségével telepíti a függőségeket, majd az npm run build segítségével elkészíti a production buildet. Az eredmény a dist mappába kerül.
4. A build eredményére később frontend-build néven lehet hivatkozni.

### A teljes összerakott projekt

A frontend build eredményét egy PHP-Apache konténerben szolgálja ki:

1. A php:8.2-apache docker image-ből indul ki.
2. Engedélyezi az Apache rewrite modult (a2enmod rewrite), ami a React Router helyes működéséhez szükséges (az SPA útvonalak kezelése).
3. Átmásolja a frontend-build-ből a dist mappa tartalmát az Apache alapértelmezett dokumentumgyökérbe (/var/www/html/).
4. Átmásolja a PHP backend fájlokat (api mappa) a megfelelő helyre.
5. A konténer a 80-as porton szolgálja ki az alkalmazást.
