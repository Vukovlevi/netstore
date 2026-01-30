# Közpotni szerver fejlesztői dokumentáció

## A központi szerver feladata

A központi szerver feladata, hogy képes legyen fogadni a hozzá csatlakozó klienseket. Ezekkel egy TCP alapú (saját fejlesztésű) kommunikációs protokollal kommunikál. Lehetőséget biztosít számukra a többi csatlakozott kliens felé indított termékkeresésekre. Ehhez biztosítja az összeköttetést, fogadja a kéréseket, válaszokat és átjátszó állomásként működik.

## Felhasznált technológiák

1. Go programozási nyelv
2. google/uuid könyvtár
3. joho/godotenv könyvtár
4. TCP alapú kommunikációs protokoll (lásd: Protokoll.md)

## Packagek

### Config

#### A package célja

A szerver működéséhez szükséges konfiguráció beolvasását tartalmazza. A konfiguráció környezeti változókon keresztül történik.

#### Környezeti változók

1. IP: az ip cím, amin a központi szerver figyel (megadás nélkül alapértelmezetten ez a 0.0.0.0)
2. PORT: az a port, amin a központi szerver figyel (megadás nélkül alapértelmezetten ez a 42069)

#### A package struktúrái

1. Config
   - ip: string
   - port: string
   - ToAddress(): string függvény: visszaadja a konfiguráció által tartalmazott ip-ből és portból készített címet, amire a go beépített net package hallgatni tud

#### A package függvényei

1. LoadConfig(): \*Config függvény: beolvassa a környezeti változókban tárolt konfigurációs paramétereket, és visszaad egy pointert, az ezeket tartalmazó Config struktúrához

#### A package működése, használata

A szerver indulásakor a LoadConfig() függvény meghívásra kerül, majd a visszakapott Config struktúrát címmé alakítva (ToAddress() függvény) elindítani a TCP szervert a beépített net package használatával.

### Queue

#### A package célja

A központi szerver által több szálon fogadott keresési kérések összefésülése egy queue adattípusba (first in first out), majd azok továbbítása feldolgozásra egyszerre egy alapon (feldolgozás: kérés kiküldése az összes többi kliensek, beérkezett válaszok üsszegyűjtése és visszajuttatása az eredeti kérdezőnek).

#### Konstansok

1. STATUS_CAN_SEARCH = 1: azt mutatja, hogy nincs jelenleg feldolgozás alatt álló kérés, tehát meg lehet kezdeni a következő kérés feldolgozását
2. STATUS_ANSWERING = 2: azt mutatja, hogy van jelenleg feldolgozás alatt álló kérés, tehát nem indítható a következő kérés feldolgozása

#### A package struktúrái

1. SearchRequestNode
   - Next: \*SearchRequestNode (a következő keresési kérés node-jára mutató pointer -> a queue adattípus megvalósításához, egy singly linked list a mögöttes adatszerkezet)
   - SearchParam: []byte (az adott node-hoz tartozó komplex szűrés végpontra küldendő JSON adat byte-jai)
   - FullAnswerChan: chan []byte (az a csatorna, ahova a szerver által begyűjtött összesített válasz JSON objektumának byte-jait küldeni kell -> erre hallgat az eredeti kérdezőt kezelő kapcsolat)
   - ClientId: string (a kérést indító kliens kapcsolatának UUID azonosítója string formában)

2. SearchRequestQueue
   - Head: \*SearchRequestNode (az első keresési kérés node-jára mutató pointer, ha nincs ilyen, akkor nil)
   - Tail: \*SearchRequestNode (az utolsó keresési kérés node-jára mutató pointer, ha nincs ilyen, akkor nil)
   - Status: int (a két konstans státusz értékét veheti fel, ez mutatja a queue számára, hogy kezdheti-e a következő kérés feldolgozását)
   - SearchRequestChan: chan \*SearchRequestNode (ezen a csatornán hallgatja a queue a beérkező keresési kéréseket, a kapcsolatokat kezelő szálak ide küldik a beérkezett keresési kéréseket node-dá alakítva)
   - ProcessCallBack: func(\*SearchRequestNode) (egy callback függvény, ami paraméterként átvesz egy feldolgozásra szánt keresési kérés node-ot, ezt hívja meg a queue, amikor megkezdi egy keresési kérés feldolgozását)
   - IsTesting: bool (azt mutatja, hogy tesztelés alatt áll-e a kód, ahhoz kell, hogy a mögöttes queue adattípus unit tesztelése közben ne kezdjen el automatikusan kéréseket feldolgozni a rendszer)
   - mutex: \*sync.Mutex (egy mutex, ami biztosítja, hogy ne legyen race condition, mivel több szál is hozzá fog férni a struktúra mezőihez)

   - HandleSearchRequest() függvény: ez a függvény hallgatja a beérkező keresési kéréseket a SearchRequestChan csatornáról, majd berakja őket a queue-ba az Enqueue() függvény segítségével (külön szálon futtatandó, mert blokkolja az adott szál futását egy üzenet érkezéséig a csatornán)
   - Enqueue(\*SearchRequestNode) függvény: paraméterként átvesz egy keresési kérés node-ot, majd thread-safe módon behelyezi a queue-ba (kezeli a mögöttes singly linked list adatstruktúrát is), valamint meghívja a Process() függvényt, hogy az megkezdje a következő kérés feldolgozását, amennyiben lehetséges (STATUS_CAN_SEARCH)
   - Dequeue(): \*SearchRequestNode függvény: visszaadja a következő keresési kérés node-jának pointerét (vagy nil-t, ha nincs ilyen), kezeli a mögöttes singly linked list adatstruktúrát is (fontos: nem thread-safe, csak thread-safe körülmények között hívható)
   - Process() függvény: thread-safe módon megkezdi a következő keresési kérés feldolgozását, amennyiben ez lehetséges (STATUS_CAN_SEARCH), a következő kérés megállapításához a Dequeue() függvényt használja (ezt megteheti, mert thread-safe a körülmény), a feldolgozás indításához használja a ProcessCallback() függvényt, és beállítja a queue Status mezőjét STATUS_ANSWERING-re
   - FinishProcess() függvény: egy kérés feldolgozásának befejezésekor hívódik meg, thread-safe módon beállítja a queue Status mezőjét STATUN_CAN_SEARCH-re, valamint a Process() függvényt meghívva megkezdi a következő kérés feldolgozását (amennyiben van)

#### A package függvényei

1. NewSearchRequestQueue(processCallback func(*SearchRequestNode)): *SearchRequestQueue: paraméterkén átveszi a feldolgozás megkezdésekor meghívandó callback függvényt, visszaad egy SearchRequestQueue pointert, amiben az összes mező inicializálva van

#### A package működése, használata

A NewSearchRequestQueue() függvény segítségével kérni kell egy új queue-t (egy ilyen példány lesz, ami a szerver indulásakor jön létre). Ezután külön szálon el kell indítani a HandleSearchRequest() függvényt, hogy hallgassa beérkező keresési kéréseket, és kezelje őket. Amint egy beérkezik, bekerül a mögöttes adatstruktúrákba, és a státusz függvényében megfelelő thread-safe módon megkezdődik a kérés(ek) feldolgozása. A feldolgozás kivitelezése a szerver feladata, a queue csak a sorrendet és az egyszerre csak egy kérés feldolgozását biztosítja az által, hogy a queue indítja a feldolgozás folyamatát, amit viszont a szerver végez (ezért kell a callback függvény).

### Tcp

#### A package célja

Ez a package felelős a TCP kapcsolatok kezeléséért, valamint a kommunikációs protokollnak megfelelő üzenetértelmezésért.

#### Header

##### Konstansok

    - VERSION = 1: a protokoll verziója
    - VERSION_SIZE = 1: a verziót jelző byte-ok száma
    - MESSAGE_LENGTH_SIZE = 4: a payload hosszúságát jelző byte-ok száma
    - HEADER_SIZE: a protokoll szerinti fejléc hossza (ez a verzóból, valamint a payload hosszúságából áll)
    - ErrVersionMismatch: nem megfelelő protokoll verzió fogadásakor használandó hiba

##### A header struktúrái

1. TcpHeader
   - Version: byte (a kapott fejlécben található protokoll verzó)
   - MsgLen: uint32 (a kapott fejlécben található payload hosszúság)

   - ValidateHeader(): error függvény: ellenőrzi, hogy a kapott fejléc a megfelelő protokoll verziót használja-e, ha nem, akkor ErrVersionMismatch hibát ad vissza, egyébként nil-t

##### A header függvényei

1. CreateHeaderFromBuffer([]byte) \*TcpHeader függvény: paraméterként átvesz egy 5 byte-ból álló tömböt, aminek adataiból létrehozza a TcpHeader struktúrát és visszaadja hozzá a pointert

2. CreateHeaderForPayload([]byte): []byte függvény: paraméterként átvesz egy protokoll szerinti payload-ot, és létrehozza hozzá a protokoll szerinti fejlécet (egy byte tömböt, aminek első byte-ja a protokoll verzió, a maradék 4 byte pedig uint32-ként értelmezve a payload hossza), majd visszaadja a fejlécet alkotó byte-okat byte tömbként

#### Message

##### Konstansok

    - a protokollban meghatározott üzenet típusok
    - a protokollban meghatározott EOF byte
    - UUID_LENGTH = 36: egy UUID string alakjának hossza (szükséges, mert egyes üzenetek a protokoll szerint tartalmaznak UUID-t)
    - AuthenticationError: egy hibaüzenet, ami nem egyező PSK-k esetén kerül felhasználásra

##### A message struktúrái és interface-ei

1. Message (interface)
   - ToMessageBytes(): []byte függvény: ez a függvény az egyes struktúrákon van implementálva, és az adott üzenetet alakítja át a protokollnak megfelelően byte tömbbé, ami küldhető a TCP kapcsolaton

2. TcpMessage (struktúra)
   - MessageType: byte (az üzenet típusa a protokoll szerint)
   - Content: []byte (az üzenet protokoll szerinti hasznos adatait tartalmazó byte tömb -> nincs benne a protokoll szerint a payload részét képző első byte, ami az üzenet típusa, valamint az utolsó byte, ami az EOF)

   - ToMessageBytes(): []byte függvény(): a Message interface implementálása, a struktúrában található adatokból összeállít egy a protokollnak megfelelő byte tömböt (header generálása, content payloaddá alakítása -> első byte message type, utolsó EOF, majd ezek összefűzése egy tömmbé)
   - ToAuthenticationMessage(): \*AuthenticationMessage: egy általános TcpMessage struktúrának adatait protokoll szerinti AuthenticationMessage-ként értelmezi, és egy olyan struktúra pointerét adja vissza
   - ToSearchMessage(): \*SearchMessage: egy általános TcpMessage struktúrának adatait protokoll szerinti SearchMessage-ként értelmezi, és egy olyan struktúra pointerét adja vissza
   - ToAnswerMessage(): \*AnswerMessage: egy általános TcpMessage struktúráinak adatait protokoll szerinti AnswerMessage-ként értelmezi, és egy olyan struktúra pointerét adja vissza

3. AuthenticationMessage (struktúra)
   - \*TcpMessage (a go féle "öröklődés"), TcpMessage tulajdonságaivel rendelkezik
   - Psk: string (a protokoll szerinti authentication message-ben küldött psk)

   - Authenticate(): error függvény: ellenőrzi a kapott psk-t a szerver által elfogadottal szemben, visszaad egy AuthenticationError-t ha hibás a psk, nil-t ha helyes

4. SearchMessage (struktúra)
   - \*TcpMessage (a go féle "öröklődés"), TcpMessage tulajdonságaivel rendelkezik
   - FullAnswerChan: chan []byte (az adott keresési kérelemhez tartozó csatorna, amin a kapcsolat hallgatja a szerver által összegyűjtött teljes választ)

5. AnswerMessage (struktúra)
   - \*TcpMessage (a go féle "öröklődés"), TcpMessage tulajdonságaivel rendelkezik
   - AnswerId: string (annak a kérésnek a UUID-ja, amihez a válasz tartozik)

6. ClientSearchMessage (struktúra)
   - \*TcpMessage (a go féle "öröklődés"), TcpMessage tulajdonságaivel rendelkezik
   - ClientId: a keresést indító kliens UUID-ja
   - AnswerId: az adott keresés UUID-ja
   - SingleAnswerChan: chan \*AnswerMessage (a csatorna, ahova a kapcsolatok a kérésre kapott válaszukat küldhetik, a szerver ezen hallgatja a válaszokat, amiket összesítve létrehozza a teljes választ)
   - SearchParam: []byte (a komplex szűrés végpontra szánt JSON adat byte tömb formában)

   - ToMessageBytes(): []byte függvény: a Message interface implementálása, a protokoll szerinti hasznos tartalom előállításáért felel az általános TcpMessage Content mezőjében (az AnswerId után fűzött SearchParam), majd visszaadja az ezen adatokból már megfelelően működő általános TcpMessage ToMessageBytes() függvényének eredményét

7. ClientAnswerMessage (struktúra)
   - \*TcpMessage (a go féle "öröklődés"), TcpMessage tulajdonságaivel rendelkezik

8. ErrorMessage (struktúra)
   - \*TcpMessage (a go féle "öröklődés"), TcpMessage tulajdonságaivel rendelkezik
   - Msg: string (a hibaüzenet)

   - ToMessageBytes(): []byte függvény: a Message interface implementálása, a protokoll szerinti hasznos tartalom előállításáért felel az általános TcpMessage Content mezőjében (a hibaüzenet byte tömbként), majd visszaadja az ezen adatokból már megfelelően működő általános TcpMessage ToMessageBytes() függvényének eredményét

9. AuthenticationSuccessMessage (struktúra)
   - \*TcpMessage (a go féle "öröklődés"), TcpMessage tulajdonságaivel rendelkezik

##### A message függvényei

1. CreateTcpMessageFromPayload([]byte): \*TcpMessage függvény: paraméterként átveszi a protokoll szerinti payload-ot, majd abból külön mezőbe kiveszi az üzenet típusát, majd a maradék byteokról leveszi az EOF-et, és a hasznos adatot a Content mezőben tárolva visszaadja az általános TcpMessage struktúra pointerét

2. CreateClientSearchMessage(string, string, []byte): \*ClientSearchMessage függvény: paraméterként átveszi a kliens UUID string verzióját, az adott keresés UUID string verzióját, valamint a komplex szűrés végpontra szánt JSON adatok byte tömb verzióját, majd ezekből előállít egy struktúrát, és visszaadja hozzá a pointert

3. CreateClientAnswerMessage([]byte): \*ClientAnswerMessage függvény: paraméterként átveszi a teljes válasz JSON adatainak byte tömb verzióját, majd létrehoz egy struktúrát az adatokkal, és visszaadja hozzá a pointert

4. CreateErrorMessage(string): \*ErrorMessage függvény: paraméterként átveszi a hibaüzenetet, majd létrehoz egy struktúrát az adatokkal, és visszaadja hozzá a pointert

5. CreateAuthenticationSuccessMessage(): \*AuthenticationSuccessMessage függvény: visszaad egy pointert egy sikeres authentikációt jelző üzenethez
