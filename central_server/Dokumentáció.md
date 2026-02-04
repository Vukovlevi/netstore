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

#### Connection

##### Konstansok

Három konstans található a fájlban, amik hibaüzeneteket tartalmaznak.

##### A connection struktúrái

1. Connection
   - Id: uuid.UUID (a kapcsolat egyedi azonosítója, a kapcsolat létrehozásakor jön létre)
   - Conn: net.Conn (a kapcsolat mögött álló beépített TCP kapcsolat)
   - SearchRequestChan: chan \*queue.SearchRequestNode (a csatorna, ahova a beérkező keresési kéréseket el kell küldenie a kapcsolatnak, hogy a queue fogadja és feldolgozza)
   - AnswerChan: chan \*AnswerMessage (a csatorna, ahova kiküldött keresési kérésre érkező válaszokat kell küldeni, hogy a szerver összegyűjthesse)
   - CurrentAnswerId: string (a jelenleg kiküldött keresési kérés egyedi azonosítója, azért kell, hogy ne kerülhessen másik keresési kérésre adott válasz a szerver által jelenleg várthoz)
   - IsAuthenticated: bool (mutatja, hogy az adott kapcsolat authentikált-e, ha nem, akkor csak authentikációs üzenet fogadható tőle)
   - ConnChan: chan \*Connection (a csatorna, ahova a kapcsolat a saját pointerét küldi, ha authentikált vagy bezárult, hogy értesítse ezen eseményekről a szervert, ami ezáltal módosítja a kezelt kapcsolatok mapját)
   - ReturnError: error (a hiba, ami miatt a kapcsolat read loopja megszakadt, logoláshoz kell)
   - mutex: \*sync.Mutex (egy mutex, ami biztosítja, hogy a több szálról is kezelt erőforrások ne okozzanak problémát)

   - ReadLoop() függvény: a ReadHeader() függvény és a ReadPayload() függvény segítségével olvas üzenetet, saját hiba esetén folytatja az olvasást és kezeli azt a protokollban meghatározott módon, network hiba esetén zárja a kapcsolatot (külön szálon futtatandó, mert blokkolja a futást, amíg egy üzenetre vár)
   - ReadHeader(): (\*TcpHeader, error, error) függvény: vár egy üzenetre a klienstől, majd kiolvas belőle egy protokoll által meghatározott fejlécet, visszaadja a fejléc struktúra pointerét ha van, egy saját errort (ami protokolltól való eltérést jelent), és egy network errort (ami egy beépített hiba, ha pl.: bezáródik a TCP kapcsolat)
   - ReadPayload(uint32): \*TcpMessage függvény: megkapja a várt payload hosszát, majd annyi adatot olvas ki, amiből csinál egy TcpMessage struktúrát, amihez a pointert visszaadja, ha hibát észlel (pl.: nem EOF az utolsó byte -> nem megfelelő message length, azaz hibás message), akkor nilt ad vissza
   - HandleMessage(\*TcpMessage) függvény: megkap egy általános TcpMessage struktúrát, majd a MessageType alapján meghívja a megfelelő függvényt az üzenet kezeléséhez
   - Authenticate(\*AuthenticationMessage) függvény: átvesz egy authentication üzenetet, majd elvégzi az authentikációt, ha sikeres, jelzi a szervernek az új kapcsolat felvételét a kapcsolatok mapba (ConnChan segítségével) és visszaküld egy sikeres authentikáció üzenetet, ha nem sikerül az authentikáció, akkor ennek megfelelő hibaüzenetet küld
   - EnqueueSearch(\*SearchMessage) függvény: átvesz egy SearchMessage-t, amiből csinál egy SearchRequestNode-ot, amit a SearchRequestChan segítségével elküld a queuenak feldolgozásra, valamint egy új szálon elindítja a várakozást a keresés válaszára a WaitForAnswer() függvény segítségével
   - WaitForAnswer(chan []byte) függvény: átvesz egy csatornát, amin egy keresésre érkezett teljes válaszának byte-jai érkeznek, amikor ez megtörténik, elküldi azt a kliensnek (külön szálon futtatandó, mert a teljes válasz megérkezéséig blokkolja a futást)
   - GiveAnswer(\*AnswerMessage) függvény: átvesz egy AnswerMessage struktúrát, amit az AnswerChan-be küld, hogy eljuttasa ezen kliens válaszát a szervernek, amennyiben a kapott válasz azonosítója megegyezik a jelenleg futó kérés azonosítójával (CurrentAnswerId), ha nem, akkor eldobja az üzenetet, ezt thread-safe módon teszi
   - SendErrorMessage(string): error függvény: átvesz egy hibaüzenetet, létrehoz belőle egy struktúrát, majd a SendMessage() függvény segítségével elküldi azt a kliensnek, visszaadja a közben felmerülő hibákat
   - SendMessage(Message): error függvény: átvesz egy Message interfészt implementáló struktúrát, majd a write() függvény segítségével elküldi a kliensnek, visszaadja a közben felmerülő hibákat
   - write([]byte): error függvény: átvesz egy byte tömböt, amit a mögöttes kapcsolaton elküld a kliensnek, visszaadja a közben felmerülő hibákat
   - SendClientSearch(\*ClientSearchMessage): error függvény: átvesz egy kliensnek szánt keresési kérés üzenetet, beállítja thread-safe módon a kapcsolat AnswerChan-jét és CurrentAnswerId-ját, valamint elküldi a kliensnek az üzenetet a SendMessage() függvény segítségével

##### A connection függvényei

1. CreateConnection(net.Conn, chan *SearchRequestNode, chan *Connection): \*Connection függvény: átvesz egy beépített mögöttes TCP kapcsolatot, a csatornát, amin keresztül értesíteni tudja a queuet egy új keresési kérésről, valamint a csatornát, amin keresztül a szerver fele kommunikálhatja a státuszát (authentikált ezért használható a hálózat részeként, vagy bezárt ezért törlendő onnan -> kezelt kapcsolatok mapja), ezek alapján a többi paraméter inicializálásával létrehoz egy Connection struktúrát, amihez visszaadja a pointert
