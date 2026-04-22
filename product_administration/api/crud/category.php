<?php
function handleCategory($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['check_deleted'])) {
                $needle = trim((string)$_GET['check_deleted']);
                if ($needle === '') {
                    echo json_encode(null, JSON_UNESCAPED_UNICODE);
                    break;
                }
                $row = getData("SELECT id, name, deleted_at FROM category WHERE LOWER(TRIM(name)) = LOWER(?) AND deleted_at IS NOT NULL ORDER BY id DESC LIMIT 1", 's', [$needle]);
                echo json_encode(!empty($row) ? $row[0] : null, JSON_UNESCAPED_UNICODE);
                break;
            }
            if (isset($_GET['id'])) {
                $category = getData("SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

                if (empty($category)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Kategória nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return;
                }
                echo json_encode($category[0], JSON_UNESCAPED_UNICODE);
            } else {
                $categories = getData("SELECT id, name FROM category WHERE deleted_at IS NULL ORDER BY name");

                echo json_encode($categories, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            $nameInput = isset($body['name']) ? trim((string)$body['name']) : '';

            if (!empty($body['restore']) && !empty($body['id']) && $nameInput !== '') {
                $deletedRow = getData("SELECT id FROM category WHERE id = ? AND deleted_at IS NOT NULL", 'i', [$body['id']]);

                if (empty($deletedRow)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Nincs visszaállítható kategória ezzel az azonosítóval!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $nameClash = getData("SELECT id FROM category WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

                if (!empty($nameClash)) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív kategória!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $restored = changeData("UPDATE category SET name = ?, deleted_at = NULL WHERE id = ?", 'si', [$nameInput, $body['id']]);

                if (!$restored) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $row = getData("SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
                http_response_code(200);
                echo json_encode($row[0], JSON_UNESCAPED_UNICODE);
                break;
            }

            if ($nameInput === '') {
                http_response_code(400);
                echo json_encode(['message' => 'Hiányzó adat: a kategória neve kötelező!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existing = getData("SELECT id, deleted_at FROM category WHERE LOWER(TRIM(name)) = LOWER(?) ORDER BY (deleted_at IS NULL) DESC LIMIT 1", 's', [$nameInput]);

            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen nevű kategória!'], JSON_UNESCAPED_UNICODE);
                    return;
                } else {
                    http_response_code(409);
                    echo json_encode([
                        'message' => 'Létezik egy korábban törölt kategória ezen a néven. Szeretné visszaállítani?',
                        'restorable' => true,
                        'id' => (int)$existing[0]['id']
                    ], JSON_UNESCAPED_UNICODE);
                    return;
                }
            }

            $success = changeData("INSERT INTO category (name) VALUES (?)", 's', [$nameInput]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $newCategory = getData("SELECT id, name FROM category WHERE name = ? AND deleted_at IS NULL", 's', [$nameInput]);

            http_response_code(201);
            echo json_encode($newCategory[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító és a név megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $nameInput = trim((string)$body['name']);

            $categoryToUpdate = getData("SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);

            if (empty($categoryToUpdate)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú kategória, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existingName = getData("SELECT id FROM category WHERE LOWER(TRIM(name)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$nameInput, $body['id']]);

            if (!empty($existingName)) {
                http_response_code(409);
                echo json_encode(['message' => 'Már létezik ilyen nevű kategória!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE category SET name = ? WHERE id = ?", 'si', [$nameInput, $body['id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $updatedCategory = getData("SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedCategory[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $categoryToDelete = getData("SELECT id FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (empty($categoryToDelete)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú kategória, vagy már törölve lett!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $subCategories = getData("SELECT id FROM sub_category WHERE category_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (!empty($subCategories)) {
                http_response_code(409);
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ez a kategória aktív alkategóriá(ka)t tartalmaz. Előbb törölje az alkategóriákat.'
                ], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE category SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);

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
