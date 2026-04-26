<?php
function handleProductType($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['check_deleted'])) {
                $needle = trim((string)$_GET['check_deleted']);
                if ($needle === '') {
                    echo json_encode(null, JSON_UNESCAPED_UNICODE);
                    break;
                }
                $row = getData("SELECT id, name, description, sub_id, storing_condition_id, deleted_at FROM product_type WHERE LOWER(TRIM(name)) = LOWER(?) AND deleted_at IS NOT NULL ORDER BY id DESC LIMIT 1", 's', [$needle]);
                echo json_encode(!empty($row) ? $row[0] : null, JSON_UNESCAPED_UNICODE);
                break;
            }
            if (isset($_GET['id'])) {
                $type = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

                if (empty($type)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Terméktípus nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return;
                }
                echo json_encode($type[0], JSON_UNESCAPED_UNICODE);
            } else {
                $types = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE deleted_at IS NULL ORDER BY name");

                echo json_encode($types, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            $nameInput = isset($body['name']) ? trim((string)$body['name']) : '';

            if (!empty($body['restore']) && !empty($body['id']) && $nameInput !== '' && !empty($body['description']) && !empty($body['sub_id']) && !empty($body['storing_condition_id'])) {
                $deletedRow = getData("SELECT id FROM product_type WHERE id = ? AND deleted_at IS NOT NULL", 'i', [$body['id']]);

                if (empty($deletedRow)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Nincs visszaállítható terméktípus ezzel az azonosítóval!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $nameClash = getData("SELECT id FROM product_type WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

                if (!empty($nameClash)) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív terméktípus!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $restored = changeData("UPDATE product_type SET name = ?, description = ?, sub_id = ?, storing_condition_id = ?, deleted_at = NULL WHERE id = ?", 'ssiii', [$nameInput, $body['description'], $body['sub_id'], $body['storing_condition_id'], $body['id']]);

                if (!$restored) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $row = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
                http_response_code(200);
                echo json_encode($row[0], JSON_UNESCAPED_UNICODE);
                break;
            }

            if ($nameInput === '' || empty($body['description']) || empty($body['sub_id']) || empty($body['storing_condition_id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: név, leírás, alkategória és tárolási körülmény megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existing = getData("SELECT id, deleted_at FROM product_type WHERE LOWER(TRIM(name)) = LOWER(?) ORDER BY (deleted_at IS NULL) DESC LIMIT 1", 's', [$nameInput]);

            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű terméktípus!'], JSON_UNESCAPED_UNICODE);
                    return;
                } else {
                    http_response_code(409);
                    echo json_encode([
                        'message' => 'Létezik egy korábban törölt terméktípus ezen a néven. Szeretné visszaállítani?',
                        'restorable' => true,
                        'id' => (int)$existing[0]['id']
                    ], JSON_UNESCAPED_UNICODE);
                    return;
                }
            }

            $success = changeData("INSERT INTO product_type (name, description, sub_id, storing_condition_id) VALUES (?, ?, ?, ?)", 'ssii', [$nameInput, $body['description'], $body['sub_id'], $body['storing_condition_id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $newType = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE name = ? AND deleted_at IS NULL", 's', [$nameInput]);

            http_response_code(201);
            echo json_encode($newType[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name']) || empty($body['description']) || empty($body['sub_id']) || empty($body['storing_condition_id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: azonosító, név, leírás, alkategória és tárolási körülmény megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $nameInput = trim((string)$body['name']);

            $typeToUpdate = getData("SELECT id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);

            if (empty($typeToUpdate)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú terméktípus, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existingName = getData("SELECT id FROM product_type WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

            if (!empty($existingName)) {
                http_response_code(409);
                echo json_encode(['message' => 'Már létezik ilyen nevű terméktípus!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE product_type SET name = ?, description = ?, sub_id = ?, storing_condition_id = ? WHERE id = ?", 'ssiii', [$nameInput, $body['description'], $body['sub_id'], $body['storing_condition_id'], $body['id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $updatedType = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedType[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $typeToDelete = getData("SELECT id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (empty($typeToDelete)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú terméktípus, vagy már törölve lett!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $products = getData("SELECT id FROM product WHERE type_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (!empty($products)) {
                http_response_code(409);
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ez a terméktípus aktív terméke(ke)t tartalmaz. Előbb törölje a termékeket.'
                ], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE product_type SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);

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
