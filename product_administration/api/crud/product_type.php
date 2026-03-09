<?php
function handleProductType($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['id'])) {
                $type = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
                
                if (empty($type)) {
                    echo json_encode(['message' => 'Terméktípus nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(404);
                }
                echo json_encode($type[0], JSON_UNESCAPED_UNICODE);
            } else {
                $types = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE deleted_at IS NULL ORDER BY name");
                
                echo json_encode($types, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            if (empty($body['name']) || empty($body['description']) || empty($body['sub_id']) || empty($body['storing_condition_id'])) {
                echo json_encode(['message' => "Hiányzó adat: 'name', 'description', 'sub_id' és 'storing_condition_id' megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $existing = getData("SELECT id, deleted_at FROM product_type WHERE name = ?", 's', [$body['name']]);
            
            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív terméktípus!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                } else {
                    echo json_encode(['message' => 'Már létezik ilyen nevű terméktípus, de törölve lett. Kérem, használjon más nevet!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                }
            }

            $success = changeData("INSERT INTO product_type (name, description, sub_id, storing_condition_id) VALUES (?, ?, ?, ?)", 'ssii', [$body['name'], $body['description'], $body['sub_id'], $body['storing_condition_id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $newType = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE name = ? AND deleted_at IS NULL", 's', [$body['name']]);
            
            echo json_encode($newType[0], JSON_UNESCAPED_UNICODE);
            http_response_code(201);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name']) || empty($body['description']) || empty($body['sub_id']) || empty($body['storing_condition_id'])) {
                echo json_encode(['message' => "Hiányzó adat: 'id', 'name', 'description', 'sub_id' és 'storing_condition_id' megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $typeToUpdate = getData("SELECT id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            
            if (empty($typeToUpdate)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú terméktípus, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $existingName = getData("SELECT id FROM product_type WHERE name = ? AND id != ?", 'si', [$body['name'], $body['id']]);
            
            if (!empty($existingName)) {
                echo json_encode(['message' => 'Ez a név már foglalt egy másik (esetleg törölt) terméktípus által!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE product_type SET name = ?, description = ?, sub_id = ?, storing_condition_id = ? WHERE id = ?", 'ssiii', [$body['name'], $body['description'], $body['sub_id'], $body['storing_condition_id'], $body['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $updatedType = getData("SELECT id, name, description, sub_id, storing_condition_id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedType[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                echo json_encode(['message' => "Hiányzó adat: 'id' query paraméter megadása kötelező!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $typeToDelete = getData("SELECT id FROM product_type WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (empty($typeToDelete)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú terméktípus, vagy már törölve lett!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $products = getData("SELECT id FROM product WHERE type_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (!empty($products)) {
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ez a terméktípus aktív terméke(ke)t tartalmaz. Előbb törölje a termékeket.'
                ], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE product_type SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);
            
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