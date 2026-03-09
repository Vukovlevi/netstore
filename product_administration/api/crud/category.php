<?php
function handleCategory($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['id'])) {
                $category = getData("SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
                
                if (empty($category)) {
                    echo json_encode(['message' => 'Kategória nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(404);
                }
                echo json_encode($category[0], JSON_UNESCAPED_UNICODE);
            } else {
                $categories = getData("SELECT id, name FROM category WHERE deleted_at IS NULL ORDER BY name");
                
                echo json_encode($categories, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            if (empty($body['name'])) {
                echo json_encode(['message' => 'Hiányzó adat!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $existing = getData("SELECT id, deleted_at FROM category WHERE name = ?", 's', [$body['name']]);
            
            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív kategória!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                } else {
                    echo json_encode(['message' => 'Már létezik ilyen nevű kategória, de törölve lett. Kérem, használjon más nevet!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                }
            }

            $success = changeData("INSERT INTO category (name) VALUES (?)", 's', [$body['name']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $newCategory = getData("SELECT id, name FROM category WHERE name = ? AND deleted_at IS NULL", 's', [$body['name']]);
            
            echo json_encode($newCategory[0], JSON_UNESCAPED_UNICODE);
            http_response_code(201);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name'])) {
                echo json_encode(['message' => "Hiányzó adat: 'id' és 'name' megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $categoryToUpdate = getData("SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            
            if (empty($categoryToUpdate)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú kategória, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $existingName = getData("SELECT id FROM category WHERE name = ? AND id != ?", 'si', [$body['name'], $body['id']]);
            
            if (!empty($existingName)) {
                echo json_encode(['message' => 'Ez a név már foglalt egy másik (esetleg törölt) kategória által!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE category SET name = ? WHERE id = ?", 'si', [$body['name'], $body['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $updatedCategory = getData("SELECT id, name FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedCategory[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                echo json_encode(['message' => "Hiányzó adat: 'id' query paraméter megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $categoryToDelete = getData("SELECT id FROM category WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (empty($categoryToDelete)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú kategória, vagy már törölve lett!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $subCategories = getData("SELECT id FROM sub_category WHERE category_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (!empty($subCategories)) {
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ez a kategória aktív alkategóriá(ka)t tartalmaz. Előbb törölje az alkategóriákat.'
                ], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE category SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);
            
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