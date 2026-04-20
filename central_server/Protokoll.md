# A Netstore alkalmazás által használt TCP alapú bináris protokoll dokumentációja - Szabó-Vukov Levente

## A protokoll alapvető működése

A protokoll egy TCP-re épülő, üzenet alapú kommunikációs módszer, amit a rendszer arra használ, hogy megvalósítsa a hálózatos keresés funkciót (aminek leírása a Funkcionális specifikációban található).

## A protokoll fogalmai

HEADER: az üzenet fejléce, az üzenettel kapcsolatos információkat tartalmaz  
PAYLOAD: az üzenet teljes tartalma ([MessageType] és [EOF] byte-okkal együtt)  
CONTENT: az üzenet által küldött lényegi információ, a [MessageType] és az [EOF] byte-ok között foglal helyet a PAYLOAD részeként  
[EOF] byte: minden üzenet lezáró byte-ja, a protokollban ez a 0x4E  
PSK: a kliensek és szerver közötti authentikációs célra szolgáló, előre megbeszélt kulcs (pre-shared key), jellemzően a konfigurációkban kell megadni  
[SearchParam]: az a JSON blob (byte tömbként reprezentált JSON), ami a termék adminisztráció komplex szűrés végpontjának meghívásakor átadandó a HTTP kérés BODY részében. A szűrési paramétereket tartalmazza.  
[SingleAnswer]: az a JSON blob (byte tömbként reprezentált JSON), ami egy kliens válaszát tartalmazza egy keresési kérésre  
[FullAnswer]: az a JSON blob (byte tömbként reprezentált JSON), ami az összes kliens összesített válaszait tartalmazza egy keresési kérésre

## Az üzenetek felépítése

Minden üzenet alapvetően 2 részből áll: HEADER + PAYLOAD

### HEADER

A header hossza 5 byte.
Ebből az első byte: protokoll verzió
A maradék 4 byte: a PAYLOAD hossza

      1 byte                         4 byte

| - - - - - - - - | - - - - - - - - | - - - - - - - - | - - - - - - - - | - - - - - - - - |

      verzó                          payload hossz

### PAYLOAD

A PAYLOAD első és utolsó byte-ja fix.
Az első byte: [MessageType]
Az utolsó byte: [EOF]
Ezen byte-okat minden üzenet PAYLOAD-ja tartalmazza, még azok is, ahol a CONTENT egyébként üres.

Az ezen két byte között található adatmennyiség a CONTENT, ami a [MessageType]-tól függően különböző módokon értelmezendő, kezelendő.

### MessageType

    - 1 (Authentication): Amennyiben a [MessageType] byte értéke 1, az üzenet authentikációs információkat tartalmaz {client -> server} (az üzenet CONTENT része később kifejtve)
    - 7 (AuthenticationSuccess): Amennyiben a [MessageType] byte értéke 7, az üzenet jelentése sikeres authentikáció {server -> client} (az üzenetnek nincs CONTENT része)
    - 2 (Search): Amennyiben a [MessageType] byte értéke 2, az üzenet jelentése keresési kérés {client -> server} (az üzenet CONTENT része később kifejtve)
    - 3 (Answer): Amennyiben a [MessageType] byte értéke 3, az üzenet jelentése válasz egy keresési kérésre {client -> server} (az üzenet CONTENT része később kifejtve)
    - 4 (ClientSearch): Amennyiben a [MessageType] byte értéke 4, az üzenet jelentése megválaszolandó keresési kérés {server -> client} (az üzenet CONTENT része később kifejtve)
    - 5 (ClientAnswer): Amennyiben a [MessageType] byte értéke 5, az üzenet jelentése egy teljesen összeállított válasz egy keresési kérésre {server -> client} (az üzenet CONTENT része később kifejtve)
    - 6 (Error): Amennyiben a [MessageType] byte értéke 6, az üzenet jelentése hiba {server -> client} (az üzenet CONTENT része később kifejtve)

#### 1: Authentication

Authentication üzenet kinézete:
| HEADER | 1 | PSK | EOF |

Ebből a CONTENT rész: PSK
A HEADER-ben található PAYLOAD hossz ezen esetben: PSK hossza + 2 (+2 a [MessageType] és [EOF] miatt)

#### 7: AuthenticationSuccess

AuthenticationSuccess üzenet kinézete:
| HEADER | 7 | EOF |

Ennek az üzenetnek nincs CONTENT része.
A HEADER-ben található PAYLOAD hossz ezen esetben: 2

#### 2: Search

Search üzenet kinézete:
| HEADER | 2 | [SearchParam] | EOF |

Ebből a CONTENT rész: [SearchParam]
A HEADER-ben található PAYLOAD hossz ezen esetben: [SearchParam] hossza + 2 (+2 a [MessageType] és [EOF] miatt)

#### 3: Answer

Answer üzenet kinézete:
| HEADER | 3 | UUID | [SingleAnswer] | EOF |

Ebből a CONTENT rész: UUID + [SingleAnswer]
A HEADER-ben található PAYLOAD hossz ezen esetben: 36 (UUID string hossza) + [SingleAnswer] hossza + 2 (+2 a [MessageType] és [EOF] miatt)

A UUID egy, a szerver által generált UUID, ami azonosítja azt a keresési kérést, amire a válasz érkezik. A UUID-t a kliens a ClientSearch üzenetben kapja meg, és a válasza küldésekor beépíti azt ebbe az üzenetbe.
Ha a szerver olyan UUID-val rendelkező Answer üzenetet kap, amit éppen nem vár (pl.: már lejárt a várakozási idő, és másik kérés van feldolgozása alatt), azt eldobja.

#### 4: ClientSearh

ClientSearch üzenet kinézete:
| HEADER | 4 | UUID | [SearchParam] | EOF |

Ebből a CONTENT rész: UUID + [SearchParam]
A HEADER-ben található PAYLOAD hossz ezen esetben: 36 (UUID string hossza) + [SearchParam] hossza + 2 (+2 a [MessageType] és [EOF] miatt)

A UUID egy, a szerver által generált UUID, ami azonosítja azt a keresési kérést, amire a szerver a választ várja. Ezen üzenet részeként küldi el a klienseknek, amik a válaszukban ezt visszaküldik.
Ha a szerver olyan UUID-val rendelkező Answer üzenetet kap, amit éppen nem vár (pl.: már lejárt a várakozási idő, és másik kérés van feldolgozása alatt), azt eldobja.

#### 5: ClientAnswer

ClientAnswer üzenet kinézete:
| HEADER | 5 | [FullAnswer] | EOF |

Ebből a CONTENT rész: [FullAnswer]
A HEADER-ben található PAYLOAD hossz ezen esetben: [FullAnswer] hossza + 2 (+2 a [MessageType] és [EOF] miatt)

#### 6: Error

Error üzenet kinézete:
| HEADER | 6 | hibaüzenet | EOF |

Ebből a CONTENT rész: hibaüzenet
A HEADER-ben található PAYLOAD hossz ezen esetben: hibaüzenet hossza + 2 (+2 a [MessageType] és [EOF] miatt)

## Kapcsolatfelvétel és kapcsolattartás egy kliens és a szerver között

1. Lezajlik az authentikációs folyamat (a további lépések csak sikeres authentikáció után futhatnak)
2. Hálózatos keresés küldése
3. Hálózatos keresés megválaszolása
4. Hibás formátumú üzenetek

### Az authentikációs folyamat

1. A kliens Authentication üzenetet küld a szerver részére.
2. A szerver az üzenetben érkező PSK-t ellenőrzi a saját konfigurációja alapján, majd két módon kezelheti: 2. a, Egyező PSK: a szerver egy AuthenticationSuccess üzenetet küld a kliens részére, a kapcsolat a két állomás között innentől kezdve használható egyéb jellegű (nem authentikációs célú) kommunikációra. 2. b, Nem egyező PSK: a szerver egy Error üzenetet küld a kliensnek, a kapcsolat a két állomás között továbbra is csak authentikációs célra alkalmazható.
3. Ha nem Authentication üzenet érkezik a klienstől, miközben a kapcsolat csak authentikációs célra alkalmazható, akkor a szerver bontja a kapcsolatot.

### Hálózatos keresés küldése

1. A kliens Search üzenetet küld a szerver részére
2. A szerver ezt továbbítja az összes többi kliensnek megválaszolásra (ClientSearch üzenet).
3. A szerver időlimittel összegyűjti a válaszokat (Answer üzenet).
4. A szerver addig vár, amíg minden kiküldött üzenetre nem érkezik válasz, vagy lejár az idő. Az idő lejárása esetén az addig megérkezett válaszokat továbbítja, a később érkezőket eldobja.
5. A szerver elküldi a válaszokat a kliensnek (ClientAnswer üzenet).

### Hálózatos keresés megválaszolása

1. A szerver ClientSearch üzenetet üzenetet küld a kliensnek.
2. A kliens lekérdezi a keresési eredményeket, majd visszaküldi a szervernek (Answer üzenet).
3. Ha a kliens túllépi a szerver várakozási idejét a válaszadással, a beérkezett választ a szerver el fogja dobni.

### Hibás formátumú üzenetek

A szerver a hibás formátummal rendelkező üzeneteket eldobja.
Erre példa, ha nem kap elég byte-ot, hogy összeállítson egy HEADER-t, vagy nem egyezik a tényleges beérkezett üzenet hossza a HEADER-ben szereplővel.
A szerver ilyen hibákra választ nem ad, de a beérkezett adatokkal nem foglalkozik.
