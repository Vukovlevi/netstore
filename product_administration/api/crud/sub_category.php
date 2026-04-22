<?php
function handleSubCategory($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['check_deleted'])) {
                $needle = trim((string)$_GET['check_deleted']);
                if ($needle === '') {
                    echo json_encode(null, JSON_UNESCAPED_UNICODE);
                    break;
                }
                $row = getData("SELECT id, name, category_id, deleted_at FROM sub_category WHERE LOWER(TRIM(name)) = LOWER(?) AND deleted_at IS NOT NULL ORDER BY id DESC LIMIT 1", 's', [$needle]);
                echo json_encode(!empty($row) ? $row[0] : null, JSON_UNESCAPED_UNICODE);
                break;
            }
            if (isset($_GET['id'])) {
                $sub = getData("SELECT sc.id, sc.name, sc.category_id, c.name as category_name FROM sub_category sc JOIN category c ON sc.category_id = c.id WHERE sc.id = ? AND sc.deleted_at IS NULL", 'i', [$_GET['id']]);

                if (empty($sub)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Alkategória nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return;
                }
                echo json_encode($sub[0], JSON_UNESCAPED_UNICODE);
            } else {
                $subs = getData("SELECT sc.id, sc.name, sc.category_id, c.name as category_name FROM sub_category sc JOIN category c ON sc.category_id = c.id WHERE sc.deleted_at IS NULL ORDER BY sc.name");
                echo json_encode($subs, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            $nameInput = isset($body['name']) ? trim((string)$body['name']) : '';

            if (!empty($body['restore']) && !empty($body['id']) && $nameInput !== '' && !empty($body['category_id'])) {
                $deletedRow = getData("SELECT id FROM sub_category WHERE id = ? AND deleted_at IS NOT NULL", 'i', [$body['id']]);

                if (empty($deletedRow)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Nincs visszaállítható alkategória ezzel az azonosítóval!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $nameClash = getData("SELECT id FROM sub_category WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

                if (!empty($nameClash)) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív alkategória!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $restored = changeData("UPDATE sub_category SET name = ?, category_id = ?, deleted_at = NULL WHERE id = ?", 'sii', [$nameInput, $body['category_id'], $body['id']]);

                if (!$restored) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $row = getData("SELECT sc.id, sc.name, sc.category_id, c.name as category_name FROM sub_category sc JOIN category c ON sc.category_id = c.id WHERE sc.id = ? AND sc.deleted_at IS NULL", 'i', [$body['id']]);
                http_response_code(200);
                echo json_encode($row[0], JSON_UNESCAPED_UNICODE);
                break;
            }

            if ($nameInput === '' || empty($body['category_id'])) {
                http_response_code(400);
                echo json_encode(['message' => 'Hiányzó adat: név és kategória megadása kötelező!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existing = getData("SELECT id, deleted_at FROM sub_category WHERE LOWER(TRIM(name)) = LOWER(?) ORDER BY (deleted_at IS NULL) DESC LIMIT 1", 's', [$nameInput]);

            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű alkategória!'], JSON_UNESCAPED_UNICODE);
                    return;
                } else {
                    http_response_code(409);
                    echo json_encode([
                        'message' => 'Létezik egy korábban törölt alkategória ezen a néven. Szeretné visszaállítani?',
                        'restorable' => true,
                        'id' => (int)$existing[0]['id']
                    ], JSON_UNESCAPED_UNICODE);
                    return;
                }
            }

            $success = changeData("INSERT INTO sub_category (name, category_id) VALUES (?, ?)", 'si', [$nameInput, $body['category_id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $newSub = getData("SELECT sc.id, sc.name, sc.category_id, c.name as category_name FROM sub_category sc JOIN category c ON sc.category_id = c.id WHERE sc.name = ? AND sc.deleted_at IS NULL", 's', [$nameInput]);

            http_response_code(201);
            echo json_encode($newSub[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name']) || empty($body['category_id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: azonosító, név és kategória megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $nameInput = trim((string)$body['name']);

            $subToUpdate = getData("SELECT id FROM sub_category WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);

            if (empty($subToUpdate)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú alkategória, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existingName = getData("SELECT id FROM sub_category WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

            if (!empty($existingName)) {
                http_response_code(409);
                echo json_encode(['message' => 'Már létezik ilyen nevű alkategória!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE sub_category SET name = ?, category_id = ? WHERE id = ?", 'sii', [$nameInput, $body['category_id'], $body['id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $updatedSub = getData("SELECT sc.id, sc.name, sc.category_id, c.name as category_name FROM sub_category sc JOIN category c ON sc.category_id = c.id WHERE sc.id = ? AND sc.deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedSub[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $subToDelete = getData("SELECT id FROM sub_category WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (empty($subToDelete)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú alkategória!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $types = getData("SELECT id FROM product_type WHERE sub_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (!empty($types)) {
                http_response_code(409);
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ehhez az alkategóriához aktív terméktípusok tartoznak.'
                ], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE sub_category SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);

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
