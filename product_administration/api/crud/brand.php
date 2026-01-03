<?php
function handleBrand($method, $body) {
    switch ($method) {
        case 'GET':
            if (isset($_GET['id'])) {
                $brand = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
                
                if (empty($brand)) {
                    echo json_encode(['message' => 'Márka nem található vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(404);
                }
                echo json_encode($brand[0], JSON_UNESCAPED_UNICODE);
            } else {
                $brands = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE deleted_at IS NULL ORDER BY name");
                echo json_encode($brands, JSON_UNESCAPED_UNICODE);
            }
            break;
        case 'POST':
            if (empty($body['name'])) {
                echo json_encode(['message' => 'Hiányzó adat: név kötelező!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $existing = getData("SELECT id, deleted_at FROM brand WHERE name = ?", 's', [$body['name']]);
            
            if (!empty($existing)) {
                if (is_null($existing[0]['deleted_at'])) {
                    echo json_encode(['message' => 'Már létezik ilyen nevű aktív márka!'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                } else {
                    echo json_encode(['message' => 'Már létezik ilyen márka, de törölve lett.'], JSON_UNESCAPED_UNICODE);
                    return http_response_code(409);
                }
            }

            $isOwn = isset($body['is_own']) ? (int)$body['is_own'] : 0;
            $isTemporary = isset($body['is_temporary']) ? (int)$body['is_temporary'] : 0;

            $success = changeData("INSERT INTO brand (name, is_own, is_temporary) VALUES (?, ?, ?)", 'sii', [$body['name'], $isOwn, $isTemporary]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $newBrand = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE name = ? AND deleted_at IS NULL", 's', [$body['name']]);
            
            echo json_encode($newBrand[0], JSON_UNESCAPED_UNICODE);
            http_response_code(201);
            break;
        case 'PUT':
            if (empty($body['id']) || empty($body['name'])) {
                echo json_encode(['message' => "Hiányzó adat!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $brandToUpdate = getData("SELECT id FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            
            if (empty($brandToUpdate)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú márka, vagy törölve lett!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $existingName = getData("SELECT id FROM brand WHERE name = ? AND id != ?", 'si', [$body['name'], $body['id']]);
            
            if (!empty($existingName)) {
                echo json_encode(['message' => 'Ez a név már foglalt!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $isOwn = isset($body['is_own']) ? (int)$body['is_own'] : 0;
            $isTemporary = isset($body['is_temporary']) ? (int)$body['is_temporary'] : 0;

            $success = changeData("UPDATE brand SET name = ?, is_own = ?, is_temporary = ? WHERE id = ?", 'siii', [$body['name'], $isOwn, $isTemporary, $body['id']]);
            
            if (!$success) {
                echo json_encode(['message' => 'Sikertelen művelet, adatbázis hiba!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(500);
            }

            $updatedBrand = getData("SELECT id, name, is_own, is_temporary FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$body['id']]);
            echo json_encode($updatedBrand[0], JSON_UNESCAPED_UNICODE);
            break;
        case 'DELETE':
            if (empty($_GET['id'])) {
                echo json_encode(['message' => "Hiányzó adat: 'id'!"], JSON_UNESCAPED_UNICODE);
                return http_response_code(400);
            }

            $brandToDelete = getData("SELECT id FROM brand WHERE id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (empty($brandToDelete)) {
                echo json_encode(['message' => 'Nincs ilyen azonosítójú márka!'], JSON_UNESCAPED_UNICODE);
                return http_response_code(404);
            }

            $products = getData("SELECT id FROM product WHERE brand_id = ? AND deleted_at IS NULL", 'i', [$_GET['id']]);
            
            if (!empty($products)) {
                echo json_encode([
                    'message' => 'Törlés sikertelen: Ehhez a márkához aktív termékek tartoznak.'
                ], JSON_UNESCAPED_UNICODE);
                return http_response_code(409);
            }

            $success = changeData("UPDATE brand SET deleted_at = CURDATE() WHERE id = ?", 'i', [$_GET['id']]);
            
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
