<?php
function handleStoringCondition($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['id'])) {
                $condition = getData("SELECT id, description FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
                
                if (empty($condition)) {
                    echo json_encode(['message' => 'Tárolási körülmény nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(404);
                }
                echo json_encode($condition[0], JSON_UNESCAPED_UNICODE);
            } else {
                $conditions = getData("SELECT id, description FROM storing_condition WHERE deleted_at IS NULL ORDER BY description");
                echo json_encode($conditions, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            if (empty($body['description'])) {
                echo json_encode(['message' => 'Hiányzó adat: leírás kötelező!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $existing = getData("SELECT id, deleted_at FROM storing_condition WHERE description = ?", 's', [$body['description']]);
            
            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    echo json_encode(['message' => 'Már létezik ilyen leírású aktív tárolási körülmény!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                } else {
                    echo json_encode(['message' => 'Már létezik ilyen tárolási körülmény, de törölve lett.'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                }
            }

            $success = changeData("INSERT INTO storing_condition (description) VALUES (?)", 's', [$body['description']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $newCondition = getData("SELECT id, description FROM storing_condition WHERE description = ? AND deleted_at IS NULL", 's', [$body['description']]);
            
            echo json_encode($newCondition[0], JSON_UNESCAPED_UNICODE);
            http_response_code(201);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['description'])) {
                echo json_encode(['message' => "Hiányzó adat!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $condToUpdate = getData("SELECT id FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            
            if (empty($condToUpdate)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú tárolási körülmény, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $existingDesc = getData("SELECT id FROM storing_condition WHERE description = ? AND id != ?", 'si', [$body['description'], $body['id']]);
            
            if (!empty($existingDesc)) {
                echo json_encode(['message' => 'Ez a leírás már foglalt!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE storing_condition SET description = ? WHERE id = ?", 'si', [$body['description'], $body['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $updatedCond = getData("SELECT id, description FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedCond[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                echo json_encode(['message' => "Hiányzó adat: 'id'!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $condToDelete = getData("SELECT id FROM storing_condition WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (empty($condToDelete)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú tárolási körülmény!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $types = getData("SELECT id FROM product_type WHERE storing_condition_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (!empty($types)) {
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ehhez a körülményhez aktív terméktípusok tartoznak.'
                ], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE storing_condition SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }
            
            http_response_code(204);
            break;
        default:
            return http_response_code(405);
    }
}