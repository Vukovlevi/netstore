<?php
function handleStoringCondition($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['check_deleted'])) {
                $needle = trim((string)$_GET['check_deleted']);
                if ($needle === '') {
                    echo json_encode(null, JSON_UNESCAPED_UNICODE);
                    break;
                }
                $row = getData("SELECT id, description, deleted_at FROM storing_condition WHERE LOWER(TRIM(description)) = LOWER(?) AND deleted_at IS NOT NULL ORDER BY id DESC LIMIT 1", 's', [$needle]);
                echo json_encode(!empty($row) ? $row[0] : null, JSON_UNESCAPED_UNICODE);
                break;
            }
            if (isset($_GET['id'])) {
                $condition = getData("SELECT id, description FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

                if (empty($condition)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Tárolási körülmény nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return;
                }
                echo json_encode($condition[0], JSON_UNESCAPED_UNICODE);
            } else {
                $conditions = getData("SELECT id, description FROM storing_condition WHERE deleted_at IS NULL ORDER BY description");
                echo json_encode($conditions, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            $descInput = isset($body['description']) ? trim((string)$body['description']) : '';

            if (!empty($body['restore']) && !empty($body['id']) && $descInput !== '') {
                $deletedRow = getData("SELECT id FROM storing_condition WHERE id = ? AND deleted_at IS NOT NULL", 'i', [$body['id']]);

                if (empty($deletedRow)) {
                    http_response_code(404);
                    echo json_encode(['message' => 'Nincs visszaállítható tárolási körülmény ezzel az azonosítóval!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $clash = getData("SELECT id FROM storing_condition WHERE LOWER(TRIM(description)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$descInput, $body['id']]);

                if (!empty($clash)) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen leírású aktív tárolási körülmény!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $restored = changeData("UPDATE storing_condition SET description = ?, deleted_at = NULL WHERE id = ?", 'si', [$descInput, $body['id']]);

                if (!$restored) {
                    http_response_code(500);
                    echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                    return;
                }

                $row = getData("SELECT id, description FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
                http_response_code(200);
                echo json_encode($row[0], JSON_UNESCAPED_UNICODE);
                break;
            }

            if ($descInput === '') {
                http_response_code(400);
                echo json_encode(['message' => 'Hiányzó adat: a leírás megadása kötelező!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existing = getData("SELECT id, deleted_at FROM storing_condition WHERE LOWER(TRIM(description)) = LOWER(?) ORDER BY (deleted_at IS NULL) DESC LIMIT 1", 's', [$descInput]);

            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    http_response_code(409);
                    echo json_encode(['message' => 'Már létezik ilyen leírású tárolási körülmény!'], JSON_UNESCAPED_UNICODE);
                    return;
                } else {
                    http_response_code(409);
                    echo json_encode([
                        'message' => 'Létezik egy korábban törölt tárolási körülmény ezen a leírással. Szeretné visszaállítani?',
                        'restorable' => true,
                        'id' => (int)$existing[0]['id']
                    ], JSON_UNESCAPED_UNICODE);
                    return;
                }
            }

            $success = changeData("INSERT INTO storing_condition (description) VALUES (?)", 's', [$descInput]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $newCondition = getData("SELECT id, description FROM storing_condition WHERE description = ? AND deleted_at IS NULL", 's', [$descInput]);

            http_response_code(201);
            echo json_encode($newCondition[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['description'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító és a leírás megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $descInput = trim((string)$body['description']);

            $condToUpdate = getData("SELECT id FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);

            if (empty($condToUpdate)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú tárolási körülmény, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $existingDesc = getData("SELECT id FROM storing_condition WHERE LOWER(TRIM(description)) = LOWER(?) AND id != ? AND deleted_at IS NULL", 'si', [$descInput, $body['id']]);

            if (!empty($existingDesc)) {
                http_response_code(409);
                echo json_encode(['message' => 'Már létezik ilyen leírású tárolási körülmény!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE storing_condition SET description = ? WHERE id = ?", 'si', [$descInput, $body['id']]);

            if (!$success) {
                http_response_code(500);
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $updatedCond = getData("SELECT id, description FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedCond[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                http_response_code(400);
                echo json_encode(['message' => "Hiányzó adat: az azonosító megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return;
            }

            $condToDelete = getData("SELECT id FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (empty($condToDelete)) {
                http_response_code(404);
                echo json_encode(['message' => 'Nincs ilyen azonosítójú tárolási körülmény!'], JSON_UNESCAPED_UNICODE);
                return;
            }

            $types = getData("SELECT id FROM product_type WHERE storing_condition_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);

            if (!empty($types)) {
                http_response_code(409);
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ehhez a körülményhez aktív terméktípusok tartoznak.'
                ], JSON_UNESCAPED_UNICODE);
                return;
            }

            $success = changeData("UPDATE storing_condition SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);

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
