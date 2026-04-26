<?php
function handleBrand($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['check_deleted'])) {
                $needle = trim((string)$_GET['check_deleted']);
                if ($needle === '') {
                    echo json_encode(null, JSON_UNESCAPED_UNICODE);
                    break;
                }
                $row = getData("SELECT id, name, is_own, is_temporary, deleted_at FROM brand WHERE LOWER(TRIM(name)) = LOWER(?) AND deleted_at IS NOT NULL ORDER BY id DESC LIMIT 1", 's', [$needle]);
                echo json_encode(!empty($row) ? $row[0] : null, JSON_UNESCAPED_UNICODE);
                break;
            }
            if (isset($_GET['id'])) {
                $brand = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

                if (empty($brand)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Márka nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return;
                }
                echo json_encode($brand[0], JSON_UNESCAPED_UNICODE);
            } else {
                $brands = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE deleted_at IS NULL ORDER BY name");
                echo json_encode($brands, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            $nameInput = isset($body['name']) ? trim((string)$body['name']) : '';

            if (!empty($body['restore']) && !empty($body['id']) && $nameInput !== '') {
                $deletedRow = getData("SELECT id FROM brand WHERE id = ? AND deleted_at IS NOT NULL", 'i', [$body['id']]);

                if (empty($deletedRow)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Nincs visszaállítható márka ezzel az azonosítóval!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $nameClash = getData("SELECT id FROM brand WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

                if (!empty($nameClash)) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív márka!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $isOwn = isset($body['is_own']) ? (int)$body['is_own'] : 0;
                $isTemporary = isset($body['is_temporary']) ? (int)$body['is_temporary'] : 0;

                $restored = changeData("UPDATE brand SET name = ?, is_own = ?, is_temporary = ?, deleted_at = NULL WHERE id = ?", 'siii', [$nameInput, $isOwn, $isTemporary, $body['id']]);

                if (!$restored) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $row = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
                http_response_code(200);
                echo json_encode($row[0], JSON_UNESCAPED_UNICODE);
                break;
            }

            if ($nameInput === '') {
                http_response_code(400);
                echo json_encode(['message' => 'Hiányzó adat: a márka neve kötelező!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existing = getData("SELECT id, deleted_at FROM brand WHERE LOWER(TRIM(name)) = LOWER(?) ORDER BY (deleted_at IS NULL) DESC LIMIT 1", 's', [$nameInput]);

            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű márka!'], JSON_UNESCAPED_UNICODE);
                    return;
                } else {
                    http_response_code(409);
                    echo json_encode([
                        'message' => 'Létezik egy korábban törölt márka ezen a néven. Szeretné visszaállítani?',
                        'restorable' => true,
                        'id' => (int)$existing[0]['id']
                    ], JSON_UNESCAPED_UNICODE);
                    return;
                }
            }

            $isOwn = isset($body['is_own']) ? (int)$body['is_own'] : 0;
            $isTemporary = isset($body['is_temporary']) ? (int)$body['is_temporary'] : 0;

            $success = changeData("INSERT INTO brand (name, is_own, is_temporary) VALUES (?, ?, ?)", 'sii', [$nameInput, $isOwn, $isTemporary]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $newBrand = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE name = ? AND deleted_at IS NULL", 's', [$nameInput]);

            http_response_code(201);
            echo json_encode($newBrand[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító és a név megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $nameInput = trim((string)$body['name']);

            $brandToUpdate = getData("SELECT id FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);

            if (empty($brandToUpdate)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú márka, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existingName = getData("SELECT id FROM brand WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

            if (!empty($existingName)) {
                http_response_code(409);
                echo json_encode(['message' => 'Már létezik ilyen nevű márka!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $isOwn = isset($body['is_own']) ? (int)$body['is_own'] : 0;
            $isTemporary = isset($body['is_temporary']) ? (int)$body['is_temporary'] : 0;

            $success = changeData("UPDATE brand SET name = ?, is_own = ?, is_temporary = ? WHERE id = ?", 'siii', [$nameInput, $isOwn, $isTemporary, $body['id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $updatedBrand = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedBrand[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $brandToDelete = getData("SELECT id FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (empty($brandToDelete)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú márka!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $products = getData("SELECT id FROM product WHERE brand_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (!empty($products)) {
                http_response_code(409);
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ehhez a márkához aktív termékek tartoznak.'
                ], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE brand SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            http_response_code(204);
            break;
        default:
            http_response_code(405);
            return;
    }
}
